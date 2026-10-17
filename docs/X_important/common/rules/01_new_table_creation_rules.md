# ルール: 新規テーブル追加における厳格遵守ルール

データベースに新規テーブルを追加・改修する際は、以下のルールを厳格に遵守すること。
詳細な実装フローやテンプレート活用法は [docs/04_new_tables/01_new_table_implementation_guide.md](../../../04_new_tables/01_new_table_implementation_guide.md) を参照。

---

## 遵守ルール一覧

1. **マイグレーションは必ず Up / Down 両方を作成し、スクリプトでロールバックまで検証する**
   - [database/migrate/db/migrations/](../../../../database/migrate/db/migrations) に連番で `.up.sql` と `.down.sql` を作成する。
   - [local_migration_up.sh](../../../../scripts/database/local_migration_up.sh) と [local_migration_down.sh](../../../../scripts/database/local_migration_down.sh) の両方を実行し、問題なく元に戻せることを確認する。

2. **sqlc の `schema.sql` への DDL 同期とコード再生成を必ず行う**
   - [code/cardApiCall/infoInsert/sqlc/db/schema.sql](../../../../code/cardApiCall/infoInsert/sqlc/db/schema.sql) にテーブル定義を追加し、`sqlc/queries/` にクエリを記述した上で `sqlc generate` を実行する。

3. **Entity を排し、DTO を防波堤として sqlc の型（`sql.Null*`）をカプセル化する**
   - 振る舞いのない無意味な Entity 層（`domain/`）を増やして多重詰め替えを行わない（[05_dto_design_policy.md](../../../02_go_architecture/05_dto_design_policy.md) 参照）。
   - sqlc が吐き出す `sql.NullString` や `sql.NullTime` などの内部型を UseCase 戻り値や外部 API に直接露出させず、必ず [dto/cardrecord/](../../../../code/cardApiCall/infoInsert/dto/cardrecord) でプリミティブ型やポインタに変換する。

4. **トランザクション境界は UseCase に配置し、`ExecTx` / `ExecNonTx` を使い分ける**
   - **更新系**: [01_transaction_boundary_rules.md](../../lambda/rules/01_transaction_boundary_rules.md) に従い、`transaction.ExecTx` で囲む。外部通信は Tx 外で行う。
   - **参照系**: `transaction.ExecNonTx` を用い、UseCase から直接 `sqlc_gen.New(...)` を呼ばない。
   - **Repository**: 自前で `db.Begin()` や `tx.Commit()` を呼ばず、引数の `*sqlc_gen.Queries` を用いてクエリ実行とログ・メトリクス計測に専念する。

5. **全パッケージの単体テストを作成し、必ずテストスクリプト（`test_all.sh`）で検証する**
   - 作成・改修したすべてのパッケージ（DTO / Repository / UseCase）に対して単体テスト（`*_test.go`）を作成する。
   - DB を共有するテストの競合を防ぐため、単純な `go test ./...` ではなく、必ず [scripts/code/cardApiCall/infoinsert/test_all.sh](../../../../scripts/code/cardApiCall/infoinsert/test_all.sh) を実行して全テストが PASS することを確認する。

6. **クラス図（`card_class_diagram.md`）および対応する docs ドキュメントを必ず同期更新する**
   - 実装完了後、[docs/01_data_prepare/card_class_diagram.md](../../../01_data_prepare/card_class_diagram.md) を開き、Mermaid クラス図および「実装との対応表」に新テーブルと DTO の対応行を必ず追記する。
   - コード修正に伴うドキュメント見直し・同期の全般的な規律については [02_doc_code_synchronization_rules.md](./02_doc_code_synchronization_rules.md) を遵守すること。

---

## 記述例

### ⭕ 正しい記述例

#### 1. DTO で Nullable 型を吸収する
```go
func FromSelectCardPackRow(row sqlc_gen.CardPack) CardPackDTO {
    var releaseDate *time.Time
    if row.ReleaseDate.Valid {
        releaseDate = &row.ReleaseDate.Time
    }
    return CardPackDTO{
        ID:          row.ID,
        PackCode:    row.PackCode,
        NameJa:      row.NameJa,
        NameEn:      row.NameEn.String, // NullString を安全に string 化
        ReleaseDate: releaseDate,        // NullTime を *time.Time 化
    }
}
```

#### 2. UseCase で ExecTx（更新）と ExecNonTx（参照）を一貫利用する
```go
// 更新系: ExecTx でトランザクション境界を制御
func (u *neonUseCaseImpl) InsertCardPack(ctx context.Context, dto CardPackDTO) (int64, error) {
    tr := transaction.NewTx(u.ProduceConnDB())
    return transaction.ExecTx(ctx, tr, func(q *sqlc_gen.Queries) (int64, error) {
        repo := repository.NewCardPackRepository(q)
        return repo.InsertCardPack(ctx, dto.ToInsertParams())
    })
}

// 参照系: ExecNonTx で統一的に Queries を受領
func (u *neonUseCaseImpl) GetCardPackByID(ctx context.Context, id int64) (CardPackDTO, error) {
    tr := transaction.NewTx(u.ProduceConnDB())
    return transaction.ExecNonTx(ctx, tr, func(q *sqlc_gen.Queries) (CardPackDTO, error) {
        repo := repository.NewCardPackRepository(q)
        return repo.GetCardPackByID(ctx, id)
    })
}
```

---

### ❌ 違反例（アンチパターン）

#### 1. sqlc の Raw Row を外部にそのまま返す（ノーガード戦法）
```go
// ❌ 違反: sqlc の構造体をそのまま返すと JSON に {"String": "...", "Valid": true} が漏れ出る
func (u *neonUseCaseImpl) GetCardPack(ctx context.Context, id int64) (sqlc_gen.CardPack, error) {
    ...
}
```

#### 2. 参照系で UseCase が直接 sqlc を初期化する
```go
// ❌ 違反: UseCase がインフラ具象に依存。ExecNonTx に揃えること！
func (u *neonUseCaseImpl) GetCardPack(ctx context.Context, id int64) (CardPackDTO, error) {
    repo := repository.NewCardPackRepository(sqlc_gen.New(u.ProduceConnDB()))
    return repo.GetCardPackByID(ctx, id)
}
```

#### 3. クラス図を更新せずにマージする
```text
❌ 違反: コードとマイグレーションだけ書いて card_class_diagram.md を更新しない
→ 設計書と実装が乖離し、後続の開発者がテーブル間のリレーションを追えなくなる。
```

---

## 実装チェックリスト

- [ ] `database/migrate/db/migrations/` に連番の `.up.sql` と `.down.sql` の両方を作成したか？
- [ ] ローカル環境で `local_migration_up.sh` と `local_migration_down.sh` の両方を実行して検証したか？
- [ ] `sqlc/db/schema.sql` にテーブル定義を追加し、`sqlc generate` を実行したか？
- [ ] DTO で `sql.Null*` 型を吸収し、sqlc の型を外部に露出させていないか？
- [ ] UseCase の更新系で `ExecTx`、参照系で `ExecNonTx` を使用しているか？
- [ ] 全パッケージ（DTO / Repository / UseCase）の単体テストを作成したか？
- [ ] テストスクリプト `bash scripts/code/cardApiCall/infoinsert/test_all.sh` を実行してすべて PASS したか？
- [ ] [docs/01_data_prepare/card_class_diagram.md](../../../01_data_prepare/card_class_diagram.md) のクラス図および対応表を更新したか？
