# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/TeaComponent/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_tea_component" "payments_api" {
  name           = "payments-api"
  component_name = "payments-api"
}
