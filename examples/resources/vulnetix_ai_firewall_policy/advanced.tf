# Every argument the kind supports. From config-schemas examples/AiFirewallPolicy/advanced.yaml.

resource "vulnetix_ai_firewall_policy" "default" {
  baseline = {
    enabled = true
    exclude = ["pii-phone"]
    ref     = "recommended"
  }
  capabilities = {
    approved-mcp = {
      action   = "block"
      family   = "mcp"
      list     = "allow"
      patterns = ["github", "linear*"]
    }
    banned-clients = {
      action   = "block"
      family   = "client"
      list     = "deny"
      patterns = ["unknown"]
      priority = 10
    }
  }
  guardrails = [
    {
      action      = "block"
      baseline_id = "secrets-aws-access-key"
      enabled     = true
      name        = "block-aws-keys"
      pattern     = "(?i)AKIA[0-9A-Z]{16}"
      priority    = 10
      rule_type   = "blocked_pattern"
    },
    {
      action    = "block"
      name      = "cap-conversation"
      pattern   = "200"
      priority  = 50
      rule_type = "max_messages"
    },
    {
      action    = "redact"
      name      = "redact-pii"
      priority  = 20
      rule_type = "pii_redact"
    },
    {
      action    = "block"
      name      = "sealed-blocks"
      priority  = 5
      rule_type = "delimiter_integrity"
    },
    {
      action    = "strip"
      name      = "strip-shell-for-codex"
      pattern   = "Bash"
      priority  = 30
      rule_type = "tool_deny"
      when      = "client == 'codex'"
    },
    {
      action    = "flag"
      enabled   = false
      name      = "flag-filesystem-writes"
      priority  = 40
      rule_type = "expression"
      when      = "'filesystem' in mcpServers && tools.exists(t, t.startsWith('mcp__filesystem__write'))"
    },
  ]
  hold_policy = {
    hold_timeout_action  = "deny"
    hold_timeout_minutes = 240
    hold_unknown         = true
  }
  models = [
    {
      action       = "allow"
      any_provider = false
      provider     = "anthropic"
      slug         = "claude-sonnet-4-5"
    },
    {
      action   = "allow"
      provider = "acme-vllm"
      slug     = "qwen3-coder-30b"
    },
    {
      action       = "deny"
      any_provider = true
      slug         = "gpt-4o-mini"
    },
  ]
  providers = [
    {
      action = "allow"
      key = {
        from_env = "ANTHROPIC_PROVIDER_KEY"
      }
      slug = "anthropic"
    },
    {
      action = "allow"
      key = {
        from_file = "~/.config/acme/vllm.key"
      }
      slug = "acme-vllm"
    },
    {
      action = "deny"
      slug   = "openrouter"
    },
    {
      action = "default"
      slug   = "groq"
    },
  ]
  prune = true
  settings = {
    logs_enabled = true
  }
  labels = {
    team = "platform-security"
  }
  annotations = {
    ticket = "https://tracker.example.com/SEC-2210"
  }
}
