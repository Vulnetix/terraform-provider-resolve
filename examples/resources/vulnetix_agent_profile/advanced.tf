# Every argument the kind supports. From config-schemas examples/AgentProfile/advanced.yaml.

resource "vulnetix_agent_profile" "security_reviewer" {
  name = "security-reviewer"
  body = "---\nid: 0b7e4d2c-9a1f-4c3e-8d5b-2f6a7e8c9d01\nname: security-reviewer\ndisplay_name: Security reviewer\npalette: [\"#0f172a\", \"#1e293b\", \"#38bdf8\", \"#f8fafc\"]\n---\nYou review pull requests for injection, authorization and secrets\nhandling, and cite the checklist in checklists/review.md.\n"
  files = [
    {
      body = "- Inputs are validated at the boundary.\n- Every route checks the caller's role.\n"
      path = "checklists/review.md"
    },
    {
      file = "agents/security-reviewer/secrets.md"
      path = "checklists/secrets.md"
    },
  ]
  labels = {
    team = "appsec"
  }
}

resource "vulnetix_agent_profile" "release_notes" {
  name = "release-notes"
  file = "agents/release-notes.md"
}
