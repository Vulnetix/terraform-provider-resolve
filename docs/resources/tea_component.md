---
page_title: "vulnetix_tea_component Resource - vulnetix"
subcategory: ""
description: |-
  A component the tenant publishes through the Transparency Exchange API (TEA): a part a product is assembled from, such as a service, a library or an installer. Component releases, and the distributions and artifacts under them, are published by CI pipelines and are not declared here; a product release names the component releases it is built from when CI publishes it. A component given the same name as a TeaProduct is the same subject to consumers, who then see one node carrying both identities. The server assigns the component's UUID and TEI and reports them in status. Who may read it is declared with a TeaAccessPolicy; until one exists the component is private. Removing the document deletes the component while no release has been published under it. Once CI has published a release, the component and its releases stay, and only the config management of the component ends.
  Manages the TeaComponent config kind through /config/objects/TeaComponent/{name}. Applies are recorded with the provider's manager, so a GitOps vulnetix config apply --prune leaves this object alone.
  Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (most of them, such as ApiSpecSource, RiskStrategy and VexSource) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object, or lets go of it for a kind the server never deletes (such as Organization). Some kinds wait on a person after they are applied (an invited Member, a Slack connection): status.phase says so, and Terraform does not wait.
---

# vulnetix_tea_component (Resource)

A component the tenant publishes through the Transparency Exchange API (TEA): a part a product is assembled from, such as a service, a library or an installer. Component releases, and the distributions and artifacts under them, are published by CI pipelines and are not declared here; a product release names the component releases it is built from when CI publishes it. A component given the same name as a TeaProduct is the same subject to consumers, who then see one node carrying both identities. The server assigns the component's UUID and TEI and reports them in status. Who may read it is declared with a TeaAccessPolicy; until one exists the component is private. Removing the document deletes the component while no release has been published under it. Once CI has published a release, the component and its releases stay, and only the config management of the component ends.

Manages the `TeaComponent` config kind through `/config/objects/TeaComponent/{name}`. Applies are recorded with the provider's `manager`, so a GitOps `vulnetix config apply --prune` leaves this object alone.

Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (most of them, such as `ApiSpecSource`, `RiskStrategy` and `VexSource`) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object, or lets go of it for a kind the server never deletes (such as `Organization`). Some kinds wait on a person after they are applied (an invited `Member`, a Slack connection): `status.phase` says so, and Terraform does not wait.

## Example Usage

### Minimal

The arguments the kind requires, and nothing else.

```terraform
# Only the arguments the kind requires. From config-schemas examples/TeaComponent/minimal.yaml.

resource "vulnetix_tea_component" "payments_api" {
  name           = "payments-api"
  component_name = "payments-api"
}
```

### Advanced

Every argument the kind supports. A secret is a write-only `*_wo` argument fed from an ephemeral variable, so it never reaches the plan or the state; change `*_wo_version` to send a new value.

```terraform
# Every argument the kind supports. From config-schemas examples/TeaComponent/advanced.yaml.

resource "vulnetix_tea_component" "payments_api" {
  name           = "payments-api"
  component_name = "payments-api"
  identifiers = [
    {
      id_type  = "PURL"
      id_value = "pkg:oci/payments-api?repository_url=ghcr.io/acme"
    },
  ]
  labels = {
    team = "payments"
  }
}
```

### Reference

Every argument that has a schema default, set to that default, plus the required arguments. Planned against a tenant it only sets values the tenant left unset. Every kind's reference file is also in [Vulnetix/platform-iac-examples](https://github.com/Vulnetix/platform-iac-examples), with the same document as YAML.

```terraform
# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/TeaComponent/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_tea_component" "payments_api" {
  name           = "payments-api"
  component_name = "payments-api"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `component_name` (String) The component name TEA consumers see, such as payments-api. Unique among the tenant's components. The first time a name is published the server derives the component's UUID from it, so a CI pipeline that publishes releases under the same component name writes to this component. Changing the name later keeps the UUID.
- `name` (String) The object's name (`metadata.name`), unique per kind within the tenant. Changing it replaces the object.

### Optional

- `annotations` (Map of String) Free key/value notes for people and tools (`metadata.annotations`), such as a ticket link.
- `identifiers` (Attributes Set) Identifiers for the component besides its TEI, in the order consumers see them. The server assigns the TEI itself, so it is reported in status and cannot be declared. (see [below for nested schema](#nestedatt--identifiers))
- `labels` (Map of String) Free key/value labels (`metadata.labels`). Stored with the object, never interpreted by the server.

### Read-Only

- `id` (String) The object's name, the same as `name`.
- `status` (Attributes) What the server says about the object beyond its spec. (see [below for nested schema](#nestedatt--status))

<a id="nestedatt--identifiers"></a>
### Nested Schema for `identifiers`

Required:

- `id_type` (String) PURL: a package URL. CPE: a CPE 2.3 name. COMPLIANCE_DOCUMENT: a compliance document type the TEA specification names. One of `PURL`, `CPE`, `COMPLIANCE_DOCUMENT`.
- `id_value` (String) The identifier itself.


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
# Import an existing TeaComponent by its name (metadata.name).
terraform import vulnetix_tea_component.payments_api payments-api
```
