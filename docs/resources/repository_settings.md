---
page_title: "vulnetix_repository_settings Resource - vulnetix"
subcategory: ""
description: |-
  The Software Factory's controls for one repository: the settings panel of the Factory repositories page. The repository must be one the factory can see, through a managed organisation or another signal. metadata.name is the repository lower-cased, with every run of characters other than a to z and 0 to 9 replaced by one hyphen: acme-payments-api for acme/payments-api. A field left out means its default. The asset fields (labels, owner and criticality) describe the repository as an asset of the business: they are stored per repository, shown on the Repositories page, and become the labels, owner and criticality variables of the CEL finding context, so a risk, suppression, jail or gate rule can select on them. They need the saas migration 20261228000000_repository_asset; until it is applied a document that sets one is refused. Deleting the document removes the repository's own settings, so it runs at the tenant's autonomy.default again, removes its monthly token cap and clears its asset fields; the record of the change is kept in the audit trail.
  Manages the RepositorySettings config kind through /config/objects/RepositorySettings/{name}. Applies are recorded with the provider's manager, so a GitOps vulnetix config apply --prune leaves this object alone.
  Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (most of them, such as ApiSpecSource, RiskStrategy and VexSource) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object, or lets go of it for a kind the server never deletes (such as Organization). Some kinds wait on a person after they are applied (an invited Member, a Slack connection): status.phase says so, and Terraform does not wait.
---

# vulnetix_repository_settings (Resource)

The Software Factory's controls for one repository: the settings panel of the Factory repositories page. The repository must be one the factory can see, through a managed organisation or another signal. metadata.name is the repository lower-cased, with every run of characters other than a to z and 0 to 9 replaced by one hyphen: acme-payments-api for acme/payments-api. A field left out means its default. The asset fields (labels, owner and criticality) describe the repository as an asset of the business: they are stored per repository, shown on the Repositories page, and become the `labels`, `owner` and `criticality` variables of the CEL finding context, so a risk, suppression, jail or gate rule can select on them. They need the saas migration 20261228000000_repository_asset; until it is applied a document that sets one is refused. Deleting the document removes the repository's own settings, so it runs at the tenant's autonomy.default again, removes its monthly token cap and clears its asset fields; the record of the change is kept in the audit trail.

Manages the `RepositorySettings` config kind through `/config/objects/RepositorySettings/{name}`. Applies are recorded with the provider's `manager`, so a GitOps `vulnetix config apply --prune` leaves this object alone.

Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (most of them, such as `ApiSpecSource`, `RiskStrategy` and `VexSource`) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object, or lets go of it for a kind the server never deletes (such as `Organization`). Some kinds wait on a person after they are applied (an invited `Member`, a Slack connection): `status.phase` says so, and Terraform does not wait.

## Example Usage

### Minimal

The arguments the kind requires, and nothing else.

```terraform
# Only the arguments the kind requires. From config-schemas examples/RepositorySettings/minimal.yaml.

resource "vulnetix_repository_settings" "acme_payments_api" {
  name       = "acme-payments-api"
  autonomy   = "draft_pr"
  repository = "acme/payments-api"
}
```

### Advanced

Every argument the kind supports. A secret is a write-only `*_wo` argument fed from an ephemeral variable, so it never reaches the plan or the state; change `*_wo_version` to send a new value.

```terraform
# Every argument the kind supports. From config-schemas examples/RepositorySettings/advanced.yaml.

resource "vulnetix_repository_settings" "acme_payments_api" {
  name = "acme-payments-api"
  asset_labels = {
    data = "pii"
    pci  = "in-scope"
    tier = "prod"
  }
  autonomy          = "suggest"
  criticality       = "critical"
  defaults_off      = ["lint-sweep", "test-sweep"]
  kanban_mode       = "skip_review"
  monthly_token_cap = 2000000
  owner = {
    kind = "team"
    name = "payments"
  }
  protected_branches = true
  repository         = "acme/payments-api"
  tea_publish        = true
  labels = {
    team = "payments"
  }
  annotations = {
    ticket = "https://tracker.example.com/PAY-1501"
  }
}
```

### Reference

Every argument that has a schema default, set to that default, plus the required arguments. Planned against a tenant it only sets values the tenant left unset. Every kind's reference file is also in [Vulnetix/platform-iac-examples](https://github.com/Vulnetix/platform-iac-examples), with the same document as YAML.

```terraform
# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/RepositorySettings/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_repository_settings" "acme_payments_api" {
  name = "acme-payments-api"
  asset_labels = {
  }
  autonomy           = "draft_pr"
  defaults_off       = []
  kanban_mode        = "off"
  protected_branches = true
  repository         = "acme/payments-api"
  tea_publish        = false
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `autonomy` (String) How far the factory goes on its own in this repository, from observe (the least) through suggest and draft_pr to auto_merge_low_risk (the most, shown as approve when ready). It may not be above the tenant's autonomy.max setting; the plan is refused when it is. One of `observe`, `suggest`, `draft_pr`, `auto_merge_low_risk`.
- `name` (String) The object's name (`metadata.name`), unique per kind within the tenant. Changing it replaces the object.
- `repository` (String) The repository: owner/name on github.com, or host/owner/name on another forge (codeberg.org/acme/api). Matched without regard to case.

### Optional

- `annotations` (Map of String) Free key/value notes for people and tools (`metadata.annotations`), such as a ticket link.
- `asset_labels` (Map of String) Key/value labels on the repository as an asset, such as tier: prod or data: pii. At most 50. Keys are lower case letters, digits and . _ / - (at most 63 characters, starting and ending with a letter or digit); values are at most 256 characters. In CEL: labels['tier'] == 'prod'. A label dropped from the document is removed from the repository. Schema default: `{}`.
- `criticality` (String) How much the business depends on the repository: low, medium, high or critical. Left out, the repository is unclassified and `criticality` is empty in CEL. It is the mission and well-being input of the SSVC decision. One of `low`, `medium`, `high`, `critical`.
- `defaults_off` (Set of String) The default workflows this repository opts out of, by key. Every default the tenant has on applies to the rest. Schema default: `[]`.
- `kanban_mode` (String) Which kanban cards are mirrored to issues in this repository. all: every card. skip_review: a card waits until it leaves Review. off: none. It also sets the mode of the repository's existing kanban bindings. One of `all`, `skip_review`, `off`. Schema default: `off`.
- `labels` (Map of String) Free key/value labels (`metadata.labels`). Stored with the object, never interpreted by the server.
- `monthly_token_cap` (Number) The most tokens the factory may spend on this repository in a month. Left out, there is no cap of its own. Setting it creates the repository's budget on the Budgets page; removing it removes that budget.
- `owner` (Attributes) Who owns the repository: a team, a user or an email address. In CEL, `owner` is the name. (see [below for nested schema](#nestedatt--owner))
- `protected_branches` (Boolean) Never push to the default branch or a protected branch. Only an owner or an admin may set it to false. Schema default: `true`.
- `tea_publish` (Boolean) Publish the repository's releases to TEA. Schema default: `false`.

### Read-Only

- `id` (String) The object's name, the same as `name`.
- `status` (Attributes) What the server says about the object beyond its spec. (see [below for nested schema](#nestedatt--status))

<a id="nestedatt--owner"></a>
### Nested Schema for `owner`

Required:

- `kind` (String) team: a team name. user: a user name on the forge or in Vulnetix. email: an email address. One of `team`, `user`, `email`.
- `name` (String) The team name, the user name or the email address.


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
# Import an existing RepositorySettings by its name (metadata.name).
terraform import vulnetix_repository_settings.acme_payments_api acme-payments-api
```
