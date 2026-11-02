# Every argument the kind supports. From config-schemas examples/TeaCollection/advanced.yaml.

resource "vulnetix_tea_collection" "acme_payments_2_4_0" {
  name    = "acme-payments-2-4-0"
  product = "acme-payments"
  release = "2.4.0"
  annotations = {
    note = "Embargoed until the advisory is published; see the TeaAccessPolicy of the same name."
  }
}

resource "vulnetix_tea_collection" "payments_api_2_4_0" {
  name      = "payments-api-2-4-0"
  component = "payments-api"
  release   = "2.4.0"
}
