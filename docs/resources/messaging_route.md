---
page_title: "vulnetix_messaging_route Resource - vulnetix"
subcategory: ""
description: |-
  A messaging workflow: which events go to which messaging connections, filtered how, rendered with which template, and how often. Every event the organisation's producers emit is matched against its routes; a route that takes it renders its template once per target and records each delivery in the console's delivery log. A route whose event types reach a sensitive type (security., vault., or a sensitive Pix event) is marked sensitive: it may then send to a signed webhook only with allowSensitiveWebhook, and to a GitHub repository only when the repository is private. Each change is recorded as a messaging.route.changed event, as a console edit is. Deleting it stops its deliveries; the delivery log is kept.
  An object of this kind can reference MessagingConnection, MessagingTemplate objects. Reference those resources' attributes, or add depends_on, so Terraform applies them first.
  Manages the MessagingRoute config kind through /config/objects/MessagingRoute/{name}. Applies are recorded with the provider's manager, so a GitOps vulnetix config apply --prune leaves this object alone.
  Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (most of them, such as ApiSpecSource, RiskStrategy and VexSource) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object, or lets go of it for a kind the server never deletes (such as Organization). Some kinds wait on a person after they are applied (an invited Member, a Slack connection): status.phase says so, and Terraform does not wait.
---

# vulnetix_messaging_route (Resource)

A messaging workflow: which events go to which messaging connections, filtered how, rendered with which template, and how often. Every event the organisation's producers emit is matched against its routes; a route that takes it renders its template once per target and records each delivery in the console's delivery log. A route whose event types reach a sensitive type (security.*, vault.*, or a sensitive Pix event) is marked sensitive: it may then send to a signed webhook only with allowSensitiveWebhook, and to a GitHub repository only when the repository is private. Each change is recorded as a messaging.route.changed event, as a console edit is. Deleting it stops its deliveries; the delivery log is kept.

An object of this kind can reference `MessagingConnection`, `MessagingTemplate` objects. Reference those resources' attributes, or add `depends_on`, so Terraform applies them first.

Manages the `MessagingRoute` config kind through `/config/objects/MessagingRoute/{name}`. Applies are recorded with the provider's `manager`, so a GitOps `vulnetix config apply --prune` leaves this object alone.

Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (most of them, such as `ApiSpecSource`, `RiskStrategy` and `VexSource`) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object, or lets go of it for a kind the server never deletes (such as `Organization`). Some kinds wait on a person after they are applied (an invited `Member`, a Slack connection): `status.phase` says so, and Terraform does not wait.

## Example Usage

### Minimal

The arguments the kind requires, and nothing else.

```terraform
# Only the arguments the kind requires. From config-schemas examples/MessagingRoute/minimal.yaml.

resource "vulnetix_messaging_route" "broken_builds" {
  name        = "broken-builds"
  event_types = ["scan.build_broken"]
}
```

### Advanced

Every argument the kind supports. A secret is a write-only `*_wo` argument fed from an ephemeral variable, so it never reaches the plan or the state; change `*_wo_version` to send a new value.

