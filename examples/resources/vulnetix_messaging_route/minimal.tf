# Only the arguments the kind requires. From config-schemas examples/MessagingRoute/minimal.yaml.

resource "vulnetix_messaging_route" "broken_builds" {
  name        = "broken-builds"
  event_types = ["scan.build_broken"]
}
