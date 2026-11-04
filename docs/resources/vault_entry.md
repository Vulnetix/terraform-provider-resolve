---
page_title: "vulnetix_vault_entry Resource - vulnetix"
subcategory: ""
description: |-
  One secret in the tenant's vault: an AI provider key, an environment variable for Pix sandboxes and Belai hosts, or a provider management key, with the grants that say where it may be used. The value is encrypted with KMS before it is stored, is sent once by the client, and is never exported, shown in a plan or returned by any endpoint. No grant means no access, and inside one entry a deny beats an allow. Two active entries that fill the same slot (a provider for provider keys, a variable name for environment variables) may not cover the same subject: a plan that would create such an overlap is refused. Deleting the document revokes the entry: its value is destroyed and it covers nothing, and the vault keeps a stub with its history. Connector secrets (forge tokens, the Wiz client secret) are not VaultEntry documents: they are created through their integration kinds. Automatic rotation is set up in the console and an apply leaves it as it is.
  VaultEntry decides who can act for the tenant, so it can only be applied by an owner or admin of the organisation, and a role above the caller's own cannot be granted. A member is refused with HTTP 403.
  Manages the VaultEntry config kind through /config/objects/VaultEntry/{name}. Applies are recorded with the provider's manager, so a GitOps vulnetix config apply --prune leaves this object alone.
  Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (most of them, such as ApiSpecSource, RiskStrategy and VexSource) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object, or lets go of it for a kind the server never deletes (such as Organization). Some kinds wait on a person after they are applied (an invited Member, a Slack connection): status.phase says so, and Terraform does not wait.
---

# vulnetix_vault_entry (Resource)

One secret in the tenant's vault: an AI provider key, an environment variable for Pix sandboxes and Belai hosts, or a provider management key, with the grants that say where it may be used. The value is encrypted with KMS before it is stored, is sent once by the client, and is never exported, shown in a plan or returned by any endpoint. No grant means no access, and inside one entry a deny beats an allow. Two active entries that fill the same slot (a provider for provider keys, a variable name for environment variables) may not cover the same subject: a plan that would create such an overlap is refused. Deleting the document revokes the entry: its value is destroyed and it covers nothing, and the vault keeps a stub with its history. Connector secrets (forge tokens, the Wiz client secret) are not VaultEntry documents: they are created through their integration kinds. Automatic rotation is set up in the console and an apply leaves it as it is.

VaultEntry decides who can act for the tenant, so it can only be applied by an owner or admin of the organisation, and a role above the caller's own cannot be granted. A member is refused with HTTP 403.

Manages the `VaultEntry` config kind through `/config/objects/VaultEntry/{name}`. Applies are recorded with the provider's `manager`, so a GitOps `vulnetix config apply --prune` leaves this object alone.

Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (most of them, such as `ApiSpecSource`, `RiskStrategy` and `VexSource`) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object, or lets go of it for a kind the server never deletes (such as `Organization`). Some kinds wait on a person after they are applied (an invited `Member`, a Slack connection): `status.phase` says so, and Terraform does not wait.

## Example Usage

### Minimal

The arguments the kind requires, and nothing else.

```terraform
# Only the arguments the kind requires. From config-schemas examples/VaultEntry/minimal.yaml.

resource "vulnetix_vault_entry" "npm_token" {
  name       = "npm-token"
  entry_name = "NPM_TOKEN"
  type       = "env"
}
```

### Advanced

Every argument the kind supports. A secret is a write-only `*_wo` argument fed from an ephemeral variable, so it never reaches the plan or the state; change `*_wo_version` to send a new value.

```terraform
# Every argument the kind supports. From config-schemas examples/VaultEntry/advanced.yaml.

variable "anthropic_production_value" {
  type      = string
  sensitive = true
  ephemeral = true
}
variable "npm_token_value" {
  type      = string
  sensitive = true
  ephemeral = true
}
variable "openai_admin_value" {
  type      = string
  sensitive = true
  ephemeral = true
}

resource "vulnetix_vault_entry" "anthropic_production" {
  name       = "anthropic-production"
  entry_name = "Anthropic production"
  expires_at = "2027-06-30T00:00:00Z"
  grants = [
    {
      all     = true
      effect  = "allow"
      surface = "ai_firewall"
    },
    {
      effect  = "allow"
      ids     = ["3f1c2a9e-8b4d-4c7a-9e21-5d6f7a8b9c01", "7a2e4b6c-1d3f-4e5a-8b7c-9d0e1f2a3b4c"]
      surface = "pix_sandbox"
    },
  ]
  provider_slug    = "anthropic"
  type             = "provider_key"
  value_wo         = var.anthropic_production_value
  value_wo_version = 1
  labels = {
    team = "platform"
  }
  annotations = {
    ticket = "https://tracker.example.com/SEC-218"
  }
}

resource "vulnetix_vault_entry" "npm_token" {
  name       = "npm-token"
  entry_name = "NPM_TOKEN"
  grants = [
    {
      all     = true
      surface = "belai_host"
    },
    {
      effect  = "deny"
      ids     = ["0b1c2d3e-4f5a-4b6c-8d7e-9f0a1b2c3d4e"]
      surface = "belai_host"
    },
  ]
  type             = "env"
  value_wo         = var.npm_token_value
  value_wo_version = 1
}

resource "vulnetix_vault_entry" "openai_admin" {
  name             = "openai-admin"
  entry_name       = "OpenAI admin"
  provider_slug    = "openai"
  type             = "management"
  value_wo         = var.openai_admin_value
  value_wo_version = 1
}
```

