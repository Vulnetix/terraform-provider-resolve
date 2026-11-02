# Only the arguments the kind requires. From config-schemas examples/PackageList/minimal.yaml.

resource "vulnetix_package_list" "banned" {
  name = "banned"
  entries = [
    {
      purl = "pkg:npm/event-stream@3.3.6"
    },
  ]
  list = "deny"
}
