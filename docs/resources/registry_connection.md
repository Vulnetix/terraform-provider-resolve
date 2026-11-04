---
page_title: "vulnetix_registry_connection Resource - vulnetix"
subcategory: ""
description: |-
  A private package or container registry that Vulnetix scans may read, and how to authenticate to it: an npm registry, a PyPI index, a Maven repository or an OCI registry. The credential is kept in the tenant's secrets vault (a connector entry), never in the connection; a scanner asks the platform to resolve the connection for a registry host and a package, and the credential is opened and the open recorded in the vault activity log. A connection with a username authenticates with basic authentication and the secret as the password; without one, the secret is a bearer token. Apply checks the credential against the registry before storing it and refuses a login the registry rejects. This is separate from the paid container registry logins on the Integrations page, which the crawler uses. Needs the saas migration 20261228020000_scan_registry_connection; until it is applied, plan refuses a document naming it.
  Manages the RegistryConnection config kind through /config/objects/RegistryConnection/{name}. Applies are recorded with the provider's manager, so a GitOps vulnetix config apply --prune leaves this object alone.
  Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (most of them, such as ApiSpecSource, RiskStrategy and VexSource) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object, or lets go of it for a kind the server never deletes (such as Organization). Some kinds wait on a person after they are applied (an invited Member, a Slack connection): status.phase says so, and Terraform does not wait.
---

# vulnetix_registry_connection (Resource)

A private package or container registry that Vulnetix scans may read, and how to authenticate to it: an npm registry, a PyPI index, a Maven repository or an OCI registry. The credential is kept in the tenant's secrets vault (a connector entry), never in the connection; a scanner asks the platform to resolve the connection for a registry host and a package, and the credential is opened and the open recorded in the vault activity log. A connection with a username authenticates with basic authentication and the secret as the password; without one, the secret is a bearer token. Apply checks the credential against the registry before storing it and refuses a login the registry rejects. This is separate from the paid container registry logins on the Integrations page, which the crawler uses. Needs the saas migration 20261228020000_scan_registry_connection; until it is applied, plan refuses a document naming it.

Manages the `RegistryConnection` config kind through `/config/objects/RegistryConnection/{name}`. Applies are recorded with the provider's `manager`, so a GitOps `vulnetix config apply --prune` leaves this object alone.

Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (most of them, such as `ApiSpecSource`, `RiskStrategy` and `VexSource`) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object, or lets go of it for a kind the server never deletes (such as `Organization`). Some kinds wait on a person after they are applied (an invited `Member`, a Slack connection): `status.phase` says so, and Terraform does not wait.

## Example Usage

### Minimal

The arguments the kind requires, and nothing else.

```terraform
# Only the arguments the kind requires. From config-schemas examples/RegistryConnection/minimal.yaml.

resource "vulnetix_registry_connection" "github_packages" {
  name      = "github-packages"
  ecosystem = "npm"
  url       = "https://npm.pkg.github.com/"
}
```

### Advanced

Every argument the kind supports. A secret is a write-only `*_wo` argument fed from an ephemeral variable, so it never reaches the plan or the state; change `*_wo_version` to send a new value.

```terraform
# Every argument the kind supports. From config-schemas examples/RegistryConnection/advanced.yaml.

variable "acme_artifactory_npm_secret" {
  type      = string
  sensitive = true
  ephemeral = true
}
variable "acme_registry_secret" {
  type      = string
  sensitive = true
  ephemeral = true
}

resource "vulnetix_registry_connection" "acme_artifactory_npm" {
  name              = "acme-artifactory-npm"
  ecosystem         = "npm"
  enabled           = true
  scope             = ["@acme", "@acme-internal"]
  secret_wo         = var.acme_artifactory_npm_secret
  secret_wo_version = 1
  url               = "https://artifactory.acme.example/artifactory/api/npm/npm-local/"
  username          = "vulnetix-scanner"
  labels = {
    team = "platform"
  }
  annotations = {
    ticket = "https://tracker.example.com/PLAT-482"
  }
}

resource "vulnetix_registry_connection" "acme_registry" {
  name              = "acme-registry"
  ecosystem         = "oci"
  scope             = ["acme/payments", "acme/platform"]
  secret_wo         = var.acme_registry_secret
  secret_wo_version = 1
  url               = "registry.acme.example"
}
```

### Reference

Every argument that has a schema default, set to that default, plus the required arguments. Planned against a tenant it only sets values the tenant left unset. Every kind's reference file is also in [Vulnetix/platform-iac-examples](https://github.com/Vulnetix/platform-iac-examples), with the same document as YAML.

```terraform
# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/RegistryConnection/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_registry_connection" "github_packages" {
  name      = "github-packages"
  ecosystem = "npm"
  enabled   = true
  scope     = []
  url       = "https://npm.pkg.github.com/"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `ecosystem` (String) oci: a container registry. npm: an npm registry. pypi: a PyPI simple index. maven: a Maven repository. One of `oci`, `npm`, `pypi`, `maven`.
- `name` (String) The object's name (`metadata.name`), unique per kind within the tenant. Changing it replaces the object.
- `url` (String) For oci, the registry host with an optional port (registry.example.com, ghcr.io, registry.example.com:5000). For the others, the https URL of the registry or index (https://npm.pkg.github.com/, https://pypi.example.com/simple/, https://maven.example.com/releases/). It is matched against the registry a scan reaches by host.

### Optional

> **NOTE**: [Write-only arguments](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments) are supported in Terraform 1.11 and later.

- `annotations` (Map of String) Free key/value notes for people and tools (`metadata.annotations`), such as a ticket link.
- `enabled` (Boolean) False keeps the connection and its credential but stops scans using it. Schema default: `true`.
- `labels` (Map of String) Free key/value labels (`metadata.labels`). Stored with the object, never interpreted by the server.
- `scope` (Set of String) The packages or repositories this connection serves. Left out, every one at the registry. npm: a scope (@acme). oci: a repository path prefix (acme/payments). maven: a group id prefix (com.acme). pypi: a project name prefix (acme-). When several connections serve one registry, the most specific scope wins. Schema default: `[]`.
- `secret_wo` (String, Sensitive, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) The password or token, read on the client and sent once. Required when the connection is created; left out afterwards, the stored credential is kept. It is stored in the secrets vault and never exported or shown in a plan. Changing it stores a new version of the vault entry. Write-only: never stored in state. Needs Terraform 1.11 or later. Change `secret_wo_version` to send a new value.
- `secret_wo_version` (Number) Any number. Change it to send `secret_wo` again; Terraform cannot see a write-only value change.
- `username` (String) The user for basic authentication, with the secret as the password. Left out, the secret is a bearer token (an npm automation token, a registry access token).

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
# Import an existing RegistryConnection by its name (metadata.name).
terraform import vulnetix_registry_connection.github_packages github-packages
```
