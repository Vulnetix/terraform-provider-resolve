# Every argument the kind supports. From config-schemas examples/Crew/advanced.yaml.

resource "vulnetix_crew" "release_crew" {
  name         = "release-crew"
  description  = "Reviews, tests and writes the notes for every release branch."
  display_name = "Release crew"
  members = [
    {
      profile  = "reviewer"
      replicas = 2
    },
    {
      profile = "security-reviewer"
    },
  ]
  one_per_repo = true
  annotations = {
    owner = "platform-team"
  }
}
