# スクリプトを介さない手動テスト実行の禁止および管理スクリプト経由の検証徹底ルール

本ルールは、AI エージェントがコードの編集・追加・改修を行う際に遵守すべき絶対原則を定義する。

---

## 1. 義務原則（Mandatory Rules）

1. **手動（Ad-hoc）な `go test` 直打ちの禁止**:
   - ターミナルから直接 `go test ./...` などの ad-hoc コマンドを手動で直接叩いてテストを済ませることを禁止する。
   - テストの検証は、**必ず `scripts/` 配下に管理された公式スクリプトを実行すること**。

2. **管理テストスクリプトの実行**:
   - 個別機能の改修時:
     - 画像同期: `./scripts/code/cardApiCall/syncstorageimage/test_all.sh`
     - カード登録: `./scripts/code/cardApiCall/cardinfoinsert/test_all.sh`
   - 作業完了前・全体検証時:
     - 総合テスト: `./scripts/code/cardApiCall/test_all.sh`

3. **スクリプト設計規約の遵守**:
   - スクリプトの仕様・構成は [docs/06_scripts_general/01_script_architecture_and_design_policy.md](docs/06_scripts_general/01_script_architecture_and_design_policy.md) に準拠すること。
   - `set -Ceu`、リポジトリルート自動解決、DB テストの直列（`-p 1`）実行を厳守する。

4. **テスト未実行・失敗状態での完了報告の禁止**:
   - スクリプトが 100% PASS（終了コード 0）した実行結果ログを必ず確認し、テスト未実行のまま作業完了を報告してはならない。
