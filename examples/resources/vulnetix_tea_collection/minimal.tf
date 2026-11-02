# Only the arguments the kind requires. From config-schemas examples/TeaCollection/minimal.yaml.

resource "vulnetix_tea_collection" "acme_payments_2_4_0" {
  name    = "acme-payments-2-4-0"
  product = "acme-payments"
  release = "2.4.0"
}
