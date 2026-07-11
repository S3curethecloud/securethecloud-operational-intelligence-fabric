CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE IF NOT EXISTS assets (
  id UUID PRIMARY KEY,
  name TEXT NOT NULL,
  asset_type TEXT NOT NULL,
  environment TEXT NOT NULL,
  owner TEXT,
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS runtime_events (
  id UUID PRIMARY KEY,
  asset_id UUID REFERENCES assets(id),
  event_type TEXT NOT NULL,
  severity TEXT NOT NULL,
  source TEXT NOT NULL,
  message TEXT NOT NULL,
  raw_event JSONB NOT NULL,
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS policy_decisions (
  id UUID PRIMARY KEY,
  runtime_event_id UUID REFERENCES runtime_events(id),
  policy_name TEXT NOT NULL,
  decision TEXT NOT NULL,
  reason TEXT NOT NULL,
  opa_input JSONB NOT NULL,
  opa_result JSONB NOT NULL,
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS incidents (
  id UUID PRIMARY KEY,
  title TEXT NOT NULL,
  severity TEXT NOT NULL,
  status TEXT NOT NULL,
  risk_score INTEGER NOT NULL,
  summary TEXT,
  created_at TIMESTAMPTZ DEFAULT now(),
  updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS incident_events (
  incident_id UUID REFERENCES incidents(id),
  runtime_event_id UUID REFERENCES runtime_events(id),
  PRIMARY KEY (incident_id, runtime_event_id)
);

CREATE TABLE IF NOT EXISTS investigation_runs (
  id UUID PRIMARY KEY,
  incident_id UUID REFERENCES incidents(id),
  status TEXT NOT NULL,
  ai_summary TEXT,
  evidence JSONB NOT NULL,
  model_provider TEXT NOT NULL,
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS approvals (
  id UUID PRIMARY KEY,
  incident_id UUID REFERENCES incidents(id),
  recommendation TEXT NOT NULL,
  decision TEXT NOT NULL,
  reviewer TEXT NOT NULL,
  rationale TEXT NOT NULL,
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS audit_log (
  id UUID PRIMARY KEY,
  actor TEXT NOT NULL,
  action TEXT NOT NULL,
  target_type TEXT NOT NULL,
  target_id UUID NOT NULL,
  metadata JSONB NOT NULL,
  created_at TIMESTAMPTZ DEFAULT now()
);
