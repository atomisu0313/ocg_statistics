# 共通開発ルール索引 (Common Rules Index)

本リポジトリにおける全体共通の開発・アーキテクチャ遵守ルールです。

| No | ルール | 概要 | 詳細ドキュメント |
| :--- | :--- | :--- | :--- |
| **01** | [01_new_table_creation_rules.md](./01_new_table_creation_rules.md) | 新規テーブル追加時の厳格遵守ルール（Up/Down検証、sqlc同期、DTO防波堤、ExecTx/ExecNonTx、test_all.sh、クラス図同期） | [docs/04_new_tables](../../../04_new_tables/01_new_table_implementation_guide.md) |
| **02** | [02_doc_code_synchronization_rules.md](./02_doc_code_synchronization_rules.md) | コード修正時における docs ドキュメント見直し・同期原則（コード変更箇所の対応 docs 見直し、Doc-Code Drift 防止） | [docs/](../../../README.md) |
| **03** | [03_security_and_credential_rules.md](./03_security_and_credential_rules.md) | リモート認証情報・パスワード保護原則（コードやcompose.yamlへの本番パスワード直接記述・フォールバック記載の厳格禁止） | [docs/X_important/common/rules](./03_security_and_credential_rules.md) |

---

## 関連ドキュメント

- 新規テーブル追加・実装完全ガイド: [docs/04_new_tables/01_new_table_implementation_guide.md](../../../04_new_tables/01_new_table_implementation_guide.md)
- DTO 設計方針: [docs/02_go_architecture/05_dto_design_policy.md](../../../02_go_architecture/05_dto_design_policy.md)
- カードクラス図・実装対応表: [docs/01_data_prepare/card_class_diagram.md](../../../01_data_prepare/card_class_diagram.md)
- トランザクション境界原則: [docs/X_important/lambda/rules/01_transaction_boundary_rules.md](../../lambda/rules/01_transaction_boundary_rules.md)
