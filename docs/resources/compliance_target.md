---
page_title: "vulnetix_compliance_target Resource - vulnetix"
subcategory: ""
description: |-
  A compliance obligation the organisation holds itself to, and how often Vulnetix measures it. For the report frameworks (eu_ai_act, nist_ai_rmf, iso_42001, owasp_top_10, pci_dss, coverage, exposure, eu_cra, iso_27001, owasp_dsomm) each run creates a compliance report for the period, scoped to the repositories named, through the same machinery as a report made by hand, so the OSCAL export, the evidence and the assessor workflow are unchanged. For the three evaluated targets there is no report: each run stores an evaluation. slsa checks the evidence for a SLSA build level; ssdf checks the NIST SP 800-218 practices Vulnetix holds evidence for; crypto-policy evaluates a cryptography policy, written in CEL over the components of the Cryptography Bill of Materials, against the newest CBOM of each repository in scope. The first run is due when the target is created and later runs follow the cadence. Needs the saas migration 20261228040000_compliance_target; until it is applied, plan refuses a document naming it.
  Manages the ComplianceTarget config kind through /config/objects/ComplianceTarget/{name}. Applies are recorded with the provider's manager, so a GitOps vulnetix config apply --prune leaves this object alone.
  Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (most of them, such as ApiSpecSource, RiskStrategy and VexSource) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object, or lets go of it for a kind the server never deletes (such as Organization). Some kinds wait on a person after they are applied (an invited Member, a Slack connection): status.phase says so, and Terraform does not wait.
---

# vulnetix_compliance_target (Resource)

A compliance obligation the organisation holds itself to, and how often Vulnetix measures it. For the report frameworks (eu_ai_act, nist_ai_rmf, iso_42001, owasp_top_10, pci_dss, coverage, exposure, eu_cra, iso_27001, owasp_dsomm) each run creates a compliance report for the period, scoped to the repositories named, through the same machinery as a report made by hand, so the OSCAL export, the evidence and the assessor workflow are unchanged. For the three evaluated targets there is no report: each run stores an evaluation. slsa checks the evidence for a SLSA build level; ssdf checks the NIST SP 800-218 practices Vulnetix holds evidence for; crypto-policy evaluates a cryptography policy, written in CEL over the components of the Cryptography Bill of Materials, against the newest CBOM of each repository in scope. The first run is due when the target is created and later runs follow the cadence. Needs the saas migration 20261228040000_compliance_target; until it is applied, plan refuses a document naming it.

Manages the `ComplianceTarget` config kind through `/config/objects/ComplianceTarget/{name}`. Applies are recorded with the provider's `manager`, so a GitOps `vulnetix config apply --prune` leaves this object alone.

Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (most of them, such as `ApiSpecSource`, `RiskStrategy` and `VexSource`) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object, or lets go of it for a kind the server never deletes (such as `Organization`). Some kinds wait on a person after they are applied (an invited `Member`, a Slack connection): `status.phase` says so, and Terraform does not wait.

## Example Usage

### Minimal

The arguments the kind requires, and nothing else.

```terraform
# Only the arguments the kind requires. From config-schemas examples/ComplianceTarget/minimal.yaml.

resource "vulnetix_compliance_target" "quarterly_owasp" {
  name      = "quarterly-owasp"
  framework = "owasp_top_10"
}
```

### Advanced

Every argument the kind supports. A secret is a write-only `*_wo` argument fed from an ephemeral variable, so it never reaches the plan or the state; change `*_wo_version` to send a new value.

