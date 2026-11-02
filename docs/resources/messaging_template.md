---
page_title: "vulnetix_messaging_template Resource - vulnetix"
subcategory: ""
description: |-
  An organisation's own message template, which a MessagingRoute uses instead of the system template for the event type. One template renders for every provider: Slack, Teams, Discord and Google Chat get the title, summary, body, fields, actions and accent; Jira and GitHub Issues get the issue fields; a signed webhook gets the rendered text and, when asked, the event payload. Every text field is a Go text/template run in a sandbox: only the messaging functions are allowed, every value is escaped, and sizes are capped. Plan compiles each field and refuses a template with errors. Every change to content is saved as a new version, noted as applied as code, as a console edit is; earlier versions stay in the template's history. Deleting it is refused while a route outside this apply still uses it.
  Manages the MessagingTemplate config kind through /config/objects/MessagingTemplate/{name}. Applies are recorded with the provider's manager, so a GitOps vulnetix config apply --prune leaves this object alone.
  Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (most of them, such as ApiSpecSource, RiskStrategy and VexSource) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object, or lets go of it for a kind the server never deletes (such as Organization). Some kinds wait on a person after they are applied (an invited Member, a Slack connection): status.phase says so, and Terraform does not wait.
---

# vulnetix_messaging_template (Resource)

An organisation's own message template, which a MessagingRoute uses instead of the system template for the event type. One template renders for every provider: Slack, Teams, Discord and Google Chat get the title, summary, body, fields, actions and accent; Jira and GitHub Issues get the issue fields; a signed webhook gets the rendered text and, when asked, the event payload. Every text field is a Go text/template run in a sandbox: only the messaging functions are allowed, every value is escaped, and sizes are capped. Plan compiles each field and refuses a template with errors. Every change to content is saved as a new version, noted as applied as code, as a console edit is; earlier versions stay in the template's history. Deleting it is refused while a route outside this apply still uses it.

Manages the `MessagingTemplate` config kind through `/config/objects/MessagingTemplate/{name}`. Applies are recorded with the provider's `manager`, so a GitOps `vulnetix config apply --prune` leaves this object alone.

Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (most of them, such as `ApiSpecSource`, `RiskStrategy` and `VexSource`) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object, or lets go of it for a kind the server never deletes (such as `Organization`). Some kinds wait on a person after they are applied (an invited `Member`, a Slack connection): `status.phase` says so, and Terraform does not wait.

## Example Usage

### Minimal

The arguments the kind requires, and nothing else.

```terraform
# Only the arguments the kind requires. From config-schemas examples/MessagingTemplate/minimal.yaml.

resource "vulnetix_messaging_template" "plain_title" {
  name = "plain-title"
  content = {
    title = "{{.event.title}}"
  }
}
```

### Advanced

Every argument the kind supports. A secret is a write-only `*_wo` argument fed from an ephemeral variable, so it never reaches the plan or the state; change `*_wo_version` to send a new value.

