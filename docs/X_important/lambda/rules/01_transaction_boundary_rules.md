# ルール: DB変更操作におけるトランザクション境界の配置原則

DB変更操作（INSERT / UPDATE / DELETE）を実装・改修する際は、以下のルールを厳格に遵守すること。
詳細な技術メカニズムや背景は [docs/2_go_architecture/03_generic_transaction_pattern.md](../../../2_go_architecture/03_generic_transaction_pattern.md) を参照。

---

## 遵守ルール一覧

1. **トランザクション境界は必ず UseCase レイヤに配置する**（複数リポジトリ操作を不可分な1単位にまとめるため）
2. **Repository レイヤでトランザクション（Begin / Commit / Rollback）を制御しない**（引数の `*sqlc_gen.Queries` を用いてクエリ実行に専念する）
3. **外部 I/O 処理はトランザクションの外で実行する**（HTTP 通信、スクレイピング、S3 アクセス等を Tx ブロック内に含めない）
4. **トランザクション実行には `transaction.ExecTx` を使用する**（コミット失敗検知・自動ロールバック保証）

---

## 記述例

### ⭕ 正しい記述例（UseCase レイヤ）

外部通信などの前処理は `transaction.ExecTx` の **外** で完了させ、DB変更操作のみをコールバック内に集約する。

```go
func (u *neonMonsterUseCase) InsertMonster(ctx context.Context, monsterCard dto.MonsterCard) (int64, error) {
    // 1. 外部通信・前処理・バリデーションは Tx の「外」で実行
    // （※HTTPリクエスト等を Tx 内に入れるとコネクション枯渇の原因になる）

    // 2. トランザクション境界の開始（UseCase で配置）
    return transaction.ExecTx(ctx, u.tr, func(q *sqlc_gen.Queries) (int64, error) {
        // 同一トランザクション (q) を各リポジトリに引き渡して実行
        cardID, err := u.cardRepo.InsertCard(ctx, q, monsterCard.Card)
        if err != nil {
            return 0, err // エラー時は自動的に Rollback
        }

        monsterCard.Monster.CardID = cardID
        err = u.monsterRepo.InsertMonster(ctx, q, monsterCard.Monster)
        if err != nil {
            return 0, err // エラー時は自動的に Rollback
        }

        return cardID, nil // 正常終了時は自動的に Commit
    })
}
```

### ❌ 違反例（アンチパターン）

#### 1. 外部通信をトランザクション内に入れてしまう
```go
// ❌ 違反: DBロックを保持したまま外部通信を行い、プール枯渇・タイムアウトを引き起こす
transaction.ExecTx(ctx, u.tr, func(q *sqlc_gen.Queries) (int64, error) {
    res, err := http.Get("https://api.example.com/data") // 外部通信は Tx 外で行うこと！
    ...
})
```

#### 2. Repository 側で個別にコミット・ロールバックしてしまう
```go
// ❌ 違反: Repository 内で勝手にコミットすると、後続のリポジトリ操作失敗時に元に戻せない
func (r *cardRepository) InsertCard(ctx context.Context, card Card) error {
    tx, _ := r.db.Begin()
    defer tx.Commit() // Repository で Tx 制御をしてはならない！
    ...
}
```

---

## 実装チェックリスト

- [ ] DB変更操作を UseCase レイヤの `transaction.ExecTx` で囲んでいるか？
- [ ] Repository 側で `db.Begin()` や `tx.Commit()` を呼んでいないか？
- [ ] 外部通信や重い前処理が `transaction.ExecTx` の外側に配置されているか？
