# The tenant's not_affected statements about one package, newest first.
data "vulnetix_vex_statements" "lodash" {
  status = "not_affected"
  purl   = "pkg:npm/lodash"
  limit  = 200
}

output "lodash_not_affected" {
  value = { for s in data.vulnetix_vex_statements.lodash.statements : s.vulnerability_id => s.justification }
}
