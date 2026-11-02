# Every argument the kind supports. From config-schemas examples/VaultEntry/advanced.yaml.

variable "anthropic_production_value" {
  type      = string
  sensitive = true
  ephemeral = true
}
variable "npm_token_value" {
  type      = string
  sensitive = true
  ephemeral = true
}
variable "openai_admin_value" {
  type      = string
  sensitive = true
  ephemeral = true
}

resource "vulnetix_vault_entry" "anthropic_production" {
  name       = "anthropic-production"
  entry_name = "Anthropic production"
  expires_at = "2027-06-30T00:00:00Z"
  grants = [
    {
      all     = true
      effect  = "allow"
      surface = "ai_firewall"
    },
    {
      effect  = "allow"
      ids     = ["3f1c2a9e-8b4d-4c7a-9e21-5d6f7a8b9c01", "7a2e4b6c-1d3f-4e5a-8b7c-9d0e1f2a3b4c"]
      surface = "pix_sandbox"
    },
  ]
  provider_slug    = "anthropic"
  type             = "provider_key"
  value_wo         = var.anthropic_production_value
  value_wo_version = 1
  labels = {
    team = "platform"
  }
  annotations = {
    ticket = "https://tracker.example.com/SEC-218"
  }
}

resource "vulnetix_vault_entry" "npm_token" {
  name       = "npm-token"
  entry_name = "NPM_TOKEN"
  grants = [
    {
      all     = true
      surface = "belai_host"
    },
    {
      effect  = "deny"
      ids     = ["0b1c2d3e-4f5a-4b6c-8d7e-9f0a1b2c3d4e"]
      surface = "belai_host"
    },
  ]
  type             = "env"
  value_wo         = var.npm_token_value
  value_wo_version = 1
}

resource "vulnetix_vault_entry" "openai_admin" {
  name             = "openai-admin"
  entry_name       = "OpenAI admin"
  provider_slug    = "openai"
  type             = "management"
  value_wo         = var.openai_admin_value
  value_wo_version = 1
}
