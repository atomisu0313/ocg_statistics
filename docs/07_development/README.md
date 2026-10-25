# 開発運用・標準プロセス (Development Process & Workflow)

本ディレクトリは、`ocg_statistics` プロジェクトにおける開発プロセス、GitHub Issue 運用、テスト・品質保証、およびプルリクエスト運用の標準規約・手順書をまとめたものです。

---

## 1. ドキュメント一覧

| ドキュメント | 概要 |
| :--- | :--- |
| **[01_issue_driven_development_workflow.md](./01_issue_driven_development_workflow.md)** | **Issue 駆動開発（Issue-Driven Development）標準ワークフロー**<br>Issue 着手から要件分析、ベースライン検証、影響範囲調査、コード・ドキュメント同期、公式テスト実行、PR 作成と Issue 連携までの一連の標準手順。 |

---

## 2. 開発時のクイックリファレンス

1. **ブランチ作成とプッシュ**: 必ず `main` から `feature/issue-<番号>-<概要>` または `fix/issue-<番号>-<概要>` を切る。**`main` への直接プッシュは禁止であり、必ず作成した feature / fix ブランチに push することを徹底する**。
2. **テスト実行**: 手動の `go test` 直打ちは禁止。必ず [`scripts/code/cardApiCall/test_all.sh`](../../scripts/code/cardApiCall/test_all.sh) などの管理スクリプトを使用する。
3. **ドキュメント同期**: コードの削除・変更時は、必ず `docs/` 配下の設計書・ツリー図も同時に更新する。
4. **パス記法**: ドキュメント内にローカル絶対パス（`/Users/...` 等）は記述せず、必ず相対パスを用いる。
5. **PR 作成 & Issue 連携**: PR 本文に `Closes #<Issue番号>` を記載し、Issue コメントにも完了報告を投稿する。