```terraform
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `event_types` (Set of String) The event types the route takes: a type (scan.completed), a family (pix.*), or * for every event. A type that was renamed is stored as its new name.
- `name` (String) The object's name (`metadata.name`), unique per kind within the tenant. Changing it replaces the object.

### Optional

- `actor_filter` (Attributes) Narrows the route by who acted. (see [below for nested schema](#nestedatt--actor_filter))
- `allow_sensitive_webhook` (Boolean) Let a sensitive route send to a signed webhook. Schema default: `false`.
- `annotations` (Map of String) Free key/value notes for people and tools (`metadata.annotations`), such as a ticket link.
- `digest` (Attributes) How events are batched. immediate sends each event as it comes; window sends one digest every 15 or 60 minutes; daily sends one digest a day at hour in tz. tiers overrides that per delivery tier. (see [below for nested schema](#nestedatt--digest))
- `display_name` (String) The name the console shows. Defaults to metadata.name.
- `enabled` (Boolean) A disabled route matches nothing. Schema default: `true`.
- `filters` (Attributes) Narrows the route to one scope, and for scan events to the findings that pass the thresholds. An unset condition lets everything through. (see [below for nested schema](#nestedatt--filters))
- `include_pii` (Boolean) Let personal details (a member's email address) into messages. Off, they are masked. Schema default: `false`.
- `labels` (Map of String) Free key/value labels (`metadata.labels`). Stored with the object, never interpreted by the server.
- `max_per_hour` (Number) At most this many messages an hour per target; the rest are suppressed and logged. Schema default: `60`.
- `mention_principals` (Set of String) Members of the organisation, by principal uuid, to @-mention on every target that supports mentions. Each is mentioned by the identity the messaging centre holds for them on that provider. At most 25.
- `mention_static` (Attributes) Provider user ids to @-mention that are not members, keyed by provider: a Slack user id (U0123ABCD), a Discord user id, a Google Chat user id, a Teams user id or UPN, a Jira account id, or a GitHub login. (see [below for nested schema](#nestedatt--mention_static))
- `min_severity` (String) The quietest event severity the route takes. One of `info`, `notice`, `warning`, `critical`. Schema default: `info`.
- `mode` (String) live sends. shadow renders and logs every delivery without sending, to try a route out. One of `live`, `shadow`. Schema default: `live`.
- `system_template` (String) The key of a system template to render every event with (generic, notice, findings.report and so on), instead of the one each event type has. Not with template.
- `targets` (Attributes Set) Where the route sends: one entry per connection, each connection at most once. At most 20. (see [below for nested schema](#nestedatt--targets))
- `template` (String) The metadata.name of a MessagingTemplate to render with. Leave it and systemTemplate out to use the system template for each event type.
- `template_version` (Number) Pin the template to one of its versions. Left out, the route follows the template's latest version. Only with template.

### Read-Only

- `id` (String) The object's name, the same as `name`.
- `status` (Attributes) What the server says about the object beyond its spec. (see [below for nested schema](#nestedatt--status))

<a id="nestedatt--actor_filter"></a>
### Nested Schema for `actor_filter`

Optional:

- `kinds` (Set of String) Actor kinds, such as principal, service_account, agent or system. At most 10.
- `mode` (String) any: everyone. include: only the principals and actor kinds listed. exclude: everyone but them. One of `any`, `include`, `exclude`. Schema default: `any`.
- `principals` (Set of String) Principal uuids. At most 100.


<a id="nestedatt--digest"></a>
### Nested Schema for `digest`

Optional:

- `hour` (Number) The hour of a daily digest, 0 to 23.
- `minutes` (Number) The window, for mode window, or for an immediate route with a tier set to digest (60 when left out). One of `15`, `60`.
- `mode` (String) One of `immediate`, `window`, `daily`. Schema default: `immediate`.
- `tiers` (Attributes) Per delivery tier: instant, digest or off. A tier left out follows mode. (see [below for nested schema](#nestedatt--digest--tiers))
- `tz` (String) The IANA time zone of a daily digest, such as Australia/Sydney.

<a id="nestedatt--digest--tiers"></a>
### Nested Schema for `digest.tiers`

Optional:

- `error` (String) One of `instant`, `digest`, `off`.
- `info` (String) One of `instant`, `digest`, `off`.
- `needs_you` (String) One of `instant`, `digest`, `off`.
- `warning` (String) One of `instant`, `digest`, `off`.



<a id="nestedatt--filters"></a>
### Nested Schema for `filters`

Optional:

- `findings` (Attributes) Thresholds a scan event's findings must pass. The event is sent with the findings that pass, and not at all when none do. (see [below for nested schema](#nestedatt--filters--findings))
- `scope` (Attributes) Only events about this application, version, upload or repository. Each field that is set must match the event. (see [below for nested schema](#nestedatt--filters--scope))

<a id="nestedatt--filters--findings"></a>
### Nested Schema for `filters.findings`

Optional:

- `cess_threshold` (Number) The lowest CESS score.
- `cvss_severities` (Set of String) Only these CVSS severities.
- `cvss_threshold` (Number) The lowest CVSS base score.
- `eol_days` (Number) Only packages that reach end of life within this many days.
- `epss_threshold` (Number) The lowest EPSS probability.
- `exploit_threshold` (Number) The fewest known exploits.
- `has_exploits` (Set of String) Only findings whose exploit maturity is one of these.
- `ignore_test` (Boolean) Leave out test dependencies. Schema default: `false`.
- `is_eol` (Boolean) Only end-of-life packages (true) or only supported ones (false).
- `is_exploitable` (Boolean) Only exploitable (true) or only not exploitable (false).
- `is_malware` (Boolean) Only malware (true) or only not malware (false).
- `is_reachable` (Boolean) Only reachable (true) or only unreachable (false) code.
- `new_only` (Boolean) Only findings new in this scan. Schema default: `false`.
- `only_direct` (Boolean) Only direct dependencies. Schema default: `false`.
- `only_production` (Boolean) Only production dependencies. Schema default: `false`.


<a id="nestedatt--filters--scope"></a>
### Nested Schema for `filters.scope`

Optional:

- `app_name` (String) The application name.
- `app_version` (String) The application version.
- `repo` (String) One repository, owner/name.
- `upload_uuid` (String) One upload.



<a id="nestedatt--mention_static"></a>
### Nested Schema for `mention_static`

Optional:

- `discord` (Set of String)
- `github` (Set of String)
- `google_chat` (Set of String)
- `jira` (Set of String)
- `slack` (Set of String)
- `teams` (Set of String)


<a id="nestedatt--targets"></a>
### Nested Schema for `targets`

Required:

- `connection` (String) The metadata.name of the MessagingConnection to send to. A GitHub Issues connection, which the GitHub App installation makes, is named as export names it.

Optional:

- `enabled` (Boolean) A disabled target is skipped. Schema default: `true`.
- `options` (Attributes) Per provider. slack: channelId overrides the connection's channel. jira: projectKey and issueType override the connection's defaults, labels are added. GitHub Issues: owner and repo (required), labels, assignMentions. Other providers take none. (see [below for nested schema](#nestedatt--targets--options))

<a id="nestedatt--targets--options"></a>
### Nested Schema for `targets.options`

Optional:

- `assign_mentions` (Boolean) Assign the GitHub issue to the mentioned people.
- `channel_id` (String) A Slack channel id.
- `issue_type` (String) A Jira issue type.
- `labels` (Set of String) Labels to add to the issue. At most 10.
- `owner` (String) The GitHub repository owner.
- `project_key` (String) A Jira project key.
- `repo` (String) The GitHub repository name.



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
# Import an existing MessagingRoute by its name (metadata.name).
terraform import vulnetix_messaging_route.broken_builds broken-builds
```
