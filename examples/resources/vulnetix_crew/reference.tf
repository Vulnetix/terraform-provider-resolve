# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/Crew/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_crew" "review_crew" {
  name = "review-crew"
  members = [
    {
      profile = "reviewer"
    },
  ]
}
