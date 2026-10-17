# ルール: 関連フィールドに関するユニットテスト実装の原則

DTO、DBエンティティ、外部APIレスポンス等の構造体にフィールドを追加・変更した場合、あるいはそれらを取り扱うロジック（Repository / UseCase / DTOマッピング）を実装・改修する際は、**関連するフィールドすべてに対して網羅的な検証（アサーション）を行うユニットテストを必ず実装すること**。

「フィールドを追加したのにテストコードで値が検証されていない」「マッピング漏れでゼロ値のまま素通りしている」といった不具合を未然に防ぎ、データ整合性を担保する。

---

## 遵守ルール一覧

1. **フィールドを追加・変更した場合は、関連するテストコードに必ずアサーションを追加する**  
   構造体・DTO・DBテーブル・SQLクエリに新しいフィールドを追加したら、既存または新規のユニットテストにおいて、そのフィールドが期待通りに設定・更新・マッピングされているかを明示的に検証（`assert.Equal` 等）すること。

2. **テストデータには意味のある値（非ゼロ値・一意な値）を設定する**  
   ゼロ値（空文字 `""`、数値 `0`、`false`、`nil` など）をテストデータに使うと、マッピング漏れや代入漏れが発生していてもテストをパスしてしまう。必ず意図を持った非ゼロ値を設定し、正しく値が伝播・保存されていることを保証すること。

3. **各レイヤでの値の伝播・永続化を漏れなく検証する**  
   - **Repository レイヤ**: INSERT / UPDATE 時にDBへ正しく保存され、取得クエリで同じ値が返ることを検証する。Null許容型（`sql.NullString`, `sql.NullInt64` 等）は有効値（`Valid: true`）とNULL（`Valid: false`）の両方を考慮する。
   - **UseCase レイヤ**: 外部データ（APIレスポンス等）からドメインモデル・DTOへの変換時に対象フィールドが欠落せず受け渡されていることを検証する。
   - **Client / Parser レイヤ**: 外部APIのJSONやスクレイピング結果から構造体へ正しくパースされていることを検証する。

4. **構造体の丸ごと比較だけに依存せず、追加・変更した関連フィールドは個別にも検証する**  
   自動生成IDやタイムスタンプ（`CreatedAt` 等）が含まれる場合、構造体同士の一括比較（`assert.Equal(t, a, b)`）が失敗しやすくなるため、対象フィールドを個別アサーション（`assert.Equal(t, expected.Field, actual.Field)`）で確実に検証すること。

---

## 記述例

### ⭕ 正しい記述例

#### 1. Repository レイヤでのテスト（登録値・取得値の各フィールド検証）

追加・変更したフィールド（例: `ImageUrlSmallApi`, `OcgApiImageID` 等）が、引数で渡した通りにDBへ永続化され、返却されているかをすべてアサーションする。

```go
func TestCardImageRepository(t *testing.T) {
    t.Run("カード画像メタデータが全フィールド正しく保存・取得できること", func(t *testing.T) {
        ...
        insertParams := sqlc_gen.InsertCardImageParams{
            OcgApiImageID:      46986414,
            OcgApiCardID:       46986414,
            CardID:             sql.NullInt64{Int64: card.ID, Valid: true},
            ImageUrlSmallApi:   "https://images.example.com/cards_small/46986414.jpg",
            ImageUrlLargeApi:   "https://images.example.com/cards/46986414.jpg",
            ImageUrlCroppedApi: "https://images.example.com/cards_cropped/46986414.jpg",
        }

        img, err := cardImageRepo.InsertCardImage(ctx, insertParams)
        assert.NoError(t, err)

        // ⭕ 遵守: 関連フィールドをすべて明示的にアサーションして検証する
        assert.NotZero(t, img.ID)
        assert.Equal(t, insertParams.OcgApiImageID, img.OcgApiImageID)
        assert.Equal(t, insertParams.OcgApiCardID, img.OcgApiCardID)
        assert.Equal(t, insertParams.CardID.Int64, img.CardID.Int64)
        assert.Equal(t, insertParams.ImageUrlSmallApi, img.ImageUrlSmallApi)
        assert.Equal(t, insertParams.ImageUrlLargeApi, img.ImageUrlLargeApi)
        assert.Equal(t, insertParams.ImageUrlCroppedApi, img.ImageUrlCroppedApi)
        assert.False(t, img.ImageUrlLargeThisApp.Valid) // 未設定フィールドが意図通りNULLかも確認
    })
}
```

