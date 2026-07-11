# GitHub Issues

## Phase 0 - Create repository foundation

- [ ] Add README.md
- [ ] Add repo structure
- [ ] Add .env.example
- [ ] Add docker-compose.yml
- [ ] Add Makefile
- [ ] Add docs folder
- [ ] Add GitHub topics
- [ ] Add milestones

Acceptance:

- Repo can be cloned.
- Local services can start with `make up`.

## Phase 1 - Local services and schema

- [ ] Start PostgreSQL, Redis, NATS, and OPA with Docker Compose
- [ ] Add operational intelligence schema
- [ ] Add sample runtime event
- [ ] Add event simulation script

Acceptance:

- `make up` works.
- `make event` posts sample event after API is running.

## Phase 4 - OPA policy correlation

- [ ] Add runtime policy
- [ ] Add approval policy
- [ ] API calls OPA for event context
- [ ] Store or return policy reason

Acceptance:

- Suspicious payment-api shell event returns high-risk policy context.

## Phase 5 - Go operational API

- [ ] Implement runtime event ingestion
- [ ] Implement incident list/detail
- [ ] Implement investigation trigger
- [ ] Implement approval endpoint
- [ ] Implement audit log
- [ ] Implement evidence replay endpoint

Acceptance:

- User can drive full flow with curl.

## Phase 6 - Python AI investigation service

- [ ] Add FastAPI service
- [ ] Add deterministic mock provider
- [ ] Generate evidence-backed summary
- [ ] Return recommendation with human approval requirement

Acceptance:

- Investigation never claims enforcement or authorization.

## Phase 7 - Human approval workflow

- [ ] Add approval decisions
- [ ] Add reviewer and rationale
- [ ] Add audit log record
- [ ] Include approval in evidence replay

Acceptance:

- Every accepted recommendation has a human reviewer and rationale.

## Phase 8 - Next.js dashboard

- [ ] Add overview page
- [ ] Add incident queue
- [ ] Add incident detail
- [ ] Add evidence replay panel
- [ ] Add approval panel

Acceptance:

- Demo can be run from browser.

## Phase 9 - Evidence replay

- [ ] Return ordered evidence chain
- [ ] Include event, policy, investigation, approval, and audit
- [ ] Render timeline in UI

Acceptance:

- Reviewer can reconstruct what happened and why.
