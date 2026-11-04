---
page_title: "vulnetix_factory_trigger Resource - vulnetix"
subcategory: ""
description: |-
  Something that starts a Software Factory workflow: a GitHub webhook event, a finding, an issue label, a kanban card move, a firewall sighting, or a person pressing Run. A trigger is a trigger node in the workflow's graph whose id is metadata.name (at most 40 characters), so applying one saves and publishes a new version of the workflow, the same as adding the trigger in the workflow builder; the workflow must be published and must not hold an unpublished draft. A trigger node the workflow already has under that id is adopted. A schedule is a FactorySchedule. Conditions are typed fields: the console takes no expression for a trigger. A trigger is on while its workflow is on. Removing the document removes the node and publishes the workflow again.
  An object of this kind can reference FactoryWorkflow objects. Reference those resources' attributes, or add depends_on, so Terraform applies them first.
  Manages the FactoryTrigger config kind through /config/objects/FactoryTrigger/{name}. Applies are recorded with the provider's manager, so a GitOps vulnetix config apply --prune leaves this object alone.
  Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (most of them, such as ApiSpecSource, RiskStrategy and VexSource) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object, or lets go of it for a kind the server never deletes (such as Organization). Some kinds wait on a person after they are applied (an invited Member, a Slack connection): status.phase says so, and Terraform does not wait.
---

# vulnetix_factory_trigger (Resource)

Something that starts a Software Factory workflow: a GitHub webhook event, a finding, an issue label, a kanban card move, a firewall sighting, or a person pressing Run. A trigger is a trigger node in the workflow's graph whose id is metadata.name (at most 40 characters), so applying one saves and publishes a new version of the workflow, the same as adding the trigger in the workflow builder; the workflow must be published and must not hold an unpublished draft. A trigger node the workflow already has under that id is adopted. A schedule is a FactorySchedule. Conditions are typed fields: the console takes no expression for a trigger. A trigger is on while its workflow is on. Removing the document removes the node and publishes the workflow again.

An object of this kind can reference `FactoryWorkflow` objects. Reference those resources' attributes, or add `depends_on`, so Terraform applies them first.

Manages the `FactoryTrigger` config kind through `/config/objects/FactoryTrigger/{name}`. Applies are recorded with the provider's `manager`, so a GitOps `vulnetix config apply --prune` leaves this object alone.

Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (most of them, such as `ApiSpecSource`, `RiskStrategy` and `VexSource`) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object, or lets go of it for a kind the server never deletes (such as `Organization`). Some kinds wait on a person after they are applied (an invited `Member`, a Slack connection): `status.phase` says so, and Terraform does not wait.

## Example Usage

### Minimal

The arguments the kind requires, and nothing else.

```terraform
# Only the arguments the kind requires. From config-schemas examples/FactoryTrigger/minimal.yaml.

resource "vulnetix_factory_trigger" "run_now" {
  name     = "run-now"
  start_at = ["scan"]
  type     = "manual"
  workflow = "security-patch-train"
}
```

### Advanced

Every argument the kind supports. A secret is a write-only `*_wo` argument fed from an ephemeral variable, so it never reaches the plan or the state; change `*_wo_version` to send a new value.

