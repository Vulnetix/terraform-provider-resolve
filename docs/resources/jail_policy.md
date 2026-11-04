---
page_title: "vulnetix_jail_policy Resource - vulnetix"
subcategory: ""
description: |-
  A jail policy: the ordered rules vulnetix jail evaluates for a repository over everything ever reported for it (vulnerabilities past their SLA, dependencies past end of life, migrations that have not landed, hygiene categories that regressed). A policy without repository is the tenant's default; a policy with one is a per-repository override that replaces the default outright. Separate from QualityGate, which judges a single scan.
  Manages the JailPolicy config kind through /config/objects/JailPolicy/{name}. Applies are recorded with the provider's manager, so a GitOps vulnetix config apply --prune leaves this object alone.
  Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (most of them, such as ApiSpecSource, RiskStrategy and VexSource) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object, or lets go of it for a kind the server never deletes (such as Organization). Some kinds wait on a person after they are applied (an invited Member, a Slack connection): status.phase says so, and Terraform does not wait.
---

# vulnetix_jail_policy (Resource)

A jail policy: the ordered rules vulnetix jail evaluates for a repository over everything ever reported for it (vulnerabilities past their SLA, dependencies past end of life, migrations that have not landed, hygiene categories that regressed). A policy without repository is the tenant's default; a policy with one is a per-repository override that replaces the default outright. Separate from QualityGate, which judges a single scan.

Manages the `JailPolicy` config kind through `/config/objects/JailPolicy/{name}`. Applies are recorded with the provider's `manager`, so a GitOps `vulnetix config apply --prune` leaves this object alone.

Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (most of them, such as `ApiSpecSource`, `RiskStrategy` and `VexSource`) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object, or lets go of it for a kind the server never deletes (such as `Organization`). Some kinds wait on a person after they are applied (an invited `Member`, a Slack connection): `status.phase` says so, and Terraform does not wait.

## Example Usage

### Minimal

The arguments the kind requires, and nothing else.

```terraform
# Only the arguments the kind requires. From config-schemas examples/JailPolicy/minimal.yaml.

resource "vulnetix_jail_policy" "org_default" {
  name = "org-default"
}
```

### Advanced

Every argument the kind supports. A secret is a write-only `*_wo` argument fed from an ephemeral variable, so it never reaches the plan or the state; change `*_wo_version` to send a new value.

```terraform
# Every argument the kind supports. From config-schemas examples/JailPolicy/advanced.yaml.

resource "vulnetix_jail_policy" "payments_api" {
  name                   = "payments-api"
  default_staleness_days = 7
  description            = "Payments ships with no overdue criticals, no live secrets, and moment.js on its way out."
  enabled                = true
  enforcement_mode       = "enforce"
  on_stale               = "fail"
  repository             = "acme/payments-api"
  rules = [
    {
      aggregate      = "COUNT_PAST_SLA"
      category       = "SCA"
      ecosystem      = "npm"
      enabled        = true
      flags          = ["reachable", "hasFix"]
      kind           = "VULN"
      label          = "No reachable critical npm finding past its SLA"
      on_breach      = "BLOCK"
      on_stale       = "warn"
      operator       = "GT"
      ratchet        = false
      severity_floor = "critical"
      staleness_days = 3
      threshold_low  = 0
      tool_name      = "vulnetix"
      when           = "fixAvailability == 'available' && (kev || epss >= 0.1)"
    },
    {
      aggregate      = "COUNT"
      deadline       = "2027-03-31"
      kind           = "GOAL"
      label          = "Off moment.js"
      on_breach      = "WARN"
      operator       = "GT"
      ratchet        = true
      target_pattern = "pkg:npm/moment@*"
      threshold_low  = 0
    },
    {
      aggregate      = "COUNT"
      category       = "SECRETS"
      kind           = "HYGIENE"
      label          = "Secrets stay between 0 and 2"
      on_stale       = "pass"
      operator       = "OUTSIDE"
      threshold_high = 2
      threshold_low  = 0
    },
  ]
  labels = {
    team = "payments"
  }
  annotations = {
    ticket = "https://tracker.example.com/SEC-1980"
  }
}
```

### Reference

Every argument that has a schema default, set to that default, plus the required arguments. Planned against a tenant it only sets values the tenant left unset. Every kind's reference file is also in [Vulnetix/platform-iac-examples](https://github.com/Vulnetix/platform-iac-examples), with the same document as YAML.

