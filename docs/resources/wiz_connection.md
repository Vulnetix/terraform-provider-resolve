---
page_title: "vulnetix_wiz_connection Resource - vulnetix"
subcategory: ""
description: |-
  The Wiz service account the tenant authorised. One per tenant. Wiz is connected as an addition to AWS: each claim Wiz makes is checked from inside the AWS account, so a plan is refused until AWS is connected, and on a plan without cloud integrations. Apply checks the credential against Wiz before it is stored, including that the service account can read audit data (admin:audit). The client secret is kept in the tenant's vault and is never exported, shown in a plan or returned by any endpoint. A connection is created and deleted, never edited: any change replaces it (the old connection is deleted and its stored secret revoked, then the new one is created and checked). Deleting the document disconnects Wiz and stops the sync; findings already ingested from Wiz, and their triage history, are kept.
  Manages the WizConnection config kind through /config/objects/WizConnection/{name}. Applies are recorded with the provider's manager, so a GitOps vulnetix config apply --prune leaves this object alone.
  Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (most of them, such as ApiSpecSource, RiskStrategy and VexSource) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object, or lets go of it for a kind the server never deletes (such as Organization). Some kinds wait on a person after they are applied (an invited Member, a Slack connection): status.phase says so, and Terraform does not wait.
---

# vulnetix_wiz_connection (Resource)

The Wiz service account the tenant authorised. One per tenant. Wiz is connected as an addition to AWS: each claim Wiz makes is checked from inside the AWS account, so a plan is refused until AWS is connected, and on a plan without cloud integrations. Apply checks the credential against Wiz before it is stored, including that the service account can read audit data (admin:audit). The client secret is kept in the tenant's vault and is never exported, shown in a plan or returned by any endpoint. A connection is created and deleted, never edited: any change replaces it (the old connection is deleted and its stored secret revoked, then the new one is created and checked). Deleting the document disconnects Wiz and stops the sync; findings already ingested from Wiz, and their triage history, are kept.

Manages the `WizConnection` config kind through `/config/objects/WizConnection/{name}`. Applies are recorded with the provider's `manager`, so a GitOps `vulnetix config apply --prune` leaves this object alone.

Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (most of them, such as `ApiSpecSource`, `RiskStrategy` and `VexSource`) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object, or lets go of it for a kind the server never deletes (such as `Organization`). Some kinds wait on a person after they are applied (an invited `Member`, a Slack connection): `status.phase` says so, and Terraform does not wait.

## Example Usage

### Minimal

The arguments the kind requires, and nothing else.

```terraform
# Only the arguments the kind requires. From config-schemas examples/WizConnection/minimal.yaml.

resource "vulnetix_wiz_connection" "default" {
  client_id         = "4ba3c9f2d7e14a6b8c0d1e2f3a4b5c6d"
  tenant_datacenter = "us17"
}
```

### Advanced

Every argument the kind supports. A secret is a write-only `*_wo` argument fed from an ephemeral variable, so it never reaches the plan or the state; change `*_wo_version` to send a new value.

```terraform
# Every argument the kind supports. From config-schemas examples/WizConnection/advanced.yaml.

variable "default_client_secret" {
  type      = string
  sensitive = true
  ephemeral = true
}

resource "vulnetix_wiz_connection" "default" {
  client_id                = "4ba3c9f2d7e14a6b8c0d1e2f3a4b5c6d"
  client_secret_wo         = var.default_client_secret
  client_secret_wo_version = 1
  tenant_datacenter        = "us17"
  labels = {
    owner = "cloud-security"
  }
  annotations = {
    note = "The service account needs the admin:audit permission."
  }
}
```

### Reference

Every argument that has a schema default, set to that default, plus the required arguments. Planned against a tenant it only sets values the tenant left unset. Every kind's reference file is also in [Vulnetix/platform-iac-examples](https://github.com/Vulnetix/platform-iac-examples), with the same document as YAML.

```terraform
# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/WizConnection/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_wiz_connection" "default" {
  client_id         = "4ba3c9f2d7e14a6b8c0d1e2f3a4b5c6d"
  tenant_datacenter = "us17"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `client_id` (String) The service account's client id, from Settings, Service Accounts in the Wiz console.
- `tenant_datacenter` (String) The datacenter segment of the tenant's API host: us17 for api.us17.wiz.io. It is the part of the Wiz console URL before .wiz.io.

### Optional

> **NOTE**: [Write-only arguments](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments) are supported in Terraform 1.11 and later.

- `annotations` (Map of String) Free key/value notes for people and tools (`metadata.annotations`), such as a ticket link.
- `client_secret_wo` (String, Sensitive, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) Where the client reads the service account's client secret. Required when the connection is created; left out afterwards, the stored secret is kept, including when another change replaces the connection. Write-only: never stored in state. Needs Terraform 1.11 or later. Change `client_secret_wo_version` to send a new value.
- `client_secret_wo_version` (Number) Any number. Change it to send `client_secret_wo` again; Terraform cannot see a write-only value change.
- `labels` (Map of String) Free key/value labels (`metadata.labels`). Stored with the object, never interpreted by the server.
- `name` (String) The object's name. WizConnection is one per tenant, so it is always `default`.

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
# Import an existing WizConnection by its name (metadata.name).
terraform import vulnetix_wiz_connection.default default
```