```terraform
# Every argument the kind supports. From config-schemas examples/FactoryTrigger/advanced.yaml.

resource "vulnetix_factory_trigger" "pr_opened" {
  name = "pr-opened"
  conditions = {
    actions     = ["opened", "synchronize", "reopened"]
    author_type = "user"
    authors = {
      allow = ["alice", "bob"]
      deny  = ["dependabot*"]
    }
    base_branch = ["@default", "release/*"]
    draft       = "exclude"
    head_branch = ["feature/**"]
    labels = {
      all = ["ready"]
      any = ["security"]
    }
    origin = "same_repo"
    properties = jsonencode({
      team = ["payments"]
    })
    subject = "pr"
  }
  debounce_seconds  = 60
  event             = "pull_request"
  max_runs_per_hour = 30
  overlap           = "cancel_previous"
  priority          = 2
  start_at          = ["scan"]
  title             = "Pull request opened or updated"
  type              = "webhook"
  workflow          = "security-patch-train"
  labels = {
    team = "appsec"
  }
}

resource "vulnetix_factory_trigger" "push_default_branch" {
  name = "push-default-branch"
  conditions = {
    base_branch  = ["@default"]
    paths        = ["package.json", "**/go.mod"]
    paths_ignore = ["docs/**"]
  }
  event    = "push"
  ref      = "default"
  start_at = ["scan"]
  type     = "webhook"
  workflow = "security-patch-train"
}

resource "vulnetix_factory_trigger" "ci_failed" {
  name = "ci-failed"
  conditions = {
    conclusion = ["failure", "timed_out"]
  }
  event    = "check_run.completed"
  start_at = ["scan"]
  type     = "webhook"
  workflow = "security-patch-train"
}

resource "vulnetix_factory_trigger" "dependabot_critical" {
  name = "dependabot-critical"
  conditions = {
    severity = ["critical", "high"]
  }
  event    = "dependabot_alert.created"
  start_at = ["scan"]
  type     = "webhook"
  workflow = "security-patch-train"
}

resource "vulnetix_factory_trigger" "production_deploy" {
  name = "production-deploy"
  conditions = {
    conclusion  = ["failure"]
    environment = ["production"]
  }
  event    = "deployment_status"
  start_at = ["scan"]
  type     = "webhook"
  workflow = "security-patch-train"
}

resource "vulnetix_factory_trigger" "image_published" {
  name = "image-published"
  conditions = {
    package_type = ["container"]
  }
  event    = "registry_package.published"
  start_at = ["scan"]
  type     = "webhook"
  workflow = "security-patch-train"
}

resource "vulnetix_factory_trigger" "changes_requested" {
  name = "changes-requested"
  conditions = {
    review_state = ["changes_requested"]
  }
  event    = "pull_request_review.submitted"
  start_at = ["scan"]
  type     = "webhook"
  workflow = "security-patch-train"
}

resource "vulnetix_factory_trigger" "security_label" {
  name              = "security-label"
  debounce_seconds  = 30
  label             = "belai:security"
  max_runs_per_hour = 20
  start_at          = ["scan"]
  type              = "label"
  workflow          = "security-patch-train"
}

resource "vulnetix_factory_trigger" "critical_finding" {
  name     = "critical-finding"
  severity = ["critical", "high"]
  start_at = ["scan"]
  type     = "finding"
  workflow = "security-patch-train"
}

resource "vulnetix_factory_trigger" "card_to_ready" {
  name      = "card-to-ready"
  from_list = "backlog"
  start_at  = ["scan"]
  to_list   = "ready"
  type      = "kanban"
  workflow  = "security-patch-train"
}

resource "vulnetix_factory_trigger" "new_package" {
  name              = "new-package"
  decision          = ["block"]
  ecosystems        = ["npm", "pypi"]
  events            = ["package.first_seen", "package.blocked"]
  max_runs_per_hour = 10
  new_to            = "org"
  source            = "package-firewall"
  start_at          = ["scan"]
  ticket_repo       = "acme/security-reviews"
  type              = "firewall"
  workflow          = "security-patch-train"
}

resource "vulnetix_factory_trigger" "new_mcp_server" {
  name     = "new-mcp-server"
  event    = "ai.capability_first_seen"
  kinds    = ["mcp_server", "skill"]
  source   = "ai-firewall"
  start_at = ["scan"]
  type     = "firewall"
  workflow = "security-patch-train"
}

resource "vulnetix_factory_trigger" "sandbox_egress" {
  name      = "sandbox-egress"
  event     = "egress.host_first_seen"
  sandboxes = ["build-sandbox"]
  source    = "egress"
  start_at  = ["scan"]
  type      = "firewall"
  workflow  = "security-patch-train"
}

resource "vulnetix_factory_trigger" "guardian_dns" {
  name       = "guardian-dns"
  categories = ["dns"]
  event      = "siem.fired"
  providers  = ["aws"]
  source     = "siem"
  start_at   = ["scan"]
  tables     = ["dns_events"]
  triggers   = ["suspicious-dns"]
  type       = "firewall"
  workflow   = "security-patch-train"
}
```

### Reference

Every argument that has a schema default, set to that default, plus the required arguments. Planned against a tenant it only sets values the tenant left unset. Every kind's reference file is also in [Vulnetix/platform-iac-examples](https://github.com/Vulnetix/platform-iac-examples), with the same document as YAML.

