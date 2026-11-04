---
page_title: "vulnetix_scheduled_report Resource - vulnetix"
subcategory: ""
description: |-
  A Guardian scheduled send: run a SavedSearch on a calendar and send its rows to messaging connections or email. Each send is charged to the daily query budget like any search, and the plan caps how many are enabled.
  An object of this kind can reference SavedSearch, MessagingConnection objects. Reference those resources' attributes, or add depends_on, so Terraform applies them first.
  Manages the ScheduledReport config kind through /config/objects/ScheduledReport/{name}. Applies are recorded with the provider's manager, so a GitOps vulnetix config apply --prune leaves this object alone.
  Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (most of them, such as ApiSpecSource, RiskStrategy and VexSource) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object, or lets go of it for a kind the server never deletes (such as Organization). Some kinds wait on a person after they are applied (an invited Member, a Slack connection): status.phase says so, and Terraform does not wait.
---

# vulnetix_scheduled_report (Resource)

A Guardian scheduled send: run a SavedSearch on a calendar and send its rows to messaging connections or email. Each send is charged to the daily query budget like any search, and the plan caps how many are enabled.

An object of this kind can reference `SavedSearch`, `MessagingConnection` objects. Reference those resources' attributes, or add `depends_on`, so Terraform applies them first.

Manages the `ScheduledReport` config kind through `/config/objects/ScheduledReport/{name}`. Applies are recorded with the provider's `manager`, so a GitOps `vulnetix config apply --prune` leaves this object alone.

Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (most of them, such as `ApiSpecSource`, `RiskStrategy` and `VexSource`) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object, or lets go of it for a kind the server never deletes (such as `Organization`). Some kinds wait on a person after they are applied (an invited `Member`, a Slack connection): `status.phase` says so, and Terraform does not wait.

## Example Usage

### Minimal

The arguments the kind requires, and nothing else.

```terraform
# Only the arguments the kind requires. From config-schemas examples/ScheduledReport/minimal.yaml.

resource "vulnetix_scheduled_report" "daily_dns_digest" {
  name = "daily-dns-digest"
  destination = {
    connections = ["secops-slack"]
    kind        = "messaging"
  }
  search_ref = "dynamic-dns-lookups"
}
```

### Advanced

Every argument the kind supports. A secret is a write-only `*_wo` argument fed from an ephemeral variable, so it never reaches the plan or the state; change `*_wo_version` to send a new value.

```terraform
# Every argument the kind supports. From config-schemas examples/ScheduledReport/advanced.yaml.

resource "vulnetix_scheduled_report" "weekly_console_sign_ins" {
  name = "weekly-console-sign-ins"
  at   = "08:00"
  destination = {
    emails  = ["secops@acme.example", "ciso@acme.example"]
    kind    = "email"
    subject = "Console sign-ins this week"
  }
  display_name = "Weekly console sign-ins"
  enabled      = true
  format = {
    csv    = true
    layout = "table_chart"
    link   = true
    query  = true
    rows   = 20
  }
  repeat         = "weekly"
  rows_window_ms = 604800000
  search_ref     = "console-sign-ins"
  weekday        = 1
  when_empty     = "all_clear"
  zone           = "Australia/Melbourne"
  labels = {
    team = "secops"
  }
  annotations = {
    ticket = "https://jira.example.com/browse/SEC-2216"
  }
}

resource "vulnetix_scheduled_report" "monthly_dns_digest" {
  name = "monthly-dns-digest"
  at   = "09:30"
  destination = {
    connections = ["secops-slack", "soc-teams"]
    kind        = "messaging"
    subject     = "Dynamic DNS lookups this month"
  }
  enabled = true
  format = {
    layout = "summary"
    link   = true
  }
  monthday       = 1
  repeat         = "monthly"
  search_ref     = "dynamic-dns-lookups"
  since_previous = true
  when_empty     = "skip"
  zone           = "UTC"
}
```

### Reference

Every argument that has a schema default, set to that default, plus the required arguments. Planned against a tenant it only sets values the tenant left unset. Every kind's reference file is also in [Vulnetix/platform-iac-examples](https://github.com/Vulnetix/platform-iac-examples), with the same document as YAML.

```terraform
# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/ScheduledReport/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_scheduled_report" "daily_dns_digest" {
  name = "daily-dns-digest"
  at   = "08:00"
  destination = {
    connections = ["secops-slack"]
    kind        = "messaging"
  }
  enabled = true
  format = {
    csv    = true
    layout = "table_chart"
    link   = true
    query  = false
    rows   = 20
  }
  monthday       = 1
  repeat         = "daily"
  rows_window_ms = 86400000
  search_ref     = "dynamic-dns-lookups"
  since_previous = false
  weekday        = 1
  when_empty     = "skip"
  zone           = "UTC"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `destination` (Attributes) (see [below for nested schema](#nestedatt--destination))
- `name` (String) The object's name (`metadata.name`), unique per kind within the tenant. Changing it replaces the object.
- `search_ref` (String) The metadata.name of the SavedSearch to run.

### Optional

- `annotations` (Map of String) Free key/value notes for people and tools (`metadata.annotations`), such as a ticket link.
- `at` (String) The time of day, 24-hour, in zone. Schema default: `08:00`.
- `display_name` (String) The name the console shows, and the default subject. Defaults to metadata.name.
- `enabled` (Boolean) Schema default: `true`.
- `format` (Attributes) (see [below for nested schema](#nestedatt--format))
- `labels` (Map of String) Free key/value labels (`metadata.labels`). Stored with the object, never interpreted by the server.
- `monthday` (Number) For monthly. At most 28, so every month has the day. Schema default: `1`.
- `repeat` (String) One of `hourly`, `every_6h`, `daily`, `weekly`, `monthly`. Schema default: `daily`.
- `rows_window_ms` (Number) How far back each send reads, in milliseconds: 1 hour to 31 days. Schema default: `86400000`.
- `since_previous` (Boolean) Read from the previous send instead of rowsWindowMs, so no row is reported twice. Schema default: `false`.
- `weekday` (Number) For weekly: 0 is Sunday. Schema default: `1`.
- `when_empty` (String) skip sends nothing when there are no rows; all_clear sends a message saying so. One of `skip`, `all_clear`. Schema default: `skip`.
- `zone` (String) An IANA time zone, such as Australia/Sydney. Schema default: `UTC`.

### Read-Only

- `id` (String) The object's name, the same as `name`.
- `status` (Attributes) What the server says about the object beyond its spec. (see [below for nested schema](#nestedatt--status))

<a id="nestedatt--destination"></a>
### Nested Schema for `destination`

Required:

- `kind` (String) One of `messaging`, `email`.

Optional:

- `connections` (Set of String)
- `emails` (Set of String)
- `subject` (String) Defaults to displayName.


<a id="nestedatt--format"></a>
### Nested Schema for `format`

Optional:

- `csv` (Boolean) Attach a CSV of up to 1,000 rows. Email only. Schema default: `true`.
- `layout` (String) table_chart: a bar chart of the first measure and the table. table: the table only. summary: the row count only. One of `table_chart`, `table`, `summary`. Schema default: `table_chart`.
- `link` (Boolean) Link to the search in the console. Schema default: `true`.
- `query` (Boolean) Include the query text. Schema default: `false`.
- `rows` (Number) Rows shown. Schema default: `20`.


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
# Import an existing ScheduledReport by its name (metadata.name).
terraform import vulnetix_scheduled_report.daily_dns_digest daily-dns-digest
```
