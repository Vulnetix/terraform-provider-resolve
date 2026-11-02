# Every argument the kind supports. From config-schemas examples/MessagingConnection/advanced.yaml.

variable "platform_teams_url" {
  type      = string
  sensitive = true
  ephemeral = true
}
variable "appsec_jira_jira_api_token" {
  type      = string
  sensitive = true
  ephemeral = true
}
variable "siem_webhook_url" {
  type      = string
  sensitive = true
  ephemeral = true
}

resource "vulnetix_messaging_connection" "security_slack" {
  name         = "security-slack"
  display_name = "Security alerts (Slack)"
  paused       = false
  slack = {
    channel_id   = "C0123456789"
    channel_name = "sec-alerts"
  }
  type = "slack"
  labels = {
    team = "security"
  }
  annotations = {
    note = "Install the Slack app from status.authorizationUrl after the first apply."
  }
}

resource "vulnetix_messaging_connection" "platform_teams" {
  name           = "platform-teams"
  display_name   = "Platform team (Teams)"
  type           = "teams"
  url_wo         = var.platform_teams_url
  url_wo_version = 1
}

resource "vulnetix_messaging_connection" "appsec_jira" {
  name         = "appsec-jira"
  display_name = "AppSec Jira"
  jira = {
    api_token_wo         = var.appsec_jira_jira_api_token
    api_token_wo_version = 1
    default_issue_type   = "Bug"
    default_project_key  = "SEC"
    email                = "appsec-bot@acme.example"
    site_url             = "https://acme.atlassian.net"
  }
  type = "jira"
}

resource "vulnetix_messaging_connection" "siem_webhook" {
  name           = "siem-webhook"
  display_name   = "SIEM intake"
  paused         = true
  type           = "webhook"
  url_wo         = var.siem_webhook_url
  url_wo_version = 1
}
