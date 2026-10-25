#!/bin/bash
set -Ceu

# リポジトリルートを特定して code/cardApiCall へ移動
REPO_ROOT="$(git rev-parse --show-toplevel)"
cd "${REPO_ROOT}/code/cardApiCall"

echo "=== Running cardInfoInsert Unit Tests ==="
# DB不要な独立テスト（並列実行）
go test ./internal/util/... ./internal/htmlget/... ./internal/http/... ./internal/usecase/tcgapi/... ./cmd/aws_lambda/cardInfoInsert/...

echo "=== Running Database Tests ==="
# 同一DBを共有するテスト（直列実行で競合回避）
go test -p 1 ./internal/repository/... ./internal/usecase/neon/... ./internal/usecase/master/...
