---
page_title: "vulnetix_tea_access_policy Resource - vulnetix"
subcategory: ""
description: |-
  Who may read one published TEA object: a product, a component or a release's collection. This is the Vulnetix publication overlay's access policy; the TEA consumption API itself does not say who may fetch what. The policy in force for any object is the narrowest one declared anywhere from the object up to its product, so a private product keeps everything beneath it private whatever those objects declare, and widening the product back restores each one to its own declaration. Plan refuses a policy wider than the one its parent enforces. An object no policy covers is private. Removing the document removes the declaration, and the object inherits from its parent again; for a product or component that means private. Mirroring to other TEA servers (publishTo in the publisher API) is not part of this kind, because this server does not mirror yet.
  An object of this kind can reference TeaProduct, TeaComponent, TeaCollection objects. Reference those resources' attributes, or add depends_on, so Terraform applies them first.
  Manages the TeaAccessPolicy config kind through /config/objects/TeaAccessPolicy/{name}. Applies are recorded with the provider's manager, so a GitOps vulnetix config apply --prune leaves this object alone.
  Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (most of them, such as ApiSpecSource, RiskStrategy and VexSource) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object, or lets go of it for a kind the server never deletes (such as Organization). Some kinds wait on a person after they are applied (an invited Member, a Slack connection): status.phase says so, and Terraform does not wait.
---

# vulnetix_tea_access_policy (Resource)

Who may read one published TEA object: a product, a component or a release's collection. This is the Vulnetix publication overlay's access policy; the TEA consumption API itself does not say who may fetch what. The policy in force for any object is the narrowest one declared anywhere from the object up to its product, so a private product keeps everything beneath it private whatever those objects declare, and widening the product back restores each one to its own declaration. Plan refuses a policy wider than the one its parent enforces. An object no policy covers is private. Removing the document removes the declaration, and the object inherits from its parent again; for a product or component that means private. Mirroring to other TEA servers (publishTo in the publisher API) is not part of this kind, because this server does not mirror yet.

An object of this kind can reference `TeaProduct`, `TeaComponent`, `TeaCollection` objects. Reference those resources' attributes, or add `depends_on`, so Terraform applies them first.

Manages the `TeaAccessPolicy` config kind through `/config/objects/TeaAccessPolicy/{name}`. Applies are recorded with the provider's `manager`, so a GitOps `vulnetix config apply --prune` leaves this object alone.

Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (most of them, such as `ApiSpecSource`, `RiskStrategy` and `VexSource`) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object, or lets go of it for a kind the server never deletes (such as `Organization`). Some kinds wait on a person after they are applied (an invited `Member`, a Slack connection): `status.phase` says so, and Terraform does not wait.

## Example Usage

### Minimal

The arguments the kind requires, and nothing else.

```terraform
# Only the arguments the kind requires. From config-schemas examples/TeaAccessPolicy/minimal.yaml.

resource "vulnetix_tea_access_policy" "acme_payments" {
  name       = "acme-payments"
  product    = "acme-payments"
  visibility = "private"
}
```

### Advanced

Every argument the kind supports. A secret is a write-only `*_wo` argument fed from an ephemeral variable, so it never reaches the plan or the state; change `*_wo_version` to send a new value.

```terraform
# Every argument the kind supports. From config-schemas examples/TeaAccessPolicy/advanced.yaml.

resource "vulnetix_tea_access_policy" "acme_payments" {
  name           = "acme-payments"
  confirm_public = true
  product        = "acme-payments"
  visibility     = "public"
  annotations = {
    approved-by = "https://tracker.example.com/SEC-210"
  }
}

resource "vulnetix_tea_access_policy" "payments_api" {
  name      = "payments-api"
  component = "payments-api"
  shared_with = [
    {
      expires           = "2027-06-30T00:00:00Z"
      organization      = "6f1c2a8e-3b4d-4e5f-8a9b-0c1d2e3f4a5b"
      organization_name = "Globex Corporation"
    },
  ]
  visibility = "shared"
}

resource "vulnetix_tea_access_policy" "acme_payments_2_4_0" {
  name       = "acme-payments-2-4-0"
  collection = "acme-payments-2-4-0"
  visibility = "private"
}
```

### Reference

Every argument that has a schema default, set to that default, plus the required arguments. Planned against a tenant it only sets values the tenant left unset. Every kind's reference file is also in [Vulnetix/platform-iac-examples](https://github.com/Vulnetix/platform-iac-examples), with the same document as YAML.

```terraform
# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/TeaAccessPolicy/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_tea_access_policy" "acme_payments" {
  name       = "acme-payments"
  product    = "acme-payments"
  visibility = "private"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The object's name (`metadata.name`), unique per kind within the tenant. Changing it replaces the object.
- `visibility` (String) private: only this organisation may read the object. shared: this organisation and the organisations in sharedWith may read it, and readers must authenticate. public: anyone may read it without authenticating. Public cannot be taken back in practice, because whatever was fetched while it was public stays fetched, so it also needs confirmPublic. One of `private`, `shared`, `public`.

### Optional

- `annotations` (Map of String) Free key/value notes for people and tools (`metadata.annotations`), such as a ticket link.
- `collection` (String) The metadata.name of a TeaCollection the policy applies to. It covers the artifacts in the collection, and is narrowed by the policy of the release's product or component.
- `component` (String) The metadata.name of a TeaComponent the policy applies to.
- `confirm_public` (Boolean) Required, and true, when visibility is public: the document's author confirms the object is meant for the open web. The publisher API asks for the same confirmation with ?confirm=public.
- `labels` (Map of String) Free key/value labels (`metadata.labels`). Stored with the object, never interpreted by the server.
- `product` (String) The metadata.name of a TeaProduct the policy applies to.
- `shared_with` (Attributes Set) The organisations entitled to read the object while visibility is shared. Allowed only with shared, so that a later change to private cannot leave a grant standing. Where several levels of the tree declare shared, a reader needs a grant at every level and the earliest expiry applies. (see [below for nested schema](#nestedatt--shared_with))

### Read-Only

- `id` (String) The object's name, the same as `name`.
- `status` (Attributes) What the server says about the object beyond its spec. (see [below for nested schema](#nestedatt--status))

<a id="nestedatt--shared_with"></a>
### Nested Schema for `shared_with`

Required:

- `organization` (String) The reading organisation's UUID on this server, the identity it authenticates as.

Optional:

- `expires` (String) When the grant lapses, such as the end of an evaluation or a contract. Left out, the grant lasts until it is removed.
- `organization_name` (String) The organisation's display name, for whoever reviews the grants. Advisory only: the UUID is the identity, so a rename never moves a grant.


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
# Import an existing TeaAccessPolicy by its name (metadata.name).
terraform import vulnetix_tea_access_policy.acme_payments acme-payments
```
