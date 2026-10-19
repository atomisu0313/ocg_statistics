package neon_test

import (
	"context"
	"testing"

	"atomisu.com/ocg-statics/cardApiCall/cmd/aws_lambda/cardInfoInsert/app"
	"atomisu.com/ocg-statics/cardApiCall/internal/config"
	"atomisu.com/ocg-statics/cardApiCall/internal/dto/cardrecord"
	"atomisu.com/ocg-statics/cardApiCall/internal/usecase/neon"
	"github.com/samber/do"
	"github.com/stretchr/testify/assert"
)

func TestInsertCardInfo(t *testing.T) {
	t.Run("モンスターカードが正しく挿入される", func(t *testing.T) {
		config.BeforeEachForUnitTest()
		defer config.AfterEachForUnitTest()

		injector := app.SetupDIContainer()
		do.Override(injector, config.TestDbConnection)
		neonUseCase := do.MustInvoke[neon.NeonUseCase](injector)

		// isMonsterCardのロジック (TypeLinesに"Monster"が含まれる) を満たすようにテストデータを作成
		sampleData := cardrecord.StandardCard{
			NameEn:    "Mokey Mokey",
			NameJa:    "もけもけ",
			DescEn:    "An outcast angel.",
			DescJa:    "天使のはみだし者",
			NeuronID:  6018,
			TcgID:     27288416,
			Def:       100,
			Atk:       300,
			Type:      "Normal Monster",
			Level:     1,
			Race:      "Fairy",
			Attribute: "LIGHT",
			TypeLines: []string{"Monster", "Normal"}, // "Monster" を追加
			CardImages: []cardrecord.CardImage{
				{
					ID:              27288416,
					ImageUrlSmall:   "https://images.ygoprodeck.com/images/cards_small/27288416.jpg",
					ImageUrlLarge:   "https://images.ygoprodeck.com/images/cards/27288416.jpg",
					ImageUrlCropped: "https://images.ygoprodeck.com/images/cards_cropped/27288416.jpg",
				},
			},
		}

		cardID, err := neonUseCase.InsertCardInfo(context.Background(), sampleData)
		assert.NoError(t, err)
		assert.NotEqual(t, int64(0), cardID)

		// モンスターテーブルに正しく挿入されたか確認
		resultsFull, err := neonUseCase.GetMonsterCardExtendedByID(context.Background(), cardID)
		assert.NoError(t, err)
		assert.NotNil(t, resultsFull)
		resultsGet := resultsFull.MonsterCardSelectResult

		assert.Equal(t, sampleData.NameEn, resultsGet.NameEn)
		assert.Equal(t, sampleData.NameJa, resultsGet.NameJa)
		assert.Equal(t, "天使族", resultsGet.RaceNameJa)
		assert.Equal(t, "光", resultsGet.AttributeNameJa)

		// カード画像テーブルに正しく挿入されたか確認
		images, err := neonUseCase.GetCardImagesByCardID(context.Background(), cardID)
		assert.NoError(t, err)
		assert.Len(t, images, 1)
		assert.Equal(t, int64(27288416), images[0].OcgApiImageID)
		assert.Equal(t, "https://images.ygoprodeck.com/images/cards/27288416.jpg", images[0].ImageUrlLargeApi)

		// 他のテーブルに挿入されていないことを確認
		_, err = neonUseCase.GetSpellCardByID(context.Background(), cardID)
		assert.Error(t, err) // Should be an error (not found)
		_, err = neonUseCase.GetTrapCardByID(context.Background(), cardID)
		assert.Error(t, err) // Should be an error (not found)

		cardPattern, err := neonUseCase.GetCardPatternByCardID(context.Background(), cardID)
		assert.True(t, cardPattern.IsMonster)
		assert.False(t, cardPattern.IsSpell)
		assert.False(t, cardPattern.IsTrap)
	})

	t.Run("魔法カードが正しく挿入される", func(t *testing.T) {
		config.BeforeEachForUnitTest()
		defer config.AfterEachForUnitTest()

		injector := app.SetupDIContainer()
		do.Override(injector, config.TestDbConnection)
		neonUseCase := do.MustInvoke[neon.NeonUseCase](injector)

		sampleData := cardrecord.StandardCard{
			NameEn:   "Raigeki",
			NameJa:   "サンダー・ボルト",
			DescEn:   "Sample Text.",
			DescJa:   "サンプルテキスト。",
			Type:     "Spell Card",
			NeuronID: 4343,
			TcgID:    12580477,
			Race:     "Normal",
		}

		cardID, err := neonUseCase.InsertCardInfo(context.Background(), sampleData)
		assert.NoError(t, err)
		assert.NotEqual(t, int64(0), cardID)

		// 魔法テーブルに正しく挿入されたか確認
		result, err := neonUseCase.GetSpellCardByID(context.Background(), cardID)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, sampleData.NameEn, result.NameEn)
		assert.Equal(t, "通常", result.SpellTypeNameJa)

		cardPattern, err := neonUseCase.GetCardPatternByCardID(context.Background(), cardID)
		assert.False(t, cardPattern.IsMonster)
		assert.True(t, cardPattern.IsSpell)
		assert.False(t, cardPattern.IsTrap)

		// 他のテーブルに挿入されていないことを確認
		_, err = neonUseCase.GetMonsterCardExtendedByID(context.Background(), cardID)
		assert.Error(t, err) // Should be an error (not found)
		_, err = neonUseCase.GetTrapCardByID(context.Background(), cardID)
		assert.Error(t, err) // Should be an error (not found)
	})

	t.Run("罠カードが正しく挿入される", func(t *testing.T) {
		config.BeforeEachForUnitTest()
		defer config.AfterEachForUnitTest()

		injector := app.SetupDIContainer()
		do.Override(injector, config.TestDbConnection)
		neonUseCase := do.MustInvoke[neon.NeonUseCase](injector)

		sampleData := cardrecord.StandardCard{
			NameEn:   "Solemn Judgment",
			NameJa:   "神の宣告",
			DescEn:   "Sample Text.",
			DescJa:   "サンプルテキスト。",
			Type:     "Trap Card",
			NeuronID: 8916,
			TcgID:    41420027,
			Race:     "Counter",
		}

		cardID, err := neonUseCase.InsertCardInfo(context.Background(), sampleData)
		assert.NoError(t, err)
		assert.NotEqual(t, int64(0), cardID)

		// 罠テーブルに正しく挿入されたか確認
		result, err := neonUseCase.GetTrapCardByID(context.Background(), cardID)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, sampleData.NameEn, result.NameEn)
		assert.Equal(t, "カウンター", result.TrapTypeNameJa)

		// 他のテーブルに挿入されていないことを確認
		_, err = neonUseCase.GetMonsterCardExtendedByID(context.Background(), cardID)
		assert.Error(t, err) // Should be an error (not found)
		_, err = neonUseCase.GetSpellCardByID(context.Background(), cardID)
		assert.Error(t, err) // Should be an error (not found)

		cardPattern, err := neonUseCase.GetCardPatternByCardID(context.Background(), cardID)
		assert.False(t, cardPattern.IsMonster)
		assert.False(t, cardPattern.IsSpell)
		assert.True(t, cardPattern.IsTrap)
	})

	t.Run("無効なカードタイプはエラーを返す", func(t *testing.T) {
		config.BeforeEachForUnitTest()
		defer config.AfterEachForUnitTest()

		injector := app.SetupDIContainer()
		do.Override(injector, config.TestDbConnection)
		neonUseCase := do.MustInvoke[neon.NeonUseCase](injector)

		// どのタイプにも一致しないデータ
		sampleData := cardrecord.StandardCard{
			NameEn:    "Invalid Card",
			NameJa:    "無効なカード",
			Type:      "Invalid Type",
			TypeLines: []string{"Invalid"},
		}

		_, err := neonUseCase.InsertCardInfo(context.Background(), sampleData)
		assert.Error(t, err)
		assert.Equal(t, "invalid card type", err.Error())
	})
}
