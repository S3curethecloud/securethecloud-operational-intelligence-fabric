const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

async function getIncident(id: string) {
  const res = await fetch(`${API_URL}/v1/incidents/${id}`, { cache: "no-store" });
  if (!res.ok) return null;
  return res.json();
}

async function getEvidence(id: string) {
  const res = await fetch(`${API_URL}/v1/incidents/${id}/evidence`, { cache: "no-store" });
  if (!res.ok) return { chain: [] };
  return res.json();
}

export default async function IncidentPage({ params }: { params: { id: string } }) {
  const incident = await getIncident(params.id);
  const evidence = await getEvidence(params.id);

  if (!incident) {
    return <main style={{ padding: 32 }}>Incident not found.</main>;
  }

  return (
    <main style={{ padding: 32, maxWidth: 1100, margin: "0 auto" }}>
      <a className="muted" href="/">← Back</a>
      <h1>{incident.title}</h1>
      <p className="badge">Risk Score: {incident.risk_score}</p>
      <p className="badge" style={{ marginLeft: 8 }}>Status: {incident.status}</p>

      <section className="grid" style={{ marginTop: 24 }}>
        <div className="card">
          <h2>Policy Decision</h2>
          <p>{incident.policy_decision.reason}</p>
          <p className="muted">Severity: {incident.policy_decision.severity}</p>
        </div>
        <div className="card">
          <h2>AI Investigation</h2>
          {incident.investigation ? (
            <>
              <p>{incident.investigation.summary}</p>
              <p className="muted">Recommendation: {incident.investigation.recommended_action}</p>
            </>
          ) : (
            <p className="muted">Run investigation from CLI: ./scripts/run_investigation.sh {incident.id}</p>
          )}
        </div>
      </section>

      <section className="card" style={{ marginTop: 24 }}>
        <h2>Evidence Replay</h2>
        <ol>
          {evidence.chain.map((item: any) => (
            <li key={`${item.step}-${item.id}`} style={{ marginBottom: 12 }}>
              <strong>{item.type}</strong>: {item.summary}
              <div className="muted">Source: {item.source}</div>
            </li>
          ))}
        </ol>
      </section>

      <section className="card" style={{ marginTop: 24 }}>
        <h2>Human Approval</h2>
        <p className="muted">
          Record approval from CLI: ./scripts/approve_incident.sh {incident.id} approved
        </p>
        {incident.approvals?.length ? (
          incident.approvals.map((approval: any) => (
            <div className="card" key={approval.id}>
              <strong>{approval.decision}</strong>
              <p>{approval.rationale}</p>
              <p className="muted">Reviewer: {approval.reviewer}</p>
            </div>
          ))
        ) : (
          <p>No approval recorded yet.</p>
        )}
      </section>
    </main>
  );
}
