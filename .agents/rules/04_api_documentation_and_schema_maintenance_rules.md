# API スキーマコード保守およびドキュメント同期ルール (API Documentation & Schema Maintenance Rules)

本ルールは、AI エージェントが AWS Lambda 関数をはじめとする API の実装・修正・追加を行う際に遵守すべき絶対原則を定義する。

---

## 1. 原則（Mandatory Principles）

### ① API 定義ファイル（`code/api/lambda_specs.yaml`）を「単一の真実（Single Source of Truth）」とする
- Lambda 関数の引数仕様、レスポンス形式、エラー構造、バリデーション制約は、すべて [code/api/lambda_specs.yaml](code/api/lambda_specs.yaml)（OpenAPI 3.1 スキーマ）をコードとして管理する原本（Single Source of Truth）とする。

### ② API・Lambda 改修時のスキーマ同期義務
- Lambda のハンドラー（`HandleRequest`）、DTO 構造体、バリデーションロジック、エラーハンドリングを変更・追加・削除した際は、**必ずセットで `code/api/lambda_specs.yaml` を更新・メンテナンスすること**。
- Go コードのみを変更し、API 定義ファイルの更新を怠ることを厳禁とする。

### ③ ドキュメント側への詳細二重管理の禁止（ポインタ原則）
- `docs/` 配下のドキュメント（特に `docs/00_infra/` など）に、リクエスト/レスポンスのフィールド詳細一覧やテーブル定義を重複して書き込んではならない。
- ドキュメント側は、**API 定義ファイル（`code/api/lambda_specs.yaml`）へのリンク・参照（ポインタ）とインフラ/運用の高レベルな概要のみに留める**こと。
- 仕様の詳細を一箇所（API 定義ファイル）に集約し、ドキュメントとの乖離（Specification Drift）およびメンテナンスコストの倍増を防止する。
