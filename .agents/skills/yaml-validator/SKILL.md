---
name: yaml-validator
description: YAML および OpenAPI 3.1 仕様書の構文検証・スキーマ抽出・整合性チェックを行うスキル。PEP 723 仕様に準拠し、uv run で依存関係（pyyaml）を自動解決して実行する。
---

# YAML & OpenAPI Validator Skill (PEP 723 / uv 対応)

本スキルは、リポジトリ内の YAML ファイル（特に `code/api/lambda_specs.yaml` などの OpenAPI 3.1 仕様書）の構文妥当性検証およびスキーマ構造の抽出を安全・隔離環境で行うための標準スキルです。

PEP 723（Inline Script Metadata）に準拠しており、ホスト環境への追加ライブラリ導入なしに `uv` 経由で実行可能です。

---

## 1. 提供スクリプト

- **スクリプトパス**: `.agents/skills/yaml-validator/scripts/validate_yaml.py`
- **仕様**: PEP 723 (`# /// script`, `pyyaml>=6.0`, `requires-python = ">=3.11"`)
- **主な機能**:
  - YAML 構文（インデント・クォート・マルチライン記法）のパースエラー検知
  - OpenAPI 仕様書の場合、バージョン、エンドポイント一覧（メソッド・OperationId・サマリー）、定義スキーマ（Schemas）一覧の自動抽出・サマリー表示

---

## 2. 実行手順 (`uv run`)

リポジトリルートから以下のコマンドで実行します（`uv` が自動的に隔離環境を作成し、`pyyaml` を一時インストールして実行します）：

```bash
# デフォルト (code/api/lambda_specs.yaml) を検証
uv run .agents/skills/yaml-validator/scripts/validate_yaml.py

# 任意の YAML ファイルを指定して検証
uv run .agents/skills/yaml-validator/scripts/validate_yaml.py path/to/target.yaml
```

---

## 3. 終了コード仕様

- `0`: 構文正常（Syntax OK）
- `1`: ファイル未存在、またはパースエラー（Syntax Error）
