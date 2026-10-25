# スクリプト体系設計書およびテスト実行スクリプト標準規約 (Scripts Architecture and Testing Standards)

本ドキュメントは、リポジトリ（`ocg_statistics`）におけるすべての自動化スクリプト（`scripts/` 配下）のディレクトリ構造、設計規約、および「テスト実行の完全スクリプト化」の標準仕様を定義する。

---

## 1. 背景と目的

開発作業において、ターミナル上で手動（Ad-hoc）に `go test ./...` やコマンドライン引数を打ち込んでテストを行うと、以下の深刻な問題が発生する：

1. **テスト実行条件のブレ（環境差異・競合）**:
   - テスト用 DB（`ocg-db-local`）に対して直列実行フラグ（`-p 1`）を付け忘れてテーブル Truncate が並行で走り、原因不明のフレイキーテスト（Flaky Tests）が頻発する。
2. **実行対象パッケージの抜け漏れ**:
   - 「手元の特定パッケージだけテストして OK だったが、共通層や他 Lambda が壊れていた」というサイレントデグレを見落とす。
3. **再現性と作業証跡の欠如**:
   - 誰がどのようなパラメータや環境変数で検証したのかが残らず、CI/CD パイプラインとの動作乖離を生む。

このため、本プロジェクトでは **「手動でのワンライナー実行を排し、管理されたシェルスクリプト経由でのテスト実行を徹底する」** ことを絶対規律とする。

---

## 2. スクリプト全体のディレクトリ体系

`scripts/` ディレクトリは対象レイヤー・責務ごとに厳格に分類・配置されている。

```text
scripts/
├── code/                              # アプリケーションコード（Go）関連スクリプト
│   └── cardApiCall/
│       ├── test_all.sh                # 【総合テスト】全モジュールの一括全件検証
│       ├── cardinfoinsert/            # cardInfoInsert Lambda 関連
│       │   ├── test_all.sh            # 個別機能テストスクリプト
│       │   ├── local_test.sh          # ローカルコンテナ（ポート9002）へのリクエスト検証
│       │   └── build_and_push_to_ecr.sh
│       ├── syncstorageimage/          # syncStorageImageWithDB Lambda 関連
│       │   ├── test_all.sh            # 個別機能テストスクリプト
│       │   ├── local_test.sh          # ローカルコンテナ（ポート9003）へのリクエスト検証
│       │   └── build_and_push_to_ecr.sh
│       ├── cardidcheck/               # cardIdCheck Lambda 関連
│       │   ├── local_test.sh          # ローカルコンテナ（ポート9001）へのリクエスト検証
│       │   └── build_and_push_to_ecr.sh
│       └── internal/                  # 共通コード生成等
│           └── sqlc_generate.sh
├── database/                          # DB マイグレーション関連
│   ├── local_migration_up.sh
│   ├── local_migration_down.sh
│   ├── remote_migration_up.sh
│   └── remote_migration_down.sh
├── actions/                           # act を用いた GitHub Actions ローカル実行
│   ├── run_cardinsert.sh
│   ├── run_syncstorageimage.sh
│   └── run_idcheck.sh
├── airflow/                           # MWAA / Airflow 関連
│   ├── init_local_airflow.sh
│   ├── check_dag_list.sh
│   └── test_dag.sh
└── aws/                               # AWS インフラ・SSM パラメータ操作
    └── preparemaster/
        ├── get_ssm_parameters.sh
        ├── update_current_id.sh
        └── update_delta_id.sh
```

---

## 3. テスト実行スクリプトの設計規約

すべての `test_all.sh` スクリプトは、以下の規約を遵守して作成されなければならない。

### 3.1 厳格なシェルオプション（`set -Ceu`）
スクリプト冒頭で必ず以下を宣言する：
```bash
#!/bin/bash
set -Ceu
```
- `-e`: コマンドが 1 つでもエラー（終了コード非 0）を返した時点でスクリプトを即時停止する。
- `-u`: 未定義変数の参照をエラーとして検知する。
- `-C`: 既存ファイルへの上書きリダイレクトを禁止する。

### 3.2 実行ディレクトリ非依存（リポジトリルート自動解決）
カレントディレクトリがどこであっても（リポジトリルート、スクリプト直下、サブディレクトリなど）確実に正しく動作するよう、`git rev-parse --show-toplevel` を用いてルートを特定する：
```bash
REPO_ROOT="$(git rev-parse --show-toplevel)"
cd "${REPO_ROOT}/code/cardApiCall"
```

