package repository_test

import (
	"context"
	"testing"

	"atomisu.com/ocg-statics/cardApiCall/internal/config"
	"atomisu.com/ocg-statics/cardApiCall/internal/dto/kind"
	"atomisu.com/ocg-statics/cardApiCall/internal/repository"
	"atomisu.com/ocg-statics/cardApiCall/internal/sqlc_gen"
	"atomisu.com/ocg-statics/cardApiCall/internal/transaction"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

// テスト用の共通処理
func testSpellTypeRetrieval(t *testing.T, dbConn *config.DbConn, target kind.SpellKind, getFunc func(ctx context.Context, repo repository.SpellTypeRepository) (kind.SpellKind, error)) {
	ctx := context.Background()
	tr := transaction.NewTx(dbConn.DB)

	result, err := transaction.ExecTx(ctx, tr, func(q *sqlc_gen.Queries) (kind.SpellKind, error) {
		spellTypeRepo := repository.NewSpellTypeRepository(q)
		return getFunc(ctx, spellTypeRepo)
	})

	assert.NoError(t, err)
	assert.Equal(t, target.ID, result.ID, "IDが一致しません")
	assert.Equal(t, target.NameJa, result.NameJa, "名前（和名）が一致しません")
	assert.Equal(t, target.NameEn, result.NameEn, "名前（英名）が一致しません")
}

// TestForSpellType tests the SpellTypeRepository
func TestForSpellType(t *testing.T) {

	testCases := []struct {
		name   string
		target kind.SpellKind
	}{
		{"通常魔法Enumの取得", kind.SpellTypeNormal},
		{"永続魔法Enumの取得", kind.SpellTypeContinuous},
		{"装備魔法Enumの取得", kind.SpellTypeEquip},
		{"フィールド魔法Enumの取得", kind.SpellTypeField},
		{"速攻魔法Enumの取得", kind.SpellTypeQuickPlay},
		{"儀式魔法Enumの取得", kind.SpellTypeRitual},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// セットアップ
			dbConn, _, cleanup := setupTest(t)
			defer cleanup()

			// 和名で取得テスト
			testSpellTypeRetrieval(t, dbConn, tc.target, func(ctx context.Context, repo repository.SpellTypeRepository) (kind.SpellKind, error) {
				return repo.GetSpellTypeByNameJa(ctx, tc.target.NameJa)
			})

			// 英名で取得テスト
			testSpellTypeRetrieval(t, dbConn, tc.target, func(ctx context.Context, repo repository.SpellTypeRepository) (kind.SpellKind, error) {
				return repo.GetSpellTypeByNameEn(ctx, tc.target.NameEn)
			})

			// IDで取得テスト
			testSpellTypeRetrieval(t, dbConn, tc.target, func(ctx context.Context, repo repository.SpellTypeRepository) (kind.SpellKind, error) {
				return repo.GetSpellTypeById(ctx, tc.target.ID)
			})
		})
	}
}
