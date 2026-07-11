package securethecloud.approvals

default requires_human_approval := false

requires_human_approval if {
  input.risk_score >= 70
}

requires_human_approval if {
  input.severity == "critical"
}

requires_human_approval if {
  input.asset_type == "kubernetes_workload"
  input.environment == "dev-banking"
}
