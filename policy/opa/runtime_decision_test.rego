package securethecloud.runtime

test_secret_access_is_critical if {
  result := decision with input as {
    "event_type": "suspicious_process_exec",
    "asset": {"name": "payment-api"},
    "process": {"name": "sh", "cmdline": "sh -c cat /etc/secrets/payment.env"}
  }
  result.allow == false
  result.severity == "critical"
}
