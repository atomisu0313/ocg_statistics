#!/bin/bash
set -Ceu

startId=${1:-100}
delta=${2:-10}

echo "Invoking cardInfoInsert on port 9002 (startId: ${startId}, delta: ${delta})..."
curl -s -X POST "http://localhost:9002/2015-03-31/functions/function/invocations" \
  -H "Content-Type: application/json" \
  -d "{\"startId\": ${startId}, \"delta\": ${delta}}"
echo ""
