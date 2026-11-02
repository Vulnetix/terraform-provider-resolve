# Every argument the kind supports. From config-schemas examples/Suppression/advanced.yaml.

resource "vulnetix_suppression" "lodash_template_not_reachable" {
  name          = "lodash-template-not-reachable"
  expires       = "2027-06-30T00:00:00Z"
  justification = "vulnerable_code_not_in_execute_path"
  reason        = "lodash.template is never called with user input; the build only renders static templates."
  scope = {
    branch     = "main"
    category   = "sca"
    line       = 1204
    line_range = "1204-1210"
    path       = "web/package-lock.json"
    purl       = "pkg:npm/lodash@4.17.20"
    repository = "acme/storefront"
    rule       = "CVE-2021-23337"
    target     = "pkg:npm/lodash"
  }
  type = "false_positive"
  when = "severity != 'critical' && !kev && repo == 'acme/storefront'"
  labels = {
    team = "payments"
  }
  annotations = {
    ticket = "https://jira.example.com/browse/SEC-1234"
  }
}
