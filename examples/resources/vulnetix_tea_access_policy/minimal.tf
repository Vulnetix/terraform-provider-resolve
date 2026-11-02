# Only the arguments the kind requires. From config-schemas examples/TeaAccessPolicy/minimal.yaml.

resource "vulnetix_tea_access_policy" "acme_payments" {
  name       = "acme-payments"
  product    = "acme-payments"
  visibility = "private"
}
