# AWS SDK and Go Quality Rules

## 1. AWS SDK Go v2 S3 Transfer Manager 標準化
- **非推奨パッケージの禁止**:
  - `github.com/aws/aws-sdk-go-v2/feature/s3/manager` は非推奨（deprecated: superceded by `feature/s3/transfermanager`）のため使用禁止。
- **推奨パッケージ**:
  - `github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager` を使用すること。
  - クライアント初期化: `transfermanager.New(s3Client)`
  - オブジェクト転送: `transferManager.UploadObject(ctx, &transfermanager.UploadObjectInput{...})`

## 2. 静的解析（Linter / Staticcheck）とエラーハンドリング
- **非推奨 API (SA1019) の完全排除**:
  - ライブラリ更新時や新機能実装時は、非推奨警告が出ている API をそのまま導入・放置しない。
- **未使用エラー変数 (SA4006) の排除**:
  - 関数実行の戻り値 `err`（特に `transaction.ExecTx` など）をチェックせずに後続の処理で上書き（`err = ...`）しない。
  - テストコードであっても、必ず `require.NoError(t, err)` または `assert.NoError(t, err)` でエラー検証を行うこと。
