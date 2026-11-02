# Every argument the kind supports. From config-schemas examples/ApiSpecSource/advanced.yaml.

resource "vulnetix_api_spec_source" "payments_openapi" {
  name     = "payments-openapi"
  owner    = "payments-team"
  schedule = "daily"
  source = {
    repo = {
      path       = "api/openapi.yaml"
      repository = "acme/payments"
    }
  }
  tags = ["payments", "public"]
  labels = {
    team = "payments"
  }
  annotations = {
    ticket = "https://tracker.example.com/API-311"
  }
}

resource "vulnetix_api_spec_source" "public_api_spec" {
  name = "public-api-spec"
  probe = {
    mode = "exact"
  }
  schedule = "hourly"
  source = {
    url = {
      url = "https://api.acme.com/openapi.json"
    }
  }
}

resource "vulnetix_api_spec_source" "staging_host" {
  name = "staging-host"
  probe = {
    authorize = true
    mode      = "authorized"
  }
  schedule = "weekly"
  source = {
    host = {
      host = "staging.acme.com:8443"
    }
  }
}

resource "vulnetix_api_spec_source" "edge_gateway" {
  name     = "edge-gateway"
  owner    = "platform-team"
  schedule = "6h"
  source = {
    connector = {
      kind = "aws_apigw"
      ref  = "a1b2c3d4e5"
    }
  }
  tags = ["edge"]
}
