# Every argument the kind supports. From config-schemas examples/ServiceAccount/advanced.yaml.

resource "vulnetix_service_account" "ci_scanner" {
  name = "ci-scanner"
  labels = {
    pipeline = "github-actions"
  }
  annotations = {
    owner = "platform-team"
  }
}

resource "vulnetix_service_account" "release_signer" {
  name = "release-signer"
}
