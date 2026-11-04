# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/DetectionRule/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_detection_rule" "sigmahq_aws_root_login" {
  name               = "sigmahq-aws-root-login"
  enabled            = true
  interval_sec       = 60
  lateness_sec       = 300
  max_alerts_per_run = 50
  sigmahq            = "rules/cloud/aws/cloudtrail/aws_root_account_usage.yml"
  suppress_sec       = 900
}
