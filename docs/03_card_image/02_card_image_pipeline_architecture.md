# カード画像パイプライン アーキテクチャ設計書

## 1. システム全体概要

本アーキテクチャは、遊戯王カードの画像メタデータ（URL群）をDBに永続化し、実体画像ファイルを外部CDNからダウンロードして自前 AWS S3 バケットへ配置、さらに配信URL（`image_url_*_this_app`）を更新するまでの一連の処理基盤を定義する。

### 主な要件
1. **DBスキーマ拡張**: クラス図に基づく `card_images` テーブルの新設と、カード（`cards`）との `1:N` 関連づけ。
2. **画像URLの登録（メタデータ収集）**: `infoInsert`（Go Lambda）内で TCG API から返却される `card_images` 配列を受け取り、`card_images` テーブルに INSERT する。
3. **画像ダウンロード & S3 配置**: 外部 CDN（`images.ygoprodeck.com`）から画像バイナリを非同期でダウンロードし、S3 にアップロード。
4. **自前配信 URL の更新**: S3 / CloudFront の URL を `card_images` の `image_url_*_this_app` に書き戻す。

---

## 2. アーキテクチャ構成図

```mermaid
flowchart TD
    subgraph DataIngestion["1. メタデータ収集フェーズ (Go Lambda: infoInsert)"]
        TCGAPI["YGOPRODeck API"] -->|"card_images[] メタデータ"| Lambda["infoInsert (Go Lambda)"]
        Lambda -->|"親カード INSERT"| CardsTable[("cards テーブル")]
        Lambda -->|"画像URL群 INSERT (this_appはNULL)"| ImagesTable[("card_images テーブル")]
    end

    subgraph ImageStoragePipeline["2. 画像同期パイプライン (S3 アップロード)"]
        Trigger["Airflow / EventBridge"] -->|"未同期画像を検出して実行"| SyncWorker["画像同期ワーカー<br/>(Python/Airflow Task or Lambda)"]
        ImagesTable -.->|"未ダウンロード画像を取得<br/>(where image_url_large_this_app IS NULL)"| SyncWorker
        SyncWorker -->|"HTTP GET (並行ダウンロード)"| ExtCDN["YGOPRODeck CDN"]
        SyncWorker -->|"PUT Object"| S3Bucket[("AWS S3 バケット<br/>ocg-card-images")]
        SyncWorker -->|"UPDATE image_url_*_this_app"| ImagesTable
    end

    subgraph Delivery["3. 画像配信フェーズ (CDN & Frontend)"]
        S3Bucket --> CloudFront["AWS CloudFront CDN"]
        ClientApp["Web UI / App"] -->|"高速・安定キャッシュ画像取得"| CloudFront
    end
```

---

## 3. DB スキーマ定義 (`card_images`)

### 3.1 テーブル定義 SQL (`sqlc/db/schema.sql` に追加)

```sql
CREATE TABLE card_images (
    id BIGSERIAL PRIMARY KEY,
    ocg_api_image_id BIGINT UNIQUE NOT NULL,
    ocg_api_card_id BIGINT NOT NULL,
    card_id BIGINT,
    image_url_small_api VARCHAR(1024) NOT NULL,
    image_url_large_api VARCHAR(1024) NOT NULL,
    image_url_cropped_api VARCHAR(1024) NOT NULL,
    image_url_small_this_app VARCHAR(1024),
    image_url_large_this_app VARCHAR(1024),
    image_url_cropped_this_app VARCHAR(1024),
    dataowner VARCHAR(255) DEFAULT 'system',
    regist_date TIMESTAMP DEFAULT current_timestamp,
    enable_start_date TIMESTAMP DEFAULT '1970-01-01 00:00:00',
    enable_end_date TIMESTAMP DEFAULT '9999-12-31 23:59:59',
    version BIGINT DEFAULT 1,
    FOREIGN KEY (card_id) REFERENCES cards(id) ON DELETE CASCADE
);

CREATE INDEX idx_card_images_card_id ON card_images(card_id);
CREATE INDEX idx_card_images_ocg_api_card_id ON card_images(ocg_api_card_id);
CREATE INDEX idx_card_images_sync_pending ON card_images(id) WHERE image_url_large_this_app IS NULL;
```

