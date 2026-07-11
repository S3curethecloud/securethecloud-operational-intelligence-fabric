# SecureTheCloud Operational Intelligence Fabric

## Palantir-Style AI Operations Lab

SecureTheCloud Operational Intelligence Fabric is a hands-on portfolio lab for building a governed AI operations platform. It simulates a regulated enterprise environment with runtime telemetry, Kubernetes-style operational events, OPA policy correlation, AI-assisted investigations, evidence replay, and human-in-the-loop approvals.

This lab is designed to show enterprise-grade operational architecture thinking rather than a toy chatbot demo.

## Scenario

**Secure AI Banking Operations**

You are building an internal operational intelligence platform for a regulated banking environment. The platform helps security and operations teams investigate suspicious runtime activity, correlate telemetry, evaluate policy context, generate evidence-backed AI summaries, and require human review before an operational recommendation is accepted.

## MVP Build Order

Start with this sequence:

```text
Phase 0 -> Phase 1 -> Phase 4 -> Phase 5 -> Phase 6 -> Phase 7 -> Phase 8 -> Phase 9
```

Leave real Falco, Kubernetes, OpenTelemetry, Grafana, and Terraform for the second pass.

## What This MVP Builds

| Phase | Outcome |
|---|---|
| Phase 0 | GitHub repo foundation |
| Phase 1 | Local Docker services and project scaffold |
| Phase 4 | OPA policy correlation |
| Phase 5 | Go operational API |
| Phase 6 | Python AI investigation service |
| Phase 7 | Human approval workflow |
| Phase 8 | Next.js AI operations workspace |
| Phase 9 | Evidence replay |

## Architecture

```text
+-------------------------------------------------------------+
|                 AI Operations Workspace                     |
|                  Next.js Investigator UI                    |
+--------------------------+----------------------------------+
                           |
                           v
+-------------------------------------------------------------+
|                         API Layer                           |
|              Go API Gateway + Python AI Service             |
+------------+-------------------+----------------------------+
             |                   |
             v                   v
+---------------------+   +-----------------------------+
| Operational Store   |   | AI Investigation Workflow   |
| PostgreSQL + Redis  |   | FastAPI + Mock LLM First    |
+---------------------+   +-----------------------------+
             |                   |
             v                   v
+-------------------------------------------------------------+
|                  Evidence + Governance Layer                |
|       OPA Decisions | Approval Records | Audit Chain         |
+-------------------------------------------------------------+
             ^
             |
+-------------------------------------------------------------+
|                  Runtime Detection Fabric                   |
| Simulated Runtime Events | OPA | Future Falco/K8s/Otel      |
+-------------------------------------------------------------+
```

## Governance Boundary

The AI layer is an evidence summarization and recommendation layer only.

It may:

- summarize evidence
- correlate signals
- explain policy context
- recommend next review steps
- cite evidence IDs

It may not:

- authorize runtime actions
- claim production enforcement occurred
- mutate infrastructure
- bypass OPA, SENTINEL, or human approval
- issue tokens or create runtime sessions
- silently approve remediation

## Quick Start

### 1. Start infrastructure

```bash
cp .env.example .env
make up
```

### 2. Start the Go API

```bash
make api
```

If port 8080 is already in use, run:

```bash
cd api
PORT=18080 go run ./cmd/server
```

### 3. Start the AI service

In another terminal:

```bash
make ai
```

Give the AI service a few seconds to install dependencies and start before running the investigation step.

### 4. Start the frontend

In another terminal:

```bash
make frontend
```

### 5. Simulate a suspicious runtime event

```bash
./scripts/simulate_runtime_event.sh
```

### 6. Open the dashboard

```text
http://localhost:3000
```

## Demo Flow

1. Open the dashboard.
2. Show no active incidents or show the sample queue.
3. Run `./scripts/simulate_runtime_event.sh`.
4. Open the incident queue.
5. Select the new high-risk incident.
6. Review runtime evidence.
7. Review OPA policy context.
8. Run AI investigation.
9. Review evidence-backed AI recommendation.
10. Approve, reject, or request more evidence.
11. Replay the evidence timeline.
12. Show audit log.

## MVP Acceptance Criteria

The MVP is complete when:

- `make up` starts local services.
- `make api` starts the Go API.
- `make ai` starts the Python AI service.
- `make frontend` starts the Next.js dashboard.
- event simulator posts a runtime event.
- API stores runtime event in memory for MVP demo.
- API calls OPA or falls back to local policy context if OPA is unavailable.
- high-risk event creates an incident.
- AI service generates an evidence-backed summary.
- human reviewer can approve, reject, or request more evidence.
- incident evidence replay shows event, policy, investigation, approval, and audit chain.

## Second Pass Roadmap

After the MVP works, add:

- PostgreSQL persistence
- Redis-backed investigation jobs
- NATS streaming ingestion
- real Falco event ingestion
- Fluent Bit forwarding
- kind Kubernetes workload simulation
- OpenTelemetry traces
- Grafana dashboards
- Terraform local infrastructure polish
- pgvector evidence retrieval
- LangGraph workflow implementation

## Recruiter-Facing Summary

SecureTheCloud Operational Intelligence Fabric is a governed AI operations lab that simulates how regulated enterprises combine telemetry, runtime detection, policy correlation, evidence-backed AI investigation, and human approval workflows into a single operational intelligence platform.

## Resume Bullet

Built SecureTheCloud Operational Intelligence Fabric, a governed AI operations platform for telemetry correlation, runtime investigation, policy-aware agent orchestration, evidence-backed workflows, and human-in-the-loop operational approvals across Kubernetes and cloud-native environments.
