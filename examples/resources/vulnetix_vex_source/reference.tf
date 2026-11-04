# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/VexSource/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_vex_source" "local_vex" {
  name     = "local-vex"
  format   = "auto"
  on_error = "keep-last"
  path = {
    path = ".vulnetix/vex.openvex.json"
  }
  trust = {
    require_signature = false
  }
}
