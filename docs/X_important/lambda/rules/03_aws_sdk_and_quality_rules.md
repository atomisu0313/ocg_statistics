# AWS SDK & 品質管理ルール (AWS SDK & Quality Rules)

本ドキュメントは、AWS SDK Go v2 の利用方針および静的解析（Linter / Staticcheck）における品質管理原則を定義する。

---

## 1. AWS SDK Go v2 S3 Transfer Manager 標準化

### 1.1 非推奨パッケージ（`feature/s3/manager`）の禁止
- `github.com/aws/aws-sdk-go-v2/feature/s3/manager` は非推奨（deprecated: superceded by `feature/s3/transfermanager`）である。
- 新規コードおよび既存改修において `manager.NewUploader` や `manager.NewDownloader` を使用してはならない。

### 1.2 推奨パッケージ（`feature/s3/transfermanager`）の採用
- 高レベルの S3 オブジェクト転送（ストリーミングアップロード、マルチパート転送など）には `github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager` を使用する。
- **実装例**:
  ```go
  import "github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager"

  // クライアント生成
  tm := transfermanager.New(s3Client)

  // ストリーミングアップロード
  output, err := tm.UploadObject(ctx, &transfermanager.UploadObjectInput{
      Bucket:      &bucketName,
      Key:         &key,
      Body:        streamReader,
      ContentType: &contentType,
  })
  ```

---

## 2. 静的解析（Linter / Staticcheck）とエラーハンドリング原則

### 2.1 非推奨 API（SA1019）の完全排除
- パッケージ追加時や改修時には静的解析（`staticcheck` や `go vet`）を実施し、非推奨 API の警告（`SA1019`）を放置しない。

### 2.2 未使用エラー変数（SA4006）およびエラー握りつぶしの禁止
- DB トランザクション（`transaction.ExecTx`）や外部 I/O の戻り値 `err` をチェックせずに、直後の別操作で上書き（`err = ...`）してはならない。
- テストコードにおいても、すべてのエラー戻り値に対して `assert.NoError(t, err)` または `require.NoError(t, err)` による明示的な検証を行うこと。

---

## 3. 詳細技術ナレッジ

本ルールの策定経緯、AWS SDK Go v2 における Transfer Manager 刷新のアーキテクチャ背景、および `staticcheck` による品質一掃の具体例については、以下の技術検証ナレッジを参照のこと。

- [技術検証ナレッジ: AWS SDK Go v2 Transfer Manager 移行と静的解析（Staticcheck）品質一掃 (knowledge_aws_sdk_s3_transfermanager_migration_and_linter_sweep.md)](../../../knowledge_aws_sdk_s3_transfermanager_migration_and_linter_sweep.md)

