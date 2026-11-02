# Only the arguments the kind requires. From config-schemas examples/AwsAccount/minimal.yaml.

resource "vulnetix_aws_account" "production" {
  name       = "production"
  account_id = "123456789012"
}