```terraform
# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/JailPolicy/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_jail_policy" "org_default" {
  name                   = "org-default"
  default_staleness_days = 7
  enabled                = true
  enforcement_mode       = "enforce"
  on_stale               = "fail"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The object's name (`metadata.name`), unique per kind within the tenant. Changing it replaces the object.

### Optional

- `annotations` (Map of String) Free key/value notes for people and tools (`metadata.annotations`), such as a ticket link.
- `default_staleness_days` (Number) How old the evidence a rule reads may be, in days, for rules that do not set stalenessDays. Schema default: `7`.
- `description` (String)
- `enabled` (Boolean) Schema default: `true`.
- `enforcement_mode` (String) enforce applies each rule's onBreach. warn downgrades every breach to a warning, which disarms the jail without a CLI release. One of `enforce`, `warn`. Schema default: `enforce`.
- `labels` (Map of String) Free key/value labels (`metadata.labels`). Stored with the object, never interpreted by the server.
- `on_stale` (String) What a rule does when the evidence it needs is stale or absent, for rules that do not set their own. fail is the default: a gate that passes on missing evidence is worse than no gate. One of `fail`, `warn`, `pass`. Schema default: `fail`.
- `repository` (String) The repository this policy overrides the default for. Leave it out for the tenant default; a tenant has at most one default and one policy per repository.
- `rules` (Attributes List) The rules, in order. Every enabled rule is evaluated on every run; order decides presentation and which rule gives the headline reason, and never short-circuits. Rules are matched to stored rules by label, so labels must be unique within a policy. (see [below for nested schema](#nestedatt--rules))

### Read-Only

- `id` (String) The object's name, the same as `name`.
- `status` (Attributes) What the server says about the object beyond its spec. (see [below for nested schema](#nestedatt--status))

<a id="nestedatt--rules"></a>
### Nested Schema for `rules`

Required:

- `aggregate` (String) COUNT: how many. COUNT_PAST_SLA: how many are past the remediation window of the tenant's RemediationPolicy (VULN only). MAX_AGE_DAYS: the age of the oldest, in days. PERCENT: the share of the population, 0 to 100. One of `COUNT`, `COUNT_PAST_SLA`, `MAX_AGE_DAYS`, `PERCENT`.
- `kind` (String) VULN: vulnerability findings. EOL: end-of-life dependencies. GOAL: a migration target, selected by targetPattern. HYGIENE: findings from the other scanner categories. One of `VULN`, `EOL`, `GOAL`, `HYGIENE`.
- `label` (String) The rule's name in jail output and the console, and its key within the policy.
- `operator` (String) How the aggregate is compared with the thresholds. BETWEEN and OUTSIDE use both thresholds, inclusive. One of `GT`, `GTE`, `LT`, `LTE`, `EQ`, `NEQ`, `BETWEEN`, `OUTSIDE`.
- `threshold_low` (Number) The threshold, or the lower bound for BETWEEN and OUTSIDE.

Optional:

- `category` (String) Only findings from this scanner category. A HYGIENE rule without one draws from SECRETS, LICENSE, IAC, OCI and MALWARE. One of `SCA`, `SAST`, `SECRETS`, `IAC`, `OCI`, `LICENSE`, `MALWARE`.
- `deadline` (String) GOAL only. Until this time a breach is downgraded to a warning, whatever onBreach says.
- `ecosystem` (String) Only dependencies from this package ecosystem.
- `enabled` (Boolean) Schema default: `true`.
- `flags` (Set of String) Every flag listed must hold for a finding to count. kev: in a known-exploited catalogue. malicious: known malware. hasFix: a fixed version exists. reachable: the vulnerable code is called. direct: a direct dependency.
- `on_breach` (String) One of `BLOCK`, `WARN`. Schema default: `BLOCK`.
- `on_stale` (String) Overrides the policy's onStale for this rule. One of `fail`, `warn`, `pass`.
- `ratchet` (Boolean) GOAL only. Breach when the value rose since the last assessment of the same rule on the same branch, so a migration holds its ground before the deadline. With no earlier assessment the value is recorded and the rule passes. Schema default: `false`.
- `severity_floor` (String) Only findings at or above this severity. One of `low`, `medium`, `high`, `critical`.
- `staleness_days` (Number) Overrides the policy's defaultStalenessDays for this rule.
- `target_pattern` (String) GOAL only, and required there: a path.Match glob over the identifier of each fact (a purl, a package name, a rule id).
- `threshold_high` (Number) The upper bound for BETWEEN and OUTSIDE; required by them. It must not be below thresholdLow.
- `tool_name` (String) Only findings reported by this tool.
- `when` (String) An extra condition each fact must meet to be counted, over the finding variables. A when that does not compile, or that uses a finding variable the facts of the rule kind do not carry yet, skips the rule with a warning. A fact it fails on leaves the rule indeterminate (verdict indeterminate, CLI exit 3). The factory gate fails closed instead: the rule breaches. A CEL expression, evaluated in the `finding` environment.


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
# Import an existing JailPolicy by its name (metadata.name).
terraform import vulnetix_jail_policy.org_default org-default
```
