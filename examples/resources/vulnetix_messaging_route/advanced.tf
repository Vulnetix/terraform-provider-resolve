# Every argument the kind supports. From config-schemas examples/MessagingRoute/advanced.yaml.

resource "vulnetix_messaging_route" "critical_findings" {
  name = "critical-findings"
  actor_filter = {
    kinds      = ["service_account"]
    mode       = "exclude"
    principals = ["7d3f2a10-4b1c-4e8a-9f20-6a5b4c3d2e1f"]
  }
  allow_sensitive_webhook = false
  digest = {
    minutes = 15
    mode    = "immediate"
    tiers = {
      error     = "instant"
      info      = "off"
      needs_you = "instant"
      warning   = "digest"
    }
  }
  display_name = "Critical findings to security"
  enabled      = true
  event_types  = ["scan.completed", "scan.build_broken"]
  filters = {
    findings = {
      cess_threshold    = 5
      cvss_severities   = ["critical", "high"]
      cvss_threshold    = 7
      eol_days          = 90
      epss_threshold    = 0.1
      exploit_threshold = 1
      has_exploits      = ["functional", "weaponized", "in_the_wild"]
      ignore_test       = true
      is_eol            = false
      is_exploitable    = true
      is_malware        = false
      is_reachable      = true
      new_only          = true
      only_direct       = false
      only_production   = true
    }
    scope = {
      app_name    = "payments-api"
      app_version = "2.4.0"
      repo        = "acme/payments-api"
      upload_uuid = "0b6f6c1e-8f4a-4c55-9a0e-3c2d1b7a9e10"
    }
  }
  include_pii        = false
  max_per_hour       = 120
  mention_principals = ["2c9e8d7f-1a2b-4c3d-8e4f-5a6b7c8d9e0f"]
  mention_static = {
    discord     = ["123456789012345678"]
    github      = ["octocat"]
    google_chat = ["112233445566"]
    jira        = ["5b10ac8d82e05b22cc7d4ef5"]
    slack       = ["U0123ABCD"]
    teams       = ["secops@acme.example"]
  }
  min_severity = "warning"
  mode         = "live"
  targets = [
    {
      connection = "security-slack"
      enabled    = true
      options = {
        channel_id = "C0123456789"
      }
    },
    {
      connection = "appsec-jira"
      options = {
        issue_type  = "Bug"
        labels      = ["vulnetix", "critical"]
        project_key = "SEC"
      }
    },
    {
      connection = "github-acme"
      options = {
        assign_mentions = true
        labels          = ["vulnetix"]
        owner           = "acme"
        repo            = "security-alerts"
      }
    },
  ]
  template         = "security-findings"
  template_version = 1
  labels = {
    team = "security"
  }
  annotations = {
    runbook = "https://wiki.example.com/security/alerts"
  }
}

resource "vulnetix_messaging_route" "nightly_digest" {
  name = "nightly-digest"
  digest = {
    hour = 8
    mode = "daily"
    tz   = "Australia/Sydney"
  }
  event_types     = ["pix.*"]
  mode            = "shadow"
  system_template = "generic"
}
