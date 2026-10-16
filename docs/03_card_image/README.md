# カード画像ドキュメント索引 (Card Image Documentation Index)

本ディレクトリは、カード画像データの取得・モデル設計・保存・配信アーキテクチャに関する技術仕様を管理する。

| No | ドキュメント名 | 概要 | 関連ファイル |
| :--- | :--- | :--- | :--- |
| **01** | [01_card_image_architecture.md](./01_card_image_architecture.md) | **カード画像データ管理設計**<br/>・クラス図変更（カードとカード画像の 1:N 集約）<br/>・TCG API の画像フィールド定義と DTO（`CardImage`）<br/>・絵違い（イラスト違い）対応<br/>・自前ストレージ配信への移行フェーズ設計 | [card_class_diagram.md](../01_data_prepare/card_class_diagram.md)<br/>[carddto.go](../../code/cardApiCall/infoInsert/dto/cardrecord/carddto.go) |
| **02** | [02_card_image_pipeline_architecture.md](./02_card_image_pipeline_architecture.md) | **カード画像パイプライン アーキテクチャ設計書**<br/>・全体アーキテクチャ構成図（Mermaid）<br/>・`card_images` テーブル定義 SQL & 部分インデックス設計<br/>・S3 キー体系・CloudFront 配信設計<br/>・同期ワーカー選定比較（Airflow vs Lambda） | [schema.sql](../../code/cardApiCall/infoInsert/sqlc/db/schema.sql)<br/>[db_insert.py](../../pipeline/dags/ocg/db_insert.py) |
| **03** | [03_implementation_guide.md](./03_implementation_guide.md) | **カード画像テーブル作成・S3配置 実装作業手順書**<br/>・フェーズ1: DBスキーマ・sqlc・リポジトリ新設<br/>・フェーズ2: Go DTO & UseCase 改修（Tx内画像登録）<br/>・フェーズ3: AWS S3 & IAM 権限構築<br/>・フェーズ4: Airflow 画像ダウンロード＆S3アップロード実装<br/>・フェーズ5: 検証・テスト手順 | - |
