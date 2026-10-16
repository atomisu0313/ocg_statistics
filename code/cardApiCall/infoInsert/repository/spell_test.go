package repository_test

import (
	"context"
	"testing"

	"atomisu.com/ocg-statics/infoInsert/dto/cardrecord"
	"atomisu.com/ocg-statics/infoInsert/repository"
	"atomisu.com/ocg-statics/infoInsert/sqlc_gen"
	"atomisu.com/ocg-statics/infoInsert/transaction"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestForSpell tests the SpellRepository
func TestForSpell(t *testing.T) {

	t.Run("正常系01 魔法カードの新規登録処理", func(t *testing.T) {
		// セットアップ
		dbConn, card, cleanup := setupTest(t)
		defer cleanup()

		// Test data
		spellTypeID := int32(1) // 通常魔法のIDを便宜上セット

		// トランザクションの整備
		ctx := context.Background()
		tr := transaction.NewTx(dbConn.DB)

		// トランザクション境界の中で実行(useCaseではこの中にbaseCard挿入処理を入れる)
		insertedSpell, err := transaction.ExecTx(ctx, tr, func(q *sqlc_gen.Queries) (sqlc_gen.Spell, error) {
			spellRepo := repository.NewSpellRepository(q)
			return spellRepo.InsertSpell(ctx, card.ID, spellTypeID)
		})
		require.NoError(t, err, "Transaction should execute without error")

		// Verification
		assert.NotZero(t, insertedSpell.CardID, "Inserted spell should have a non-zero card ID")
		assert.Equal(t, spellTypeID, insertedSpell.SpellTypeID.Int32, "The spell's type ID should match the input")

		// データの取得
		ctx2 := context.Background()
		fetchedSpell, err := transaction.ExecTx(ctx2, tr, func(q *sqlc_gen.Queries) (cardrecord.SpellCardSelectResult, error) {
			spellRepo := repository.NewSpellRepository(q)
			return spellRepo.GetSpellByCardID(ctx2, insertedSpell.CardID)
		})

		assert.Equal(t, insertedSpell.CardID, fetchedSpell.ID, "Fetched spell card ID should match the inserted one")
		assert.Equal(t, "通常", fetchedSpell.SpellTypeNameJa, "Fetched spell type ID should match the inserted one")
	})
}
