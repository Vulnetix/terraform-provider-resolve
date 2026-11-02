# Only the arguments the kind requires. From config-schemas examples/ApiSpecSource/minimal.yaml.

resource "vulnetix_api_spec_source" "payments_openapi" {
  name = "payments-openapi"
  source = {
    url = {
      url = "https://api.acme.com/openapi.json"
    }
  }
}
