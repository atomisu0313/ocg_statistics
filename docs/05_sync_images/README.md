# カード画像 S3 同期設計・仕様 (Card Image Sync Architecture & Specs)

本ディレクトリは、外部 CDN 上のカード画像を自前 AWS S3 バケットにストリーミング同期し、データベース（Neon / PostgreSQL）の配信 URL レコードを更新する Lambda 関数 `syncStorageImageWithDB` の仕様およびインフラ要件をまとめたドキュメント群です。

---

## 1. ドキュメント構成

| No | ドキュメント | 概要 |
| :--- | :--- | :--- |
| **01** | [01_sync_storage_images_api_specification.md](./01_sync_storage_images_api_specification.md) | **API 仕様書 (Go 実装仕様)**<br>Lambda I/O 仕様、既存登録情報の突合、ダウンロード対象 DTO の事前計算、Stream による S3 プッシュ、成功 ID 集計、DB レコード更新フロー、パッケージ構成。 |
| **02** | [02_infrastructure_and_prerequisites_guide.md](./02_infrastructure_and_prerequisites_guide.md) | **インフラ準備・要件定義書**<br>S3 バケット仕様（命名規則、暗号化、ブロックパブリックアクセス）、IAM ロール・最小権限ポリシー、環境変数定義、OpenTofu (IaC) 設計、ローカル検証（MinIO）方針。 |
| **03** | [03_local_minio_testing_guide.md](./03_local_minio_testing_guide.md) | **ローカル MinIO テスト環境構築・検証手順書**<br>MinIO コンテナ連携アーキテクチャ、バケット自動作成、S3 クライアント DI の拡張仕様、local_test.sh による結合検証手順。 |

---

## 2. 関連ドキュメント

- カード画像データ管理設計: [docs/03_card_image/01_card_image_architecture.md](../03_card_image/01_card_image_architecture.md)
- カード画像パイプライン アーキテクチャ設計書: [docs/03_card_image/02_card_image_pipeline_architecture.md](../03_card_image/02_card_image_pipeline_architecture.md)
- マルチ Lambda 共通基盤パッケージ構成規約: [docs/02_go_architecture/06_multi_lambda_package_architecture.md](../02_go_architecture/06_multi_lambda_package_architecture.md)
- トランザクション境界原則: [docs/X_important/lambda/rules/01_transaction_boundary_rules.md](../X_important/lambda/rules/01_transaction_boundary_rules.md)
- リモート認証情報・パスワード保護原則: [docs/X_important/common/rules/03_security_and_credential_rules.md](../X_important/common/rules/03_security_and_credential_rules.md)
