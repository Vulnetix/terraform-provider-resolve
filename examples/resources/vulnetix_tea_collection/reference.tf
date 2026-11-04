# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/TeaCollection/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_tea_collection" "acme_payments_2_4_0" {
  name    = "acme-payments-2-4-0"
  product = "acme-payments"
  release = "2.4.0"
}
