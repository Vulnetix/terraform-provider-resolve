# Only the arguments the kind requires. From config-schemas examples/TeaComponent/minimal.yaml.

resource "vulnetix_tea_component" "payments_api" {
  name           = "payments-api"
  component_name = "payments-api"
}
