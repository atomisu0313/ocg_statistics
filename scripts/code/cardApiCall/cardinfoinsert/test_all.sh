#!/bin/bash
set -Ceu

cd code/cardApiCall

# DB不要な独立テスト（並列実行）
go test ./internal/util/... ./internal/htmlget/... ./internal/http/... ./internal/usecase/tcgapi/...

# 同一DBを共有するテスト（直列実行で競合回避）
go test -p 1 ./internal/repository/... ./internal/usecase/neon/... ./internal/usecase/master/...
