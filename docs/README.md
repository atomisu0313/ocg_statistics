# ocg_statistics ドキュメント総合案内 (Documentation Map)

本リポジトリ（`ocg_statistics`）の設計書、アーキテクチャ資料、運用手順書、および重要ルールの総合インデックスです。

---

## 1. 最重要ルール（必読）

開発・改修作業に着手する前に、必ず以下の重要ルールを確認してください。

- **[docs/X_important/](./X_important/README.md)**:
  - **[共通ルール](./X_important/common/rules/README.md)**:
    - [01_new_table_creation_rules.md](./X_important/common/rules/01_new_table_creation_rules.md): 新規テーブル追加時の厳格ルール
    - [02_doc_code_synchronization_rules.md](./X_important/common/rules/02_doc_code_synchronization_rules.md): コード修正時の docs ドキュメント見直し・同期原則
    - [03_security_and_credential_rules.md](./X_important/common/rules/03_security_and_credential_rules.md): リモート認証情報・パスワード保護原則
    - [04_mandatory_test_verification_rules.md](./X_important/common/rules/04_mandatory_test_verification_rules.md): コード編集時におけるテスト全件実行・検証義務原則
  - **[Lambda / Go アプリ開発ルール](./X_important/lambda/rules/README.md)**:
    - [01_transaction_boundary_rules.md](./X_important/lambda/rules/01_transaction_boundary_rules.md): トランザクション境界配置原則
    - [02_related_fields_unit_test_rules.md](./X_important/lambda/rules/02_related_fields_unit_test_rules.md): 関連フィールド単体テスト原則
  - **[インフラ運用ルール](./X_important/infra/rules/README.md)**:
    - [01_script_only_infra_deployment_rules.md](./X_important/infra/rules/01_script_only_infra_deployment_rules.md): スクリプト外手動デプロイの完全禁止原則
  - **[ドキュメント作成ルール](./X_important/document/rules/README.md)**:
    - [01_prohibit_absolute_path_rules.md](./X_important/document/rules/01_prohibit_absolute_path_rules.md): ドキュメント内での絶対パス（フルパス）記述の厳格禁止原則

---

## 2. ドキュメント体系一覧

| ディレクトリ | 名称 | 主な内容 |
|---|---|---|
| **[00_infra/](./00_infra/README.md)** | インフラ・環境運用 | AWS、Neon DB、OpenTofu 構成、主要デプロイスクリプト一覧 |
| **[01_data_prepare/](./01_data_prepare)** | データ準備・モデリング | [カードクラス図 (card_class_diagram.md)](./01_data_prepare/card_class_diagram.md)、実装対応表、データ準備アクティビティ図 |
| **[02_go_architecture/](./02_go_architecture)** | Go アーキテクチャ | レイヤー責務、[DTO 設計方針 (05_dto_design_policy.md)](./02_go_architecture/05_dto_design_policy.md)、[マルチ Lambda パッケージ構成 (06_multi_lambda_package_architecture.md)](./02_go_architecture/06_multi_lambda_package_architecture.md)、トランザクションパターン、DI パターン、[コードテンプレート集](./02_go_architecture/templates/README.md) |
| **[03_card_image/](./03_card_image/README.md)** | カード画像機能 | 画像取得・S3 保存パイプラインアーキテクチャ、詳細実装ガイド |
| **[04_new_tables/](./04_new_tables/README.md)** | 新規テーブル追加 | [新規テーブル追加・実装完全ガイド (01_new_table_implementation_guide.md)](./04_new_tables/01_new_table_implementation_guide.md) |
| **[05_sync_images/](./05_sync_images/README.md)** | カード画像 S3 同期 | [S3 画像同期 API 仕様書 (01_sync_storage_images_api_specification.md)](./05_sync_images/01_sync_storage_images_api_specification.md)、[インフラ準備・要件定義書 (02_infrastructure_and_prerequisites_guide.md)](./05_sync_images/02_infrastructure_and_prerequisites_guide.md) |
| **[06_scripts_general/](./06_scripts_general/README.md)** | スクリプト総合設計 | [スクリプト体系設計書およびテスト実行スクリプト標準規約 (01_script_architecture_and_design_policy.md)](./06_scripts_general/01_script_architecture_and_design_policy.md) |
| **[07_development/](./07_development/README.md)** | 開発標準プロセス | [Issue 駆動開発標準ワークフロー (01_issue_driven_development_workflow.md)](./07_development/01_issue_driven_development_workflow.md) |
| **[X_important/](./X_important/README.md)** | 重要遵守ルール | 開発者が遵守すべき共通・Lambda・インフラの厳格ルール集 |

---

## 3. 技術検証ナレッジ集 (Technical Knowledge)

システムのアーキテクチャ変遷、トラブルシューティング、調査検証結果を記録した完全独立ナレッジ集です。

- **[AWS SDK Go v2 Transfer Manager 移行と静的解析（Staticcheck）品質一掃](./knowledge_aws_sdk_s3_transfermanager_migration_and_linter_sweep.md)**
  - `feature/s3/manager` の非推奨化と `feature/s3/transfermanager` への統合刷新、および静的解析（SA1019/SA4006）によるエラー握りつぶし撲滅の完全記録。

