# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/AwsAccount/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_aws_account" "production" {
  name       = "production"
  account_id = "123456789012"
  deploy     = false
  enabled    = true
}
