# Only the arguments the kind requires. From config-schemas examples/TeaProduct/minimal.yaml.

resource "vulnetix_tea_product" "acme_payments" {
  name         = "acme-payments"
  product_name = "Acme Payments"
}
