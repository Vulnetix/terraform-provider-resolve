# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/AwsConnection/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_aws_connection" "default" {
  advanced_networking = false
  region              = "ap-southeast-2"
  sso = {
    start_url = "https://acme.awsapps.com/start"
  }
}
