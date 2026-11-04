# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/PackageList/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_package_list" "banned" {
  name = "banned"
  entries = [
    {
      purl = "pkg:npm/event-stream@3.3.6"
    },
  ]
  list = "deny"
}
