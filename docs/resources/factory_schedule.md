---
page_title: "vulnetix_factory_schedule Resource - vulnetix"
subcategory: ""
description: |-
  A cron schedule that starts a Software Factory workflow: a cron trigger node in the workflow's graph whose id is metadata.name (at most 40 characters). Applying one saves and publishes a new version of the workflow, the same as adding the schedule on the Schedules tab, and creates or updates the timer that fires it; the workflow must be published and must not hold an unpublished draft. A cron trigger node the workflow already has under that id is adopted. Removing the document removes the node and its timer and publishes the workflow again.
  An object of this kind can reference FactoryWorkflow objects. Reference those resources' attributes, or add depends_on, so Terraform applies them first.
  Manages the FactorySchedule config kind through /config/objects/FactorySchedule/{name}. Applies are recorded with the provider's manager, so a GitOps vulnetix config apply --prune leaves this object alone.
  Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (most of them, such as ApiSpecSource, RiskStrategy and VexSource) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object, or lets go of it for a kind the server never deletes (such as Organization). Some kinds wait on a person after they are applied (an invited Member, a Slack connection): status.phase says so, and Terraform does not wait.
---

# vulnetix_factory_schedule (Resource)

A cron schedule that starts a Software Factory workflow: a cron trigger node in the workflow's graph whose id is metadata.name (at most 40 characters). Applying one saves and publishes a new version of the workflow, the same as adding the schedule on the Schedules tab, and creates or updates the timer that fires it; the workflow must be published and must not hold an unpublished draft. A cron trigger node the workflow already has under that id is adopted. Removing the document removes the node and its timer and publishes the workflow again.

An object of this kind can reference `FactoryWorkflow` objects. Reference those resources' attributes, or add `depends_on`, so Terraform applies them first.

Manages the `FactorySchedule` config kind through `/config/objects/FactorySchedule/{name}`. Applies are recorded with the provider's `manager`, so a GitOps `vulnetix config apply --prune` leaves this object alone.

Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (most of them, such as `ApiSpecSource`, `RiskStrategy` and `VexSource`) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object, or lets go of it for a kind the server never deletes (such as `Organization`). Some kinds wait on a person after they are applied (an invited `Member`, a Slack connection): `status.phase` says so, and Terraform does not wait.

## Example Usage

### Minimal

The arguments the kind requires, and nothing else.

```terraform
# Only the arguments the kind requires. From config-schemas examples/FactorySchedule/minimal.yaml.

resource "vulnetix_factory_schedule" "nightly" {
  name     = "nightly"
  cron     = "0 3 * * *"
  start_at = ["scan"]
  workflow = "security-patch-train"
}
```

### Advanced

Every argument the kind supports. A secret is a write-only `*_wo` argument fed from an ephemeral variable, so it never reaches the plan or the state; change `*_wo_version` to send a new value.

```terraform
# Every argument the kind supports. From config-schemas examples/FactorySchedule/advanced.yaml.

resource "vulnetix_factory_schedule" "weekday_sweep" {
  name     = "weekday-sweep"
  cron     = "0 3 * * 1-5"
  enabled  = true
  overlap  = "skip_if_running"
  priority = 3
  runs_on  = "aws"
  start_at = ["scan"]
  timezone = "Australia/Melbourne"
  title    = "Weekday dependency sweep"
  workflow = "security-patch-train"
  annotations = {
    ticket = "https://tracker.acme.com/SEC-142"
  }
}
```

### Reference

Every argument that has a schema default, set to that default, plus the required arguments. Planned against a tenant it only sets values the tenant left unset. Every kind's reference file is also in [Vulnetix/platform-iac-examples](https://github.com/Vulnetix/platform-iac-examples), with the same document as YAML.

```terraform
# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/FactorySchedule/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_factory_schedule" "nightly" {
  name     = "nightly"
  cron     = "0 3 * * *"
  enabled  = true
  runs_on  = "aws"
  start_at = ["scan"]
  timezone = "UTC"
  workflow = "security-patch-train"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `cron` (String) When the workflow runs: five cron fields (minute, hour, day of month, month, day of week) of numbers, ranges, lists and steps. Day of month and day of week cannot both be restricted.
- `name` (String) The object's name (`metadata.name`), unique per kind within the tenant. Changing it replaces the object.
- `start_at` (Set of String) The ids of the steps in the workflow's graph the schedule leads to: one edge from the trigger to each, in this order.
- `workflow` (String) The metadata.name of a FactoryWorkflow: the workflow the schedule starts.

### Optional

- `annotations` (Map of String) Free key/value notes for people and tools (`metadata.annotations`), such as a ticket link.
- `enabled` (Boolean) Whether the schedule fires. A schedule of a workflow that is switched off fires but starts nothing. Schema default: `true`.
- `labels` (Map of String) Free key/value labels (`metadata.labels`). Stored with the object, never interpreted by the server.
- `overlap` (String) What happens when the previous run is still going at the next fire. Defaults to skip_if_running. One of `skip_if_running`, `queue`, `cancel_previous`.
- `priority` (Number) The run's priority, 1 (highest) to 5, when runs wait for capacity.
- `runs_on` (String) The timer that fires the schedule: aws (an EventBridge Scheduler schedule) or cloudflare (a Worker cron trigger, for workflows that run in a Pix Sandbox). Defaults to aws. One of `aws`, `cloudflare`. Schema default: `aws`.
- `timezone` (String) The IANA time zone the cron fields are read in, such as Australia/Melbourne. Defaults to UTC. Schema default: `UTC`.
- `title` (String) The trigger node's title in the builder. Left out, the node is titled with metadata.name.

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
# Import an existing FactorySchedule by its name (metadata.name).
terraform import vulnetix_factory_schedule.nightly nightly
```
