# 外部 API 連携仕様・ナレッジ (External APIs Integration)

本ディレクトリは、本プロジェクト（`ocg_statistics`）において連携する外部サービス（公式遊戯王 Neuron、YGOPRODeck API 等）との通信仕様、パラメータエンコーディング方針、および過去の不具合・知見ナレッジをまとめたものです。

---

## 1. ドキュメント一覧

| ドキュメント | 概要 |
| :--- | :--- |
| **[01_ygoprodeck_api_query_encoding_and_character_sanitization.md](./01_ygoprodeck_api_query_encoding_and_character_sanitization.md)** | **外部 API 連携におけるクエリエンコードと文字サニタイズの責務分離（YGOPRODeck API）**<br>「&」を含むカード（Droll & Lock Bird 等）で発生した API 取得失敗問題の分析、ブラックリスト方式文字削除の撤廃、および `url.QueryEscape` による RFC 3986 準拠パーセントエンコードの完全記録。 |
| **[02_card_id_nomenclature_and_airflow_neuron_id_mapping.md](./02_card_id_nomenclature_and_airflow_neuron_id_mapping.md)** | **遊戯王カード識別子体系とAirflow連番パイプラインにおけるNeuron CIDの責務定義**<br>Airflow の連番（`current_id`）が公式 Neuron の CID（Card ID）であることを明確化し、YGOPRODeck のパスコード（`ocg_api_id`）との混同による画像同期失敗メカニズム、JOIN 検索による是正、およびシステム全体の命名規律を記録。 |

---

## 2. 外部 API 連携の基本方針

1. **識別子体系（Neuron CID vs パスコード）の厳格な分離**: Airflow が生成・管理する連番は「公式遊戯王 Neuron の CID」であり、YGOPRODeck API のカード ID（8桁パスコード）と混同してはならない。変数・DTO・DB カラムは `neuron_id` と `ocg_api_id` を明確に区別する。
2. **ドメイン層での文字改変禁止**: カード名などの固有識別子から、URL の都合で記号（`&` や `#`）を削除してはならない。
3. **通信層での厳格な URL エンコーディング**: クエリパラメータ値には必ず `net/url` の `url.QueryEscape` を適用し、安全かつ可逆的に送信する。
4. **完全一致仕様の尊重**: 外部プロバイダ（YGOPRODeck 等）の検索仕様に合わせ、正確な名称でリクエストを構築する。

