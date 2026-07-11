#!/usr/bin/env bash
set -euo pipefail

API_URL="${API_URL:-http://localhost:8080}"
INCIDENT_ID="${1:-}"
DECISION="${2:-approved}"

if [ -z "${INCIDENT_ID}" ]; then
  INCIDENT_ID=$(curl -s "${API_URL}/v1/incidents" | jq -r '.[0].id')
fi

curl -s -X POST "${API_URL}/v1/incidents/${INCIDENT_ID}/approvals" \
  -H "Content-Type: application/json" \
  -d "{\"decision\":\"${DECISION}\",\"reviewer\":\"demo.reviewer@securethecloud.dev\",\"rationale\":\"Reviewed evidence chain and approved for lab-only containment workflow.\"}" | jq .
