package repository_test

import (
	"context"
	"database/sql"
	"testing"

	"atomisu.com/ocg-statics/infoInsert/config"
	"atomisu.com/ocg-statics/infoInsert/repository"
	"atomisu.com/ocg-statics/infoInsert/sqlc_gen"
	"github.com/samber/do"
	"github.com/stretchr/testify/assert"
)

func TestCardImageRepository(t *testing.T) {
	t.Run("カード画像の登録と各種取得・更新・削除のテスト", func(t *testing.T) {
		config.BeforeEachForUnitTest()
		defer config.AfterEachForUnitTest()

		injector := do.New()
		do.Provide(injector, config.TestDbConnection)
		dbConn := do.MustInvoke[*config.DbConn](injector)

		ctx := context.Background()
		q := sqlc_gen.New(dbConn)
		cardRepo := repository.NewCardRepository(q)
		cardImageRepo := repository.NewCardImageRepository(q)

		// 1. 親カードを作成
		card, err := cardRepo.InsertCard(ctx, sqlc_gen.InsertCardParams{
			NameEn:     sql.NullString{String: "Dark Magician", Valid: true},
			NameJa:     sql.NullString{String: "ブラック・マジシャン", Valid: true},
			CardTextJa: sql.NullString{String: "魔法使いとしては、攻撃力・守備力ともに最高クラス。", Valid: true},
			CardTextEn: sql.NullString{String: "The ultimate wizard in terms of attack and defense.", Valid: true},
			NeuronID:   sql.NullInt64{Int64: 4005, Valid: true},
			OcgApiID:   sql.NullInt64{Int64: 46986414, Valid: true},
		})
		assert.NoError(t, err)

		// 2. カード画像の挿入（絵違い2枚）
		img1, err := cardImageRepo.InsertCardImage(ctx, sqlc_gen.InsertCardImageParams{
			OcgApiImageID:      46986414,
			OcgApiCardID:       46986414,
			CardID:             sql.NullInt64{Int64: card.ID, Valid: true},
			ImageUrlSmallApi:   "https://images.ygoprodeck.com/images/cards_small/46986414.jpg",
			ImageUrlLargeApi:   "https://images.ygoprodeck.com/images/cards/46986414.jpg",
			ImageUrlCroppedApi: "https://images.ygoprodeck.com/images/cards_cropped/46986414.jpg",
		})
		assert.NoError(t, err)
		assert.NotZero(t, img1.ID)
		assert.Equal(t, int64(46986414), img1.OcgApiImageID)
		assert.False(t, img1.ImageUrlLargeThisApp.Valid) // 初期はNULL

		img2, err := cardImageRepo.InsertCardImage(ctx, sqlc_gen.InsertCardImageParams{
			OcgApiImageID:      46986415, // 別イラスト
			OcgApiCardID:       46986414,
			CardID:             sql.NullInt64{Int64: card.ID, Valid: true},
			ImageUrlSmallApi:   "https://images.ygoprodeck.com/images/cards_small/46986415.jpg",
			ImageUrlLargeApi:   "https://images.ygoprodeck.com/images/cards/46986415.jpg",
			ImageUrlCroppedApi: "https://images.ygoprodeck.com/images/cards_cropped/46986415.jpg",
		})
		assert.NoError(t, err)
		assert.NotZero(t, img2.ID)

		// 3. card_id で一覧取得（2枚取得できること）
		imagesByCardID, err := cardImageRepo.ListCardImagesByCardID(ctx, card.ID)
		assert.NoError(t, err)
		assert.Len(t, imagesByCardID, 2)

		// 4. ocg_api_card_id で一覧取得
		imagesByApiID, err := cardImageRepo.ListCardImagesByOcgApiCardID(ctx, 46986414)
		assert.NoError(t, err)
		assert.Len(t, imagesByApiID, 2)

		// 5. 未同期リスト（this_app が NULL のもの）の取得
		pending, err := cardImageRepo.ListPendingSyncCardImages(ctx, 10)
		assert.NoError(t, err)
		assert.Len(t, pending, 2)

		// 6. 自前配信用 URL の更新
		updated, err := cardImageRepo.UpdateCardImageThisAppUrls(ctx, sqlc_gen.UpdateCardImageThisAppUrlsParams{
			ID:                     img1.ID,
			ImageUrlSmallThisApp:   sql.NullString{String: "https://images.example.com/cards/small/46986414.jpg", Valid: true},
			ImageUrlLargeThisApp:   sql.NullString{String: "https://images.example.com/cards/large/46986414.jpg", Valid: true},
			ImageUrlCroppedThisApp: sql.NullString{String: "https://images.example.com/cards/cropped/46986414.jpg", Valid: true},
		})
		assert.NoError(t, err)
		assert.True(t, updated.ImageUrlLargeThisApp.Valid)
		assert.Equal(t, "https://images.example.com/cards/large/46986414.jpg", updated.ImageUrlLargeThisApp.String)

		// 更新後、未同期リストは1枚に減るはず
		pendingAfter, err := cardImageRepo.ListPendingSyncCardImages(ctx, 10)
		assert.NoError(t, err)
		assert.Len(t, pendingAfter, 1)

		// 7. card_id による削除
		err = cardImageRepo.DeleteCardImagesByCardID(ctx, card.ID)
		assert.NoError(t, err)

		imagesAfterDelete, err := cardImageRepo.ListCardImagesByCardID(ctx, card.ID)
		assert.NoError(t, err)
		assert.Empty(t, imagesAfterDelete)
	})
}