### 3.2 外部キーとインデックス設計の要点
- `card_id` は `cards(id)` を参照する外部キー。カード本体が削除された場合はカスケード削除。
- `image_url_large_this_app IS NULL` に対する部分インデックス（Partial Index）を作成し、**S3 未同期の画像を高速に抽出**可能にする。

---

## 4. S3 ストレージ設計と命名規則

### 4.1 S3 キー体系（ディレクトリ構造）
```text
s3://<bucket-name>/cards/
  ├── small/
  │   └── {ocg_api_image_id}.jpg      # サムネイル (カード全体)
  ├── large/
  │   └── {ocg_api_image_id}.jpg      # 高解像度 (カード全体)
  └── cropped/
      └── {ocg_api_image_id}.jpg      # イラスト部分のみ
```

- **キーに `ocg_api_image_id` を採用する理由**:
  - 同一カードの絵違い（オルタナティブアート）ごとにユニークなファイル名となり、上書き衝突を完全に防止できる。

### 4.2 CloudFront 配信 URL フォーマット
```text
https://images.example.com/cards/small/{ocg_api_image_id}.jpg
https://images.example.com/cards/large/{ocg_api_image_id}.jpg
https://images.example.com/cards/cropped/{ocg_api_image_id}.jpg
```

---

## 5. 同期ワーカー（画像ダウンロード & S3 配置）の選定とAirflow連携構成

### 5.1 方式の選定経緯
画像同期処理は、ストリーミングアップロードと高速な並行処理が求められるため、専用の Go Lambda（`syncStorageImageWithDB`）として実装されている。
Airflow（`ocg_db_insert` DAG）では、**DB インサートタスク（`invoke_lambda_function`）の直後に画像同期タスク（`invoke_sync_image_lambda`）を直列にキックする方針**を採用している。

```mermaid
flowchart LR
    Start(["start"]) --> Fetch["start_fetch_param"]
    Fetch --> GetParams["get_current_id<br/>get_delta_id<br/>get_max_id"]
    GetParams --> GenIDs["generate_card_ids<br/>(ocgCardIds配列生成)"]
    GetParams --> Branch{"choose_branch_path<br/>(上限超過チェック)"}
    Branch -->|"超過"| Alarm["alarm_excess_id"]
    Branch -->|"正常"| InsertLambda["invoke_lambda_function<br/>(cardInfoInsert Lambda)"]
    GenIDs -.->|"ocgCardIds"| InsertLambda
    InsertLambda --> SyncLambda["invoke_sync_image_lambda<br/>(syncStorageImageWithDB Lambda)"]
    GenIDs -.->|"neuronCardIds"| SyncLambda
    SyncLambda --> UpdateID["update_current_id<br/>(SSM更新)"]
    Alarm --> EndUpdate["end_update_task"]
    UpdateID --> EndUpdate
    EndUpdate --> End(["end"])
```

### 5.2 各タスクの責務
1. **`generate_card_ids`**: SSM から取得した `current_id` と `delta_id` より、処理対象のカード ID 配列（Neuron CID 整数配列）を一意に生成し XCom へ渡す。
2. **`invoke_lambda_function`**: `cardInfoInsert` Lambda を呼び出し、Neuron DB / YGOPRODeck API からカード基本情報・詳細情報および画像メタデータ（URL）を Neon DB へ INSERT する。
3. **`invoke_sync_image_lambda`**: インサート直後に `syncStorageImageWithDB` Lambda を呼び出し、同一の Neuron CID 配列（`neuronCardIds`）を対象として外部 CDN から実体画像を S3 へストリーム保存し、DB の配信 URL（`image_url_*_this_app`）を更新する。
4. **`update_current_id`**: 画像同期まで完了した後に、次回開始 ID（`current_id + delta_id`）を SSM Parameter Store に書き戻す。

