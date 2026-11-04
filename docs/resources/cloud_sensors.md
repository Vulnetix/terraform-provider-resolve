---
page_title: "vulnetix_cloud_sensors Resource - vulnetix"
subcategory: ""
description: |-
  Which AWS-native sources and Vulnetix passes are read, for one account or for every account. The sources listed are turned on; every other source is turned off for that scope. Once a source is on, everything it reports is ingested: granularity stops at the source. Turning one off stops reading it and keeps the findings already ingested, with their triage history. Each source's consent record is a named inline policy on the execution role, attached by redeploying the stack, so set deploy to attach it in the same apply. Turning a source on needs a connected AWS account. Not stored yet for one account: CloudNativeSource.accountId exists and is keyed, but no route writes an account id and the stack deploy reads the sources for the whole tenant, so only account "*" applies today. A document naming an account validates, and applying it needs that write path and a per-account read in the deploy.
  An object of this kind can reference AwsConnection, AwsAccount objects. Reference those resources' attributes, or add depends_on, so Terraform applies them first.
  Manages the CloudSensors config kind through /config/objects/CloudSensors/{name}. Applies are recorded with the provider's manager, so a GitOps vulnetix config apply --prune leaves this object alone.
  Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (most of them, such as ApiSpecSource, RiskStrategy and VexSource) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object, or lets go of it for a kind the server never deletes (such as Organization). Some kinds wait on a person after they are applied (an invited Member, a Slack connection): status.phase says so, and Terraform does not wait.
---

# vulnetix_cloud_sensors (Resource)

Which AWS-native sources and Vulnetix passes are read, for one account or for every account. The sources listed are turned on; every other source is turned off for that scope. Once a source is on, everything it reports is ingested: granularity stops at the source. Turning one off stops reading it and keeps the findings already ingested, with their triage history. Each source's consent record is a named inline policy on the execution role, attached by redeploying the stack, so set deploy to attach it in the same apply. Turning a source on needs a connected AWS account. Not stored yet for one account: CloudNativeSource.accountId exists and is keyed, but no route writes an account id and the stack deploy reads the sources for the whole tenant, so only account "*" applies today. A document naming an account validates, and applying it needs that write path and a per-account read in the deploy.

An object of this kind can reference `AwsConnection`, `AwsAccount` objects. Reference those resources' attributes, or add `depends_on`, so Terraform applies them first.

Manages the `CloudSensors` config kind through `/config/objects/CloudSensors/{name}`. Applies are recorded with the provider's `manager`, so a GitOps `vulnetix config apply --prune` leaves this object alone.

Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (most of them, such as `ApiSpecSource`, `RiskStrategy` and `VexSource`) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object, or lets go of it for a kind the server never deletes (such as `Organization`). Some kinds wait on a person after they are applied (an invited `Member`, a Slack connection): `status.phase` says so, and Terraform does not wait.

## Example Usage

### Minimal

The arguments the kind requires, and nothing else.

```terraform
# Only the arguments the kind requires. From config-schemas examples/CloudSensors/minimal.yaml.

resource "vulnetix_cloud_sensors" "estate" {
  name    = "estate"
  sources = ["inspector", "security-hub"]
}
```

### Advanced

Every argument the kind supports. A secret is a write-only `*_wo` argument fed from an ephemeral variable, so it never reaches the plan or the state; change `*_wo_version` to send a new value.

```terraform
# Every argument the kind supports. From config-schemas examples/CloudSensors/advanced.yaml.

resource "vulnetix_cloud_sensors" "estate" {
  name    = "estate"
  account = "*"
  deploy  = true
  sources = ["inspector", "security-hub", "guardduty", "config", "macie", "access-analyzer", "waf", "compute", "compute-artifacts", "network-logs", "eks-workloads"]
  labels = {
    scope = "all-accounts"
  }
  annotations = {
    note = "Only account \"*\" applies today; a per-account document validates but is not stored yet."
  }
}
```

### Reference

Every argument that has a schema default, set to that default, plus the required arguments. Planned against a tenant it only sets values the tenant left unset. Every kind's reference file is also in [Vulnetix/platform-iac-examples](https://github.com/Vulnetix/platform-iac-examples), with the same document as YAML.

```terraform
# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/CloudSensors/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_cloud_sensors" "estate" {
  name    = "estate"
  account = "*"
  deploy  = false
  sources = ["inspector", "security-hub"]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The object's name (`metadata.name`), unique per kind within the tenant. Changing it replaces the object.
- `sources` (Set of String) inspector: Amazon Inspector. security-hub: AWS Security Hub. guardduty, config, macie and access-analyzer: those products, read through Security Hub (one grant, filtered by product). waf: the AWS WAF configuration in front of the account. compute: SSM package inventory, EMR, EKS node groups, running ECS images, CIS AWS Foundations and CloudFormation template checks. compute-artifacts: read inside ECR image layers and Lambda code. network-logs: match Resolver and flow logs against malicious indicators. eks-workloads: list EKS pods and the images each runs; each cluster also needs an access entry you create. compute-artifacts, network-logs and eks-workloads have an effect only with compute on.

### Optional

- `account` (String) An AWS account id, or * for every account, including those connected later. One document per value. Schema default: `*`.
- `annotations` (Map of String) Free key/value notes for people and tools (`metadata.annotations`), such as a ticket link.
- `deploy` (Boolean) Redeploy the pipeline stack after a change, so the execution role carries exactly the policies the sources turned on need. Schema default: `false`.
- `labels` (Map of String) Free key/value labels (`metadata.labels`). Stored with the object, never interpreted by the server.

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
# Import an existing CloudSensors by its name (metadata.name).
terraform import vulnetix_cloud_sensors.estate estate
```
