package syncimage

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"atomisu.com/ocg-statics/cardApiCall/internal/config"
	"atomisu.com/ocg-statics/cardApiCall/internal/dto"
	appHttp "atomisu.com/ocg-statics/cardApiCall/internal/http"
	"atomisu.com/ocg-statics/cardApiCall/internal/repository"
	"atomisu.com/ocg-statics/cardApiCall/internal/sqlc_gen"
	"atomisu.com/ocg-statics/cardApiCall/internal/storage"
	"atomisu.com/ocg-statics/cardApiCall/internal/transaction"
	"github.com/samber/do"
)

// SyncStorageImageUseCase はカード画像の外部ストレージ同期ユースケースのインターフェースです。
type SyncStorageImageUseCase interface {
	SyncStorageImagesByNeuronIDs(ctx context.Context, neuronCardIds []int64) (*dto.SyncImagesOutput, error)
	SyncStorageImagesByCardIDs(ctx context.Context, neuronCardIds []int64) (*dto.SyncImagesOutput, error)
	SyncStorageImages(ctx context.Context, startId int64, delta int64) ([]int64, error)
}

type syncStorageImageUseCaseImpl struct {
	dbConn        *config.DbConn
	cardImageRepo repository.CardImageRepository
	storageClient storage.StorageClient
	streamClient  appHttp.StreamClient
	concurrency   int
}

// NewSyncStorageImageUseCase はユースケースを初期化する純粋コンストラクタです。
func NewSyncStorageImageUseCase(
	dbConn *config.DbConn,
	cardImageRepo repository.CardImageRepository,
	storageClient storage.StorageClient,
	streamClient appHttp.StreamClient,
) SyncStorageImageUseCase {
	return &syncStorageImageUseCaseImpl{
		dbConn:        dbConn,
		cardImageRepo: cardImageRepo,
		storageClient: storageClient,
		streamClient:  streamClient,
		concurrency:   5, // 外部CDNレートリミット対策のためのデフォルト同時並行数
	}
}

// NewSyncStorageImageUseCaseForDI は samber/do 用のプロバイダ関数です。
func NewSyncStorageImageUseCaseForDI(i *do.Injector) (SyncStorageImageUseCase, error) {
	dbConn, err := do.Invoke[*config.DbConn](i)
	if err != nil {
		return nil, err
	}

	var cardImageRepo repository.CardImageRepository
	repo, err := do.Invoke[repository.CardImageRepository](i)
	if err == nil {
		cardImageRepo = repo
	} else {
		cardImageRepo = repository.NewCardImageRepository(sqlc_gen.New(dbConn.DB))
	}

	storageClient, err := do.Invoke[storage.StorageClient](i)
	if err != nil {
		return nil, err
	}

	streamClient, err := do.Invoke[appHttp.StreamClient](i)
	if err != nil {
		return nil, err
	}

	return NewSyncStorageImageUseCase(dbConn, cardImageRepo, storageClient, streamClient), nil
}

// SyncStorageImagesByCardIDs は SyncStorageImagesByNeuronIDs へのエイリアスです。
func (u *syncStorageImageUseCaseImpl) SyncStorageImagesByCardIDs(ctx context.Context, neuronCardIds []int64) (*dto.SyncImagesOutput, error) {
	return u.SyncStorageImagesByNeuronIDs(ctx, neuronCardIds)
}

