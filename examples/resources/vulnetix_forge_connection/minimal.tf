# Only the arguments the kind requires. From config-schemas examples/ForgeConnection/minimal.yaml.

resource "vulnetix_forge_connection" "codeberg" {
  name = "codeberg"
  host = "codeberg.org"
  type = "forgejo"
}
