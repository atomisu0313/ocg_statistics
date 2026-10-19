# リモート認証情報・パスワード保護原則 (Security & Credential Rules)

本ドキュメントは、リモートデータベース（Neon 等）やクラウド環境（AWS 等）の接続情報・パスワードをコードおよび設定ファイルに直接含めることを厳格に禁止するルールです。

---

## 1. 根本原則（Strict Prohibition）

1. **リモート環境（本番 / ステージング）のパスワード・シークレットをソースコード内に直接記述（ハードコード）することを固く禁ずる**。
2. **フォールバック値（デフォルト引数など）にリモート環境の実パスワードを設定することを禁ずる**。
3. **`compose.yaml` やリポジトリ配下の設定ファイルにリモートの資格情報を混入させてはならない**。
   - `compose.yaml` はローカルコンテナ環境（`ocg_db_local` 等）専用の設定であり、常にローカル用のダミー値・デフォルト値（`postgres`/`postgres` 等）を維持する。

---

## 2. 認証情報の管理規約

| 区分 | 許容される配置場所 | 禁止事項 |
| :--- | :--- | :--- |
| **リモートDB (Neon等)** | Lambdaの環境変数（IaC/OpenTofuやCIから注入）、`.env`（ローカル用・Git管理外）、GitHub Secrets | コード内ハードコード、`database.go` 等のデフォルト引数への記載、`compose.yaml` への混入 |
| **AWS 資格情報** | `~/.aws/credentials`、IAMロール、`act_secretfile`（Git管理外）、GitHub Secrets | コード内ハードコード、公開リポジトリへのコミット |
| **ローカルDB** | `compose.yaml`（`postgres`/`postgres` などの開発用ローカル値） | 本番・リモート認証情報の混入 |

---

## 3. 実装上の遵守事項

### 3.1 Go アプリケーションコード
接続情報は必ず環境変数（`os.Getenv` や `env.GetAsString`）から取得し、フォールバック値にはローカル開発用のデフォルト値（例: `localhost`, `postgres`）のみを使用する。

```go
// 正しい例: リモートパスワードを含めず、環境変数から取得
func NewDbConnection(i *do.Injector) (*DbConn, error) {
	dbURI := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=require",
		env.GetAsString("PG_DB_USER", "postgres"),
		env.GetAsString("PG_DB_PASSWORD", "postgres"),
		env.GetAsString("PG_DB_HOST_PORT", "localhost:5555"),
		env.GetAsString("PG_DB_NAME", "ocg_statics"),
	)
    // ...
}
```

```go
// 違反例: リモートのパスワードや接続先をハードコードする
dbURI := fmt.Sprintf("postgresql://neondb_owner:npg_xxxx@ep-xxxx.aws.neon.tech/neondb?sslmode=require") // 絶対禁止
```

### 3.2 誤混入時の対応
万が一リモートのパスワードがコードや作業用ファイルに混入した場合は、直ちに削除し、認証情報（DBパスワード等）のローテーションを実施すること。
