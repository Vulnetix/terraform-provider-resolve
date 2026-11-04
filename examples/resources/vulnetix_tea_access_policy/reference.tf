# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/TeaAccessPolicy/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_tea_access_policy" "acme_payments" {
  name       = "acme-payments"
  product    = "acme-payments"
  visibility = "private"
}
