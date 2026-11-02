# Every argument the kind supports. From config-schemas examples/TeaProduct/advanced.yaml.

resource "vulnetix_tea_product" "acme_payments" {
  name = "acme-payments"
  identifiers = [
    {
      id_type  = "PURL"
      id_value = "pkg:github/acme/payments-api"
    },
    {
      id_type  = "CPE"
      id_value = "cpe:2.3:a:acme:payments:*:*:*:*:*:*:*:*"
    },
  ]
  product_name = "Acme Payments"
  labels = {
    team = "payments"
  }
  annotations = {
    note = "Releases are published by the payments-api release workflow."
  }
}