### Reference

Every argument that has a schema default, set to that default, plus the required arguments. Planned against a tenant it only sets values the tenant left unset. Every kind's reference file is also in [Vulnetix/platform-iac-examples](https://github.com/Vulnetix/platform-iac-examples), with the same document as YAML.

```terraform
# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/VaultEntry/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_vault_entry" "npm_token" {
  name       = "npm-token"
  entry_name = "NPM_TOKEN"
  type       = "env"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `entry_name` (String) For env, the variable name: upper case letters, digits and underscore, starting with a letter or underscore. Names such as PATH, HOME and LD_PRELOAD, and names starting LD_, BELAI_, VULNETIX_, VX_, BASH_ or SSH_, are reserved and refused at plan. A variable's name cannot change: declare a new entry instead. For provider_key and management, the label the console shows: 1 to 64 letters, digits, spaces or . _ : / @ + -.
- `name` (String) The object's name (`metadata.name`), unique per kind within the tenant. Changing it replaces the object.
- `type` (String) provider_key: an AI provider's API key, handed to the AI Firewall gateway, Pix sandboxes and Belai hosts that hold a grant. env: an environment variable set in Pix sandboxes and on Belai hosts that hold a grant. management: a provider admin key that the vault uses to rotate provider keys; it is never granted to a surface. The type cannot change: declare a new entry instead. One of `provider_key`, `env`, `management`.

### Optional

> **NOTE**: [Write-only arguments](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments) are supported in Terraform 1.11 and later.

- `annotations` (Map of String) Free key/value notes for people and tools (`metadata.annotations`), such as a ticket link.
- `expires_at` (String) When the value stops being handed out. Left out, it never expires. A new expiry must be at least a minute and at most 10 years away; an expired entry is refused with vault_key_expired rather than decrypted.
- `grants` (Attributes Set) Where the value may be used. Each grant names a surface and either every subject on it or a list of subject ids. Left out or empty, the entry is stored but nothing may use it. At most 12 grants. (see [below for nested schema](#nestedatt--grants))
- `labels` (Map of String) Free key/value labels (`metadata.labels`). Stored with the object, never interpreted by the server.
- `provider_slug` (String) The AI Firewall provider slug the key is for, a built-in provider or one of the tenant's AiProvider objects. Required for provider_key and management, and not allowed for env. The provider must exist and be enabled when the plan runs. It cannot change: declare a new entry instead.
- `value_wo` (String, Sensitive, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) Where the client reads the secret value. Required when the entry is created; left out, the stored value is kept. Supplying the value the entry already holds is no change; a different value makes a new version and keeps the old one's fingerprint in the entry's history. An env value is at least 8 characters, so that it can be scrubbed from output reliably, and at most 4096; a key has no line breaks. Write-only: never stored in state. Needs Terraform 1.11 or later. Change `value_wo_version` to send a new value.
- `value_wo_version` (Number) Any number. Change it to send `value_wo` again; Terraform cannot see a write-only value change.

### Read-Only

- `id` (String) The object's name, the same as `name`.
- `status` (Attributes) What the server says about the object beyond its spec. (see [below for nested schema](#nestedatt--status))

<a id="nestedatt--grants"></a>
### Nested Schema for `grants`

Required:

- `surface` (String) pix_sandbox: Pix sandboxes, by their canonical sandbox id, which survives relaunches. belai_host: self-hosted Belai hosts, by host id. ai_firewall: AI Firewall callers, by the id of the principal that owns the Vulnetix API key making the request. One of `pix_sandbox`, `belai_host`, `ai_firewall`.

Optional:

- `all` (Boolean) true covers every subject on the surface, now and later. A grant is all or a list of ids, not both.
- `effect` (String) allow lets the subjects use the value; deny takes subjects away from an allow, and wins over it. One of `allow`, `deny`. Schema default: `allow`.
- `ids` (Set of String) The subjects covered, as UUIDs. 1 to 100 ids; duplicates are dropped and the list is kept sorted.


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
# Import an existing VaultEntry by its name (metadata.name).
terraform import vulnetix_vault_entry.npm_token npm-token
```
