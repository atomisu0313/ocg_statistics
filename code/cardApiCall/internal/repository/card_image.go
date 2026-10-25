package repository

import (
	"context"
	"database/sql"
	"time"

	"atomisu.com/ocg-statics/cardApiCall/internal/sqlc_gen"
	"go.uber.org/zap"
)

type CardImageRepository interface {
	Repository
	InsertCardImage(ctx context.Context, arg sqlc_gen.InsertCardImageParams) (sqlc_gen.CardImage, error)
	ListCardImagesByCardID(ctx context.Context, cardID int64) ([]sqlc_gen.CardImage, error)
	ListCardImagesByNeuronID(ctx context.Context, neuronID int64) ([]sqlc_gen.CardImage, error)
	ListPendingSyncCardImages(ctx context.Context, limit int32) ([]sqlc_gen.CardImage, error)
	UpdateCardImageThisAppUrls(ctx context.Context, arg sqlc_gen.UpdateCardImageThisAppUrlsParams) (sqlc_gen.CardImage, error)
	DeleteCardImagesByCardID(ctx context.Context, cardID int64) error
}

type cardImageRepositoryImpl struct {
	*repository
	queries *sqlc_gen.Queries
}

// NewCardImageRepository creates a new instance of CardImageRepository.
func NewCardImageRepository(q *sqlc_gen.Queries) CardImageRepository {
	return NewRepository(func(r *repository) CardImageRepository {
		return &cardImageRepositoryImpl{
			repository: r,
			queries:    q,
		}
	})
}

func (r *cardImageRepositoryImpl) InsertCardImage(ctx context.Context, arg sqlc_gen.InsertCardImageParams) (sqlc_gen.CardImage, error) {
	start := time.Now()
	defer r.logDBOperation("InsertCardImage", start, zap.Int64("ocg_api_image_id", arg.OcgApiImageID))

	cardImage, err := r.queries.InsertCardImage(ctx, arg)
	if err != nil {
		r.logDBError("InsertCardImage", err, zap.Int64("ocg_api_image_id", arg.OcgApiImageID))
		return sqlc_gen.CardImage{}, err
	}

	r.logDBResult("InsertCardImage", cardImage, zap.Int64("ocg_api_image_id", arg.OcgApiImageID))
	return cardImage, nil
}

func (r *cardImageRepositoryImpl) ListCardImagesByCardID(ctx context.Context, cardID int64) ([]sqlc_gen.CardImage, error) {
	start := time.Now()
	defer r.logDBOperation("ListCardImagesByCardID", start, zap.Int64("card_id", cardID))

	cardImages, err := r.queries.ListCardImagesByCardID(ctx, sql.NullInt64{Int64: cardID, Valid: true})
	if err != nil {
		r.logDBError("ListCardImagesByCardID", err, zap.Int64("card_id", cardID))
		return nil, err
	}

	r.logDBResult("ListCardImagesByCardID", cardImages, zap.Int64("card_id", cardID))
	return cardImages, nil
}

func (r *cardImageRepositoryImpl) ListCardImagesByNeuronID(ctx context.Context, neuronID int64) ([]sqlc_gen.CardImage, error) {
	start := time.Now()
	defer r.logDBOperation("ListCardImagesByNeuronID", start, zap.Int64("neuron_id", neuronID))

	cardImages, err := r.queries.ListCardImagesByNeuronID(ctx, sql.NullInt64{Int64: neuronID, Valid: true})
	if err != nil {
		r.logDBError("ListCardImagesByNeuronID", err, zap.Int64("neuron_id", neuronID))
		return nil, err
	}

	r.logDBResult("ListCardImagesByNeuronID", cardImages, zap.Int64("neuron_id", neuronID))
	return cardImages, nil
}

func (r *cardImageRepositoryImpl) ListPendingSyncCardImages(ctx context.Context, limit int32) ([]sqlc_gen.CardImage, error) {
	start := time.Now()
	defer r.logDBOperation("ListPendingSyncCardImages", start, zap.Int32("limit", limit))

	cardImages, err := r.queries.ListPendingSyncCardImages(ctx, limit)
	if err != nil {
		r.logDBError("ListPendingSyncCardImages", err, zap.Int32("limit", limit))
		return nil, err
	}

	r.logDBResult("ListPendingSyncCardImages", cardImages, zap.Int32("limit", limit))
	return cardImages, nil
}

func (r *cardImageRepositoryImpl) UpdateCardImageThisAppUrls(ctx context.Context, arg sqlc_gen.UpdateCardImageThisAppUrlsParams) (sqlc_gen.CardImage, error) {
	start := time.Now()
	defer r.logDBOperation("UpdateCardImageThisAppUrls", start, zap.Int64("id", arg.ID))

	cardImage, err := r.queries.UpdateCardImageThisAppUrls(ctx, arg)
	if err != nil {
		r.logDBError("UpdateCardImageThisAppUrls", err, zap.Int64("id", arg.ID))
		return sqlc_gen.CardImage{}, err
	}

	r.logDBResult("UpdateCardImageThisAppUrls", cardImage, zap.Int64("id", arg.ID))
	return cardImage, nil
}

func (r *cardImageRepositoryImpl) DeleteCardImagesByCardID(ctx context.Context, cardID int64) error {
	start := time.Now()
	defer r.logDBOperation("DeleteCardImagesByCardID", start, zap.Int64("card_id", cardID))

	err := r.queries.DeleteCardImagesByCardID(ctx, sql.NullInt64{Int64: cardID, Valid: true})
	if err != nil {
		r.logDBError("DeleteCardImagesByCardID", err, zap.Int64("card_id", cardID))
		return err
	}

	return nil
}
