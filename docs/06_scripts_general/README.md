# スクリプト総合設計・運用ガイド (Scripts General Documentation)

本ディレクトリは、リポジトリ（`ocg_statistics`）におけるすべての運用・テスト・デプロイスクリプト（`scripts/` 配下）の体系設計書および規準書を管理する。

---

## 1. ドキュメント一覧

| No | ドキュメント名 | 概要 |
| :--- | :--- | :--- |
| **01** | [01_script_architecture_and_design_policy.md](./01_script_architecture_and_design_policy.md) | スクリプト全体のディレクトリ体系設計、テストスクリプト標準規約（`set -Ceu`、リポジトリルート自動解決、直列 `-p 1` / 並列制御、標準テンプレート） |

---

## 2. 関連する重要ルール

- **[04_mandatory_test_verification_rules.md](../X_important/common/rules/04_mandatory_test_verification_rules.md)**:
  テスト実行時の Ad-hoc 手動コマンド禁止と、スクリプト経由での実行徹底ルール
- **[01_script_only_infra_deployment_rules.md](../X_important/infra/rules/01_script_only_infra_deployment_rules.md)**:
  スクリプトおよび IaC を介さない手動インフラデプロイの完全禁止ルール
