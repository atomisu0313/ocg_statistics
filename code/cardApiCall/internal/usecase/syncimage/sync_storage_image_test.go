package syncimage_test

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"io"
	"strings"
	"sync"
	"testing"

	"atomisu.com/ocg-statics/cardApiCall/internal/repository"
	"atomisu.com/ocg-statics/cardApiCall/internal/sqlc_gen"
	"atomisu.com/ocg-statics/cardApiCall/internal/usecase/syncimage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// モック: CardImageRepository
type mockCardImageRepository struct {
	repository.Repository
	mu           sync.Mutex
	records      map[int64][]sqlc_gen.CardImage
	updateCalls  []sqlc_gen.UpdateCardImageThisAppUrlsParams
	listErr      error
	updateErr    error
}

func newMockCardImageRepository() *mockCardImageRepository {
	return &mockCardImageRepository{
		records: make(map[int64][]sqlc_gen.CardImage),
	}
}

func (m *mockCardImageRepository) InsertCardImage(ctx context.Context, arg sqlc_gen.InsertCardImageParams) (sqlc_gen.CardImage, error) {
	return sqlc_gen.CardImage{}, nil
}

func (m *mockCardImageRepository) ListCardImagesByCardID(ctx context.Context, cardID int64) ([]sqlc_gen.CardImage, error) {
	return nil, nil
}

func (m *mockCardImageRepository) ListCardImagesByOcgApiCardID(ctx context.Context, ocgApiCardID int64) ([]sqlc_gen.CardImage, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.listErr != nil {
		return nil, m.listErr
	}
	return m.records[ocgApiCardID], nil
}

func (m *mockCardImageRepository) ListPendingSyncCardImages(ctx context.Context, limit int32) ([]sqlc_gen.CardImage, error) {
	return nil, nil
}

func (m *mockCardImageRepository) UpdateCardImageThisAppUrls(ctx context.Context, arg sqlc_gen.UpdateCardImageThisAppUrlsParams) (sqlc_gen.CardImage, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.updateErr != nil {
		return sqlc_gen.CardImage{}, m.updateErr
	}
	m.updateCalls = append(m.updateCalls, arg)
	return sqlc_gen.CardImage{ID: arg.ID}, nil
}

func (m *mockCardImageRepository) DeleteCardImagesByCardID(ctx context.Context, cardID int64) error {
	return nil
}

// モック: StorageClient
type mockStorageClient struct {
	mu           sync.Mutex
	existingKeys map[string]bool
	uploadedKeys map[string][]byte
	headErr      error
	uploadErr    error
}

func newMockStorageClient() *mockStorageClient {
	return &mockStorageClient{
		existingKeys: make(map[string]bool),
		uploadedKeys: make(map[string][]byte),
	}
}

func (m *mockStorageClient) ObjectExists(ctx context.Context, key string) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.headErr != nil {
		return false, m.headErr
	}
	return m.existingKeys[key], nil
}

func (m *mockStorageClient) UploadStream(ctx context.Context, key string, body io.Reader, contentType string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.uploadErr != nil {
		return "", m.uploadErr
	}
	data, err := io.ReadAll(body)
	if err != nil {
		return "", err
	}
	m.uploadedKeys[key] = data
	return m.GetDestinationURL(key), nil
}

func (m *mockStorageClient) GetDestinationURL(key string) string {
	cleanKey := strings.TrimPrefix(key, "/")
	return fmt.Sprintf("https://test-bucket.s3.ap-northeast-1.amazonaws.com/%s", cleanKey)
}

// モック: StreamClient
type mockStreamClient struct {
	mu        sync.Mutex
	responses map[string][]byte
	streamErr error
}

func newMockStreamClient() *mockStreamClient {
	return &mockStreamClient{
		responses: make(map[string][]byte),
	}
}

func (m *mockStreamClient) GetStream(ctx context.Context, url string) (io.ReadCloser, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.streamErr != nil {
		return nil, m.streamErr
	}
	data, ok := m.responses[url]
	if !ok {
		// デフォルトでダミーバイト列を返す
		data = []byte("fake-jpeg-data")
	}
	return io.NopCloser(bytes.NewReader(data)), nil
}

