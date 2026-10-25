# /// script
# requires-python = ">=3.11"
# dependencies = [
#     "pyyaml>=6.0",
# ]
# ///
"""
YAML / OpenAPI 3.1 仕様書 構文検証・スキーマ抽出スクリプト (PEP 723 / uv 対応)
"""
import sys
import os
import argparse
import yaml

def validate_yaml(file_path):
    if not os.path.exists(file_path):
        print(f"✗ Error: File not found: {file_path}", file=sys.stderr)
        return False

    try:
        with open(file_path, "r", encoding="utf-8") as f:
            data = yaml.safe_load(f)

        if not isinstance(data, dict):
            print(f"✗ Error: Root element is not a mapping/dict in {file_path}", file=sys.stderr)
            return False

        print(f"✓ Valid YAML syntax: {file_path}")

        # OpenAPI 固有の検証とサマリー抽出
        if "openapi" in data:
            print(f"  - OpenAPI Version: {data.get('openapi')}")
            info = data.get("info", {})
            print(f"  - Title: {info.get('title')} (v{info.get('version')})")

            paths = data.get("paths", {})
            print(f"  - Defined Paths ({len(paths)}):")
            for path, methods in paths.items():
                for method, op in methods.items():
                    op_id = op.get("operationId", "no-op-id")
                    summary = op.get("summary", "")
                    print(f"    • [{method.upper()}] {path} -> {op_id} ({summary})")

            schemas = data.get("components", {}).get("schemas", {})
            print(f"  - Defined Schemas ({len(schemas)}):")
            for schema_name in schemas.keys():
                print(f"    • {schema_name}")

        return True

    except Exception as e:
        print(f"✗ YAML Parse Error in {file_path}: {e}", file=sys.stderr)
        return False

def main():
    parser = argparse.ArgumentParser(description="Validate YAML / OpenAPI specification files with uv")
    parser.add_argument("file", nargs="?", default="code/api/lambda_specs.yaml", help="Path to YAML file")
    args = parser.parse_args()

    success = validate_yaml(args.file)
    sys.exit(0 if success else 1)

if __name__ == "__main__":
    main()
