# Only the arguments the kind requires. From config-schemas examples/DetectionRule/minimal.yaml.

resource "vulnetix_detection_rule" "sigmahq_aws_root_login" {
  name    = "sigmahq-aws-root-login"
  sigmahq = "rules/cloud/aws/cloudtrail/aws_root_account_usage.yml"
}
