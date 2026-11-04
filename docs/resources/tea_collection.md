---
page_title: "vulnetix_tea_collection Resource - vulnetix"
subcategory: ""
description: |-
  The collection of one published release: the versioned set of artifacts (SBOMs, VEX documents, attestations) TEA serves for it. TEA gives a collection the UUID of its release, so a collection exists only for a release a CI pipeline has already published, and plan fails until it has been. Declaring the collection makes sure it exists and gives a TeaAccessPolicy something to name, such as an embargo on one release. The artifacts are runtime data: CI adds them, and each addition moves the collection's version on, which status reports and which is not a change to this document. A collection this document creates starts at version 1 with the update reason INITIAL_RELEASE; one CI created first is adopted as it is. Removing the document deletes the collection while it has never held an artifact; otherwise the collection stays with its release and only the config management of it ends.
  An object of this kind can reference TeaProduct, TeaComponent objects. Reference those resources' attributes, or add depends_on, so Terraform applies them first.
  Manages the TeaCollection config kind through /config/objects/TeaCollection/{name}. Applies are recorded with the provider's manager, so a GitOps vulnetix config apply --prune leaves this object alone.
  Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (most of them, such as ApiSpecSource, RiskStrategy and VexSource) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object, or lets go of it for a kind the server never deletes (such as Organization). Some kinds wait on a person after they are applied (an invited Member, a Slack connection): status.phase says so, and Terraform does not wait.
---

# vulnetix_tea_collection (Resource)

The collection of one published release: the versioned set of artifacts (SBOMs, VEX documents, attestations) TEA serves for it. TEA gives a collection the UUID of its release, so a collection exists only for a release a CI pipeline has already published, and plan fails until it has been. Declaring the collection makes sure it exists and gives a TeaAccessPolicy something to name, such as an embargo on one release. The artifacts are runtime data: CI adds them, and each addition moves the collection's version on, which status reports and which is not a change to this document. A collection this document creates starts at version 1 with the update reason INITIAL_RELEASE; one CI created first is adopted as it is. Removing the document deletes the collection while it has never held an artifact; otherwise the collection stays with its release and only the config management of it ends.

An object of this kind can reference `TeaProduct`, `TeaComponent` objects. Reference those resources' attributes, or add `depends_on`, so Terraform applies them first.

Manages the `TeaCollection` config kind through `/config/objects/TeaCollection/{name}`. Applies are recorded with the provider's `manager`, so a GitOps `vulnetix config apply --prune` leaves this object alone.

Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (most of them, such as `ApiSpecSource`, `RiskStrategy` and `VexSource`) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object, or lets go of it for a kind the server never deletes (such as `Organization`). Some kinds wait on a person after they are applied (an invited `Member`, a Slack connection): `status.phase` says so, and Terraform does not wait.

## Example Usage

### Minimal

The arguments the kind requires, and nothing else.

```terraform
# Only the arguments the kind requires. From config-schemas examples/TeaCollection/minimal.yaml.

resource "vulnetix_tea_collection" "acme_payments_2_4_0" {
  name    = "acme-payments-2-4-0"
  product = "acme-payments"
  release = "2.4.0"
}
```

### Advanced

Every argument the kind supports. A secret is a write-only `*_wo` argument fed from an ephemeral variable, so it never reaches the plan or the state; change `*_wo_version` to send a new value.

```terraform
# Every argument the kind supports. From config-schemas examples/TeaCollection/advanced.yaml.

resource "vulnetix_tea_collection" "acme_payments_2_4_0" {
  name    = "acme-payments-2-4-0"
  product = "acme-payments"
  release = "2.4.0"
  annotations = {
    note = "Embargoed until the advisory is published; see the TeaAccessPolicy of the same name."
  }
}

resource "vulnetix_tea_collection" "payments_api_2_4_0" {
  name      = "payments-api-2-4-0"
  component = "payments-api"
  release   = "2.4.0"
}
```

### Reference

Every argument that has a schema default, set to that default, plus the required arguments. Planned against a tenant it only sets values the tenant left unset. Every kind's reference file is also in [Vulnetix/platform-iac-examples](https://github.com/Vulnetix/platform-iac-examples), with the same document as YAML.

```terraform
# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/TeaCollection/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_tea_collection" "acme_payments_2_4_0" {
  name    = "acme-payments-2-4-0"
  product = "acme-payments"
  release = "2.4.0"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The object's name (`metadata.name`), unique per kind within the tenant. Changing it replaces the object.
- `release` (String) The release's version, exactly as CI published it, such as 2.4.0.

### Optional

- `annotations` (Map of String) Free key/value notes for people and tools (`metadata.annotations`), such as a ticket link.
- `component` (String) The metadata.name of a TeaComponent, for the collection of one of its component releases.
- `labels` (Map of String) Free key/value labels (`metadata.labels`). Stored with the object, never interpreted by the server.
- `product` (String) The metadata.name of a TeaProduct, for the collection of one of its product releases.

### Read-Only

- `id` (String) The object's name, the same as `name`.
- `status` (Attributes) What the server says about the object beyond its spec. (see [below for nested schema](#nestedatt--status))

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
# Import an existing TeaCollection by its name (metadata.name).
terraform import vulnetix_tea_collection.acme_payments_2_4_0 acme-payments-2-4-0
```
