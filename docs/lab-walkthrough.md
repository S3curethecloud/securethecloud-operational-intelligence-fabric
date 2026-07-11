# Lab Walkthrough

## 1. Start services

```bash
cp .env.example .env
make up
```

## 2. Start API

```bash
make api
```

## 3. Start AI service

```bash
make ai
```

## 4. Start frontend

```bash
make frontend
```

## 5. Generate event

```bash
make event
```

## 6. Inspect API manually

```bash
curl -s http://localhost:8080/v1/events | jq .
curl -s http://localhost:8080/v1/incidents | jq .
curl -s http://localhost:8080/v1/audit-log | jq .
```

## 7. Run AI investigation

Get an incident ID:

```bash
INCIDENT_ID=$(curl -s http://localhost:8080/v1/incidents | jq -r '.[0].id')
```

Run investigation:

```bash
curl -s -X POST http://localhost:8080/v1/incidents/${INCIDENT_ID}/investigate | jq .
```

## 8. Approve recommendation

```bash
curl -s -X POST http://localhost:8080/v1/incidents/${INCIDENT_ID}/approvals \
  -H 'Content-Type: application/json' \
  -d '{"decision":"approved","reviewer":"demo.reviewer@securethecloud.dev","rationale":"Approved for lab isolation only after evidence review."}' | jq .
```

## 9. Replay evidence

```bash
curl -s http://localhost:8080/v1/incidents/${INCIDENT_ID}/evidence | jq .
```
