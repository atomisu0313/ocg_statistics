#!/bin/bash
set -Ceu

startId=${1:-100}
delta=${2:-10}

echo "Invoking syncStorageImage on port 9003 (startId: ${startId}, delta: ${delta})..."
curl -s -X POST "http://localhost:9003/2015-03-31/functions/function/invocations" \
  -H "Content-Type: application/json" \
  -d "{\"startId\": ${startId}, \"delta\": ${delta}}"
echo ""