// SyncStorageImagesByNeuronIDs は対象 Neuron CID のリストを受け取り、画像の突合、S3ストリーミング同期、DB更新を行います。
func (u *syncStorageImageUseCaseImpl) SyncStorageImagesByNeuronIDs(ctx context.Context, neuronCardIds []int64) (*dto.SyncImagesOutput, error) {
	if len(neuronCardIds) == 0 {
		return nil, fmt.Errorf("neuronCardIds must not be empty")
	}

	output := &dto.SyncImagesOutput{
		TotalProcessedCards:     len(neuronCardIds),
		SuccessfulNeuronCardIds: make([]int64, 0),
		FailedNeuronCardIds:     make([]int64, 0),
		SyncedImageCount:        0,
		SkippedImageCount:       0,
		Errors:                  make([]dto.SyncError, 0),
	}

	for _, neuronCardId := range neuronCardIds {
		records, err := u.cardImageRepo.ListCardImagesByNeuronID(ctx, neuronCardId)
		if err != nil {
			output.FailedNeuronCardIds = append(output.FailedNeuronCardIds, neuronCardId)
			output.Errors = append(output.Errors, dto.SyncError{
				NeuronCardID: neuronCardId,
				Message:      fmt.Sprintf("failed to list card images: %v", err),
			})
			continue
		}

		if len(records) == 0 {
			output.FailedNeuronCardIds = append(output.FailedNeuronCardIds, neuronCardId)
			output.Errors = append(output.Errors, dto.SyncError{
				NeuronCardID: neuronCardId,
				Message:      "card images not found in database",
			})
			continue
		}

		// 事前計算と突合
		type recordResolvedUrls struct {
			recordID   int64
			smallURL   string
			largeURL   string
			croppedURL string
		}

		var pendingTasks []dto.ImageSyncTaskDTO
		resolvedRecords := make([]recordResolvedUrls, 0, len(records))
		var cardErrors []dto.SyncError
		cardSkippedCount := 0

		for _, record := range records {
			variants := []struct {
				variant        dto.ImageVariant
				sourceURL      string
				existingAppURL sql.NullString
				s3Key          string
			}{
				{
					variant:        dto.ImageVariantSmall,
					sourceURL:      record.ImageUrlSmallApi,
					existingAppURL: record.ImageUrlSmallThisApp,
					s3Key:          fmt.Sprintf("cards/small/%d.jpg", record.OcgApiImageID),
				},
				{
					variant:        dto.ImageVariantLarge,
					sourceURL:      record.ImageUrlLargeApi,
					existingAppURL: record.ImageUrlLargeThisApp,
					s3Key:          fmt.Sprintf("cards/large/%d.jpg", record.OcgApiImageID),
				},
				{
					variant:        dto.ImageVariantCropped,
					sourceURL:      record.ImageUrlCroppedApi,
					existingAppURL: record.ImageUrlCroppedThisApp,
					s3Key:          fmt.Sprintf("cards/cropped/%d.jpg", record.OcgApiImageID),
				},
			}

			var resolvedSmall, resolvedLarge, resolvedCropped string

			for _, v := range variants {
				exists, err := u.storageClient.ObjectExists(ctx, v.s3Key)
				if err != nil {
					cardErrors = append(cardErrors, dto.SyncError{
						NeuronCardID:  neuronCardId,
						OcgApiImageID: record.OcgApiImageID,
						Message:       fmt.Sprintf("failed to check existence for %s: %v", v.s3Key, err),
					})
					continue
				}

				destURL := u.storageClient.GetDestinationURL(v.s3Key)

				// S3 に存在し、かつ DB の URL も設定済みの場合はスキップ
				if exists && v.existingAppURL.Valid && v.existingAppURL.String != "" {
					cardSkippedCount++
					destURL = v.existingAppURL.String
				} else {
					task := dto.ImageSyncTaskDTO{
						CardImageID:    record.ID,
						NeuronID:       neuronCardId,
						OcgApiImageID:  record.OcgApiImageID,
						Variant:        v.variant,
						SourceURL:      v.sourceURL,
						S3Key:          v.s3Key,
						DestinationURL: destURL,
					}
					pendingTasks = append(pendingTasks, task)
				}

				switch v.variant {
				case dto.ImageVariantSmall:
					resolvedSmall = destURL
				case dto.ImageVariantLarge:
					resolvedLarge = destURL
				case dto.ImageVariantCropped:
					resolvedCropped = destURL
				}
			}

			resolvedRecords = append(resolvedRecords, recordResolvedUrls{
				recordID:   record.ID,
				smallURL:   resolvedSmall,
				largeURL:   resolvedLarge,
				croppedURL: resolvedCropped,
			})
		}

		if len(cardErrors) > 0 {
			output.FailedNeuronCardIds = append(output.FailedNeuronCardIds, neuronCardId)
			output.Errors = append(output.Errors, cardErrors...)
			continue
		}

		// ダウンロードおよびストリーミングアップロード（トランザクション外）
		cardSyncedCount := 0
		if len(pendingTasks) > 0 {
			sem := make(chan struct{}, u.concurrency)
			var wg sync.WaitGroup
			var mu sync.Mutex

			for _, t := range pendingTasks {
				wg.Add(1)
				go func(task dto.ImageSyncTaskDTO) {
					defer wg.Done()
					sem <- struct{}{}
					defer func() { <-sem }()

					stream, err := u.streamClient.GetStream(ctx, task.SourceURL)
					if err != nil {
						mu.Lock()
						cardErrors = append(cardErrors, dto.SyncError{
							NeuronCardID:  task.NeuronID,
							OcgApiImageID: task.OcgApiImageID,
							Message:       fmt.Sprintf("failed to get stream from %s: %v", task.SourceURL, err),
						})
						mu.Unlock()
						return
					}
					defer stream.Close()

					_, err = u.storageClient.UploadStream(ctx, task.S3Key, stream, "image/jpeg")
					if err != nil {
						mu.Lock()
						cardErrors = append(cardErrors, dto.SyncError{
							NeuronCardID:  task.NeuronID,
							OcgApiImageID: task.OcgApiImageID,
							Message:       fmt.Sprintf("failed to stream image to S3: %v", err),
						})
						mu.Unlock()
						return
					}

					mu.Lock()
					cardSyncedCount++
					mu.Unlock()
				}(t)
			}
			wg.Wait()
		}

		if len(cardErrors) > 0 {
			output.FailedNeuronCardIds = append(output.FailedNeuronCardIds, neuronCardId)
			output.Errors = append(output.Errors, cardErrors...)
			continue
		}

		// DBレコードの更新（トランザクション境界: 外部I/O完了後に実施）
		updateArgs := make([]sqlc_gen.UpdateCardImageThisAppUrlsParams, 0, len(resolvedRecords))
		for _, rr := range resolvedRecords {
			updateArgs = append(updateArgs, sqlc_gen.UpdateCardImageThisAppUrlsParams{
				ID:                     rr.recordID,
				ImageUrlSmallThisApp:   sql.NullString{String: rr.smallURL, Valid: rr.smallURL != ""},
				ImageUrlLargeThisApp:   sql.NullString{String: rr.largeURL, Valid: rr.largeURL != ""},
				ImageUrlCroppedThisApp: sql.NullString{String: rr.croppedURL, Valid: rr.croppedURL != ""},
			})
		}

		var updateErr error
		if u.dbConn != nil && u.dbConn.DB != nil {
			tr := transaction.NewTx(u.dbConn.DB)
			_, updateErr = transaction.ExecTx(ctx, tr, func(q *sqlc_gen.Queries) (struct{}, error) {
				txRepo := repository.NewCardImageRepository(q)
				for _, arg := range updateArgs {
					if _, err := txRepo.UpdateCardImageThisAppUrls(ctx, arg); err != nil {
						return struct{}{}, err
					}
				}
				return struct{}{}, nil
			})
		} else {
			for _, arg := range updateArgs {
				if _, err := u.cardImageRepo.UpdateCardImageThisAppUrls(ctx, arg); err != nil {
					updateErr = err
					break
				}
			}
		}

		if updateErr != nil {
			output.FailedNeuronCardIds = append(output.FailedNeuronCardIds, neuronCardId)
			output.Errors = append(output.Errors, dto.SyncError{
				NeuronCardID: neuronCardId,
				Message:      fmt.Sprintf("failed to update card images in database: %v", updateErr),
			})
			continue
		}

		output.SuccessfulNeuronCardIds = append(output.SuccessfulNeuronCardIds, neuronCardId)
		output.SyncedImageCount += cardSyncedCount
		output.SkippedImageCount += cardSkippedCount
	}

	return output, nil
}

// SyncStorageImages は後方互換性のためのメソッドです。
func (u *syncStorageImageUseCaseImpl) SyncStorageImages(ctx context.Context, startId int64, delta int64) ([]int64, error) {
	ids := make([]int64, 0, delta)
	for i := int64(0); i < delta; i++ {
		ids = append(ids, startId+i)
	}
	out, err := u.SyncStorageImagesByNeuronIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	return out.FailedNeuronCardIds, nil
}

