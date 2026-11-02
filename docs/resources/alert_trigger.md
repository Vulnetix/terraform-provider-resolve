---
page_title: "vulnetix_alert_trigger Resource - vulnetix"
subcategory: ""
description: |-
  A Guardian trigger: run a SavedSearch or a DetectionRule on a schedule and, when its rows pass the levers, start a Software Factory workflow, post to messaging connections or send an email, with the rows as context. No rows, nothing starts. Name exactly one of searchRef or ruleRef. The plan caps how many triggers are enabled at once.
  An object of this kind can reference SavedSearch, DetectionRule, MessagingConnection objects. Reference those resources' attributes, or add depends_on, so Terraform applies them first.
  Manages the AlertTrigger config kind through /config/objects/AlertTrigger/{name}. Applies are recorded with the provider's manager, so a GitOps vulnetix config apply --prune leaves this object alone.
  Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (most of them, such as ApiSpecSource, RiskStrategy and VexSource) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object, or lets go of it for a kind the server never deletes (such as Organization). Some kinds wait on a person after they are applied (an invited Member, a Slack connection): status.phase says so, and Terraform does not wait.
---

# vulnetix_alert_trigger (Resource)

A Guardian trigger: run a SavedSearch or a DetectionRule on a schedule and, when its rows pass the levers, start a Software Factory workflow, post to messaging connections or send an email, with the rows as context. No rows, nothing starts. Name exactly one of searchRef or ruleRef. The plan caps how many triggers are enabled at once.

An object of this kind can reference `SavedSearch`, `DetectionRule`, `MessagingConnection` objects. Reference those resources' attributes, or add `depends_on`, so Terraform applies them first.

Manages the `AlertTrigger` config kind through `/config/objects/AlertTrigger/{name}`. Applies are recorded with the provider's `manager`, so a GitOps `vulnetix config apply --prune` leaves this object alone.

Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (most of them, such as `ApiSpecSource`, `RiskStrategy` and `VexSource`) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object, or lets go of it for a kind the server never deletes (such as `Organization`). Some kinds wait on a person after they are applied (an invited `Member`, a Slack connection): `status.phase` says so, and Terraform does not wait.

## Example Usage

### Minimal

The arguments the kind requires, and nothing else.

```terraform
# Only the arguments the kind requires. From config-schemas examples/AlertTrigger/minimal.yaml.

resource "vulnetix_alert_trigger" "dns_investigate" {
  name = "dns-investigate"
  destination = {
    kind = "workflow"
  }
  search_ref = "dynamic-dns-lookups"
}
```

### Advanced

Every argument the kind supports. A secret is a write-only `*_wo` argument fed from an ephemeral variable, so it never reaches the plan or the state; change `*_wo_version` to send a new value.

```terraform
# Every argument the kind supports. From config-schemas examples/AlertTrigger/advanced.yaml.

resource "vulnetix_alert_trigger" "dns_investigate" {
  name = "dns-investigate"
  destination = {
    kind     = "workflow"
    severity = "warning"
  }
  display_name = "Dynamic DNS lookups from one host"
  enabled      = true
  every_sec    = 300
  levers = {
    ask          = "always"
    columns      = ["ts", "src_ip", "qname"]
    cooldown_sec = 3600
    count_by     = "src_ip"
    hours        = "business"
    max_per_hour = 4
    min_rows     = 3
    rows_cap     = 25
    window_sec   = 900
    zone         = "Australia/Melbourne"
  }
  search_ref = "dynamic-dns-lookups"
  labels = {
    team = "secops"
  }
  annotations = {
    runbook = "https://wiki.example.com/secops/dynamic-dns"
  }
}

resource "vulnetix_alert_trigger" "root_login_page" {
  name = "root-login-page"
  destination = {
    connections = ["secops-slack"]
    kind        = "messaging"
    severity    = "critical"
  }
  display_name = "AWS root console login"
  every_sec    = 60
  levers = {
    cooldown_sec = 900
    hours        = "always"
    window_sec   = 300
  }
  rule_ref = "sigmahq-aws-root-login"
}

resource "vulnetix_alert_trigger" "failed_logins_digest" {
  name = "failed-logins-digest"
  destination = {
    emails   = ["secops@example.com", "oncall@example.com"]
    kind     = "email"
    severity = "notice"
  }
  display_name = "Failed SSO logins outside business hours"
  enabled      = false
  every_sec    = 3600
  levers = {
    count_by   = "user_name"
    hours      = "outside_business"
    min_rows   = 20
    rows_cap   = 50
    window_sec = 3600
    zone       = "Europe/London"
  }
  search_ref = "failed-sso-logins"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `destination` (Attributes) What a fire starts. workflow: a siem.fired source event for the Software Factory, which every published workflow whose trigger matches it takes. messaging: a post to each connection. email: a message to each address. (see [below for nested schema](#nestedatt--destination))
- `name` (String) The object's name (`metadata.name`), unique per kind within the tenant. Changing it replaces the object.

### Optional

- `annotations` (Map of String) Free key/value notes for people and tools (`metadata.annotations`), such as a ticket link.
- `display_name` (String) The name the console shows. Defaults to metadata.name.
- `enabled` (Boolean) Schema default: `true`.
- `every_sec` (Number) Seconds between runs. Schema default: `300`.
- `labels` (Map of String) Free key/value labels (`metadata.labels`). Stored with the object, never interpreted by the server.
- `levers` (Attributes) When a run fires, and what it sends. (see [below for nested schema](#nestedatt--levers))
- `rule_ref` (String) The metadata.name of the DetectionRule to run.
- `search_ref` (String) The metadata.name of the SavedSearch to run.

### Read-Only

- `id` (String) The object's name, the same as `name`.
- `status` (Attributes) What the server says about the object beyond its spec. (see [below for nested schema](#nestedatt--status))

<a id="nestedatt--destination"></a>
### Nested Schema for `destination`

Required:

- `kind` (String) One of `workflow`, `messaging`, `email`.

Optional:

- `connections` (Set of String)
- `emails` (Set of String)
- `severity` (String) One of `info`, `notice`, `warning`, `critical`. Schema default: `warning`.


<a id="nestedatt--levers"></a>
### Nested Schema for `levers`

Optional:

- `ask` (String) For workflow runs: workflow follows the workflow's asks; always asks a person before any change. One of `workflow`, `always`. Schema default: `workflow`.
- `columns` (Set of String) Columns sent. Left out, every column of the result.
- `cooldown_sec` (Number) Do not fire the same group again for this long, in seconds. Schema default: `3600`.
- `count_by` (String) Count rows per this column. Left out, the whole result is one group.
- `hours` (String) When it may fire. business is 09:00 to 17:00, Monday to Friday, in zone. One of `always`, `business`, `outside_business`. Schema default: `always`.
- `max_per_hour` (Number) At most this many fires per hour for the trigger. Schema default: `4`.
- `min_rows` (Number) Fire when a group has at least this many rows in the window. Schema default: `1`.
- `rows_cap` (Number) Rows sent as context, at most. Schema default: `25`.
- `window_sec` (Number) How far back each run looks, in seconds. Schema default: `300`.
- `zone` (String) The IANA time zone for hours. Schema default: `UTC`.


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
# Import an existing AlertTrigger by its name (metadata.name).
terraform import vulnetix_alert_trigger.dns_investigate dns-investigate
```
