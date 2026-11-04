# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/FactoryWorkflow/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_factory_workflow" "patch_on_demand" {
  name    = "patch-on-demand"
  enabled = true
  graph = {
    nodes = [
      {
        config = jsonencode({
          kind = "manual"
        })
        id    = "start"
        title = "Run by hand"
        type  = "trigger"
      },
    ]
  }
  publish = true
}
