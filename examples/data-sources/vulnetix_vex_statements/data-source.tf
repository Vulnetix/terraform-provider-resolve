# The tenant's not_affected statements, newest first.
data "vulnetix_vex_statements" "not_affected" {
  status = "not_affected"
  limit  = 200
}

output "not_affected_vulnerabilities" {
  value = distinct([for s in data.vulnetix_vex_statements.not_affected.statements : s.vulnerability_id])
}
