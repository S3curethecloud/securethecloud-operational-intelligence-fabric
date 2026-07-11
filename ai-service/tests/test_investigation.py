from fastapi.testclient import TestClient
from app.main import app


def test_investigation_requires_human_approval():
    client = TestClient(app)
    response = client.post(
        "/v1/investigate",
        json={
            "incident": {"id": "inc_test", "risk_score": 97},
            "runtime_events": [{"id": "evt_test", "asset": {"name": "payment-api"}}],
            "policy_decision": {"id": "pol_test", "reason": "critical policy context"},
        },
    )
    assert response.status_code == 200
    body = response.json()
    assert body["requires_human_approval"] is True
    assert "evt_test" in body["evidence_refs"]
