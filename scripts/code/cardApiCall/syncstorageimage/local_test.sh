#!/bin/bash
set -Ceu

cardIds=${1:-"[89631139, 46986414]"}

echo "Invoking syncStorageImage on port 9003 (ocgCardIds: ${cardIds})..."
curl -s -X POST "http://localhost:9003/2015-03-31/functions/function/invocations" \
  -H "Content-Type: application/json" \
  -d "{\"ocgCardIds\": ${cardIds}}"
echo ""
