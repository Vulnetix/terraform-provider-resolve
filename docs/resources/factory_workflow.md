---
page_title: "vulnetix_factory_workflow Resource - vulnetix"
subcategory: ""
description: |-
  A Software Factory workflow: the graph of steps a run walks, from a trigger through scans, decisions, crews, gates and people to its outputs. metadata.name is the workflow's slug (at most 60 characters; slugs that start with sys- belong to the default workflows and cannot be declared). Workflows are versioned: a change to the graph never edits a published version, it saves a new version and, with publish true, publishes it. Triggers and schedules declared as FactoryTrigger and FactorySchedule documents add their own trigger nodes to the graph; those nodes are not part of this document and a graph change keeps them. Removing the document archives the workflow: it is switched off and its triggers removed, and its versions and run history are kept.
  Manages the FactoryWorkflow config kind through /config/objects/FactoryWorkflow/{name}. Applies are recorded with the provider's manager, so a GitOps vulnetix config apply --prune leaves this object alone.
  Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (most of them, such as ApiSpecSource, RiskStrategy and VexSource) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object, or lets go of it for a kind the server never deletes (such as Organization). Some kinds wait on a person after they are applied (an invited Member, a Slack connection): status.phase says so, and Terraform does not wait.
---

# vulnetix_factory_workflow (Resource)

A Software Factory workflow: the graph of steps a run walks, from a trigger through scans, decisions, crews, gates and people to its outputs. metadata.name is the workflow's slug (at most 60 characters; slugs that start with sys- belong to the default workflows and cannot be declared). Workflows are versioned: a change to the graph never edits a published version, it saves a new version and, with publish true, publishes it. Triggers and schedules declared as FactoryTrigger and FactorySchedule documents add their own trigger nodes to the graph; those nodes are not part of this document and a graph change keeps them. Removing the document archives the workflow: it is switched off and its triggers removed, and its versions and run history are kept.

Manages the `FactoryWorkflow` config kind through `/config/objects/FactoryWorkflow/{name}`. Applies are recorded with the provider's `manager`, so a GitOps `vulnetix config apply --prune` leaves this object alone.

Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (most of them, such as `ApiSpecSource`, `RiskStrategy` and `VexSource`) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object, or lets go of it for a kind the server never deletes (such as `Organization`). Some kinds wait on a person after they are applied (an invited `Member`, a Slack connection): `status.phase` says so, and Terraform does not wait.

## Example Usage

### Minimal

The arguments the kind requires, and nothing else.

```terraform
# Only the arguments the kind requires. From config-schemas examples/FactoryWorkflow/minimal.yaml.

resource "vulnetix_factory_workflow" "patch_on_demand" {
  name = "patch-on-demand"
  graph = {
    nodes = [
      {
        config = jsonencode({
          kind = "manual"
        })
        id    = "start"
        title = "Run by hand"
        type  = "trigger"
      },
    ]
  }
}
```

### Advanced

Every argument the kind supports. A secret is a write-only `*_wo` argument fed from an ephemeral variable, so it never reaches the plan or the state; change `*_wo_version` to send a new value.

```terraform
# Every argument the kind supports. From config-schemas examples/FactoryWorkflow/advanced.yaml.

resource "vulnetix_factory_workflow" "security_patch_train" {
  name         = "security-patch-train"
  description  = "Scan, patch with belai:security, hold the change to CI and the quality gate, then open a draft pull request."
  display_name = "Security patch train"
  enabled      = true
  graph = {
    edges = [
      {
        from = "start"
        to   = "scan"
      },
      {
        from  = "scan"
        label = "findings"
        to    = "patch"
      },
      {
        from  = "patch"
        label = "branch"
        to    = "gate"
      },
      {
        from    = "gate"
        outcome = "pass"
        to      = "pr"
      },
      {
        from    = "gate"
        loop    = true
        outcome = "fail"
        to      = "patch"
      },
    ]
    limits = {
      agent_cost_usd          = 5
      agent_finalize_at       = 0.8
      agent_max_attempts      = 3
      agent_max_continuations = 2
      agent_max_iterations    = 80
      agent_max_passes        = 20
      agent_max_wall_minutes  = 120
      concurrency             = 2
      github_calls_per_run = {
        max     = 8
        typical = 4
      }
      month_budget_tokens    = 6000000
      run_budget_tokens      = 400000
      stop_the_line_attempts = 3
    }
    nodes = [
      {
        config = jsonencode({
          kind = "manual"
        })
        id = "start"
        position = {
          x = 0
          y = 0
        }
        title = "Run by hand"
        type  = "trigger"
      },
      {
        config = jsonencode({
          failOn     = ["malware", "secrets"]
          sbomFormat = "cyclonedx-1.6"
          scanners   = ["sca", "malscan", "secrets"]
          source     = "posture"
        })
        github_calls = {
          max = 1
        }
        id           = "scan"
        on_quota_low = "use_cached"
        position = {
          x = 240
          y = 0
        }
        title = "Scanners"
        type  = "scan"
      },
      {
        budget = {
          tokens = 300000
        }
        config = jsonencode({
          codeTier     = "smart"
          crew         = "belai:security"
          maxAttempts  = 3
          mode         = "goal"
          testsOnFail  = "fix"
          testsPostEnd = "goal"
        })
        id = "patch"
        position = {
          x = 480
          y = 0
        }
        title = "Patch and verify"
        type  = "crew"
      },
      {
        config = jsonencode({
          checkSource    = "mirror"
          maxGithubCalls = 2
          onFail         = "return_to_crew"
          onQuotaLow     = "use_cached"
          profile        = "org_default"
          releaseJail    = "enforce"
        })
        github_calls = {
          max = 2
        }
        id = "gate"
        position = {
          x = 720
          y = 0
        }
        title = "Quality gate and tests"
        type  = "gate"
      },
      {
        config = jsonencode({
          attach    = "cdx_slsa"
          kind      = "draft_pr"
          prState   = "draft"
          reviewers = "codeowners"
        })
        github_calls = {
          max = 2
        }
        id = "pr"
        position = {
          x = 960
          y = 0
        }
        title = "Draft pull request"
        type  = "output"
      },
    ]
    scope = {
      orgs  = ["acme"]
      repos = ["acme/payments", "acme-labs/*"]
    }
  }
  publish = true
  labels = {
    team = "appsec"
  }
  annotations = {
    owner = "appsec@acme.com"
  }
}
```

