package neon

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"atomisu.com/ocg-statics/cardApiCall/internal/dto/cardrecord"
	"atomisu.com/ocg-statics/cardApiCall/internal/repository"
	"atomisu.com/ocg-statics/cardApiCall/internal/sqlc_gen"
	"atomisu.com/ocg-statics/cardApiCall/internal/transaction"
)

func checkStandardCardIsMonster(cardInfo cardrecord.StandardCard) bool {
	return strings.Contains(cardInfo.Type, "Monster")
}

func checkStandardCardIsTrap(cardInfo cardrecord.StandardCard) bool {
	return strings.Contains(cardInfo.Type, "Trap")
}

func checkStandardCardIsSpell(cardInfo cardrecord.StandardCard) bool {
	return strings.Contains(cardInfo.Type, "Spell")
}

// StandardCardを引数として、魔法・罠・モンスターを判定して適切なテーブルに挿入する
func (n *neonUseCaseImpl) InsertCardInfo(ctx context.Context, cardInfo cardrecord.StandardCard) (int64, error) {
	var cardID int64
	var err error

	switch {
	case checkStandardCardIsMonster(cardInfo):
		cardID, err = n.InsertMonsterCardInfo(ctx, cardInfo)
	case checkStandardCardIsTrap(cardInfo):
		cardID, err = n.InsertTrapCardInfo(ctx, cardInfo)
	case checkStandardCardIsSpell(cardInfo):
		cardID, err = n.InsertSpellCardInfo(ctx, cardInfo)
	default:
		return 0, fmt.Errorf("invalid card type")
	}

	if err != nil {
		return 0, err
	}

	// 種別によらず共通の画像情報（card_images）を登録
	if err := n.InsertCardImageInfo(ctx, cardID, cardInfo.TcgID, cardInfo.CardImages); err != nil {
		return 0, fmt.Errorf("error insert card image info: %w", err)
	}

	return cardID, nil
}

// InsertCardImageInfo はカードに紐づく画像一覧を card_images テーブルに登録します。
func (n *neonUseCaseImpl) InsertCardImageInfo(ctx context.Context, cardID int64, ocgApiCardID int64, images []cardrecord.CardImage) error {
	if len(images) == 0 {
		return nil
	}

	tr := transaction.NewTx(n.ProduceConnDB())
	_, err := transaction.ExecTx(ctx, tr, func(q *sqlc_gen.Queries) (struct{}, error) {
		repo := repository.NewCardImageRepository(q)
		for _, img := range images {
			_, err := repo.InsertCardImage(ctx, sqlc_gen.InsertCardImageParams{
				OcgApiImageID:      img.ID,
				OcgApiCardID:       ocgApiCardID,
				CardID:             sql.NullInt64{Int64: cardID, Valid: true},
				ImageUrlSmallApi:   img.ImageUrlSmall,
				ImageUrlLargeApi:   img.ImageUrlLarge,
				ImageUrlCroppedApi: img.ImageUrlCropped,
			})
			if err != nil {
				return struct{}{}, err
			}
		}
		return struct{}{}, nil
	})
	return err
}

func (n *neonUseCaseImpl) GetCardPatternByCardID(ctx context.Context, cardID int64) (cardrecord.CardPatternSelectResult, error) {
	tr := transaction.NewTx(n.ProduceConnDB())

	return transaction.ExecTx(ctx, tr, func(q *sqlc_gen.Queries) (cardrecord.CardPatternSelectResult, error) {
		cardRepo := repository.NewCardRepository(q)
		cardPattern, err := cardRepo.GetCardPatternByCardID(ctx, cardID)
		if err != nil {
			return cardrecord.CardPatternSelectResult{}, fmt.Errorf("error get card pattern: %w", err)
		}
		result := cardrecord.CardPatternSelectResult{}
		return *result.FromSelectCardPatternByCardIDRow(cardPattern), nil
	})
}

func (n *neonUseCaseImpl) GetCardImagesByCardID(ctx context.Context, cardID int64) ([]sqlc_gen.CardImage, error) {
	tr := transaction.NewTx(n.ProduceConnDB())
	return transaction.ExecNonTx(ctx, tr, func(q *sqlc_gen.Queries) ([]sqlc_gen.CardImage, error) {
		repo := repository.NewCardImageRepository(q)
		return repo.ListCardImagesByCardID(ctx, cardID)
	})
}
