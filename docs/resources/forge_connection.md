---
page_title: "vulnetix_forge_connection Resource - vulnetix"
subcategory: ""
description: |-
  How the Software Factory reaches a git forge other than github.com: a bot account's access token on Gitea, Forgejo or GitLab (or GitHub Enterprise Server), or a GitHub App registered on a GitHub Enterprise Server. github.com is reached through the Vulnetix GitHub App instead (GitHubInstallation). Apply asks the forge what it is and checks the credential can list repositories before anything is stored; the credential and the webhook secret are kept in the tenant's vault as a connector secret, never in the document, an export or a plan. The type, host and account identify the connection and cannot change: declare a new connection instead. Deleting the document removes the webhooks the factory created, revokes the stored credential and marks the connection deleted; managed scope and repository settings rows are kept, read-only, as for an uninstalled App.
  Manages the ForgeConnection config kind through /config/objects/ForgeConnection/{name}. Applies are recorded with the provider's manager, so a GitOps vulnetix config apply --prune leaves this object alone.
  Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (most of them, such as ApiSpecSource, RiskStrategy and VexSource) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object, or lets go of it for a kind the server never deletes (such as Organization). Some kinds wait on a person after they are applied (an invited Member, a Slack connection): status.phase says so, and Terraform does not wait.
---

# vulnetix_forge_connection (Resource)

How the Software Factory reaches a git forge other than github.com: a bot account's access token on Gitea, Forgejo or GitLab (or GitHub Enterprise Server), or a GitHub App registered on a GitHub Enterprise Server. github.com is reached through the Vulnetix GitHub App instead (GitHubInstallation). Apply asks the forge what it is and checks the credential can list repositories before anything is stored; the credential and the webhook secret are kept in the tenant's vault as a connector secret, never in the document, an export or a plan. The type, host and account identify the connection and cannot change: declare a new connection instead. Deleting the document removes the webhooks the factory created, revokes the stored credential and marks the connection deleted; managed scope and repository settings rows are kept, read-only, as for an uninstalled App.

Manages the `ForgeConnection` config kind through `/config/objects/ForgeConnection/{name}`. Applies are recorded with the provider's `manager`, so a GitOps `vulnetix config apply --prune` leaves this object alone.

Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (most of them, such as `ApiSpecSource`, `RiskStrategy` and `VexSource`) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object, or lets go of it for a kind the server never deletes (such as `Organization`). Some kinds wait on a person after they are applied (an invited `Member`, a Slack connection): `status.phase` says so, and Terraform does not wait.

## Example Usage

### Minimal

The arguments the kind requires, and nothing else.

```terraform
# Only the arguments the kind requires. From config-schemas examples/ForgeConnection/minimal.yaml.

resource "vulnetix_forge_connection" "codeberg" {
  name = "codeberg"
  host = "codeberg.org"
  type = "forgejo"
}
```

### Advanced

Every argument the kind supports. A secret is a write-only `*_wo` argument fed from an ephemeral variable, so it never reaches the plan or the state; change `*_wo_version` to send a new value.

```terraform
# Every argument the kind supports. From config-schemas examples/ForgeConnection/advanced.yaml.

variable "gitlab_platform_oauth_client_secret" {
  type      = string
  sensitive = true
  ephemeral = true
}
variable "gitlab_platform_token" {
  type      = string
  sensitive = true
  ephemeral = true
}
variable "ghes_acme_app_private_key" {
  type      = string
  sensitive = true
  ephemeral = true
}

resource "vulnetix_forge_connection" "gitlab_platform" {
  name     = "gitlab-platform"
  account  = "platform"
  api_base = "https://gitlab.example.com/api/v4"
  host     = "gitlab.example.com"
  oauth = {
    client_id                = "6f1e2d3c4b5a69788796a5b4c3d2e1f0"
    client_secret_wo         = var.gitlab_platform_oauth_client_secret
    client_secret_wo_version = 1
  }
  token_wo         = var.gitlab_platform_token
  token_wo_version = 1
  type             = "gitlab"
  labels = {
    team = "platform"
  }
  annotations = {
    note = "The bot account is vulnetix-bot, a Developer in the platform group."
  }
}

resource "vulnetix_forge_connection" "ghes_acme" {
  name    = "ghes-acme"
  account = "acme"
  app = {
    app_id                 = 12
    app_slug               = "vulnetix"
    installation_id        = 3
    private_key_wo         = var.ghes_acme_app_private_key
    private_key_wo_version = 1
  }
  host = "github.acme.example"
  type = "ghes"
}
```