func TestSyncStorageImagesByCardIDs(t *testing.T) {
	ctx := context.Background()

	t.Run("空のカードIDリストを渡した場合はエラーを返す", func(t *testing.T) {
		repo := newMockCardImageRepository()
		storage := newMockStorageClient()
		stream := newMockStreamClient()
		uc := syncimage.NewSyncStorageImageUseCase(nil, repo, storage, stream)

		output, err := uc.SyncStorageImagesByCardIDs(ctx, []int64{})
		assert.Error(t, err)
		assert.Nil(t, output)
	})

	t.Run("DBに対象カードレコードが存在しない場合は failedOcgCardIds に記録される", func(t *testing.T) {
		repo := newMockCardImageRepository()
		storage := newMockStorageClient()
		stream := newMockStreamClient()
		uc := syncimage.NewSyncStorageImageUseCase(nil, repo, storage, stream)

		output, err := uc.SyncStorageImagesByCardIDs(ctx, []int64{99999999})
		require.NoError(t, err)
		assert.Equal(t, 1, output.TotalProcessedCards)
		assert.Empty(t, output.SuccessfulOcgCardIds)
		assert.Equal(t, []int64{99999999}, output.FailedOcgCardIds)
		assert.Equal(t, 0, output.SyncedImageCount)
		assert.Equal(t, 0, output.SkippedImageCount)
		assert.Len(t, output.Errors, 1)
		assert.Contains(t, output.Errors[0].Message, "not found")
	})

	t.Run("新規同期: S3に未存在でDB未同期の画像がストリーミング同期されDB更新される", func(t *testing.T) {
		repo := newMockCardImageRepository()
		storage := newMockStorageClient()
		stream := newMockStreamClient()
		uc := syncimage.NewSyncStorageImageUseCase(nil, repo, storage, stream)

		cardID := int64(89631139)
		imageID := int64(89631139)
		repo.records[cardID] = []sqlc_gen.CardImage{
			{
				ID:                 1,
				OcgApiCardID:       cardID,
				OcgApiImageID:      imageID,
				ImageUrlSmallApi:   "https://images.ygoprodeck.com/images/cards_small/89631139.jpg",
				ImageUrlLargeApi:   "https://images.ygoprodeck.com/images/cards/89631139.jpg",
				ImageUrlCroppedApi: "https://images.ygoprodeck.com/images/cards_cropped/89631139.jpg",
			},
		}

		output, err := uc.SyncStorageImagesByCardIDs(ctx, []int64{cardID})
		require.NoError(t, err)

		assert.Equal(t, 1, output.TotalProcessedCards)
		assert.Equal(t, []int64{cardID}, output.SuccessfulOcgCardIds)
		assert.Empty(t, output.FailedOcgCardIds)
		assert.Equal(t, 3, output.SyncedImageCount)
		assert.Equal(t, 0, output.SkippedImageCount)
		assert.Empty(t, output.Errors)

		// 3種類のバリエーションがS3にアップロードされていること
		assert.Contains(t, storage.uploadedKeys, fmt.Sprintf("cards/small/%d.jpg", imageID))
		assert.Contains(t, storage.uploadedKeys, fmt.Sprintf("cards/large/%d.jpg", imageID))
		assert.Contains(t, storage.uploadedKeys, fmt.Sprintf("cards/cropped/%d.jpg", imageID))

		// DB更新が実行されていること
		require.Len(t, repo.updateCalls, 1)
		call := repo.updateCalls[0]
		assert.Equal(t, int64(1), call.ID)
		assert.True(t, call.ImageUrlSmallThisApp.Valid)
		assert.True(t, call.ImageUrlLargeThisApp.Valid)
		assert.True(t, call.ImageUrlCroppedThisApp.Valid)
		assert.Contains(t, call.ImageUrlSmallThisApp.String, "cards/small/89631139.jpg")
		assert.Contains(t, call.ImageUrlLargeThisApp.String, "cards/large/89631139.jpg")
		assert.Contains(t, call.ImageUrlCroppedThisApp.String, "cards/cropped/89631139.jpg")
	})

	t.Run("スキップ判定: S3に存在しDBも更新済みの場合はスキップされる", func(t *testing.T) {
		repo := newMockCardImageRepository()
		storage := newMockStorageClient()
		stream := newMockStreamClient()
		uc := syncimage.NewSyncStorageImageUseCase(nil, repo, storage, stream)

		cardID := int64(89631139)
		imageID := int64(89631139)

		// S3に既に存在
		storage.existingKeys[fmt.Sprintf("cards/small/%d.jpg", imageID)] = true
		storage.existingKeys[fmt.Sprintf("cards/large/%d.jpg", imageID)] = true
		storage.existingKeys[fmt.Sprintf("cards/cropped/%d.jpg", imageID)] = true

		// DBにも既に設定済み
		repo.records[cardID] = []sqlc_gen.CardImage{
			{
				ID:                     1,
				OcgApiCardID:           cardID,
				OcgApiImageID:          imageID,
				ImageUrlSmallApi:       "https://images.ygoprodeck.com/images/cards_small/89631139.jpg",
				ImageUrlLargeApi:       "https://images.ygoprodeck.com/images/cards/89631139.jpg",
				ImageUrlCroppedApi:     "https://images.ygoprodeck.com/images/cards_cropped/89631139.jpg",
				ImageUrlSmallThisApp:   sql.NullString{String: "https://test-bucket.s3.ap-northeast-1.amazonaws.com/cards/small/89631139.jpg", Valid: true},
				ImageUrlLargeThisApp:   sql.NullString{String: "https://test-bucket.s3.ap-northeast-1.amazonaws.com/cards/large/89631139.jpg", Valid: true},
				ImageUrlCroppedThisApp: sql.NullString{String: "https://test-bucket.s3.ap-northeast-1.amazonaws.com/cards/cropped/89631139.jpg", Valid: true},
			},
		}

		output, err := uc.SyncStorageImagesByCardIDs(ctx, []int64{cardID})
		require.NoError(t, err)

		assert.Equal(t, 1, output.TotalProcessedCards)
		assert.Equal(t, []int64{cardID}, output.SuccessfulOcgCardIds)
		assert.Empty(t, output.FailedOcgCardIds)
		assert.Equal(t, 0, output.SyncedImageCount)
		assert.Equal(t, 3, output.SkippedImageCount)
		assert.Empty(t, storage.uploadedKeys)
	})

	t.Run("部分失敗: 1枚のカードが失敗しても残りのカードの処理は継続される", func(t *testing.T) {
		repo := newMockCardImageRepository()
		storage := newMockStorageClient()
		stream := newMockStreamClient()
		uc := syncimage.NewSyncStorageImageUseCase(nil, repo, storage, stream)

		cardSuccess := int64(89631139)
		cardFail := int64(46986414)

		repo.records[cardSuccess] = []sqlc_gen.CardImage{
			{
				ID:                 1,
				OcgApiCardID:       cardSuccess,
				OcgApiImageID:      cardSuccess,
				ImageUrlSmallApi:   "https://images.ygoprodeck.com/images/cards_small/89631139.jpg",
				ImageUrlLargeApi:   "https://images.ygoprodeck.com/images/cards/89631139.jpg",
				ImageUrlCroppedApi: "https://images.ygoprodeck.com/images/cards_cropped/89631139.jpg",
			},
		}

		repo.records[cardFail] = []sqlc_gen.CardImage{
			{
				ID:                 2,
				OcgApiCardID:       cardFail,
				OcgApiImageID:      cardFail,
				ImageUrlSmallApi:   "https://fail-url.com/small.jpg",
				ImageUrlLargeApi:   "https://fail-url.com/large.jpg",
				ImageUrlCroppedApi: "https://fail-url.com/cropped.jpg",
			},
		}

		// 特定URLでストリーミング失敗をシミュレート
		stream.streamErr = fmt.Errorf("connection timeout")

		output, err := uc.SyncStorageImagesByCardIDs(ctx, []int64{cardFail, cardSuccess})
		require.NoError(t, err)

		assert.Equal(t, 2, output.TotalProcessedCards)
		assert.Equal(t, 2, len(output.FailedOcgCardIds))
		assert.Contains(t, output.FailedOcgCardIds, cardFail)
		assert.NotEmpty(t, output.Errors)
	})
}
