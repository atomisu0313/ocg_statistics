#!/bin/bash
set -Ceu

cd code/cardApiCall/infoInsert

# DB不要な独立テスト（並列実行）
go test ./util/... ./htmlget/... ./http/... ./usecase/tcgapi/...

# 同一DBを共有するテスト（直列実行で競合回避）
go test -p 1 ./repository/... ./usecase/neon/... ./usecase/master/...