### Reference

Every argument that has a schema default, set to that default, plus the required arguments. Planned against a tenant it only sets values the tenant left unset. Every kind's reference file is also in [Vulnetix/platform-iac-examples](https://github.com/Vulnetix/platform-iac-examples), with the same document as YAML.

```terraform
# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/ForgeConnection/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_forge_connection" "codeberg" {
  name = "codeberg"
  host = "codeberg.org"
  type = "forgejo"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `host` (String) The forge's host name, with a port when it is not 443: git.example.com, codeberg.org, gitlab.example.com:8443. No scheme or path. It is stored lower case.
- `name` (String) The object's name (`metadata.name`), unique per kind within the tenant. Changing it replaces the object.
- `type` (String) The forge software. Gitea and Forgejo speak one API, so the server keeps whichever of the two the host answers as, and a plan does not count the other spelling as a change. One of `ghes`, `gitea`, `forgejo`, `gitlab`.

### Optional

> **NOTE**: [Write-only arguments](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments) are supported in Terraform 1.11 and later.

- `account` (String) The organisation, group or user the connection is for. Left out, the connection covers every repository the credential reaches.
- `annotations` (Map of String) Free key/value notes for people and tools (`metadata.annotations`), such as a ticket link.
- `api_base` (String) The REST base, when it is not the forge's default (https://<host>/api/v4 for GitLab, /api/v1 for Gitea and Forgejo, /api/v3 for GitHub Enterprise Server). It must be https on the same host.
- `app` (Attributes) A GitHub App registered on the GitHub Enterprise Server, instead of a bot token. ghes only. (see [below for nested schema](#nestedatt--app))
- `labels` (Map of String) Free key/value labels (`metadata.labels`). Stored with the object, never interpreted by the server.
- `oauth` (Attributes) An OAuth application made on the forge, so that members can link their own forge account to the connection: their commits from a factory page then carry an email the forge verifies, and their signing key can be added to their forge account. Each member signs in for themselves from the console. Register the redirect URI the object's status reports (status.oauthRedirectUri) on the application. Left out, the application stored on the connection is removed. (see [below for nested schema](#nestedatt--oauth))
- `token_wo` (String, Sensitive, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) Where the client reads the bot account's access token. Required when the connection is created, unless app is set; left out afterwards, the stored token is kept. The token is checked against the forge at apply. Write-only: never stored in state. Needs Terraform 1.11 or later. Change `token_wo_version` to send a new value.
- `token_wo_version` (Number) Any number. Change it to send `token_wo` again; Terraform cannot see a write-only value change.

### Read-Only

- `id` (String) The object's name, the same as `name`.
- `status` (Attributes) What the server says about the object beyond its spec. (see [below for nested schema](#nestedatt--status))

<a id="nestedatt--app"></a>
### Nested Schema for `app`

Required:

- `app_id` (Number) The App's id on the server.
- `installation_id` (Number) The App's installation on the account.

Optional:

- `app_slug` (String) The App's slug, which names the bot that commits (<slug>[bot]).
- `private_key_wo` (String, Sensitive, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) Where the client reads the App's private key (PEM). Required when the connection is created; left out afterwards, the stored key is kept. Write-only: never stored in state. Needs Terraform 1.11 or later. Change `private_key_wo_version` to send a new value.
- `private_key_wo_version` (Number) Any number. Change it to send `private_key_wo` again; Terraform cannot see a write-only value change.


<a id="nestedatt--oauth"></a>
### Nested Schema for `oauth`

Required:

- `client_id` (String) The application's client id.

Optional:

- `client_secret_wo` (String, Sensitive, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) Where the client reads the application's client secret. Required when sign-in is first set up; left out afterwards, the stored secret is kept. Write-only: never stored in state. Needs Terraform 1.11 or later. Change `client_secret_wo_version` to send a new value.
- `client_secret_wo_version` (Number) Any number. Change it to send `client_secret_wo` again; Terraform cannot see a write-only value change.


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
# Import an existing ForgeConnection by its name (metadata.name).
terraform import vulnetix_forge_connection.codeberg codeberg
```
