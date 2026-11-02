---
page_title: "vulnetix_messaging_connection Resource - vulnetix"
subcategory: ""
description: |-
  One destination the messaging centre can send to: a Slack channel, a Microsoft Teams workflow, a Discord or Google Chat incoming webhook, a Jira site, or a signed webhook. MessagingRoute targets, AlertTrigger, ScheduledReport and DetectionRule actions name it by metadata.name. Its credential (a webhook URL, a Jira API token, a Slack bot token) is kept in the tenant's secrets vault as a connector entry, never in the document, an export or a plan. A slack connection is authorized in a browser: apply returns status.phase AwaitingAuthorization and status.authorizationUrl, and the connection sends nothing until a person installs the Vulnetix Slack app from that URL. The type and a Jira site cannot change; declare a new connection under another name instead. Deleting it revokes its vault entry and removes it from every route that targets it. GitHub Issues connections come from the GitHub App installation and cannot be declared here, but a route can still target one by the name export gives it.
  Manages the MessagingConnection config kind through /config/objects/MessagingConnection/{name}. Applies are recorded with the provider's manager, so a GitOps vulnetix config apply --prune leaves this object alone.
  Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (most of them, such as ApiSpecSource, RiskStrategy and VexSource) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object, or lets go of it for a kind the server never deletes (such as Organization). Some kinds wait on a person after they are applied (an invited Member, a Slack connection): status.phase says so, and Terraform does not wait.
---

# vulnetix_messaging_connection (Resource)

One destination the messaging centre can send to: a Slack channel, a Microsoft Teams workflow, a Discord or Google Chat incoming webhook, a Jira site, or a signed webhook. MessagingRoute targets, AlertTrigger, ScheduledReport and DetectionRule actions name it by metadata.name. Its credential (a webhook URL, a Jira API token, a Slack bot token) is kept in the tenant's secrets vault as a connector entry, never in the document, an export or a plan. A slack connection is authorized in a browser: apply returns status.phase AwaitingAuthorization and status.authorizationUrl, and the connection sends nothing until a person installs the Vulnetix Slack app from that URL. The type and a Jira site cannot change; declare a new connection under another name instead. Deleting it revokes its vault entry and removes it from every route that targets it. GitHub Issues connections come from the GitHub App installation and cannot be declared here, but a route can still target one by the name export gives it.

Manages the `MessagingConnection` config kind through `/config/objects/MessagingConnection/{name}`. Applies are recorded with the provider's `manager`, so a GitOps `vulnetix config apply --prune` leaves this object alone.

Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (most of them, such as `ApiSpecSource`, `RiskStrategy` and `VexSource`) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object, or lets go of it for a kind the server never deletes (such as `Organization`). Some kinds wait on a person after they are applied (an invited `Member`, a Slack connection): `status.phase` says so, and Terraform does not wait.

## Example Usage

### Minimal

The arguments the kind requires, and nothing else.

```terraform
# Only the arguments the kind requires. From config-schemas examples/MessagingConnection/minimal.yaml.

resource "vulnetix_messaging_connection" "security_slack" {
  name = "security-slack"
  type = "slack"
}
```

### Advanced

Every argument the kind supports. A secret is a write-only `*_wo` argument fed from an ephemeral variable, so it never reaches the plan or the state; change `*_wo_version` to send a new value.