### 3.3 テスト種別ごとの並列／直列（`-p 1`）制御
テストは性質に応じて明確に分割して実行する：
1. **DB 不要の独立ユニットテスト（並列実行）**:
   - `internal/util/...`, `internal/htmlget/...`, `internal/http/...`, `internal/usecase/...`
   - DB に接続しないため、デフォルトの並列実行で高速に完了させる。
2. **DB を使用する統合テスト（直列実行: `-p 1`）**:
   - `internal/repository/...`, `internal/usecase/neon/...`, `internal/usecase/master/...`
   - 同一のテスト DB（`ocg-db-local`）に対して Truncate を伴うため、競合によるテスト失敗を防ぐため必ず `-p 1` を付与する。

---

## 4. 標準テストスクリプト仕様とテンプレート

### 4.1 総合テストスクリプト: `scripts/code/cardApiCall/test_all.sh`
カードAPI呼び出し機能全体の総合回帰テストを実行する。

```bash
#!/bin/bash
set -Ceu

REPO_ROOT="$(git rev-parse --show-toplevel)"
cd "${REPO_ROOT}/code/cardApiCall"

echo "=================================================="
echo " [1/2] Running Standalone Unit Tests (Parallel)"
echo "=================================================="
go test \
  ./internal/util/... \
  ./internal/htmlget/... \
  ./internal/http/... \
  ./internal/usecase/neuron/... \
  ./internal/usecase/tcgapi/... \
  ./internal/usecase/syncimage/... \
  ./cmd/aws_lambda/syncStorageImageWithDB/...

echo ""
echo "=================================================="
echo " [2/2] Running Database Tests (Sequential: -p 1)"
echo "=================================================="
go test -p 1 \
  ./internal/repository/... \
  ./internal/usecase/neon/... \
  ./internal/usecase/master/...

echo ""
echo "=================================================="
echo " ✓ ALL TESTS PASSED SUCCESSFULLY!"
echo "=================================================="
```

### 4.2 機能別テストスクリプト: `scripts/code/cardApiCall/<feature>/test_all.sh`
各 Lambda や機能モジュールに閉じた検証を行うスクリプト。

- **`syncstorageimage/test_all.sh`**:
  ```bash
  #!/bin/bash
  set -Ceu

  REPO_ROOT="$(git rev-parse --show-toplevel)"
  cd "${REPO_ROOT}/code/cardApiCall"

  echo "=== Running syncStorageImage Unit Tests ==="
  go test ./internal/storage/... ./internal/http/... ./internal/usecase/syncimage/... ./cmd/aws_lambda/syncStorageImageWithDB/...

  echo "=== Running Repository Tests ==="
  go test -p 1 ./internal/repository/...
  ```

- **`cardinfoinsert/test_all.sh`**:
  ```bash
  #!/bin/bash
  set -Ceu

  REPO_ROOT="$(git rev-parse --show-toplevel)"
  cd "${REPO_ROOT}/code/cardApiCall"

  echo "=== Running cardInfoInsert Unit Tests ==="
  go test ./internal/util/... ./internal/htmlget/... ./internal/http/... ./internal/usecase/tcgapi/...

  echo "=== Running Database Tests ==="
  go test -p 1 ./internal/repository/... ./internal/usecase/neon/... ./internal/usecase/master/...
  ```

---

## 5. 運用ルール・保守指針

1. **Ad-hoc コマンド手打ちの禁止**:
   - ターミナルから即興の `go test` コマンドを手打ちで実行して「テスト完了」と判断してはならない。必ず上記スクリプトを実行すること。
2. **パッケージ変更時のスクリプト同期義務**:
   - パッケージのリファクタリング、移動、追加を行った際は、対象スクリプト内の `go test` 対象パスを即座に最新化し、スクリプトが 100% 成功する状態を保つこと。
3. **新規 Lambda 追加時の義務**:
   - 新規に Lambda 関数を追加した際は、必ず `scripts/code/<module>/<feature>/` 配下に `test_all.sh`、`local_test.sh`、`build_and_push_to_ecr.sh` を整備し、総合テストスクリプト（`test_all.sh`）にも登録すること。
