# インフラアーキテクチャ・運用ドキュメント (Infrastructure Documentation)

本ディレクトリは、`ocg_statistics` プロジェクトにおけるクラウドインフラ（AWS、Neon DB、コンテナ、ネットワーク、CI/CD パイプライン等）の設計書・アーキテクチャ図・運用手順を管理するディレクトリです。

---

## 1. 最重要ルール（厳格遵守）

インフラの構築・変更・デプロイを行う際は、以下のルールを必ず遵守してください。

- **[01_script_only_infra_deployment_rules.md](../X_important/infra/rules/01_script_only_infra_deployment_rules.md)**:
  - **マネジメントコンソールからの手動変更（ClickOps）および手打ち Ad-hoc CLI の絶対禁止**
  - すべてのデプロイ・更新は、Git 管理されたスクリプト（[scripts/](../../scripts)）または IaC（[infra/opentofu/](../../infra/opentofu)）を経由すること
  - 未スクリプト化の作業は、先にスクリプトを作成・コミット・PR レビューを受けてから実行すること

---

## 2. インフラリソース管理構成

本プロジェクトのインフラは以下の 2 つの軸でコード化・自動化されています。

```
ocg_statistics/
├── infra/
│   └── opentofu/               # OpenTofu (Terraform 互換) によるリソース定義
│       ├── modules/            # 再利用可能なインフラモジュール
│       └── env/dev/            # 開発環境用設定
│           ├── base/           # ネットワーク・共通基盤
│           ├── prepare_master/ # データ準備用リソース
│           └── rdb_neon/       # Neon データベース接続・設定
│
└── scripts/                    # 運用・デプロイ自動化スクリプト
    ├── aws/                    # AWS リソース・パラメータ操作スクリプト
    ├── database/               # DB マイグレーション適用スクリプト (local / remote)
    ├── code/                   # コンテナビルド & ECR プッシュスクリプト
    ├── common/                 # 共通テスト・検証スクリプト
    └── airflow/                # Airflow 関連スクリプト
```

---

## 3. 主要デプロイスクリプト一覧

| 対象 | スクリプトパス | 概要 |
|---|---|---|
| **コンテナデプロイ (cardInfoInsert)** | [`scripts/code/cardApiCall/cardinfoinsert/build_and_push_to_ecr.sh`](../../scripts/code/cardApiCall/cardinfoinsert/build_and_push_to_ecr.sh) | cardInfoInsert コンテナの Docker ビルド・タグ付け・ECR プッシュ |
| **コンテナデプロイ (syncStorageImage)** | [`scripts/code/cardApiCall/syncstorageimage/build_and_push_to_ecr.sh`](../../scripts/code/cardApiCall/syncstorageimage/build_and_push_to_ecr.sh) | syncStorageImage コンテナの Docker ビルド・タグ付け・ECR プッシュ |
| **コンテナデプロイ (旧 infoInsert 互換)** | [`scripts/code/cardApiCall/infoinsert/build_and_push_to_ecr.sh`](../../scripts/code/cardApiCall/infoinsert/build_and_push_to_ecr.sh) | 旧 infoInsert コンテナのビルド・プッシュ（互換性維持） |
| **DB マイグレーション (Remote)** | [`scripts/database/remote_migration_up.sh`](../../scripts/database/remote_migration_up.sh) | Neon DB に対する golang-migrate の Up 適用 |
| **DB マイグレーション (Local)** | [`scripts/database/local_migration_up.sh`](../../scripts/database/local_migration_up.sh) | ローカル Docker PostgreSQL への Up 適用 |
| **パラメータ更新** | [`scripts/aws/preparemaster/update_current_id.sh`](../../scripts/aws/preparemaster/update_current_id.sh) | AWS Systems Manager (SSM) パラメータストアの更新 |

---

## 4. ドキュメント一覧

- [01_local_infrastructure_docker_compose.md](./01_local_infrastructure_docker_compose.md): **ローカルインフラ検証環境 (Docker Compose) とポート一覧**
- [02_lambda_api_specifications_and_error_handling.md](./02_lambda_api_specifications_and_error_handling.md): **AWS Lambda API 仕様・参照インデックス** (code/api/lambda_specs.yaml へのポインタ・インフラ一覧)
- [04_lambda_and_airflow_logging_investigation_guide.md](./04_lambda_and_airflow_logging_investigation_guide.md): **AWS Lambda と Airflow パイプライン連携におけるログ調査手法と IAM 落とし穴** (Airflow タスクログ、CloudWatch Logs 権限不備、aws lambda invoke 検証)
- [インフラ運用ルール索引](../X_important/infra/rules/README.md)
- [スクリプト外デプロイ禁止ルール](../X_important/infra/rules/01_script_only_infra_deployment_rules.md)

