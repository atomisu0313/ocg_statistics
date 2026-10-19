#!/bin/bash
set -Ceu

cd code/cardApiCall

# 共通およびsyncStorageImageのテスト
go test ./common/... ./syncStorageImageWithDB/...
