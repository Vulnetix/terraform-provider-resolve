---
page_title: "vulnetix_aws_connection Resource - vulnetix"
subcategory: ""
description: |-
  How Vulnetix reaches the tenant's AWS estate. One per tenant. sso: IAM Identity Center, which can reach every account the operator's session can; the accounts are AwsAccount documents. iamUser: a long-lived access key, for an estate without Identity Center; the key's own account (from sts:GetCallerIdentity, never from the document) is the connected account. An sso connection starts a device authorization that a person approves in the AWS access portal, so its first apply, and any apply after the session expires, cannot finish unattended. Deleting it disconnects AWS; the stacks deployed in each account stay until removed there.
  Manages the AwsConnection config kind through /config/objects/AwsConnection/{name}. Applies are recorded with the provider's manager, so a GitOps vulnetix config apply --prune leaves this object alone.
  Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (most of them, such as ApiSpecSource, RiskStrategy and VexSource) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object, or lets go of it for a kind the server never deletes (such as Organization). Some kinds wait on a person after they are applied (an invited Member, a Slack connection): status.phase says so, and Terraform does not wait.
---

# vulnetix_aws_connection (Resource)

How Vulnetix reaches the tenant's AWS estate. One per tenant. sso: IAM Identity Center, which can reach every account the operator's session can; the accounts are AwsAccount documents. iamUser: a long-lived access key, for an estate without Identity Center; the key's own account (from sts:GetCallerIdentity, never from the document) is the connected account. An sso connection starts a device authorization that a person approves in the AWS access portal, so its first apply, and any apply after the session expires, cannot finish unattended. Deleting it disconnects AWS; the stacks deployed in each account stay until removed there.

Manages the `AwsConnection` config kind through `/config/objects/AwsConnection/{name}`. Applies are recorded with the provider's `manager`, so a GitOps `vulnetix config apply --prune` leaves this object alone.

Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (most of them, such as `ApiSpecSource`, `RiskStrategy` and `VexSource`) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object, or lets go of it for a kind the server never deletes (such as `Organization`). Some kinds wait on a person after they are applied (an invited `Member`, a Slack connection): `status.phase` says so, and Terraform does not wait.

## Example Usage

### Minimal

The arguments the kind requires, and nothing else.

```terraform
# Only the arguments the kind requires. From config-schemas examples/AwsConnection/minimal.yaml.

resource "vulnetix_aws_connection" "default" {
  region = "ap-southeast-2"
  sso = {
    start_url = "https://acme.awsapps.com/start"
  }
}
```

### Advanced

Every argument the kind supports. A secret is a write-only `*_wo` argument fed from an ephemeral variable, so it never reaches the plan or the state; change `*_wo_version` to send a new value.

```terraform
# Every argument the kind supports. From config-schemas examples/AwsConnection/advanced.yaml.

resource "vulnetix_aws_connection" "default" {
  advanced_networking = true
  region              = "ap-southeast-2"
  sso = {
    start_url = "https://acme.awsapps.com/start"
  }
  labels = {
    owner = "platform"
  }
  annotations = {
    note = "Approve the device authorization in the AWS access portal on first apply."
  }
}

# Or, instead (one per tenant, so declare only one):
#
# variable "default_2_iam_user_secret_access_key" {
#   type      = string
#   sensitive = true
#   ephemeral = true
# }
# resource "vulnetix_aws_connection" "default_2" {
#   advanced_networking = false
#   iam_user = {
#     access_key_id                = "AKIAIOSFODNN7EXAMPLE"
#     secret_access_key_wo         = var.default_2_iam_user_secret_access_key
#     secret_access_key_wo_version = 1
#   }
#   region = "us-east-1"
# }
```

### Reference

Every argument that has a schema default, set to that default, plus the required arguments. Planned against a tenant it only sets values the tenant left unset. Every kind's reference file is also in [Vulnetix/platform-iac-examples](https://github.com/Vulnetix/platform-iac-examples), with the same document as YAML.

```terraform
# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/AwsConnection/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_aws_connection" "default" {
  advanced_networking = false
  region              = "ap-southeast-2"
  sso = {
    start_url = "https://acme.awsapps.com/start"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `region` (String) The Identity Center region for sso; the region calls are made in for iamUser. An allowlist, because the region becomes part of an endpoint host name. One of `us-east-1`, `us-east-2`, `us-west-1`, `us-west-2`, `eu-west-1`, `eu-west-2`, `eu-west-3`, `eu-central-1`, `eu-central-2`, `eu-north-1`, `eu-south-1`, `eu-south-2`, `ap-southeast-1`, `ap-southeast-2`, `ap-southeast-3`, `ap-southeast-4`, `ap-northeast-1`, `ap-northeast-2`, `ap-northeast-3`, `ap-south-1`, `ap-south-2`, `ap-east-1`, `ca-central-1`, `ca-west-1`, `sa-east-1`, `me-south-1`, `me-central-1`, `af-south-1`, `il-central-1`.

### Optional

- `advanced_networking` (Boolean) Report over the private path from the customer's accounts to Vulnetix instead of the internet. Applies to every account; the next stack deploy fills the parameters. Schema default: `false`.
- `annotations` (Map of String) Free key/value notes for people and tools (`metadata.annotations`), such as a ticket link.
- `iam_user` (Attributes) (see [below for nested schema](#nestedatt--iam_user))
- `labels` (Map of String) Free key/value labels (`metadata.labels`). Stored with the object, never interpreted by the server.
- `name` (String) The object's name. AwsConnection is one per tenant, so it is always `default`.
- `sso` (Attributes) (see [below for nested schema](#nestedatt--sso))

### Read-Only

- `id` (String) The object's name, the same as `name`.
- `status` (Attributes) What the server says about the object beyond its spec. (see [below for nested schema](#nestedatt--status))

<a id="nestedatt--iam_user"></a>
### Nested Schema for `iam_user`

Required:

- `access_key_id` (String) An identifier, not a secret: it is in CloudTrail on every call it makes.

Optional:

- `secret_access_key_wo` (String, Sensitive, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) Where the client reads the secret access key. It is proved with sts:GetCallerIdentity, kept in the tenant's vault, and never exported or shown in a plan. Write-only: never stored in state. Needs Terraform 1.11 or later. Change `secret_access_key_wo_version` to send a new value. Required by the server when the object is created; leave it out afterwards to keep the value it holds.
- `secret_access_key_wo_version` (Number) Any number. Change it to send `secret_access_key_wo` again; Terraform cannot see a write-only value change.


<a id="nestedatt--sso"></a>
### Nested Schema for `sso`

Required:

- `start_url` (String) The AWS access portal URL. Stored as https://<name>.awsapps.com/start.


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
# Import an existing AwsConnection by its name (metadata.name).
terraform import vulnetix_aws_connection.default default
```
