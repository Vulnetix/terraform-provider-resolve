---
page_title: "vulnetix_log_source Resource - vulnetix"
subcategory: ""
description: |-
  A way logs reach Guardian, backed by a Guardian ingest token of the same type. The token is minted on create and shown once, by the apply that creates it; it is never exported or shown in a plan. aws: the token for the vulnetix-siem CloudFormation stack you deploy (IngestToken), which sends VPC flow logs, Route 53 Resolver logs and CloudTrail events by Firehose. cloudflare: Guardian creates Logpush jobs with your Cloudflare API token, which is used for those calls only and never stored. events: the token for POST /v1/ingest/events, used by your own scripts and the Vulnetix runtime sensor. Guardian has no update for a token or a Logpush job, so a changed spec is applied as delete and create: the old token is revoked and a new one minted. Managing sources needs the Guardian manager role.
  Manages the LogSource config kind through /config/objects/LogSource/{name}. Applies are recorded with the provider's manager, so a GitOps vulnetix config apply --prune leaves this object alone.
  Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (most of them, such as ApiSpecSource, RiskStrategy and VexSource) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object, or lets go of it for a kind the server never deletes (such as Organization). Some kinds wait on a person after they are applied (an invited Member, a Slack connection): status.phase says so, and Terraform does not wait.
---

# vulnetix_log_source (Resource)

A way logs reach Guardian, backed by a Guardian ingest token of the same type. The token is minted on create and shown once, by the apply that creates it; it is never exported or shown in a plan. aws: the token for the vulnetix-siem CloudFormation stack you deploy (IngestToken), which sends VPC flow logs, Route 53 Resolver logs and CloudTrail events by Firehose. cloudflare: Guardian creates Logpush jobs with your Cloudflare API token, which is used for those calls only and never stored. events: the token for POST /v1/ingest/events, used by your own scripts and the Vulnetix runtime sensor. Guardian has no update for a token or a Logpush job, so a changed spec is applied as delete and create: the old token is revoked and a new one minted. Managing sources needs the Guardian manager role.

Manages the `LogSource` config kind through `/config/objects/LogSource/{name}`. Applies are recorded with the provider's `manager`, so a GitOps `vulnetix config apply --prune` leaves this object alone.

Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (most of them, such as `ApiSpecSource`, `RiskStrategy` and `VexSource`) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object, or lets go of it for a kind the server never deletes (such as `Organization`). Some kinds wait on a person after they are applied (an invited `Member`, a Slack connection): `status.phase` says so, and Terraform does not wait.

## Example Usage

### Minimal

The arguments the kind requires, and nothing else.

```terraform
# Only the arguments the kind requires. From config-schemas examples/LogSource/minimal.yaml.

resource "vulnetix_log_source" "sensor_events" {
  name = "sensor-events"
  type = "events"
}
```

### Advanced

Every argument the kind supports. A secret is a write-only `*_wo` argument fed from an ephemeral variable, so it never reaches the plan or the state; change `*_wo_version` to send a new value.

```terraform
# Every argument the kind supports. From config-schemas examples/LogSource/advanced.yaml.

variable "cloudflare_edge_cloudflare_api_token" {
  type      = string
  sensitive = true
  ephemeral = true
}

resource "vulnetix_log_source" "cloudflare_edge" {
  name = "cloudflare-edge"
  cloudflare = {
    account_id           = "0123456789abcdef0123456789abcdef"
    api_token_wo         = var.cloudflare_edge_cloudflare_api_token
    api_token_wo_version = 1
    datasets             = ["firewall_events", "http_requests", "gateway_dns"]
    zone_ids             = ["fedcba9876543210fedcba9876543210"]
  }
  label = "Cloudflare edge and Gateway"
  type  = "cloudflare"
  labels = {
    team = "netsec"
  }
  annotations = {
    note = "The Cloudflare token is used for this apply only; deleting the source needs it again."
  }
}

resource "vulnetix_log_source" "aws_prod" {
  name  = "aws-prod"
  label = "Production VPCs"
  type  = "aws"
}

resource "vulnetix_log_source" "sensor_events" {
  name  = "sensor-events"
  label = "Runtime sensors"
  type  = "events"
}
```

### Reference

Every argument that has a schema default, set to that default, plus the required arguments. Planned against a tenant it only sets values the tenant left unset. Every kind's reference file is also in [Vulnetix/platform-iac-examples](https://github.com/Vulnetix/platform-iac-examples), with the same document as YAML.

```terraform
# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/LogSource/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_log_source" "sensor_events" {
  name = "sensor-events"
  type = "events"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The object's name (`metadata.name`), unique per kind within the tenant. Changing it replaces the object.
- `type` (String) One of `aws`, `cloudflare`, `events`.

### Optional

- `annotations` (Map of String) Free key/value notes for people and tools (`metadata.annotations`), such as a ticket link.
- `cloudflare` (Attributes) The Logpush jobs to create. Running it again does not make a second job for a dataset that already sends to Guardian. Cloudflare allows four Logpush jobs per zone. (see [below for nested schema](#nestedatt--cloudflare))
- `label` (String) The token's label on the Sources page. Defaults to the type and the date.
- `labels` (Map of String) Free key/value labels (`metadata.labels`). Stored with the object, never interpreted by the server.

### Read-Only

- `id` (String) The object's name, the same as `name`.
- `status` (Attributes) What the server says about the object beyond its spec. (see [below for nested schema](#nestedatt--status))

<a id="nestedatt--cloudflare"></a>
### Nested Schema for `cloudflare`

Required:

- `datasets` (Set of String) firewall_events, http_requests and dns_logs are zone datasets, created on every zone the token can see (or those in zoneIds). The others are account datasets and need accountId.

Optional:

- `account_id` (String) The Cloudflare account id, for account datasets.
- `api_token_wo` (String, Sensitive, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) Where the client reads the Cloudflare API token. Sent for this apply only; Guardian never stores it, so deleting the source needs it again. Write-only: never stored in state. Needs Terraform 1.11 or later. Change `api_token_wo_version` to send a new value. Required by the server when the object is created; leave it out afterwards to keep the value it holds.
- `api_token_wo_version` (Number) Any number. Change it to send `api_token_wo` again; Terraform cannot see a write-only value change.
- `zone_ids` (Set of String) Only these zones. Left out, every zone the token can see (at most 1,000).


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
# Import an existing LogSource by its name (metadata.name).
terraform import vulnetix_log_source.sensor_events sensor-events
```
