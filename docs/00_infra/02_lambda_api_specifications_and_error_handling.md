# AWS Lambda API 仕様・参照インデックス (Lambda API Reference Index)

本ドキュメントは、`ocg_statistics` における各 AWS Lambda 関数の API 定義への参照先およびインフラ概要を示すインデックスです。

> [!IMPORTANT]
> **API 仕様の単一の真実 (Single Source of Truth)**:
> 各 Lambda 関数の引数仕様、レスポンス形式、バリデーション制約、エラー挙動（Function Error、部分失敗耐性）の詳細な定義は、ドキュメントに重複記載せず、すべて **[code/api/lambda_specs.yaml](../../code/api/lambda_specs.yaml)（OpenAPI 3.1 スキーマ）** でコードとして一元管理されています。
> API 仕様を確認・更新する場合は、必ず [code/api/lambda_specs.yaml](../../code/api/lambda_specs.yaml) を参照・編集してください。

---

## 1. 対象 Lambda 関数一覧と API 仕様参照先

各 Lambda 関数のインフラ基本構成および [code/api/lambda_specs.yaml](../../code/api/lambda_specs.yaml) 内の対応オペレーションは以下の通りです。

| 関数名 | エントリポイント | ローカルポート (RIE) | タイムアウト / メモリ | API 定義参照先 (OperationId) | 責務概要 |
| :--- | :--- | :---: | :---: | :--- | :--- |
| **`cardIdCheck`** | [`code/cardApiCall/cmd/aws_lambda/cardIdCheck/main.go`](../../code/cardApiCall/cmd/aws_lambda/cardIdCheck/main.go) | `9001` | 240秒 / 1024MB | `checkCardId`<br>([`#/components/schemas/CardIdCheckInput`](../../code/api/lambda_specs.yaml)) | 公式 Neuron DB をスクレイピングし、カード ID の実在判定を行う |
| **`cardInfoInsert`** | [`code/cardApiCall/cmd/aws_lambda/cardInfoInsert/main.go`](../../code/cardApiCall/cmd/aws_lambda/cardInfoInsert/main.go) | `9002` | 240秒 / 1024MB | `insertCardInfo`<br>([`#/components/schemas/CardInfoInsertInput`](../../code/api/lambda_specs.yaml)) | 外部 API からカード詳細情報を収集し、Neon DB へ並行挿入する |
| **`syncStorageImageWithDB`** | [`code/cardApiCall/cmd/aws_lambda/syncStorageImageWithDB/main.go`](../../code/cardApiCall/cmd/aws_lambda/syncStorageImageWithDB/main.go) | `9003` | 300秒 / 1024MB | `syncStorageImages`<br>([`#/components/schemas/SyncImagesInput`](../../code/api/lambda_specs.yaml)) | 外部 CDN の画像を S3 にストリーム同期し、DB の配信 URL を更新する |

---

## 2. API 定義の保守運用ルール

API 改修を行う際は、エージェント共通ルールである [.agents/rules/04_api_documentation_and_schema_maintenance_rules.md](../../.agents/rules/04_api_documentation_and_schema_maintenance_rules.md) を必ず遵守してください。

1. **スキーマコードの同期義務**:
   - Lambda ハンドラー、DTO、バリデーション、エラー挙動を変更した際は、必ず [code/api/lambda_specs.yaml](../../code/api/lambda_specs.yaml) を同時にメンテナンスすること。
2. **ドキュメント側の二重管理禁止**:
   - 本ドキュメント等の Markdown ファイル側にリクエスト/レスポンスの詳細テーブルや JSON スキーマを二重に書き込んではならない。必ず API 定義ファイルを参照させること。

---

## 3. ローカル検証用エンドポイント (クイックリファレンス)

[compose.yaml](../../compose.yaml) 起動時（Lambda RIE 経由）のローカル検証用 cURL 実行コマンドです。

```bash
# 1. cardIdCheck (ポート 9001)
curl -s -X POST "http://localhost:9001/2015-03-31/functions/function/invocations" \
  -H "Content-Type: application/json" \
  -d '{"id": "10000"}'

# 2. cardInfoInsert (ポート 9002)
curl -s -X POST "http://localhost:9002/2015-03-31/functions/function/invocations" \
  -H "Content-Type: application/json" \
  -d '{"startId": 10000, "delta": 10}'

# 3. syncStorageImageWithDB (ポート 9003)
curl -s -X POST "http://localhost:9003/2015-03-31/functions/function/invocations" \
  -H "Content-Type: application/json" \
  -d '{"ocgCardIds": [89631139, 46986414]}'
```
