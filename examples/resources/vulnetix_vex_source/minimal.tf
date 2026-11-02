# Only the arguments the kind requires. From config-schemas examples/VexSource/minimal.yaml.

resource "vulnetix_vex_source" "local_vex" {
  name = "local-vex"
  path = {
    path = ".vulnetix/vex.openvex.json"
  }
}
