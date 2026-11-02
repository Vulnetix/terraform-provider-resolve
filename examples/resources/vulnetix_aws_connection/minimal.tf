# Only the arguments the kind requires. From config-schemas examples/AwsConnection/minimal.yaml.

resource "vulnetix_aws_connection" "default" {
  region = "ap-southeast-2"
  sso = {
    start_url = "https://acme.awsapps.com/start"
  }
}
