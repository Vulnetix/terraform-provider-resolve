# Every argument the kind supports. From config-schemas examples/AwsConnection/advanced.yaml.

resource "vulnetix_aws_connection" "default" {
  advanced_networking = true
  region              = "ap-southeast-2"
  sso = {
    start_url = "https://acme.awsapps.com/start"
  }
  labels = {
    owner = "platform"
  }
  annotations = {
    note = "Approve the device authorization in the AWS access portal on first apply."
  }
}

# Or, instead (one per tenant, so declare only one):
#
# variable "default_2_iam_user_secret_access_key" {
#   type      = string
#   sensitive = true
#   ephemeral = true
# }
# resource "vulnetix_aws_connection" "default_2" {
#   advanced_networking = false
#   iam_user = {
#     access_key_id                = "AKIAIOSFODNN7EXAMPLE"
#     secret_access_key_wo         = var.default_2_iam_user_secret_access_key
#     secret_access_key_wo_version = 1
#   }
#   region = "us-east-1"
# }
