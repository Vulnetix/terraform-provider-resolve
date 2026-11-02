---
page_title: "vulnetix_detection_rule Resource - vulnetix"
subcategory: ""
description: |-
  A Guardian alert rule: a detection run on a schedule over the tenant's logs, raising alerts. Name exactly one of body, file, searchRef, countermeasure, sigmahq or ioc. body and file are the tenant's own rule; the others run a rule Guardian already has, which can be tuned (level, threshold, suppressions, timing) but not edited. A rule whose detection does not compile, or names a log source Guardian has no table for, is stored turned off. The plan caps how many rules run at once. Deleting one that an AlertTrigger still uses is refused.
  An object of this kind can reference SavedSearch, MessagingConnection objects. Reference those resources' attributes, or add depends_on, so Terraform applies them first.
  Manages the DetectionRule config kind through /config/objects/DetectionRule/{name}. Applies are recorded with the provider's manager, so a GitOps vulnetix config apply --prune leaves this object alone.
  Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (most of them, such as ApiSpecSource, RiskStrategy and VexSource) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object, or lets go of it for a kind the server never deletes (such as Organization). Some kinds wait on a person after they are applied (an invited Member, a Slack connection): status.phase says so, and Terraform does not wait.
---

# vulnetix_detection_rule (Resource)

A Guardian alert rule: a detection run on a schedule over the tenant's logs, raising alerts. Name exactly one of body, file, searchRef, countermeasure, sigmahq or ioc. body and file are the tenant's own rule; the others run a rule Guardian already has, which can be tuned (level, threshold, suppressions, timing) but not edited. A rule whose detection does not compile, or names a log source Guardian has no table for, is stored turned off. The plan caps how many rules run at once. Deleting one that an AlertTrigger still uses is refused.

An object of this kind can reference `SavedSearch`, `MessagingConnection` objects. Reference those resources' attributes, or add `depends_on`, so Terraform applies them first.

Manages the `DetectionRule` config kind through `/config/objects/DetectionRule/{name}`. Applies are recorded with the provider's `manager`, so a GitOps `vulnetix config apply --prune` leaves this object alone.

Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (most of them, such as `ApiSpecSource`, `RiskStrategy` and `VexSource`) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object, or lets go of it for a kind the server never deletes (such as `Organization`). Some kinds wait on a person after they are applied (an invited `Member`, a Slack connection): `status.phase` says so, and Terraform does not wait.

## Example Usage

### Minimal

The arguments the kind requires, and nothing else.

```terraform
# Only the arguments the kind requires. From config-schemas examples/DetectionRule/minimal.yaml.

resource "vulnetix_detection_rule" "sigmahq_aws_root_login" {
  name    = "sigmahq-aws-root-login"
  sigmahq = "rules/cloud/aws/cloudtrail/aws_root_account_usage.yml"
}
```

### Advanced

Every argument the kind supports. A secret is a write-only `*_wo` argument fed from an ephemeral variable, so it never reaches the plan or the state; change `*_wo_version` to send a new value.