### Reference

Every argument that has a schema default, set to that default, plus the required arguments. Planned against a tenant it only sets values the tenant left unset. Every kind's reference file is also in [Vulnetix/platform-iac-examples](https://github.com/Vulnetix/platform-iac-examples), with the same document as YAML.

```terraform
# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/FactoryWorkflow/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_factory_workflow" "patch_on_demand" {
  name    = "patch-on-demand"
  enabled = true
  graph = {
    nodes = [
      {
        config = jsonencode({
          kind = "manual"
        })
        id    = "start"
        title = "Run by hand"
        type  = "trigger"
      },
    ]
  }
  publish = true
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `graph` (Attributes) The workflow graph, in the shape the console's workflow builder saves and the workflow API takes. (see [below for nested schema](#nestedatt--graph))
- `name` (String) The object's name (`metadata.name`), unique per kind within the tenant. Changing it replaces the object.

### Optional

- `annotations` (Map of String) Free key/value notes for people and tools (`metadata.annotations`), such as a ticket link.
- `description` (String) What the workflow is for, shown under its name in the console.
- `display_name` (String) The name the console shows. Left out, the workflow is shown under its metadata.name.
- `enabled` (Boolean) Whether the factory starts runs of the workflow. Switching a workflow off also switches off its event triggers; a run already going is left to finish. A workflow that has never been published cannot be switched on. Schema default: `true`.
- `labels` (Map of String) Free key/value labels (`metadata.labels`). Stored with the object, never interpreted by the server.
- `publish` (Boolean) true, or left out: the graph is the workflow's published version, and new runs use it; a changed graph is saved as the next version and published, which needs a graph that passes the full publish checks (a trigger, every step reachable, no cycle, every decision outcome routed, and every step but config-apply configured, not left with an empty config). false: the graph is kept as the workflow's draft and runs keep using the last published version. Apply refuses to publish over a draft someone left unpublished in the console unless the draft already holds this graph. A draft may hold steps that are not configured yet (config: {}). Schema default: `true`.

### Read-Only

- `id` (String) The object's name, the same as `name`.
- `status` (Attributes) What the server says about the object beyond its spec. (see [below for nested schema](#nestedatt--status))

<a id="nestedatt--graph"></a>
### Nested Schema for `graph`

Required:

- `nodes` (Attributes List) The steps, at most 60. (see [below for nested schema](#nestedatt--graph--nodes))

Optional:

- `edges` (Attributes List) The connections between steps, at most 120. A run moves along an edge when the step it leaves finishes. (see [below for nested schema](#nestedatt--graph--edges))
- `limits` (Attributes) Run limits. A value left out, or 0, keeps the tenant's setting. (see [below for nested schema](#nestedatt--graph--limits))
- `scope` (Attributes) The repositories the workflow may run on. (see [below for nested schema](#nestedatt--graph--scope))

<a id="nestedatt--graph--nodes"></a>
### Nested Schema for `graph.nodes`

Required:

- `id` (String) The step's id, unique in the graph: 1 to 40 letters, digits, dashes or underscores. Edges name steps by it.
- `title` (String) The step's title in the builder and the run view. Publishing needs a title that is not empty on every step.
- `type` (String) What the step does. trigger: what starts a run. scan: run scanners or read their results. decision: ask a model a question and route on the answer. crew: hand the work to a Belai crew or agent. gate: hold the run to CI checks, the quality gate and the release jail, and ask a person to approve a config-apply plan. human: ask a person to approve. output: open a pull request, comment, publish to TEA, notify or write VEX. state: read or write a stored object. config-apply: plan, or apply an approved plan of, the Platform IaC documents in the repository's .vulnetix/config/ directory at the pushed commit (config: phase plan or apply, and prune). One of `trigger`, `scan`, `decision`, `crew`, `gate`, `human`, `output`, `state`, `config-apply`.

Optional:

- `budget` (Attributes) The step's token budget. (see [below for nested schema](#nestedatt--graph--nodes--budget))
- `config` (String) The step's settings, which depend on its type. A published workflow needs at least one on every step except config-apply, such as a trigger's kind or a gate's profile; a draft (publish: false) may leave them empty (config: {}), and so may a config-apply step, whose settings all have defaults. They are the settings the workflow builder edits, under the same names; the server checks them with the builder's own rules, so a step the builder would refuse is refused here too. A trigger step's config is the same as a FactoryTrigger or FactorySchedule spec, with the type under kind (kind: cron for a schedule, with cron, tz and runsOn). A JSON-encoded value; use jsonencode().
- `github_calls` (Attributes) The step's GitHub API call limit. (see [below for nested schema](#nestedatt--graph--nodes--github_calls))
- `on_quota_low` (String) What the step does when the GitHub quota is low: carry on with cached data, wait for the quota window, or ask a person. One of `use_cached`, `defer`, `escalate`.
- `position` (Attributes) Where the builder draws the step. It has no effect on runs. (see [below for nested schema](#nestedatt--graph--nodes--position))

<a id="nestedatt--graph--nodes--budget"></a>
### Nested Schema for `graph.nodes.budget`

Optional:

- `tokens` (Number) Tokens the step may spend in one run; 0 for no limit of its own.


<a id="nestedatt--graph--nodes--github_calls"></a>
### Nested Schema for `graph.nodes.github_calls`

Optional:

- `max` (Number) GitHub API calls the step may make in one run.


<a id="nestedatt--graph--nodes--position"></a>
### Nested Schema for `graph.nodes.position`

Optional:

- `x` (Number)
- `y` (Number)



<a id="nestedatt--graph--edges"></a>
### Nested Schema for `graph.edges`

Required:

- `from` (String) The id of the step the edge leaves.
- `to` (String) The id of the step the edge enters. Nothing may lead into a trigger.

Optional:

- `label` (String) The artifact that travels along the edge, shown on the canvas.
- `loop` (Boolean) Marks the one kind of edge allowed to go back: a gate or decision returning work to a crew, bounded by limits.stopTheLineAttempts.
- `outcome` (String) On an edge leaving a decision: the outcome that takes it, or $below_threshold or $unknown. On an edge leaving a gate: pass, or fail for the loop back to a crew.


<a id="nestedatt--graph--limits"></a>
### Nested Schema for `graph.limits`

Optional:

- `agent_cost_usd` (Number) The cost limit, in US dollars, of each card a worker this workflow starts takes on, up to 10000.
- `agent_finalize_at` (Number) The share of a limit (0 to 0.95) at which a worker starts wrapping up.
- `agent_max_attempts` (Number) Failed attempts a worker may make on one card, up to 10.
- `agent_max_continuations` (Number) How often a card may resume after its budget ran out, up to 50.
- `agent_max_iterations` (Number) Tool rounds a worker may take in one pass, up to 400.
- `agent_max_passes` (Number) Passes a worker may make on one card, up to 1000.
- `agent_max_wall_minutes` (Number) Minutes a worker may spend on one card, up to 10080 (a week).
- `concurrency` (Number) Runs of the workflow at once, 1 to 50.
- `github_calls_per_run` (Attributes) The GitHub API calls a run is expected to make, for the estimate the builder shows. (see [below for nested schema](#nestedatt--graph--limits--github_calls_per_run))
- `month_budget_tokens` (Number) Tokens the workflow may spend in a calendar month (UTC).
- `run_budget_tokens` (Number) Tokens one run may spend. At most monthBudgetTokens when both are set.
- `stop_the_line_attempts` (Number) Failed attempts after which a loop back to a crew stops the run, 1 to 10.

<a id="nestedatt--graph--limits--github_calls_per_run"></a>
### Nested Schema for `graph.limits.github_calls_per_run`

Optional:

- `max` (Number)
- `typical` (Number)



<a id="nestedatt--graph--scope"></a>
### Nested Schema for `graph.scope`

Optional:

- `orgs` (Set of String) GitHub organisations whose managed repositories the workflow may run on.
- `repos` (Set of String) Repositories as owner/name, owner/* for every repository of an owner, or * for every managed repository. Names are compared without case.



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
# Import an existing FactoryWorkflow by its name (metadata.name).
terraform import vulnetix_factory_workflow.patch_on_demand patch-on-demand
```
