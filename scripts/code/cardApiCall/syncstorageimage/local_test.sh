#!/bin/bash
set -Ceu

neuronCardIds=${1:-"[4000, 4001]"}

echo "Invoking syncStorageImage on port 9003 (neuronCardIds: ${neuronCardIds})..."
curl -s -X POST "http://localhost:9003/2015-03-31/functions/function/invocations" \
  -H "Content-Type: application/json" \
  -d "{\"neuronCardIds\": ${neuronCardIds}}"
echo ""
