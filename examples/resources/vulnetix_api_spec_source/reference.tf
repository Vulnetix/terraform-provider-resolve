# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/ApiSpecSource/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_api_spec_source" "payments_openapi" {
  name     = "payments-openapi"
  schedule = "daily"
  source = {
    url = {
      url = "https://api.acme.com/openapi.json"
    }
  }
}
