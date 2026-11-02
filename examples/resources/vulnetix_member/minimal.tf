# Only the arguments the kind requires. From config-schemas examples/Member/minimal.yaml.

resource "vulnetix_member" "alice" {
  name  = "alice"
  email = "alice@acme.com"
}