```terraform
# Every argument the kind supports. From config-schemas examples/MessagingTemplate/advanced.yaml.

resource "vulnetix_messaging_template" "security_findings" {
  name                 = "security-findings"
  base_system_template = "generic"
  content = {
    accent = "{{.event.severity}}"
    actions = [
      {
        label = "Open in Vulnetix"
        url   = "{{.event.url}}"
      },
    ]
    body = "{{range $k, $v := .data}}- **{{$k}}:** {{truncate 200 $v}}\n{{end}}{{with .mentions}}\n{{mentions .}}{{end}}"
    digest = {
      item  = "{{.event.title}}{{with .subject.label}}: {{.}}{{end}}"
      title = "{{.digest.count}} {{plural .digest.count \"event\" \"events\"}}"
    }
    fields = [
      {
        inline = true
        label  = "Severity"
        value  = "{{severityEmoji .event.severity}} {{severityLabel .event.severity}}"
      },
    ]
    issue = {
      body   = "{{.event.title}}"
      labels = ["vulnetix", "security"]
      title  = "[Vulnetix] {{.event.title}}"
    }
    summary = "{{with .subject.label}}{{.}}{{end}}"
    title   = "{{.event.title}}"
    webhook = {
      include_payload = true
    }
  }
  description  = "Scan results for the security channel, with a link back to the scan."
  display_name = "Security findings"
  event_scope  = ["scan.completed", "scan.build_broken"]
  labels = {
    team = "security"
  }
  annotations = {
    owner = "appsec@acme.example"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `content` (Attributes) The template itself. Each text field is a template over the event view (.event, .subject, .actor, .data, .mentions, and .digest for a digest). (see [below for nested schema](#nestedatt--content))
- `name` (String) The object's name (`metadata.name`), unique per kind within the tenant. Changing it replaces the object.

### Optional

- `annotations` (Map of String) Free key/value notes for people and tools (`metadata.annotations`), such as a ticket link.
- `base_system_template` (String) The key of the system template this one was copied from (generic, findings.report, scan.build_broken and so on). It is recorded so the console can say when that system template has a newer version; it does not change the content.
- `description` (String) What the template is for.
- `display_name` (String) The name the console shows, unique among the organisation's templates. Defaults to metadata.name.
- `event_scope` (Set of String) The event types the template is meant for: a type (scan.completed), a family (pix.*), or *. The console uses it to offer the template; a route may still use it for any event.
- `labels` (Map of String) Free key/value labels (`metadata.labels`). Stored with the object, never interpreted by the server.

### Read-Only

- `id` (String) The object's name, the same as `name`.
- `status` (Attributes) What the server says about the object beyond its spec. (see [below for nested schema](#nestedatt--status))

<a id="nestedatt--content"></a>
### Nested Schema for `content`

Required:

- `title` (String) The message title. Required.

Optional:

- `accent` (String) The colour bar: a severity (info, notice, warning, critical) or a #rrggbb colour, usually {{.event.severity}}.
- `actions` (Attributes List) Buttons or links. At most 10. A URL that does not render to an https URL is dropped from the message. (see [below for nested schema](#nestedatt--content--actions))
- `body` (String) The message body, in markdown-lite.
- `digest` (Attributes) How a batch of events is summed up when a route sends a digest. (see [below for nested schema](#nestedatt--content--digest))
- `fields` (Attributes List) Name and value pairs shown beside the body. At most 25. (see [below for nested schema](#nestedatt--content--fields))
- `issue` (Attributes) What an issue tracker target gets. A field left out falls back to the message's own. (see [below for nested schema](#nestedatt--content--issue))
- `summary` (String) One line under the title.
- `webhook` (Attributes) What a signed webhook target receives. (see [below for nested schema](#nestedatt--content--webhook))

<a id="nestedatt--content--actions"></a>
### Nested Schema for `content.actions`

Required:

- `label` (String) The button text.
- `url` (String) Where it goes.


<a id="nestedatt--content--digest"></a>
### Nested Schema for `content.digest`

Optional:

- `item` (String) One line per event in the batch.
- `title` (String) The digest title, over the batch view (.digest.count).


<a id="nestedatt--content--fields"></a>
### Nested Schema for `content.fields`

Required:

- `label` (String) The field's name.
- `value` (String) The field's value.

Optional:

- `inline` (Boolean) Show the field beside the previous one where the provider can. Schema default: `false`.


<a id="nestedatt--content--issue"></a>
### Nested Schema for `content.issue`

Optional:

- `body` (String) The issue body.
- `labels` (Set of String) Labels to put on the issue. At most 10.
- `title` (String) The issue title.


<a id="nestedatt--content--webhook"></a>
### Nested Schema for `content.webhook`

Optional:

- `include_payload` (Boolean) Send the event's payload, reduced to the fields the catalogue allows, beside the rendered text. Schema default: `false`.



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
# Import an existing MessagingTemplate by its name (metadata.name).
terraform import vulnetix_messaging_template.plain_title plain-title
```
