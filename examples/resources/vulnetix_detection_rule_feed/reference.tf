# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/DetectionRuleFeed/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_detection_rule_feed" "edge_waf" {
  name = "edge-waf"
  entries = [
    {
      rule       = "9c2d4e6f-1a3b-4c5d-8e7f-0a1b2c3d4e5f"
      rule_table = "virtual_patch"
      rule_type  = "AWS_WAF"
    },
  ]
  format         = "AWS_WAF"
  publish        = false
  schema_version = 1
}
