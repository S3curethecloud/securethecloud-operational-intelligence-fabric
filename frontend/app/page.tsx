const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

async function getIncidents() {
  try {
    const res = await fetch(`${API_URL}/v1/incidents`, { cache: "no-store" });
    if (!res.ok) return [];
    return res.json();
  } catch {
    return [];
  }
}

async function getAuditLog() {
  try {
    const res = await fetch(`${API_URL}/v1/audit-log`, { cache: "no-store" });
    if (!res.ok) return [];
    return res.json();
  } catch {
    return [];
  }
}

export default async function Home() {
  const incidents = await getIncidents();
  const audit = await getAuditLog();
  const highRisk = incidents.filter((i: any) => i.risk_score >= 70);

  return (
    <main style={{ padding: 32, maxWidth: 1200, margin: "0 auto" }}>
      <p className="badge">Secure AI Banking Operations</p>
      <h1>SecureTheCloud Operational Intelligence Fabric</h1>
      <p className="muted">
        Governed AI operations workspace for runtime telemetry, OPA policy context,
        evidence-backed investigation, and human approval workflows.
      </p>

      <section className="grid" style={{ marginTop: 24 }}>
        <div className="card">
          <h2>{incidents.length}</h2>
          <p className="muted">Total incidents</p>
        </div>
        <div className="card">
          <h2>{highRisk.length}</h2>
          <p className="muted">High-risk incidents</p>
        </div>
        <div className="card">
          <h2>{audit.length}</h2>
          <p className="muted">Audit records</p>
        </div>
      </section>

      <section style={{ marginTop: 24 }} className="card">
        <h2>Incident Queue</h2>
        {incidents.length === 0 ? (
          <p className="muted">No incidents yet. Run ./scripts/simulate_runtime_event.sh.</p>
        ) : (
          <div style={{ display: "grid", gap: 12 }}>
            {incidents.map((incident: any) => (
              <a className="card" key={incident.id} href={`/incidents/${incident.id}`}>
                <strong>{incident.title}</strong>
                <p className="muted">Risk {incident.risk_score} - {incident.status}</p>
              </a>
            ))}
          </div>
        )}
      </section>

      <section style={{ marginTop: 24 }} className="card">
        <h2>Governance Boundary</h2>
        <p className="muted">
          AI recommendations are evidence summaries only. Human approval is required before
          an operational recommendation is accepted. Runtime-impacting enforcement must follow
          the configured policy control path.
        </p>
      </section>
    </main>
  );
}