```terraform
# Every argument the kind supports. From config-schemas examples/ComplianceTarget/advanced.yaml.

resource "vulnetix_compliance_target" "eu_ai_act_payments" {
  name      = "eu-ai-act-payments"
  cadence   = "quarterly"
  enabled   = true
  framework = "eu_ai_act"
  report = {
    context = "The organisation is the provider of a high-risk credit scoring system (Annex III, 5b)."
    title   = "EU AI Act, payments"
  }
  scope = {
    repositories = ["acme/payments-api", "acme/fraud-model"]
  }
  labels = {
    team = "governance"
  }
  annotations = {
    ticket = "https://tracker.example.com/GRC-77"
  }
}

resource "vulnetix_compliance_target" "slsa_build_l2" {
  name      = "slsa-build-l2"
  cadence   = "monthly"
  framework = "slsa"
  slsa = {
    level = 2
  }
}

resource "vulnetix_compliance_target" "ssdf_core" {
  name      = "ssdf-core"
  cadence   = "quarterly"
  framework = "ssdf"
  ssdf = {
    min_met   = 80
    practices = ["PS", "PW", "RV"]
  }
}

resource "vulnetix_compliance_target" "pqc_readiness" {
  name    = "pqc-readiness"
  cadence = "monthly"
  crypto = {
    allow          = "quantumSafe || (primitive in ['block-cipher', 'hash', 'mac'] && classicalSecurityLevel >= 128)"
    applies_to     = "category == 'algorithm' && confidence != 'low'"
    max_violations = 0
  }
  framework = "crypto-policy"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `framework` (String) What is measured. A report framework creates a report on each run. slsa, ssdf and crypto-policy are evaluated and need their block below. One of `eu_ai_act`, `nist_ai_rmf`, `iso_42001`, `owasp_top_10`, `pci_dss`, `coverage`, `exposure`, `eu_cra`, `iso_27001`, `owasp_dsomm`, `slsa`, `ssdf`, `crypto-policy`.
- `name` (String) The object's name (`metadata.name`), unique per kind within the tenant. Changing it replaces the object.

### Optional

- `annotations` (Map of String) Free key/value notes for people and tools (`metadata.annotations`), such as a ticket link.
- `cadence` (String) How often it runs. The period of a report, or the window an evaluation looks back over, is the cadence before the run: a week, a month, a quarter or a year. One of `weekly`, `monthly`, `quarterly`, `yearly`. Schema default: `quarterly`.
- `crypto` (Attributes) For framework crypto-policy. A component of the CBOM is acceptable when allow holds for it. (see [below for nested schema](#nestedatt--crypto))
- `enabled` (Boolean) False keeps the target and its history and stops it running. Schema default: `true`.
- `labels` (Map of String) Free key/value labels (`metadata.labels`). Stored with the object, never interpreted by the server.
- `report` (Attributes) For a report framework: how the report is made. (see [below for nested schema](#nestedatt--report))
- `scope` (Attributes) What is measured. Left out, the whole organisation. (see [below for nested schema](#nestedatt--scope))
- `slsa` (Attributes) For framework slsa. Evaluated against the evidence Vulnetix holds: provenance records for the dependencies in scope, signatures and attestations on what the organisation publishes, and the Software Factory's minimum SLSA level (ScannerProfile slsaMin). Level 3 asks for what a hardened build platform leaves behind (verified provenance, witnessed signatures); Vulnetix cannot inspect the platform itself, so a pass says the evidence is there. (see [below for nested schema](#nestedatt--slsa))
- `ssdf` (Attributes) For framework ssdf. Checks the practices of NIST SP 800-218 that Vulnetix holds evidence for; a practice without evidence is not in the count. (see [below for nested schema](#nestedatt--ssdf))

### Read-Only

- `id` (String) The object's name, the same as `name`.
- `status` (Attributes) What the server says about the object beyond its spec. (see [below for nested schema](#nestedatt--status))

<a id="nestedatt--crypto"></a>
### Nested Schema for `crypto`

Required:

- `allow` (String) CEL over the cbom variables: true for a component the policy allows, such as `quantumSafe || (primitive == 'block-cipher' && classicalSecurityLevel >= 128)`. A component the expression is false for, or cannot be evaluated for, is a violation. A CEL expression, evaluated in the `cbom` environment.

Optional:

- `applies_to` (String) CEL over the cbom variables: which components the policy judges. Left out, algorithms only (`category == 'algorithm'`). Schema default: `category == 'algorithm'`. A CEL expression, evaluated in the `cbom` environment.
- `max_violations` (Number) How many violations still pass. Schema default: `0`.


<a id="nestedatt--report"></a>
### Nested Schema for `report`

Optional:

- `context` (String) The report's context text. For eu_ai_act this describes the AI system and the organisation's role (provider or deployer) and risk class, which the Act's obligations depend on.
- `title` (String) The report's name; the period is appended. Left out, the target's name.


<a id="nestedatt--scope"></a>
### Nested Schema for `scope`

Optional:

- `repositories` (Set of String) Only these repositories (owner/name), compared without regard to case. Left out or empty, every repository. Schema default: `[]`.


<a id="nestedatt--slsa"></a>
### Nested Schema for `slsa`

Required:

- `level` (Number) The SLSA Build level targeted: 1 provenance exists, 2 signed and verified provenance, 3 hardened build with witnessed signatures.


<a id="nestedatt--ssdf"></a>
### Nested Schema for `ssdf`

Optional:

- `min_met` (Number) The percentage of the checked practices that must be met for the target to pass. Schema default: `100`.
- `practices` (Set of String) The practice groups checked: PO prepare the organisation, PS protect the software, PW produce well-secured software, RV respond to vulnerabilities. Left out, all four.


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
# Import an existing ComplianceTarget by its name (metadata.name).
terraform import vulnetix_compliance_target.quarterly_owasp quarterly-owasp
```
