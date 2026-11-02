# Only the arguments the kind requires. From config-schemas examples/Crew/minimal.yaml.

resource "vulnetix_crew" "review_crew" {
  name = "review-crew"
  members = [
    {
      profile = "reviewer"
    },
  ]
}
