# Only the arguments the kind requires. From config-schemas examples/FactoryWorkflow/minimal.yaml.

resource "vulnetix_factory_workflow" "patch_on_demand" {
  name = "patch-on-demand"
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
}