```terraform
# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/FactoryTrigger/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_factory_trigger" "run_now" {
  name     = "run-now"
  start_at = ["scan"]
  type     = "manual"
  workflow = "security-patch-train"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The object's name (`metadata.name`), unique per kind within the tenant. Changing it replaces the object.
- `start_at` (Set of String) The ids of the steps in the workflow's graph the trigger leads to: one edge from the trigger to each, in this order.
- `type` (String) What starts the workflow. webhook: a GitHub delivery to the App (event, conditions). finding: a new finding (severity). label: a label added to an issue or pull request (label). kanban: a card moved on the board (fromList, toList). firewall: a sighting from a firewall, Guardian or the GitHub risk rules (source and its filters). manual: only a person pressing Run. One of `webhook`, `finding`, `label`, `kanban`, `manual`, `firewall`.
- `workflow` (String) The metadata.name of a FactoryWorkflow: the workflow the trigger starts.

### Optional

- `annotations` (Map of String) Free key/value notes for people and tools (`metadata.annotations`), such as a ticket link.
- `categories` (Set of String) firewall, siem source: the Guardian categories.
- `conditions` (Attributes) webhook: conditions the delivery must meet, all of them. A condition the delivery carries no data for does not match. Globs: * is any run of characters except /, ** includes /, ? is one character; @default is the repository's default branch. (see [below for nested schema](#nestedatt--conditions))
- `debounce_seconds` (Number) webhook and label: deliveries for the same subject inside this window start one run. 0 to 3600; 30 when left out.
- `decision` (Set of String) firewall: the firewall decisions that start the workflow, such as block or hold.
- `ecosystems` (Set of String) firewall, package-firewall source: the package ecosystems, such as npm and pypi.
- `event` (String) webhook: the GitHub event, such as push or pull_request, or one action of it, such as pull_request.opened. firewall: one event of the source, such as package.first_seen (or list several in events).
- `events` (Set of String) firewall: the source's events that start the workflow, such as package.first_seen and package.blocked. Every event of the source when neither this nor event is set.
- `from_list` (String) kanban: the board list a card leaves.
- `kinds` (Set of String) firewall, ai-firewall source: the capability kinds.
- `label` (String) label: the issue or pull request label that starts the workflow.
- `labels` (Map of String) Free key/value labels (`metadata.labels`). Stored with the object, never interpreted by the server.
- `max_runs_per_hour` (Number) webhook, label and firewall: the most runs the trigger starts in an hour, 1 to 1000; 60 when left out.
- `new_to` (String) firewall: only a first sighting in the organisation (org), or, for packages, in the repository (repo). One of `org`, `repo`.
- `overlap` (String) What happens when a run of the workflow for the same subject is already going. Defaults to skip_if_running. One of `skip_if_running`, `queue`, `cancel_previous`.
- `priority` (Number) The run's priority, 1 (highest) to 5, when runs wait for capacity.
- `providers` (Set of String) firewall, siem source: the log providers.
- `ref` (String) webhook push and release events: default for branch pushes, tags for tag pushes and releases. One of `default`, `tags`.
- `sandboxes` (Set of String) firewall, egress source: the Pix Sandboxes, by name.
- `severity` (Set of String) finding: the finding severities that start the workflow. All of them when left out.
- `source` (String) firewall: where the sighting comes from: the package firewall, the AI firewall, Pix Sandbox egress, Guardian (siem), or the GitHub risk and posture rules. One of `package-firewall`, `ai-firewall`, `egress`, `siem`, `github-risk`, `github-posture`.
- `tables` (Set of String) firewall, siem source: the Guardian tables.
- `ticket_repo` (String) firewall: the repository, owner/name, where review tickets for a sighting are opened.
- `title` (String) The trigger node's title in the builder. Left out, the node is titled with metadata.name.
- `to_list` (String) kanban: the board list a card enters.
- `triggers` (Set of String) firewall, siem source: the Guardian alert triggers by name.

### Read-Only

- `id` (String) The object's name, the same as `name`.
- `status` (Attributes) What the server says about the object beyond its spec. (see [below for nested schema](#nestedatt--status))

<a id="nestedatt--conditions"></a>
### Nested Schema for `conditions`

Optional:

- `actions` (Set of String) The event's actions that count, such as opened and synchronize for pull_request. Only when event names no action.
- `author_type` (String) Who made the change: a person (user), a bot, or any. One of `any`, `user`, `bot`.
- `authors` (Attributes) GitHub logins, as globs. (see [below for nested schema](#nestedatt--conditions--authors))
- `base_branch` (Set of String) Branch globs the pull request targets, or the pushed branch for a push.
- `conclusion` (Set of String) check_run, check_suite and workflow_run: the conclusions that count, such as failure. deployment_status: the states that count.
- `draft` (String) Draft pull requests: include them, exclude them, or only them. One of `include`, `exclude`, `only`.
- `environment` (Set of String) deployment_status: environment globs.
- `head_branch` (Set of String) Branch globs of the pull request's head, or the pushed branch for a push.
- `labels` (Attributes) Labels on the issue or pull request. (see [below for nested schema](#nestedatt--conditions--labels))
- `origin` (String) Pull requests from the same repository, from forks, or any. One of `any`, `same_repo`, `fork`.
- `package_type` (Set of String) registry_package: the package types that count, such as container and npm.
- `paths` (Set of String) push only: path globs, at least one of which a pushed commit changed.
- `paths_ignore` (Set of String) push only: path globs that do not count as a change.
- `properties` (String) Repository custom properties: each property named with the values it may hold, such as team: [payments]. A JSON-encoded value; use jsonencode().
- `review_state` (Set of String) pull_request_review: the review states that count.
- `severity` (Set of String) dependabot_alert, code_scanning_alert and posture_finding: the alert severities that count.
- `subject` (String) What the delivery must be about: a pull request (pr), an issue, nothing in particular (none), or any. One of `any`, `pr`, `issue`, `none`.

<a id="nestedatt--conditions--authors"></a>
### Nested Schema for `conditions.authors`

Optional:

- `allow` (Set of String) Only these authors.
- `deny` (Set of String) Never these authors.


<a id="nestedatt--conditions--labels"></a>
### Nested Schema for `conditions.labels`

Optional:

- `all` (Set of String) Every one of these labels.
- `any` (Set of String) At least one of these labels.



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
# Import an existing FactoryTrigger by its name (metadata.name).
terraform import vulnetix_factory_trigger.run_now run-now
```
