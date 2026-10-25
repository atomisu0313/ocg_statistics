#!/bin/bash
set -Ceu

cardIds=${1:-"[100, 101, 102]"}

echo "Invoking cardInfoInsert on port 9002 (ocgCardIds: ${cardIds})..."
curl -s -X POST "http://localhost:9002/2015-03-31/functions/function/invocations" \
  -H "Content-Type: application/json" \
  -d "{\"ocgCardIds\": ${cardIds}}"
echo ""
