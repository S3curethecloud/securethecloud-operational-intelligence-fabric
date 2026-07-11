#!/usr/bin/env bash
set -euo pipefail

EVENT_FILE="${1:-telemetry/sample-events/suspicious_exec.json}"
API_URL="${API_URL:-http://localhost:8080}"

echo "Publishing runtime event from ${EVENT_FILE} to ${API_URL}/v1/events"

curl -s -X POST "${API_URL}/v1/events" \
  -H "Content-Type: application/json" \
  --data-binary "@${EVENT_FILE}" | jq .
