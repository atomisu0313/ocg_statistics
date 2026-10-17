# カード設計図 クラス図 (Card Class Diagram)


## 1. クラス図 (Mermaid)

```mermaid
classDiagram
    direction TB

    %% 基底クラス
    class カード {
        -long ニューロンID
        -long ocg_API_ID
        -varchar 和名
        -varchar 英名
        -String カードテキスト_日
        -String カードテキスト_英
    }

    class カード画像 {
        - long カード画像ID
        - long ocg_API_image_id
        - long ocg_API_Card_ID
        - varchar image_url_small_api
        - varchar image_url_large_api
        - varchar image_url_cropped_api
        - varchar image_url_small_this_app
        - varchar image_url_large_this_app
        - varchar image_url_cropped_this_app
    }

    %% 派生カード
    class 罠 {
        -long カードID
        -int 罠タイプ
    }

    class 魔法 {
        -long カードID
        -int 魔法タイプ
    }

    class モンスター {
        -int カードID
        -int 攻撃力
        -int 守備力
        -int レベル
    }

    %% マスタ・区分値
    class 罠タイプ {
        -int ID
        -varchar 和名
        -varchar 英名
    }

    class 魔法タイプ {
        -int ID
        -varchar 和名
        -varchar 英名
    }

    class タイプ {
        -int ID
        -varchar 英名
        -varchar 和名
    }

    class 属性 {
        -int ID
        -int 和名
        -int 英名
    }

    class 種族 {
        -int ID
        -varchar 和名
        -varchar 英名
    }

    %% モンスター特殊召喚・拡張
    class ペンデュラム {
        -long カードID
        -int スケール
        -String テキスト_日
        -String テキスト_英
    }

    class エクシーズ {
        -long カードID
    }

    class シンクロ {
        -long カードID
    }

    class 融合 {
        -long カードID
    }

    class リンク {
        -int リンクマーカー
        -long カードID
    }

    %% 汎化（継承）
    カード <|-- 罠
    カード <|-- 魔法
    カード <|-- モンスター

    %% 関連・コンポジション（マスタ参照）
    罠 *--> "1" 罠タイプ
    魔法 *--> "1" 魔法タイプ
    モンスター *--> "0..*" タイプ : -タイップライン
    モンスター *--> "1" 属性 : -属性
    モンスター *--> "1" 種族 : -種族

    %% モンスター拡張テーブルとの関連
    モンスター o-- ペンデュラム
    エクシーズ --> モンスター
    シンクロ --> モンスター
    融合 --> モンスター
    リンク --> モンスター

    %% カード画像とカードの関連
    カード o-- "1..*" カード画像

    %% 注記
    note for モンスター "カードIDにはDB内で定義されたシリアライズを採用する"
```

---

## 2. 実装（DBスキーマ・Go DTO）との対応表

| クラス図のエンティティ | DB テーブル ([schema.sql](../../code/cardApiCall/infoInsert/sqlc/db/schema.sql)) | Go DTO ([dto/cardrecord/](../../code/cardApiCall/infoInsert/dto/cardrecord/)) |
|---|---|---|
| **カード** | `cards` | `AbstractCardSelectResult` ([carddto.go](../../code/cardApiCall/infoInsert/dto/cardrecord/carddto.go)) |
| **カード画像** | `card_images` | `CardImage` ([carddto.go](../../code/cardApiCall/infoInsert/dto/cardrecord/carddto.go#L88-L93)) |
| **モンスター** | `monsters` | `MonsterCardSelectResult` ([monster.go](../../code/cardApiCall/infoInsert/dto/cardrecord/monster.go)) |
| **魔法** | `spells` | `SpellCardSelectResult` ([spell.go](../../code/cardApiCall/infoInsert/dto/cardrecord/spell.go)) |
| **罠** | `traps` | `TrapCardSelectResult` ([trap.go](../../code/cardApiCall/infoInsert/dto/cardrecord/trap.go)) |
| **ペンデュラム** | `pendulum_monsters` | `MonsterCardSelectResultExtended` ([monster.go](../../code/cardApiCall/infoInsert/dto/cardrecord/monster.go)) |
| **リンク** | `link_monsters` | `MonsterCardSelectResultExtended` ([monster.go](../../code/cardApiCall/infoInsert/dto/cardrecord/monster.go)) |
| **融合** | `fusion_monsters` | - |
| **シンクロ** | `synchro_monsters` | - |
| **エクシーズ** | `xyz_monsters` | - |
| **属性** | `attributes` | `AttributeKind` ([attribute.go](../../code/cardApiCall/infoInsert/dto/kind/attribute.go)) |
| **種族** | `races` | `RaceKind` ([race.go](../../code/cardApiCall/infoInsert/dto/kind/race.go)) |
| **タイプ** | `monster_types` | `MonsterTypeKind` ([monster_type.go](../../code/cardApiCall/infoInsert/dto/kind/monster_type.go)) |
| **魔法タイプ** | `spell_types` | `SpellTypeKind` ([spell.go](../../code/cardApiCall/infoInsert/dto/kind/spell.go)) |
| **罠タイプ** | `trap_types` | `TrapTypeKind` ([trap.go](../../code/cardApiCall/infoInsert/dto/kind/trap.go)) |

---

## 3. 関連ドキュメント

- カード画像パイプライン詳細設計: [docs/03_card_image/01_card_image_architecture.md](../03_card_image/01_card_image_architecture.md)
- 新規テーブル追加手順書: [docs/04_new_tables/01_new_table_implementation_guide.md](../04_new_tables/01_new_table_implementation_guide.md)

