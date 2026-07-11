# Verified Output

This document records validated lab execution evidence for the SecureTheCloud Operational Intelligence Fabric.

## MVP Smoke Test — 2026-07-11

Validated end-to-end governed AI operations flow.

### What was tested

1. Simulated suspicious runtime event from `falco-simulated`.
2. Ingested event through the Go API at `/v1/events`.
3. Evaluated live OPA policy through `securethecloud.runtime.decision`.
4. Returned live OPA policy decision:
   - allow: false
   - severity: critical
   - reason: Secret file access attempt inside payment workload
5. Created a high-risk incident for `payment-api`.
6. Calculated risk score of `100`.
7. Generated AI-assisted investigation summary.
8. Required human approval.
9. Recorded human approval from `demo.reviewer@securethecloud.dev`.
10. Replayed the full evidence chain.

### Evidence chain

runtime_event -> policy_decision -> risk_score -> ai_investigation -> human_approval

### OPA verification

OPA fallback status:

fallback: false

Verified policy name:

securethecloud.runtime.decision

### Governance boundary verified

The lab confirms that:

- AI summarizes evidence and recommends next steps.
- OPA evaluates policy context.
- Human approval is required for high-risk actions.
- The lab does not claim production enforcement.
- The lab does not claim SOC 2 certification.
- Evidence, explanation, and packaging do not create authority.

### MVP status

Phase 4 — OPA Policy Correlation: complete.

Phase 5 — Go Event API: complete.

Phase 6 — AI Investigation Service: complete.

Phase 7 — Human Approval Workflow: complete.

Phase 9 — Evidence Replay: complete.
