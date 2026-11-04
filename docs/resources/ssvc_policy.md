---
page_title: "vulnetix_ssvc_policy Resource - vulnetix"
subcategory: ""
description: |-
  How the tenant makes its CISA SSVC decision (https://www.cisa.gov/stakeholder-specific-vulnerability-categorization-ssvc) for a finding: Track, Track*, Attend or Act. The decision table is the CISA Coordinator table 2.0.3 and is fixed: four decision points (exploitation, automatable, technical impact, mission and well-being) go in. This document configures how Vulnetix's signals become those four values (CEL over the finding variables, and the repository's criticality for mission and well-being), which table cells the organisation reads differently (mapping), which findings it decides by rule instead (overrides), and what each outcome asks of it (outcomes). A field left out takes the default, which is the policy of a tenant with no document. One per tenant.
  Manages the SsvcPolicy config kind through /config/objects/SsvcPolicy/{name}. Applies are recorded with the provider's manager, so a GitOps vulnetix config apply --prune leaves this object alone.
  Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (most of them, such as ApiSpecSource, RiskStrategy and VexSource) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object, or lets go of it for a kind the server never deletes (such as Organization). Some kinds wait on a person after they are applied (an invited Member, a Slack connection): status.phase says so, and Terraform does not wait.
---

# vulnetix_ssvc_policy (Resource)

How the tenant makes its CISA SSVC decision (https://www.cisa.gov/stakeholder-specific-vulnerability-categorization-ssvc) for a finding: Track, Track*, Attend or Act. The decision table is the CISA Coordinator table 2.0.3 and is fixed: four decision points (exploitation, automatable, technical impact, mission and well-being) go in. This document configures how Vulnetix's signals become those four values (CEL over the finding variables, and the repository's criticality for mission and well-being), which table cells the organisation reads differently (mapping), which findings it decides by rule instead (overrides), and what each outcome asks of it (outcomes). A field left out takes the default, which is the policy of a tenant with no document. One per tenant.

Manages the `SsvcPolicy` config kind through `/config/objects/SsvcPolicy/{name}`. Applies are recorded with the provider's `manager`, so a GitOps `vulnetix config apply --prune` leaves this object alone.

Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (most of them, such as `ApiSpecSource`, `RiskStrategy` and `VexSource`) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object, or lets go of it for a kind the server never deletes (such as `Organization`). Some kinds wait on a person after they are applied (an invited `Member`, a Slack connection): `status.phase` says so, and Terraform does not wait.

## Example Usage

### Minimal

The arguments the kind requires, and nothing else.

```terraform
# Only the arguments the kind requires. From config-schemas examples/SsvcPolicy/minimal.yaml.

resource "vulnetix_ssvc_policy" "default" {
}
```

### Advanced

Every argument the kind supports. A secret is a write-only `*_wo` argument fed from an ephemeral variable, so it never reaches the plan or the state; change `*_wo_version` to send a new value.

```terraform
# Every argument the kind supports. From config-schemas examples/SsvcPolicy/advanced.yaml.

resource "vulnetix_ssvc_policy" "default" {
  automatable = "cvss >= 9.0 || epss >= 0.5"
  exploitation = {
    active = "kev || exploitMaturity in ['active', 'weaponized']"
    poc    = "exploitMaturity == 'poc' || exploitCount > 0"
  }
  mapping = [
    {
      automatable        = ["no", "yes"]
      exploitation       = ["active"]
      mission_well_being = ["medium", "high"]
      outcome            = "Act"
      technical_impact   = ["total"]
    },
  ]
  mission_well_being = {
    critical     = "high"
    high         = "high"
    low          = "low"
    medium       = "medium"
    unclassified = "medium"
  }
  outcomes = {
    act = {
      due_days = 3
      priority = "immediate"
    }
    attend = {
      due_days = 14
      priority = "high"
    }
    track = {
      due_days = 90
      priority = "low"
    }
    track_star = {
      due_days = 60
      priority = "medium"
    }
  }
  overrides = [
    {
      name    = "card data in scope"
      outcome = "Act"
      when    = "labels['pci'] == 'in-scope' && severity in ['critical', 'high']"
    },
    {
      name    = "no fix and not reachable"
      outcome = "Track"
      when    = "!reachable && fixAvailability == 'no_fix'"
    },
  ]
  technical_impact = "cvss >= 9.0 || malicious"
  labels = {
    team = "platform-security"
  }
  annotations = {
    ticket = "https://tracker.example.com/SEC-3120"
  }
}
```

### Reference

Every argument that has a schema default, set to that default, plus the required arguments. Planned against a tenant it only sets values the tenant left unset. Every kind's reference file is also in [Vulnetix/platform-iac-examples](https://github.com/Vulnetix/platform-iac-examples), with the same document as YAML.

```terraform
# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/SsvcPolicy/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_ssvc_policy" "default" {
  automatable = "cvss >= 9.0 || epss >= 0.7"
  exploitation = {
    active = "kev || exploitMaturity in ['active', 'weaponized']"
    poc    = "exploitMaturity == 'poc' || exploitCount > 0"
  }
  mission_well_being = {
    critical     = "high"
    high         = "high"
    low          = "low"
    medium       = "medium"
    unclassified = "medium"
  }
  outcomes = {
    act = {
      due_days = 0
    }
    attend = {
      due_days = 0
    }
    track = {
      due_days = 0
    }
    track_star = {
      due_days = 0
    }
  }
  technical_impact = "cvss >= 9.0 || malicious"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `annotations` (Map of String) Free key/value notes for people and tools (`metadata.annotations`), such as a ticket link.
- `automatable` (String) Holds when an attacker can reliably automate steps 1 to 4 of the kill chain (reconnaissance, weaponization, delivery, exploitation) for the finding. CEL over the finding variables. Schema default: `cvss \u003e= 9.0 || epss \u003e= 0.7`. A CEL expression, evaluated in the `finding` environment.
- `exploitation` (Attributes) When a finding counts as exploited. Active wins when both conditions hold; when neither does, exploitation is none. (see [below for nested schema](#nestedatt--exploitation))
- `labels` (Map of String) Free key/value labels (`metadata.labels`). Stored with the object, never interpreted by the server.
- `mapping` (Attributes List) Reads cells of the decision table differently. A rule matches every cell whose decision points are among the values it lists; a decision point it leaves out matches any value. When several rules match a cell, the last one wins. A rule changes a cell for the whole tenant: to decide by a property of the finding or its repository, use overrides. (see [below for nested schema](#nestedatt--mapping))
- `mission_well_being` (Attributes) The mission and well-being value of a finding in a repository, by the repository's criticality (RepositorySettings spec.criticality). unclassified is a repository with no criticality. (see [below for nested schema](#nestedatt--mission_well_being))
- `name` (String) The object's name. SsvcPolicy is one per tenant, so it is always `default`.
- `outcomes` (Attributes) What each outcome asks of the organisation. The priority is what a finding with that outcome is ranked as; dueDays is its remediation window, and 0 or left out means the standard update cycle. (see [below for nested schema](#nestedatt--outcomes))
- `overrides` (Attributes List) Decides the findings a condition holds for by rule, after the table and the mapping. The first rule that holds sets the outcome, so list the most specific first. (see [below for nested schema](#nestedatt--overrides))
- `technical_impact` (String) Holds when exploiting the finding gives total control of the software or total disclosure of its information; otherwise the impact is partial. CEL over the finding variables. Schema default: `cvss \u003e= 9.0 || malicious`. A CEL expression, evaluated in the `finding` environment.

### Read-Only

- `id` (String) The object's name, the same as `name`.
- `status` (Attributes) What the server says about the object beyond its spec. (see [below for nested schema](#nestedatt--status))

<a id="nestedatt--exploitation"></a>
### Nested Schema for `exploitation`

Optional:

- `active` (String) Credible evidence of exploitation in the wild. CEL over the finding variables. Schema default: `kev || exploitMaturity in ['active', 'weaponized']`. A CEL expression, evaluated in the `finding` environment.
- `poc` (String) A public proof of concept or a well-known method exists. CEL over the finding variables. Schema default: `exploitMaturity == 'poc' || exploitCount \u003e 0`. A CEL expression, evaluated in the `finding` environment.


<a id="nestedatt--mapping"></a>
### Nested Schema for `mapping`

Required:

- `outcome` (String) The SSVC decision. Track and Track* are remediated in the standard update cycle (Track* is watched more closely), Attend sooner, Act as soon as possible. One of `Track`, `Track*`, `Attend`, `Act`.

Optional:

- `automatable` (Set of String)
- `exploitation` (Set of String)
- `mission_well_being` (Set of String)
- `technical_impact` (Set of String)


<a id="nestedatt--mission_well_being"></a>
### Nested Schema for `mission_well_being`

Optional:

- `critical` (String) Mission and well-being: the combination of mission prevalence and public well-being impact. One of `low`, `medium`, `high`. Schema default: `high`.
- `high` (String) Mission and well-being: the combination of mission prevalence and public well-being impact. One of `low`, `medium`, `high`. Schema default: `high`.
- `low` (String) Mission and well-being: the combination of mission prevalence and public well-being impact. One of `low`, `medium`, `high`. Schema default: `low`.
- `medium` (String) Mission and well-being: the combination of mission prevalence and public well-being impact. One of `low`, `medium`, `high`. Schema default: `medium`.
- `unclassified` (String) Mission and well-being: the combination of mission prevalence and public well-being impact. One of `low`, `medium`, `high`. Schema default: `medium`.


<a id="nestedatt--outcomes"></a>
### Nested Schema for `outcomes`

Optional:

- `act` (Attributes) (see [below for nested schema](#nestedatt--outcomes--act))
- `attend` (Attributes) (see [below for nested schema](#nestedatt--outcomes--attend))
- `track` (Attributes) (see [below for nested schema](#nestedatt--outcomes--track))
- `track_star` (Attributes) (see [below for nested schema](#nestedatt--outcomes--track_star))

<a id="nestedatt--outcomes--act"></a>
### Nested Schema for `outcomes.act`

Optional:

- `due_days` (Number) The remediation window in days. 0 or left out is the standard update cycle. Schema default: `0`.
- `priority` (String) low, medium, high or immediate. Defaults: Track low, Track* medium, Attend high, Act immediate. One of `low`, `medium`, `high`, `immediate`.


<a id="nestedatt--outcomes--attend"></a>
### Nested Schema for `outcomes.attend`

Optional:

- `due_days` (Number) The remediation window in days. 0 or left out is the standard update cycle. Schema default: `0`.
- `priority` (String) low, medium, high or immediate. Defaults: Track low, Track* medium, Attend high, Act immediate. One of `low`, `medium`, `high`, `immediate`.


<a id="nestedatt--outcomes--track"></a>
### Nested Schema for `outcomes.track`

Optional:

- `due_days` (Number) The remediation window in days. 0 or left out is the standard update cycle. Schema default: `0`.
- `priority` (String) low, medium, high or immediate. Defaults: Track low, Track* medium, Attend high, Act immediate. One of `low`, `medium`, `high`, `immediate`.


<a id="nestedatt--outcomes--track_star"></a>
### Nested Schema for `outcomes.track_star`

Optional:

- `due_days` (Number) The remediation window in days. 0 or left out is the standard update cycle. Schema default: `0`.
- `priority` (String) low, medium, high or immediate. Defaults: Track low, Track* medium, Attend high, Act immediate. One of `low`, `medium`, `high`, `immediate`.



<a id="nestedatt--overrides"></a>
### Nested Schema for `overrides`

Required:

- `outcome` (String) The SSVC decision. Track and Track* are remediated in the standard update cycle (Track* is watched more closely), Attend sooner, Act as soon as possible. One of `Track`, `Track*`, `Attend`, `Act`.
- `when` (String) CEL over the finding variables, including the repository's labels, owner and criticality. A CEL expression, evaluated in the `finding` environment.

Optional:

- `name` (String) Shown beside the decision when this rule set it.


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
# Import an existing SsvcPolicy by its name (metadata.name).
terraform import vulnetix_ssvc_policy.default default
```
