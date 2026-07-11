#!/usr/bin/env bash
set -euo pipefail

API_URL="${API_URL:-http://localhost:8080}"
INCIDENT_ID="${1:-}"

if [ -z "${INCIDENT_ID}" ]; then
  INCIDENT_ID=$(curl -s "${API_URL}/v1/incidents" | jq -r '.[0].id')
fi

curl -s -X POST "${API_URL}/v1/incidents/${INCIDENT_ID}/investigate" | jq .
