# Lambda 開発ルール索引 (Rules Index)

Lambda / Go アプリケーション開発におけるコーディング・アーキテクチャの遵守ルール。

| No | ルール | 概要 | 詳細ドキュメント |
| :--- | :--- | :--- | :--- |
| **01** | [01_transaction_boundary_rules.md](./01_transaction_boundary_rules.md) | DB変更時のトランザクション境界配置原則（UseCase配置、RepoでのTx禁止、外部I/O除外） | [03_generic_transaction_pattern.md](../../../02_go_architecture/03_generic_transaction_pattern.md) |
| **02** | [02_related_fields_unit_test_rules.md](./02_related_fields_unit_test_rules.md) | 関連フィールドに関するユニットテスト実装原則（追加フィールドの全網羅検証、非ゼロ値データ使用、レイヤ間伝播保証） | - |
