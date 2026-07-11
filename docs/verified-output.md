# Verified Local Output

The scaffold was smoke-tested with alternate local ports because port 8080 may already be reserved in some environments.

## Commands Used

```bash
cd api
PORT=18080 AI_SERVICE_URL=http://localhost:18081 go run ./cmd/server
```

```bash
cd ai-service
python3 -m uvicorn app.main:app --host 127.0.0.1 --port 18081
```

```bash
API_URL=http://localhost:18080 ./scripts/simulate_runtime_event.sh
API_URL=http://localhost:18080 ./scripts/run_investigation.sh <incident-id>
API_URL=http://localhost:18080 ./scripts/approve_incident.sh <incident-id> approved
API_URL=http://localhost:18080 ./scripts/replay_evidence.sh <incident-id>
```

## Expected Evidence Replay Shape

```json
{
  "incident_id": "inc_example",
  "chain": [
    {
      "step": 1,
      "type": "runtime_event",
      "source": "falco-simulated",
      "summary": "Unexpected shell execution detected inside payment-api container"
    },
    {
      "step": 2,
      "type": "policy_decision",
      "source": "opa",
      "summary": "Secret file access attempt inside payment workload"
    },
    {
      "step": 3,
      "type": "risk_score",
      "source": "api",
      "summary": "Risk score calculated as 100"
    },
    {
      "step": 4,
      "type": "ai_investigation",
      "source": "ai-service",
      "summary": "Evidence indicates suspicious runtime behavior on payment-api."
    },
    {
      "step": 5,
      "type": "human_approval",
      "source": "demo.reviewer@securethecloud.dev",
      "summary": "approved: Reviewed evidence chain and approved for lab-only containment workflow."
    }
  ]
}
```
