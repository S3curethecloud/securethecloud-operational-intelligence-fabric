from __future__ import annotations

from datetime import datetime, timezone
from typing import Any
from uuid import uuid4

from fastapi import FastAPI
from pydantic import BaseModel, Field

app = FastAPI(title="SecureTheCloud OIF AI Investigation Service")


class InvestigationRequest(BaseModel):
    incident: dict[str, Any]
    runtime_events: list[dict[str, Any]] = Field(default_factory=list)
    policy_decision: dict[str, Any] = Field(default_factory=dict)
    boundary: str | None = None


@app.get("/health")
def health() -> dict[str, str]:
    return {"status": "ok"}


@app.post("/v1/investigate")
def investigate(request: InvestigationRequest) -> dict[str, Any]:
    incident = request.incident
    events = request.runtime_events
    policy = request.policy_decision
    event_ids = [event.get("id", "unknown-event") for event in events]

    asset_name = "unknown asset"
    if events:
        asset_name = events[0].get("asset", {}).get("name", asset_name)

    policy_reason = policy.get("reason", "No policy reason provided")
    risk_score = incident.get("risk_score", "unknown")

    return {
        "id": f"inv_{uuid4().hex[:16]}",
        "incident_id": incident.get("id"),
        "status": "completed",
        "summary": (
            f"Evidence indicates suspicious runtime behavior on {asset_name}. "
            f"The incident has risk score {risk_score}. Policy context: {policy_reason}."
        ),
        "hypothesis": (
            "Likely unauthorized shell execution, misconfigured operational access, "
            "or a lab-generated runtime detection requiring reviewer validation."
        ),
        "recommended_action": (
            "Request human approval to perform lab-only containment review, collect additional workload evidence, "
            "and verify whether the activity was expected."
        ),
        "confidence": "medium",
        "evidence_refs": event_ids + [policy.get("id", "policy-decision")],
        "requires_human_approval": True,
        "created_at": datetime.now(timezone.utc).isoformat(),
    }