```terraform
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The object's name (`metadata.name`), unique per kind within the tenant. Changing it replaces the object.
- `type` (String) The provider. slack: a channel in a workspace that installs the Vulnetix Slack app. teams: a Microsoft Teams Workflows (Power Automate) incoming webhook. discord and google_chat: an incoming webhook. jira: issues in a Jira Cloud site, with an API token. webhook: an HTTPS endpoint of yours that receives each message signed with the Standard Webhooks scheme. It cannot change once the connection exists. One of `slack`, `teams`, `discord`, `google_chat`, `jira`, `webhook`.

### Optional

> **NOTE**: [Write-only arguments](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments) are supported in Terraform 1.11 and later.

- `annotations` (Map of String) Free key/value notes for people and tools (`metadata.annotations`), such as a ticket link.
- `display_name` (String) The label the console shows. Defaults to metadata.name.
- `jira` (Attributes) The Jira Cloud site a jira connection opens issues in, and the account it acts as. (see [below for nested schema](#nestedatt--jira))
- `labels` (Map of String) Free key/value labels (`metadata.labels`). Stored with the object, never interpreted by the server.
- `paused` (Boolean) Stop sending to this connection without removing it. A connection the server stopped after repeated failures, or whose vault entry was revoked, is reported in status.connectionStatus; supplying its credential again starts it. Schema default: `false`.
- `slack` (Attributes) The channel a slack connection posts to. The workspace is the one the app is installed into from status.authorizationUrl. (see [below for nested schema](#nestedatt--slack))
- `url_wo` (String, Sensitive, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) Where the client reads the URL messages are posted to, for teams, discord, google_chat and webhook. A Discord URL is on discord.com or discordapp.com, a Google Chat URL on chat.googleapis.com, and a Teams URL on logic.azure.com, environment.api.powerplatform.com or webhook.office.com; every URL is https. The URL is the credential, so it is kept in the vault and only a masked form is shown. Required when the connection is created; left out later, the stored URL is kept. Write-only: never stored in state. Needs Terraform 1.11 or later. Change `url_wo_version` to send a new value.
- `url_wo_version` (Number) Any number. Change it to send `url_wo` again; Terraform cannot see a write-only value change.

### Read-Only

- `id` (String) The object's name, the same as `name`.
- `status` (Attributes) What the server says about the object beyond its spec. (see [below for nested schema](#nestedatt--status))

<a id="nestedatt--jira"></a>
### Nested Schema for `jira`

Required:

- `email` (String) The Atlassian account the API token belongs to.
- `site_url` (String) The site, https://<site>.atlassian.net, with no path. It cannot change once the connection exists.

Optional:

- `api_token_wo` (String, Sensitive, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) Where the client reads the account's API token. Jira is asked who the token belongs to before the connection is saved. Required when the connection is created; left out later, the stored token is kept. Write-only: never stored in state. Needs Terraform 1.11 or later. Change `api_token_wo_version` to send a new value.
- `api_token_wo_version` (Number) Any number. Change it to send `api_token_wo` again; Terraform cannot see a write-only value change.
- `default_issue_type` (String) The issue type a route target uses when it names none, such as Bug or Task.
- `default_project_key` (String) The project a route target opens issues in when it names none. Stored upper case.


<a id="nestedatt--slack"></a>
### Nested Schema for `slack`

Optional:

- `channel_id` (String) The Slack channel id (C0123456789). Left out, messages need a channel on each route target.
- `channel_name` (String) The channel's name, for display only.


<a id="nestedatt--status"></a>
### Nested Schema for `status`

Read-Only:

- `authorization_url` (String) Where a person authorizes the object while the phase is `AwaitingAuthorization`.
- `extra` (Map of String) Other fields the kind reports, such as an invitation's expiry or a service account's principal. Structured values are JSON.
- `message` (String) Why the object is in its phase, when the server says.
- `object_uuid` (String) The object's id on the server.
- `once` (Map of String, Sensitive) Credentials the server shows only when it creates the object: a service account's `apiKey` and `secret`, a messaging connection's `signingSecret`, a log source's ingest token. They are kept in the state from the create, so keep the state encrypted. They cannot be read again: replace the object to get new ones.
- `phase` (String) `Applied` for most objects. `Pending` when the object waits on someone (an invited member who has not accepted), `AwaitingAuthorization` when a person has to authorize it in a browser at `authorization_url`. Terraform does not wait for either. `Unmanaged` means the object exists but no manager owns it.

## Import

Import an existing object by its name (`metadata.name`). The first read after an import fills in every field the server holds.

```shell
# Import an existing MessagingConnection by its name (metadata.name).
terraform import vulnetix_messaging_connection.security_slack security-slack
```