```terraform
# Every argument the kind supports. From config-schemas examples/DetectionRule/advanced.yaml.

resource "vulnetix_detection_rule" "dynamic_dns" {
  name               = "dynamic-dns"
  body               = "title: Dynamic DNS lookup\nid: 6f3e0c52-0d4b-4c8e-9a51-2a6b9a0e7d10\nlevel: medium\nlogsource:\n  category: dns_query\ndetection:\n  selection:\n    QueryName|endswith:\n      - .duckdns.org\n      - .no-ip.org\n  condition: selection\n"
  display_name       = "Dynamic DNS lookup"
  enabled            = true
  interval_sec       = 120
  language           = "sigma"
  lateness_sec       = 300
  level              = "medium"
  max_alerts_per_run = 50
  severity_override  = "warning"
  suppress_sec       = 1800
  suppressions = [
    {
      field = "QueryName"
      until = "2026-12-31"
      value = "updates.duckdns.org"
    },
    {
      field = "src_ip"
      value = "10.20.0.15"
    },
  ]
  target = {
    connections = ["secops-slack"]
    emails      = ["secops@example.com"]
  }
  threshold = {
    by         = ["src_ip"]
    count      = 5
    window_sec = 600
  }
  labels = {
    mitre-tactic = "command-and-control"
    team         = "secops"
  }
  annotations = {
    runbook = "https://wiki.example.com/secops/dynamic-dns"
  }
}

resource "vulnetix_detection_rule" "brute_force_file" {
  name         = "brute-force-file"
  display_name = "SSH brute force"
  enabled      = false
  file         = "siem/rules/brute-force.yml"
  language     = "sigma"
  level        = "high"
}

resource "vulnetix_detection_rule" "dns_search_rule" {
  name         = "dns-search-rule"
  interval_sec = 600
  level        = "low"
  search_ref   = "dynamic-dns-lookups"
}

resource "vulnetix_detection_rule" "log4shell_countermeasure" {
  name              = "log4shell-countermeasure"
  countermeasure    = "3b1f8a62-7d2e-4c41-9a0b-5e6f7c8d9e01"
  interval_sec      = 300
  severity_override = "critical"
  target = {
    connections = ["secops-slack"]
  }
}

resource "vulnetix_detection_rule" "sigmahq_aws_root_login" {
  name         = "sigmahq-aws-root-login"
  level        = "high"
  sigmahq      = "rules/cloud/aws/cloudtrail/aws_root_account_usage.yml"
  suppress_sec = 3600
}

resource "vulnetix_detection_rule" "indicator_domains" {
  name               = "indicator-domains"
  display_name       = "Indicator domain resolved"
  ioc                = "ioc-dns"
  lateness_sec       = 120
  max_alerts_per_run = 100
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The object's name (`metadata.name`), unique per kind within the tenant. Changing it replaces the object.

### Optional

- `annotations` (Map of String) Free key/value notes for people and tools (`metadata.annotations`), such as a ticket link.
- `body` (String) The rule text, at most 256 KiB.
- `countermeasure` (String) A Vulnetix countermeasure of kind SIGMA, by uuid: the virtual patching lane's detection for a vulnerability.
- `display_name` (String) The name the console shows. Defaults to the Sigma rule's title, or the saved search's, countermeasure's or built-in rule's name.
- `enabled` (Boolean) Whether the rule runs. A rule that cannot compile is never on, whatever this says. Schema default: `true`.
- `file` (String) A file holding the rule text, the same as body.
- `interval_sec` (Number) Seconds between runs. Schema default: `60`.
- `ioc` (String) A built-in indicator rule. ioc-dns: a domain on an indicator list resolved. ioc-net: a connection to an indicator address. One of `ioc-dns`, `ioc-net`.
- `labels` (Map of String) Free key/value labels (`metadata.labels`). Stored with the object, never interpreted by the server.
- `language` (String) The language of body or file. sigma: a Sigma rule, or several in one multi-document YAML with one correlation rule, which is the one that alerts. dsl: the Guardian search language. Guardian keeps a DSL rule only as a copy of a saved search today, so a dsl body or file has no direct write path yet; use searchRef until Guardian takes DSL rule text. One of `sigma`, `dsl`.
- `lateness_sec` (Number) How long a run waits for late rows before it reads a window, in seconds. Schema default: `300`.
- `level` (String) The Sigma level. Defaults to the rule's own level, or medium. One of `informational`, `low`, `medium`, `high`, `critical`.
- `max_alerts_per_run` (Number) Alerts one run may raise. Beyond it, one overflow alert says how many more matched. Schema default: `50`.
- `search_ref` (String) The metadata.name of a SavedSearch whose text the rule runs. Its text is copied when the rule is created: editing the search later does not change the rule.
- `severity_override` (String) The alert severity, instead of the one the level maps to (informational to info, low to notice, medium to warning, high and critical to critical). One of `info`, `notice`, `warning`, `critical`.
- `sigmahq` (String) A rule from the SigmaHQ catalogue Guardian imports weekly, by its path in the SigmaHQ repository. Only rules that compile are offered.
- `suppress_sec` (Number) Seconds during which a repeat of the same alert is merged into the open one instead of raising another. Schema default: `900`.
- `suppressions` (Attributes Set) Values that never alert. Each adds a filter on field equal to value to the rule's detection until it lapses. (see [below for nested schema](#nestedatt--suppressions))
- `target` (Attributes) Where the rule's alerts are sent, beyond the Guardian alerts list. (see [below for nested schema](#nestedatt--target))
- `threshold` (Attributes) Alert only when at least count rows match within windowSec, counted per the by fields. Wraps the detection as a Sigma event_count correlation. (see [below for nested schema](#nestedatt--threshold))

### Read-Only

- `id` (String) The object's name, the same as `name`.
- `status` (Attributes) What the server says about the object beyond its spec. (see [below for nested schema](#nestedatt--status))

<a id="nestedatt--suppressions"></a>
### Nested Schema for `suppressions`

Required:

- `field` (String)
- `value` (String)

Optional:

- `until` (String) The last day it applies is the day before this date (UTC). Left out, it never lapses.


<a id="nestedatt--target"></a>
### Nested Schema for `target`

Optional:

- `connections` (Set of String)
- `emails` (Set of String)


<a id="nestedatt--threshold"></a>
### Nested Schema for `threshold`

Required:

- `count` (Number)
- `window_sec` (Number)

Optional:

- `by` (Set of String)


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
# Import an existing DetectionRule by its name (metadata.name).
terraform import vulnetix_detection_rule.sigmahq_aws_root_login sigmahq-aws-root-login
```
