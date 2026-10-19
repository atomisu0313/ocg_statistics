package syncimage

import (
	"context"

	"atomisu.com/ocg-statics/cardApiCall/internal/repository"
	commonUsecase "atomisu.com/ocg-statics/cardApiCall/internal/usecase"
	"github.com/samber/do"
)

type SyncStorageImageUseCase interface {
	SyncStorageImages(ctx context.Context, startId int64, delta int64) ([]int64, error)
}

type syncStorageImageUseCaseImpl struct {
	*commonUsecase.UseCaseImpl
}

func NewSyncStorageImageUseCase(i *do.Injector) (SyncStorageImageUseCase, error) {
	return commonUsecase.NewUseCase(i, func(base *commonUsecase.UseCaseImpl) SyncStorageImageUseCase {
		return &syncStorageImageUseCaseImpl{
			UseCaseImpl: base,
		}
	})
}

func (u *syncStorageImageUseCaseImpl) SyncStorageImages(ctx context.Context, startId int64, delta int64) ([]int64, error) {
	// common/repository から card_images テーブルの指定範囲レコードを操作可能
	_ = repository.NewCardImageRepository
	return []int64{}, nil
}
