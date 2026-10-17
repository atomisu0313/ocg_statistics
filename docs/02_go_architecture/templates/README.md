# Go アーキテクチャ コードテンプレート集

本ディレクトリには、`code/cardApiCall/infoInsert` のアーキテクチャに準拠した新規機能・テーブル追加時のボイラープレートを生成するための **Jinja2 テンプレート** を配置しています。

---

## 1. テンプレート一覧

| テンプレートファイル | 対象レイヤー / ファイル | 責務・役割 |
|---|---|---|
| [`sqlc_query.sql.jinja`](./sqlc_query.sql.jinja) | `sqlc/queries/*.sql` | sqlc の CRUD および JOIN クエリ定義 |
| [`dto.go.jinja`](./dto.go.jinja) | `dto/cardrecord/` など | 層間データ転送構造体、DB Raw Row（`sql.Null*`）からの DTO 変換 |
| [`repository.go.jinja`](./repository.go.jinja) | `repository/*.go` | sqlc クエリのラッパー、構造化ログ（zap）、メトリクス記録、DTO 変換 |
| [`usecase.go.jinja`](./usecase.go.jinja) | `usecase/<subpackage>/` | ビジネスフローのオーケストレーション、`samber/do` DI コンストラクタ、`transaction.ExecTx` によるトランザクション制御 |

---

## 2. 共通化しやすい点 vs 共通化が難しい点

### 2.1 共通化しやすい点（テンプレートで自動化できる部分）

- **sqlc クエリの基本骨格**:
  - 主キー取得（`SELECT ... WHERE id = $1`）、特定列検索（`WHERE name_ja = $1`）、Insert、Update、Delete の標準構文。
- **リポジトリのログ・エラーハンドリング**:
  - `repository/base.go` の `logDBOperation` / `logDBError` / `logDBResult` と `defer` による所要時間計測の定型パターン。
- **DTO の Nullable 解除処理**:
  - sqlc 生成構造体（`sql.NullString`, `sql.NullInt64` 等）からプリミティブ型・値オブジェクトへの取り出しパターン（`FromSelectFull...Row`）。
- **DI コンテナ（`samber/do`）との結合部**:
  - `NewUseCase[T]` を用いたコンストラクタ登録とインジェクション。
- **トランザクション制御**:
  - `transaction.ExecTx` を用いて、コールバック内で `Queries` を受け取り Repository を初期化する流れ。

### 2.2 共通化が難しい点（個別設計・判断が必要なポイント）

1. **DTO の構造設計と粒度**:
   - 本プロジェクト（`infoInsert`）は ETL / データ連携パイプラインであるため、複雑な振る舞いを持つ DDD Entity は導入せず、DTO 中心（`dto/`）のアーキテクチャを採用しています。
   - 画面やテーブルの粒度に合わせて、どのフィールドを DTO に持たせるか、Null 許容値をどうハンドリングするかを個別にモデリングする必要があります。
2. **UseCase 内の複合トランザクション順序**:
   - 単純な「親テーブル挿入 → 子テーブル挿入（例: `neon_spell.go`）」は定型化できますが、`neon_monster.go` のように「種族マスタ解決 → 属性マスタ解決 → 複数モンスタータイプ走査 → リンクマーカー判定 → ペンデュラム判定」といった複数テーブルにまたがる複雑な依存順序・分岐ロジックは、完全なテンプレート化が難しく、ハンドコード（TODO コメント部）での肉付けが必要です。
3. **リレーション（JOIN クエリ）の多様性**:
   - 1:1、1:N、N:M（中間テーブル）や LEFT OUTER JOIN など、取得したい画面・集約モデルによって SQL の結合形態が異なるため、クエリテンプレートの `join_info` にて柔軟に指定できるようにしています。

---

## 3. 各テンプレートのパラメータ例

### 3.1 `sqlc_query.sql.jinja`

```yaml
table_name: "spells"
entity_name: "Spell"
primary_key:
  name: "card_id"
  type: "bigint"
  auto_increment: false
fields:
  - name: "card_id"
    type: "bigint"
  - name: "spell_type_id"
    type: "bigint"
search_fields:
  - name: "spell_type_id"
    pascal_name: "SpellTypeID"
join_info:
  joins:
    - table: "spell_types"
      fk: "spell_type_id"
      pk: "id"
      column: "name_ja"
      alias: "spell_type_name_ja"
    - table: "spell_types"
      fk: "spell_type_id"
      pk: "id"
      column: "name_en"
      alias: "spell_type_name_en"
```

### 3.2 `dto.go.jinja`

```yaml
package_name: "cardrecord"
entity_name: "SpellCard"
entity_var: "s"
embed_base: true
fields:
  - name: "SpellTypeNameJa"
    go_type: "string"
    sql_null_type: "sql.NullString"
    db_tag: "spell_type_name_ja"
    json_tag: "spellTypeNameJa"
  - name: "SpellTypeNameEn"
    go_type: "string"
    sql_null_type: "sql.NullString"
    db_tag: "spell_type_name_en"
    json_tag: "spellTypeNameEn"
```

### 3.3 `repository.go.jinja`

```yaml
entity_name: "Spell"
entity_var: "spell"
methods:
  - name: "GetSpellByCardID"
    args:
      - name: "cardID"
        type: "int64"
    return_type: "cardrecord.SpellCardSelectResult"
    sqlc_call: "r.queries.SelectFullSpellCardInfoByCardID(ctx, cardID)"
    zap_field: "zap.Int64(\"card_id\", cardID)"
    converts_to_dto: true
```

### 3.4 `usecase.go.jinja`

```yaml
package_name: "neon"
usecase_name: "Neon"
entity_name: "SpellCard"
entity_dto: "cardrecord.StandardCard"
result_dto: "cardrecord.SpellCardSelectResult"
id_type: "int64"
```
