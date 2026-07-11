# Doctrine Boundaries

This lab is designed to preserve SecureTheCloud governance boundaries.

## AI Boundary

The AI service is an evidence summarization and recommendation layer.

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
- issue tokens
- create runtime sessions
- silently approve remediation

## Policy Boundary

OPA provides deterministic policy context for the lab.

The MVP API records policy decisions and reasons. The AI service may explain this policy context, but the AI service does not replace policy evaluation.

## Human Approval Boundary

High-risk recommendations require a human reviewer decision before they are accepted as an operational decision.

Allowed reviewer decisions:

- approved
- rejected
- needs_more_evidence

## Evidence Boundary

Evidence replay should reconstruct what happened and why a recommendation was made. It must not claim live production enforcement unless a future governed phase explicitly implements and approves that behavior.
