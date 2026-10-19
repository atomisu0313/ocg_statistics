# ocg_statistics ドキュメント総合案内 (Documentation Map)

本リポジトリ（`ocg_statistics`）の設計書、アーキテクチャ資料、運用手順書、および重要ルールの総合インデックスです。

---

## 1. 最重要ルール（必読）

開発・改修作業に着手する前に、必ず以下の重要ルールを確認してください。

- **[docs/X_important/](./X_important/README.md)**:
  - **[共通ルール](./X_important/common/rules/README.md)**:
    - [01_new_table_creation_rules.md](./X_important/common/rules/01_new_table_creation_rules.md): 新規テーブル追加時の厳格ルール
    - [02_doc_code_synchronization_rules.md](./X_important/common/rules/02_doc_code_synchronization_rules.md): コード修正時の docs ドキュメント見直し・同期原則
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
| **[X_important/](./X_important/README.md)** | 重要遵守ルール | 開発者が遵守すべき共通・Lambda・インフラの厳格ルール集 |
