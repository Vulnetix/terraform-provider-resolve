# Every argument the kind supports. From config-schemas examples/AwsAccount/advanced.yaml.

resource "vulnetix_aws_account" "production" {
  name         = "production"
  account_id   = "123456789012"
  account_name = "acme-production"
  deploy       = true
  enabled      = true
  role_name    = "AdministratorAccess"
  labels = {
    environment = "production"
  }
  annotations = {
    owner = "platform-team@example.com"
  }
}
