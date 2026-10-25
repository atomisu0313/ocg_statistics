# 外部 API 連携仕様・ナレッジ (External APIs Integration)

本ディレクトリは、本プロジェクト（`ocg_statistics`）において連携する外部サービス（公式遊戯王 Neuron、YGOPRODeck API 等）との通信仕様、パラメータエンコーディング方針、および過去の不具合・知見ナレッジをまとめたものです。

---

## 1. ドキュメント一覧

| ドキュメント | 概要 |
| :--- | :--- |
| **[01_ygoprodeck_api_query_encoding_and_character_sanitization.md](./01_ygoprodeck_api_query_encoding_and_character_sanitization.md)** | **外部 API 連携におけるクエリエンコードと文字サニタイズの責務分離（YGOPRODeck API）**<br>「&」を含むカード（Droll & Lock Bird 等）で発生した API 取得失敗問題の分析、ブラックリスト方式文字削除の撤廃、および `url.QueryEscape` による RFC 3986 準拠パーセントエンコードの完全記録。 |

---

## 2. 外部 API 連携の基本方針

1. **ドメイン層での文字改変禁止**: カード名などの固有識別子から、URL の都合で記号（`&` や `#`）を削除してはならない。
2. **通信層での厳格な URL エンコーディング**: クエリパラメータ値には必ず `net/url` の `url.QueryEscape` を適用し、安全かつ可逆的に送信する。
3. **完全一致仕様の尊重**: 外部プロバイダ（YGOPRODeck 等）の検索仕様に合わせ、正確な名称でリクエストを構築する。
