# Architecture

## MVP Architecture

```text
Runtime Event Fixture
        |
        v
Go API /v1/events
        |
        +--> OPA policy evaluation
        |
        +--> risk scoring
        |
        +--> incident creation
        |
        +--> audit log
        |
        v
Next.js Dashboard
        |
        v
Python AI Service /v1/investigate
        |
        v
Evidence-backed investigation summary
        |
        v
Human approval workflow
```

## Core Entities

```text
Asset -> Runtime Event -> Policy Decision -> Incident -> Investigation -> Approval -> Audit Log
```

## Second Pass Integrations

- Falco forwards real runtime detections.
- Fluent Bit ships logs/events.
- NATS carries event streams.
- PostgreSQL persists records.
- pgvector retrieves related evidence.
- Redis queues investigation jobs.
- Grafana visualizes operational telemetry.
- kind runs the banking workload simulation.
