#!/bin/bash
set -Ceu

# リポジトリルートを特定して code/cardApiCall へ移動
REPO_ROOT="$(git rev-parse --show-toplevel)"
cd "${REPO_ROOT}/code/cardApiCall"

echo "=================================================="
echo " [1/2] Running Standalone Unit Tests (Parallel)"
echo "=================================================="
go test \
  ./internal/util/... \
  ./internal/htmlget/... \
  ./internal/http/... \
  ./internal/usecase/neuron/... \
  ./internal/usecase/tcgapi/... \
  ./internal/usecase/syncimage/... \
  ./cmd/aws_lambda/syncStorageImageWithDB/...

echo ""
echo "=================================================="
echo " [2/2] Running Database Tests (Sequential: -p 1)"
echo "=================================================="
go test -p 1 \
  ./internal/repository/... \
  ./internal/usecase/neon/... \
  ./internal/usecase/master/...

echo ""
echo "=================================================="
echo " ✓ ALL TESTS PASSED SUCCESSFULLY!"
echo "=================================================="
