# Only the arguments the kind requires. From config-schemas examples/JailPolicy/minimal.yaml.

resource "vulnetix_jail_policy" "org_default" {
  name = "org-default"
}