#### 2. DTO マッピング / UseCase レイヤでのテスト（値の引き継ぎ検証）

外部APIレスポンスからDTOへ変換する際、新規フィールドが脱落せずマッピングされているかを検証する。

```go
func TestGenerateStandardCard(t *testing.T) {
    t.Run("TCG API レスポンスの画像リストが StandardCard へ正しくマッピングされること", func(t *testing.T) {
        tcgCard := dto.TcgCardResponse{
            ID:   46986414,
            Name: "Dark Magician",
            CardImages: []dto.CardImage{
                {
                    ID:             46986414,
                    ImageUrlSmall:  "https://images.example.com/small.jpg",
                    ImageUrl:       "https://images.example.com/large.jpg",
                    ImageUrlCropped: "https://images.example.com/cropped.jpg",
                },
            },
        }

        standardCard := cardrecord.GenerateStandardCard(tcgCard, neuronCard)

        // ⭕ 遵守: 新設した画像リストおよびその各要素のフィールドが欠落していないことを検証
        assert.Len(t, standardCard.CardImages, 1)
        assert.Equal(t, int64(46986414), standardCard.CardImages[0].ID)
        assert.Equal(t, "https://images.example.com/small.jpg", standardCard.CardImages[0].ImageUrlSmall)
        assert.Equal(t, "https://images.example.com/large.jpg", standardCard.CardImages[0].ImageUrlLarge)
        assert.Equal(t, "https://images.example.com/cropped.jpg", standardCard.CardImages[0].ImageUrlCropped)
    })
}
```

---

## ❌ 違反例（アンチパターン）

### 1. エラー有無やIDしか確認せず、フィールドの値検証を省略している
```go
// ❌ 違反: エラーが起きないこととIDの生成しか確認しておらず、
// 各フィールドが正しく保存されたのか、あるいは空文字/NULLで保存されたのか検知できない
img, err := cardImageRepo.InsertCardImage(ctx, insertParams)
assert.NoError(t, err)
assert.NotZero(t, img.ID)
// ここでテストが終了してしまっている！
```

### 2. テストデータにゼロ値（空文字や0）を使ってしまう
```go
// ❌ 違反: テストデータが空文字・0だと、マッピング処理で代入が漏れていてもテストが通ってしまう
insertParams := sqlc_gen.InsertCardImageParams{
    OcgApiImageID:      0,  // 0 や "" は使わない！
    ImageUrlSmallApi:   "", // 実際の値を模した一意な文字列を設定すること！
}
```

### 3. DTO やテーブルにカラムを追加したのに、既存テストを更新しない
```go
// ❌ 違反: テーブル定義やDTOに new_column を追加したが、
// 既存の単体テストに new_column のアサーションを追加せず放置している
```

---

## 実装チェックリスト

- [ ] テーブル / DTO / 構造体に追加・変更したすべてのフィールドについて、テストコードでアサーション（`assert.Equal` 等）を記述したか？
- [ ] テストデータにゼロ値（`""`, `0`, `false` 等）を使わず、明確に判別可能なダミー値を設定しているか？
- [ ] Null許容型（`sql.NullString` 等）を扱う場合、有効値設定時（`Valid: true`）とNULL時（`Valid: false`）の両ケースを検証しているか？
- [ ] DTO変換・マッピング関数において、追加フィールドが抜け落ちずに伝播していることをテストしたか？
- [ ] 追加・改修したテストを含め、`go test ./...` が正常に成功することを確認したか？
