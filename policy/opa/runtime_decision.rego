package securethecloud.runtime

default decision := {
  "allow": true,
  "severity": "low",
  "reason": "No policy violation detected"
}

decision := {
  "allow": false,
  "severity": "critical",
  "reason": "Secret file access attempt inside payment workload"
} if {
  input.event_type == "suspicious_process_exec"
  input.asset.name == "payment-api"
  input.process.name == "sh"
  contains(input.process.cmdline, "secrets")
}

decision := {
  "allow": false,
  "severity": "high",
  "reason": "Unexpected shell execution in sensitive workload"
} if {
  input.event_type == "suspicious_process_exec"
  input.asset.name == "payment-api"
  input.process.name == "sh"
  not contains(input.process.cmdline, "secrets")
}
