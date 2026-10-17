# プロジェクト重要開発ルール索引 (Project Rules Index)

本ディレクトリは、`ocg_statistics` プロジェクトにおけるアーキテクチャ規律、品質担保、安全な運用のための**最重要遵守ルール**を集約・管理しています。

本プロジェクトに携わるすべての開発者は、PR 作成前に必ず該当するルールを確認・遵守してください。

---

## カテゴリ別ルール一覧

### 1. [共通ルール (Common Rules)](./common/rules/README.md)
プロジェクト全体で横断的に遵守すべき開発・運用の基本原則。

| No | ルール | 概要 | 詳細ドキュメント |
| :--- | :--- | :--- | :--- |
| **01** | [01_new_table_creation_rules.md](./common/rules/01_new_table_creation_rules.md) | 新規テーブル追加時の厳格遵守ルール（Up/Down検証、sqlc同期、DTO防波堤、ExecTx/ExecNonTx、test_all.sh、クラス図同期） | [docs/04_new_tables](../04_new_tables/01_new_table_implementation_guide.md) |
| **02** | [02_doc_code_synchronization_rules.md](./common/rules/02_doc_code_synchronization_rules.md) | コード修正時における docs ドキュメント見直し・同期原則（コード変更箇所の対応 docs 見直し、Doc-Code Drift 防止） | [docs/](../) |

---

### 2. [Lambda / Go アプリケーション開発ルール (Lambda Rules)](./lambda/rules/README.md)
`code/cardApiCall/` 配下の Go 実装、トランザクション、単体テストに関する遵守ルール。

| No | ルール | 概要 | 詳細ドキュメント |
| :--- | :--- | :--- | :--- |
| **01** | [01_transaction_boundary_rules.md](./lambda/rules/01_transaction_boundary_rules.md) | DB変更時のトランザクション境界配置原則（UseCase配置、RepoでのTx禁止、外部I/O除外） | [03_generic_transaction_pattern.md](../02_go_architecture/03_generic_transaction_pattern.md) |
| **02** | [02_related_fields_unit_test_rules.md](./lambda/rules/02_related_fields_unit_test_rules.md) | 関連フィールドに関するユニットテスト実装原則（追加フィールドの全網羅検証、非ゼロ値データ使用、レイヤ間伝播保証） | - |

---

### 3. [インフラ運用・デプロイルール (Infra Rules)](./infra/rules/README.md)
AWS、Neon DB、OpenTofu、シェルスクリプトに関する運用・デプロイの遵守ルール。

| No | ルール | 概要 | 詳細ドキュメント |
| :--- | :--- | :--- | :--- |
| **01** | [01_script_only_infra_deployment_rules.md](./infra/rules/01_script_only_infra_deployment_rules.md) | スクリプトおよび IaC を介さない手動デプロイの完全禁止（ClickOps 禁止、Ad-hoc CLI 禁止、スクリプト先行作成の原則） | [docs/00_infra](../00_infra) |

---

### 4. [ドキュメント作成ルール (Document Rules)](./document/rules/README.md)
ドキュメント内でのパス記述、ポータビリティ、プライバシー保護に関する遵守ルール。

| No | ルール | 概要 | 詳細ドキュメント |
| :--- | :--- | :--- | :--- |
| **01** | [01_prohibit_absolute_path_rules.md](./document/rules/01_prohibit_absolute_path_rules.md) | ドキュメント内における絶対パス（フルパス）記述の厳格禁止（相対パス記述、個人情報・環境差異排除の原則） | [docs/README.md](../README.md) |
