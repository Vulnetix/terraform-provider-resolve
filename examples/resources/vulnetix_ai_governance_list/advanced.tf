# Every argument the kind supports. From config-schemas examples/AiGovernanceList/advanced.yaml.

resource "vulnetix_ai_governance_list" "unapproved_mcp" {
  name     = "unapproved-mcp"
  action   = "strip"
  enabled  = true
  family   = "mcp"
  list     = "deny"
  patterns = ["*-shell", "filesystem"]
  priority = 20
  labels = {
    team = "platform-security"
  }
  annotations = {
    ticket = "https://tracker.example.com/SEC-2214"
  }
}

resource "vulnetix_ai_governance_list" "approved_skills" {
  name     = "approved-skills"
  action   = "flag"
  enabled  = false
  family   = "skill"
  list     = "allow"
  patterns = ["code-review", "acme-*"]
  priority = 100
}
