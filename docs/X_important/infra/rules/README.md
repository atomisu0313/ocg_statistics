# インフラ運用・デプロイルール索引 (Infra Rules Index)

本リポジトリにおけるクラウドインフラ（AWS、Neon DB、コンテナ環境等）の構築・変更・デプロイに関する遵守ルールです。

| No | ルール | 概要 | 詳細ドキュメント |
| :--- | :--- | :--- | :--- |
| **01** | [01_script_only_infra_deployment_rules.md](./01_script_only_infra_deployment_rules.md) | スクリプトおよび IaC を介さない手動デプロイの完全禁止（ClickOps 禁止、Ad-hoc CLI 禁止、スクリプト先行作成の原則） | [docs/00_infra](../../../00_infra) |

---

## 関連ドキュメント

- インフラ設計・アーキテクチャ・手順書: [docs/00_infra/](../../../00_infra)
- IaC コード定義: [infra/opentofu/](../../../../infra/opentofu)
- 運用自動化スクリプト集: [scripts/](../../../../scripts)
