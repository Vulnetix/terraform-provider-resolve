# Every argument the kind supports. From config-schemas examples/TeaComponent/advanced.yaml.

resource "vulnetix_tea_component" "payments_api" {
  name           = "payments-api"
  component_name = "payments-api"
  identifiers = [
    {
      id_type  = "PURL"
      id_value = "pkg:oci/payments-api?repository_url=ghcr.io/acme"
    },
  ]
  labels = {
    team = "payments"
  }
}
