#!/bin/bash
set -Ceu

REPO_ROOT="$(git rev-parse --show-toplevel)"
cd "${REPO_ROOT}/code/cardApiCall"

echo "=== Running syncStorageImage Unit Tests ==="
# DB不要な独立テスト（ストレージ、HTTPストリーム、ユースケース、Lambda）
go test ./internal/storage/... ./internal/http/... ./internal/usecase/syncimage/... ./cmd/aws_lambda/syncStorageImageWithDB/...

echo "=== Running Repository Tests ==="
# DBを使用するリポジトリテスト（直列実行で競合回避）
go test -p 1 ./internal/repository/...
