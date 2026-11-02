---
page_title: "vulnetix_factory_settings Resource - vulnetix"
subcategory: ""
description: |-
  The Software Factory settings document: the tenant-wide factory settings the console's Settings page edits. One per tenant. Each field is one setting, and its dotted path (autonomy.max) is the key the console, the API and the audit trail name it by. A field left out means the setting's default, so applying a document resets every setting it does not name, and an export lists only the settings that differ from their defaults. Applying saves a new revision of the document through the console's own writer, with one audit entry per changed setting; status.revision is the revision saved. The security scanner settings are the ScannerProfile document, the managed organisations are FactoryManagedOrg documents and the web push defaults (notify.push.defaults) are the NotificationDefaults document, so none of them appears here. Read-only settings (forge.polling.forbidden, firewall_enabled, systemone.min_confidence.adapter) are not part of the document.
  Manages the FactorySettings config kind through /config/objects/FactorySettings/{name}. Applies are recorded with the provider's manager, so a GitOps vulnetix config apply --prune leaves this object alone.
  Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (most of them, such as ApiSpecSource, RiskStrategy and VexSource) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object, or lets go of it for a kind the server never deletes (such as Organization). Some kinds wait on a person after they are applied (an invited Member, a Slack connection): status.phase says so, and Terraform does not wait.
---

# vulnetix_factory_settings (Resource)

The Software Factory settings document: the tenant-wide factory settings the console's Settings page edits. One per tenant. Each field is one setting, and its dotted path (autonomy.max) is the key the console, the API and the audit trail name it by. A field left out means the setting's default, so applying a document resets every setting it does not name, and an export lists only the settings that differ from their defaults. Applying saves a new revision of the document through the console's own writer, with one audit entry per changed setting; status.revision is the revision saved. The security scanner settings are the ScannerProfile document, the managed organisations are FactoryManagedOrg documents and the web push defaults (notify.push.defaults) are the NotificationDefaults document, so none of them appears here. Read-only settings (forge.polling.forbidden, firewall_enabled, systemone.min_confidence.adapter) are not part of the document.

Manages the `FactorySettings` config kind through `/config/objects/FactorySettings/{name}`. Applies are recorded with the provider's `manager`, so a GitOps `vulnetix config apply --prune` leaves this object alone.

Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (most of them, such as `ApiSpecSource`, `RiskStrategy` and `VexSource`) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object, or lets go of it for a kind the server never deletes (such as `Organization`). Some kinds wait on a person after they are applied (an invited `Member`, a Slack connection): `status.phase` says so, and Terraform does not wait.

## Example Usage

### Minimal

The arguments the kind requires, and nothing else.

```terraform
# Only the arguments the kind requires. From config-schemas examples/FactorySettings/minimal.yaml.

resource "vulnetix_factory_settings" "default" {
}
```

### Advanced

Every argument the kind supports. A secret is a write-only `*_wo` argument fed from an ephemeral variable, so it never reaches the plan or the state; change `*_wo_version` to send a new value.

```terraform
# Every argument the kind supports. From config-schemas examples/FactorySettings/advanced.yaml.

resource "vulnetix_factory_settings" "default" {
  agents = {
    finalize_at           = 0.8
    max_attempts          = 0
    max_continuations     = 0
    max_cost_usd_per_item = 12.5
    max_iterations        = 0
    max_passes_per_item   = 0
    max_wall_minutes      = 0
  }
  approve = {
    sensitive_paths = [".github/workflows/**", "infra/**"]
  }
  audit = {
    calls         = true
    sign          = true
    weekly_report = true
  }
  autonomy = {
    claim_ttl_minutes  = 60
    concurrency        = 8
    default            = "suggest"
    max                = "auto_merge_low_risk"
    max_attempts       = 3
    no_host_ttl_hours  = 6
    parked_concurrency = 50
    plan_approval      = true
  }
  belai = {
    a_2_a = {
      show = true
    }
    ask_permission       = false
    auto_commit_per_task = true
    guardrails           = "on"
    mode                 = "goal"
    show_internal_work   = "decisions"
    tests = {
      on_fail  = "fix"
      post_end = "goal"
    }
  }
  budget = {
    alerts        = [50, 75, 90, 100]
    factory_share = 75
    month         = 50000000
    on_exceed     = "switch_to_pix_allocation"
  }
  crews = {
    lease = "30m"
  }
  defaults = {
    approve_ready = {
      enabled = true
    }
    ci_fix = {
      enabled = true
    }
    conflict_fix = {
      enabled = true
    }
    enabled = true
    issue_to_pr = {
      enabled = true
    }
    lint_sweep = {
      enabled = false
    }
    pr_review = {
      enabled = true
    }
    test_sweep = {
      enabled = true
    }
  }
  delivery = {
    company_domain  = "acme.example"
    compare         = true
    decisions       = true
    enabled         = true
    hours_per_month = 150
    idle_gap_min    = 5
    kinds = {
      goal = {
        enabled = true
        minutes = 3
      }
      merge = {
        enabled = true
        minutes = 8
      }
      mode = {
        enabled = true
        minutes = 0.5
      }
      other = {
        enabled = true
        minutes = 2
      }
      plan = {
        enabled = true
        minutes = 4
      }
      toolcall = {
        enabled = false
        minutes = 0
      }
      triage = {
        enabled = true
        minutes = 6
      }
    }
    manual_card_days     = 5
    manual_cycle_days    = 4.5
    manual_decision_wait = 6
    manual_issue_days    = 3
    multiplier           = 2.4
    overhead_pct         = 20
    overhead_rate        = 105
    people               = true
    people_rate          = 110
    rate                 = 120
    review_min_per_100   = 4
    roles = {
      devops = {
        enabled = true
        mult    = 2
        rate    = 90
      }
      planner = {
        enabled = true
        mult    = 1.6
        rate    = 120
      }
      qa = {
        enabled = true
        mult    = 2.2
        rate    = 70
      }
      reviewer = {
        enabled = true
        mult    = 1.5
        rate    = 95
      }
      security = {
        enabled = true
        mult    = 2
        rate    = 110
      }
      sse = {
        enabled = true
        mult    = 2.4
        rate    = 95
      }
      writer = {
        enabled = true
        mult    = 1.8
        rate    = 60
      }
    }
    sources = {
      admin        = true
      chat         = true
      checkpoint   = true
      comment      = true
      commit       = true
      notification = true
      pixcall      = true
      review       = true
    }
    value = true
  }
  factory = {
    fanout_cap = 6
  }
  forge = {
    branch_prefix  = "factory/"
    calls_per_hour = 5000
    check_runs     = false
    etag           = true
    poll = {
      fallback = "conditional_poll_10m"
    }
    pr_state  = "ready"
    protected = true
    quota = {
      factory_share             = 15
      on_exhausted              = "pause_live_reads"
      proactive_writes_per_hour = 60
      reserve                   = 20
    }
    read = {
      default = "posture_mirror"
    }
    stale = {
      alerts = "1h"
      checks = "2m"
      pulls  = "10m"
      repo   = "12h"
    }
  }
  hosts = {
    run_on          = "label:linux-large"
    run_on_fallback = "none"
  }
  jail = {
    enforcement_mode = "enforce"
  }
  kanban = {
    conflict_comment = true
    mode             = "all"
    opt_out_label    = "no-factory"
  }
  notify = {
    connector = "teams"
    events    = "escalations_holds_budget"
  }
  persistence = {
    export    = "manual"
    postgres  = "on"
    retention = "400d"
  }
  pix = {
    classifier     = "clef"
    egress_mode    = "deny_except_allow"
    extension_mode = "packs_only"
    jev_daily_cap  = 2000
    maintenance = [
      {
        days  = ["sun"]
        end   = "04:00"
        start = "02:00"
        tz    = "Australia/Melbourne"
      },
    ]
    reserve_percent = 20
    size            = "medium"
    version_policy  = "safest"
  }
  providers = {
    code = "pix-smart"
  }
  schedules = {
    blackout = [
      {
        from = {
          day  = "fri"
          time = "17:00"
        }
        label = "weekend"
        to = {
          day  = "mon"
          time = "08:00"
        }
      },
    ]
    missed = "run_once"
    tz     = "Australia/Melbourne"
  }
  scope = {
    archived     = true
    dependencies = false
    interest     = "scanned_or_pix"
    sync         = "1h"
  }
  signing = {
    auto_add    = true
    email       = "primary"
    mode        = "required"
    on_missing  = "readd"
    retire_days = "30"
    rotate_days = "180"
  }
  systemone = {
    below_threshold = "escalate_human"
    fallback        = "jev"
    min_confidence  = 0.85
    model           = "clef-flash"
    on_abstain      = "escalate_human"
    override        = true
    retention       = "full_400d"
  }
  tea = {
    publish = "on_release_tag"
  }
  vault = {
    grant = {
      entries = ["ve_ci_token"]
      mode    = "specific"
    }
    on_error = "pause_affected_workflows"
  }
  labels = {
    team = "platform"
  }
  annotations = {
    "vulnetix.com/owner" = "platform@acme.example"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `agents` (Attributes) Per-card limits for workers. (see [below for nested schema](#nestedatt--agents))
- `annotations` (Map of String) Free key/value notes for people and tools (`metadata.annotations`), such as a ticket link.
- `approve` (Attributes) (see [below for nested schema](#nestedatt--approve))
- `audit` (Attributes) The audit trail. (see [below for nested schema](#nestedatt--audit))
- `autonomy` (Attributes) How far the factory goes on its own, and how many runs it keeps going. (see [below for nested schema](#nestedatt--autonomy))
- `belai` (Attributes) How crews the factory starts behave. (see [below for nested schema](#nestedatt--belai))
- `budget` (Attributes) The factory's token budget. (see [below for nested schema](#nestedatt--budget))
- `crews` (Attributes) (see [below for nested schema](#nestedatt--crews))
- `defaults` (Attributes) The default workflows. One switch per default, plus the master switch. (see [below for nested schema](#nestedatt--defaults))
- `delivery` (Attributes) The Delivery pages: how outcomes, value, decisions and people are measured and priced. (see [below for nested schema](#nestedatt--delivery))
- `factory` (Attributes) (see [below for nested schema](#nestedatt--factory))
- `forge` (Attributes) How the factory reads and writes GitHub and the other forges. (see [below for nested schema](#nestedatt--forge))
- `hosts` (Attributes) Where crew steps run. (see [below for nested schema](#nestedatt--hosts))
- `jail` (Attributes) (see [below for nested schema](#nestedatt--jail))
- `kanban` (Attributes) How kanban cards are mirrored to forge issues. (see [below for nested schema](#nestedatt--kanban))
- `labels` (Map of String) Free key/value labels (`metadata.labels`). Stored with the object, never interpreted by the server.
- `name` (String) The object's name. FactorySettings is one per tenant, so it is always `default`.
- `notify` (Attributes) Notifications. (see [below for nested schema](#nestedatt--notify))
- `persistence` (Attributes) Where factory records are kept and for how long. (see [below for nested schema](#nestedatt--persistence))
- `pix` (Attributes) The Pix Sandbox the factory works in. (see [below for nested schema](#nestedatt--pix))
- `providers` (Attributes) (see [below for nested schema](#nestedatt--providers))
- `schedules` (Attributes) Schedules and blackout windows. (see [below for nested schema](#nestedatt--schedules))
- `scope` (Attributes) Which repositories of the managed organisations the factory looks at. The organisations themselves are FactoryManagedOrg documents. (see [below for nested schema](#nestedatt--scope))
- `signing` (Attributes) Commits a person makes from factory pages. (see [below for nested schema](#nestedatt--signing))
- `systemone` (Attributes) System One, which scores the factory's decisions. (see [below for nested schema](#nestedatt--systemone))
- `tea` (Attributes) (see [below for nested schema](#nestedatt--tea))
- `vault` (Attributes) Vault access from the Pix Sandbox. (see [below for nested schema](#nestedatt--vault))

### Read-Only

- `id` (String) The object's name, the same as `name`.
- `status` (Attributes) What the server says about the object beyond its spec. (see [below for nested schema](#nestedatt--status))

<a id="nestedatt--agents"></a>
### Nested Schema for `agents`

Optional:

- `finalize_at` (Number) Share of any per-card limit at which the worker is told to wrap up. A value under 0.5 is raised to 0.5; 0 keeps the agent's own. Schema default: `0`.
- `max_attempts` (Number) Failed attempts at one stage before Belai blocks the card. 0 keeps the agent's own. Not autonomy.max_attempts, which retries a workflow step. Schema default: `0`.
- `max_continuations` (Number) Times a card may be resumed after its budget ran out with progress made, before it is blocked for a person. 0 keeps the agent's own. Schema default: `0`.
- `max_cost_usd_per_item` (Number) Estimated cost in US dollars at which a worker's claim on one card ends. 0 keeps each agent's own limit. Schema default: `0`.
- `max_iterations` (Number) Tool rounds per pass. 0 keeps the agent's own. Schema default: `0`.
- `max_passes_per_item` (Number) Goal passes per card. 0 keeps the agent's own. Schema default: `0`.
- `max_wall_minutes` (Number) Minutes per card. 0 keeps the agent's own. Schema default: `0`.


<a id="nestedatt--approve"></a>
### Nested Schema for `approve`

Optional:

- `sensitive_paths` (Set of String) Path globs that stop an approval, besides the built-in list. Schema default: `[]`.


<a id="nestedatt--audit"></a>
### Nested Schema for `audit`

Optional:

- `calls` (Boolean) Log every model call with its input hash. Schema default: `true`.
- `sign` (Boolean) Sign the decision ledger. Schema default: `true`.
- `weekly_report` (Boolean) Send the weekly report. Schema default: `true`.


<a id="nestedatt--autonomy"></a>
### Nested Schema for `autonomy`

Optional:

- `claim_ttl_minutes` (Number) Minutes a started crew or worker dispatch may go without a session before the attempt fails card_unclaimed and is retried. Schema default: `60`.
- `concurrency` (Number) Runs of the tenant that may run at once. Schema default: `4`.
- `default` (String) The autonomy of a repository of a managed organisation that has no RepositorySettings of its own. auto_merge_low_risk is shown as approve when ready. One of `observe`, `suggest`, `draft_pr`, `auto_merge_low_risk`. Schema default: `draft_pr`.
- `max` (String) The most autonomy any repository may run at. A repository's own autonomy above it is lowered to it; the stored value is kept. One of `observe`, `suggest`, `draft_pr`, `auto_merge_low_risk`. Schema default: `draft_pr`.
- `max_attempts` (Number) Attempts at one workflow step before the run fails. Schema default: `3`.
- `no_host_ttl_hours` (Number) Hours a crew step may wait for a machine before the run ends no_host_timeout. Schema default: `6`.
- `parked_concurrency` (Number) Runs of the tenant that may wait for a machine at once. They do not count against autonomy.concurrency. Schema default: `50`.
- `plan_approval` (Boolean) Whether a belai:feature plan waits for a person to approve it. Schema default: `true`.


<a id="nestedatt--belai"></a>
### Nested Schema for `belai`

Optional:

- `a_2_a` (Attributes) (see [below for nested schema](#nestedatt--belai--a_2_a))
- `ask_permission` (Boolean) Whether the workers and crews a workflow or schedule starts stop at permission asks. Schema default: `false`.
- `auto_commit_per_task` (Boolean) Commit after each task a crew finishes. Schema default: `true`.
- `guardrails` (String) Belai guardrails for crews the factory starts. off is refused for worker crews. One of `on`, `strict`, `off`. Schema default: `on`.
- `mode` (String) The mode crews the factory starts run in. One of `goal`, `agent`, `plan`, `auto`. Schema default: `goal`.
- `show_internal_work` (String) Which internal work a transcript shows. One of `decisions`, `security`, `all`, `hidden`. Schema default: `decisions`.
- `tests` (Attributes) (see [below for nested schema](#nestedatt--belai--tests))

<a id="nestedatt--belai--a_2_a"></a>
### Nested Schema for `belai.a_2_a`

Optional:

- `show` (Boolean) Whether transcripts show agent-to-agent rows (broadcasts, messages, replies, acknowledgements) by default. A viewer can override it. Schema default: `true`.


<a id="nestedatt--belai--tests"></a>
### Nested Schema for `belai.tests`

Optional:

- `on_fail` (String) What a crew does when those tests fail: fix them, diagnose only, or nothing. One of `fix`, `diagnose`, `off`. Schema default: `fix`.
- `post_end` (String) When tests run after a crew ends: after the goal, after the goal and its plan, after the session, or never. One of `goal`, `goal_plan`, `session`, `off`. Schema default: `goal`.



<a id="nestedatt--budget"></a>
### Nested Schema for `budget`

Optional:

- `alerts` (List of Number) Percents of the monthly budget at which an alert is sent, ascending. Schema default: `[50,80,100]`.
- `factory_share` (Number) Percent of the Pix allowance the factory may use. Schema default: `75`.
- `month` (Number) Tokens per month, all providers. This, budget.on_exceed and budget.alerts are the tenant's Factory budget on the Budgets page. Schema default: `30000000`.
- `on_exceed` (String) What happens when the monthly budget is used up. One of `pause_run`, `downgrade_to_fast`, `switch_to_pix_allocation`, `escalate_human`, `hard_stop`. Schema default: `switch_to_pix_allocation`.


<a id="nestedatt--crews"></a>
### Nested Schema for `crews`

Optional:

- `lease` (String) How long a crew holds a handoff lease. One of `15m`, `30m`, `2h`. Schema default: `30m`.


<a id="nestedatt--defaults"></a>
### Nested Schema for `defaults`

Optional:

- `approve_ready` (Attributes) (see [below for nested schema](#nestedatt--defaults--approve_ready))
- `ci_fix` (Attributes) (see [below for nested schema](#nestedatt--defaults--ci_fix))
- `conflict_fix` (Attributes) (see [below for nested schema](#nestedatt--defaults--conflict_fix))
- `enabled` (Boolean) The master switch of the default workflows. Off, no default is seeded and none starts. Schema default: `true`.
- `issue_to_pr` (Attributes) (see [below for nested schema](#nestedatt--defaults--issue_to_pr))
- `lint_sweep` (Attributes) (see [below for nested schema](#nestedatt--defaults--lint_sweep))
- `pr_review` (Attributes) (see [below for nested schema](#nestedatt--defaults--pr_review))
- `test_sweep` (Attributes) (see [below for nested schema](#nestedatt--defaults--test_sweep))

<a id="nestedatt--defaults--approve_ready"></a>
### Nested Schema for `defaults.approve_ready`

Optional:

- `enabled` (Boolean) Whether the approve when ready default workflow is on for the tenant. A repository can still opt out with defaultsOff. Schema default: `true`.


<a id="nestedatt--defaults--ci_fix"></a>
### Nested Schema for `defaults.ci_fix`

Optional:

- `enabled` (Boolean) Whether the CI fix default workflow is on for the tenant. A repository can still opt out with defaultsOff. Schema default: `true`.


<a id="nestedatt--defaults--conflict_fix"></a>
### Nested Schema for `defaults.conflict_fix`

Optional:

- `enabled` (Boolean) Whether the merge conflict fix default workflow is on for the tenant. A repository can still opt out with defaultsOff. Schema default: `true`.


<a id="nestedatt--defaults--issue_to_pr"></a>
### Nested Schema for `defaults.issue_to_pr`

Optional:

- `enabled` (Boolean) Whether the issue to pull request default workflow is on for the tenant. A repository can still opt out with defaultsOff. Schema default: `true`.


<a id="nestedatt--defaults--lint_sweep"></a>
### Nested Schema for `defaults.lint_sweep`

Optional:

- `enabled` (Boolean) Whether the lint sweep default workflow is on for the tenant. A repository can still opt out with defaultsOff. Schema default: `true`.


<a id="nestedatt--defaults--pr_review"></a>
### Nested Schema for `defaults.pr_review`

Optional:

- `enabled` (Boolean) Whether the pull request review default workflow is on for the tenant. A repository can still opt out with defaultsOff. Schema default: `true`.


<a id="nestedatt--defaults--test_sweep"></a>
### Nested Schema for `defaults.test_sweep`

Optional:

- `enabled` (Boolean) Whether the test sweep default workflow is on for the tenant. A repository can still opt out with defaultsOff. Schema default: `true`.



<a id="nestedatt--delivery"></a>
### Nested Schema for `delivery`

Optional:

- `company_domain` (String) An email domain such as acme.example. A contributor on it, or on a subdomain, counts as internal. Empty for none. Schema default: ``.
- `compare` (Boolean) Show the hand-worked comparison, which adds the refinement and hand-off overhead. Schema default: `true`.
- `decisions` (Boolean) Price the decisions System One scored. Schema default: `true`.
- `enabled` (Boolean) The master switch of the Delivery pages. Off, nothing is computed. Schema default: `true`.
- `hours_per_month` (Number) Productive hours per person-month. Schema default: `150`.
- `idle_gap_min` (Number) Minutes without an action that end a sitting. Schema default: `5`.
- `kinds` (Attributes) Pricing per decision kind. (see [below for nested schema](#nestedatt--delivery--kinds))
- `manual_card_days` (Number) Days a team takes per card. Schema default: `5`.
- `manual_cycle_days` (Number) Days a team takes per pull request. Schema default: `4.5`.
- `manual_decision_wait` (Number) Hours a team waits to decide, compared with the measured time. Schema default: `6`.
- `manual_issue_days` (Number) Days a team takes per issue. Schema default: `3`.
- `multiplier` (Number) Engineer hours per agent hour for agent work no role covers. Schema default: `2.4`.
- `overhead_pct` (Number) Refinement and hand-off overhead, as a percent of build and decision hours. Schema default: `20`.
- `overhead_rate` (Number) US dollars per hour for the overhead. Schema default: `105`.
- `people` (Boolean) Show the People and Contributors pages. Schema default: `true`.
- `people_rate` (Number) US dollars per hour for the hours people spent in the factory. Schema default: `110`.
- `rate` (Number) US dollars per hour for agent work no role covers. Schema default: `95`.
- `review_min_per_100` (Number) Minutes of reading per 100 changed lines of a review. Schema default: `4`.
- `roles` (Attributes) Pricing per role. (see [below for nested schema](#nestedatt--delivery--roles))
- `sources` (Attributes) Where to look for people. (see [below for nested schema](#nestedatt--delivery--sources))
- `value` (Boolean) Price agent work by role. Schema default: `true`.

<a id="nestedatt--delivery--kinds"></a>
### Nested Schema for `delivery.kinds`

Optional:

- `goal` (Attributes) (see [below for nested schema](#nestedatt--delivery--kinds--goal))
- `merge` (Attributes) (see [below for nested schema](#nestedatt--delivery--kinds--merge))
- `mode` (Attributes) (see [below for nested schema](#nestedatt--delivery--kinds--mode))
- `other` (Attributes) (see [below for nested schema](#nestedatt--delivery--kinds--other))
- `plan` (Attributes) (see [below for nested schema](#nestedatt--delivery--kinds--plan))
- `toolcall` (Attributes) (see [below for nested schema](#nestedatt--delivery--kinds--toolcall))
- `triage` (Attributes) (see [below for nested schema](#nestedatt--delivery--kinds--triage))

<a id="nestedatt--delivery--kinds--goal"></a>
### Nested Schema for `delivery.kinds.goal`

Optional:

- `enabled` (Boolean) Count goal decisions in the estimate. Schema default: `true`.
- `minutes` (Number) Minutes a person would take over one goal decision. Schema default: `3`.


<a id="nestedatt--delivery--kinds--merge"></a>
### Nested Schema for `delivery.kinds.merge`

Optional:

- `enabled` (Boolean) Count merge decisions in the estimate. Schema default: `true`.
- `minutes` (Number) Minutes a person would take over one merge decision. Schema default: `8`.


<a id="nestedatt--delivery--kinds--mode"></a>
### Nested Schema for `delivery.kinds.mode`

Optional:

- `enabled` (Boolean) Count mode change decisions in the estimate. Schema default: `true`.
- `minutes` (Number) Minutes a person would take over one mode change decision. Schema default: `0.5`.


<a id="nestedatt--delivery--kinds--other"></a>
### Nested Schema for `delivery.kinds.other`

Optional:

- `enabled` (Boolean) Count other decisions in the estimate. Schema default: `true`.
- `minutes` (Number) Minutes a person would take over one other decision. Schema default: `2`.


<a id="nestedatt--delivery--kinds--plan"></a>
### Nested Schema for `delivery.kinds.plan`

Optional:

- `enabled` (Boolean) Count plan decisions in the estimate. Schema default: `true`.
- `minutes` (Number) Minutes a person would take over one plan decision. Schema default: `4`.


<a id="nestedatt--delivery--kinds--toolcall"></a>
### Nested Schema for `delivery.kinds.toolcall`

Optional:

- `enabled` (Boolean) Count tool call guard rail decisions in the estimate. Tool calls are never in the estimate, whatever this switch says. Schema default: `false`.
- `minutes` (Number) Minutes a person would take over one tool call guard rail decision. Schema default: `0`.


<a id="nestedatt--delivery--kinds--triage"></a>
### Nested Schema for `delivery.kinds.triage`

Optional:

- `enabled` (Boolean) Count triage decisions in the estimate. Schema default: `true`.
- `minutes` (Number) Minutes a person would take over one triage decision. Schema default: `6`.



<a id="nestedatt--delivery--roles"></a>
### Nested Schema for `delivery.roles`

Optional:

- `devops` (Attributes) (see [below for nested schema](#nestedatt--delivery--roles--devops))
- `planner` (Attributes) (see [below for nested schema](#nestedatt--delivery--roles--planner))
- `qa` (Attributes) (see [below for nested schema](#nestedatt--delivery--roles--qa))
- `reviewer` (Attributes) (see [below for nested schema](#nestedatt--delivery--roles--reviewer))
- `security` (Attributes) (see [below for nested schema](#nestedatt--delivery--roles--security))
- `sse` (Attributes) (see [below for nested schema](#nestedatt--delivery--roles--sse))
- `writer` (Attributes) (see [below for nested schema](#nestedatt--delivery--roles--writer))

<a id="nestedatt--delivery--roles--devops"></a>
### Nested Schema for `delivery.roles.devops`

Optional:

- `enabled` (Boolean) Price agent work in the DevOps engineer role. Schema default: `true`.
- `mult` (Number) Engineer hours per agent hour for the DevOps engineer role. Schema default: `2`.
- `rate` (Number) US dollars per hour for the DevOps engineer role. Schema default: `90`.


<a id="nestedatt--delivery--roles--planner"></a>
### Nested Schema for `delivery.roles.planner`

Optional:

- `enabled` (Boolean) Price agent work in the planner role. Schema default: `true`.
- `mult` (Number) Engineer hours per agent hour for the planner role. Schema default: `1.6`.
- `rate` (Number) US dollars per hour for the planner role. Schema default: `120`.


<a id="nestedatt--delivery--roles--qa"></a>
### Nested Schema for `delivery.roles.qa`

Optional:

- `enabled` (Boolean) Price agent work in the QA engineer role. Schema default: `true`.
- `mult` (Number) Engineer hours per agent hour for the QA engineer role. Schema default: `2.2`.
- `rate` (Number) US dollars per hour for the QA engineer role. Schema default: `70`.


<a id="nestedatt--delivery--roles--reviewer"></a>
### Nested Schema for `delivery.roles.reviewer`

Optional:

- `enabled` (Boolean) Price agent work in the code reviewer role. Schema default: `true`.
- `mult` (Number) Engineer hours per agent hour for the code reviewer role. Schema default: `1.5`.
- `rate` (Number) US dollars per hour for the code reviewer role. Schema default: `95`.


<a id="nestedatt--delivery--roles--security"></a>
### Nested Schema for `delivery.roles.security`

Optional:

- `enabled` (Boolean) Price agent work in the security engineer role. Schema default: `true`.
- `mult` (Number) Engineer hours per agent hour for the security engineer role. Schema default: `2`.
- `rate` (Number) US dollars per hour for the security engineer role. Schema default: `110`.


<a id="nestedatt--delivery--roles--sse"></a>
### Nested Schema for `delivery.roles.sse`

Optional:

- `enabled` (Boolean) Price agent work in the software engineer role. Schema default: `true`.
- `mult` (Number) Engineer hours per agent hour for the software engineer role. Schema default: `2.4`.
- `rate` (Number) US dollars per hour for the software engineer role. Schema default: `95`.


<a id="nestedatt--delivery--roles--writer"></a>
### Nested Schema for `delivery.roles.writer`

Optional:

- `enabled` (Boolean) Price agent work in the technical writer role. Schema default: `true`.
- `mult` (Number) Engineer hours per agent hour for the technical writer role. Schema default: `1.8`.
- `rate` (Number) US dollars per hour for the technical writer role. Schema default: `60`.



<a id="nestedatt--delivery--sources"></a>
### Nested Schema for `delivery.sources`

Optional:

- `admin` (Boolean) Look for people in board and admin actions. Schema default: `true`.
- `chat` (Boolean) Look for people in chat messages. Schema default: `true`.
- `checkpoint` (Boolean) Look for people in checkpoints. Schema default: `true`.
- `comment` (Boolean) Look for people in comments. Schema default: `true`.
- `commit` (Boolean) Look for people in commits. Schema default: `true`.
- `notification` (Boolean) Look for people in notification actions. Schema default: `true`.
- `pixcall` (Boolean) Look for people in Pix calls. Schema default: `true`.
- `review` (Boolean) Look for people in reviews. Schema default: `true`.



<a id="nestedatt--factory"></a>
### Nested Schema for `factory`

Optional:

- `fanout_cap` (Number) Runs one delivery or sighting may start. Triggers past it, by priority, are skipped as fanout_cap. Schema default: `6`.


<a id="nestedatt--forge"></a>
### Nested Schema for `forge`

Optional:

- `branch_prefix` (String) The prefix of branches the factory pushes. Schema default: `belai/K-`.
- `calls_per_hour` (Number) The hourly call budget of a forge that sends no rate-limit headers, such as Gitea or Forgejo. Schema default: `5000`.
- `check_runs` (Boolean) Publish gate results as check runs. Needs the App permission checks: write on every managed installation. Schema default: `false`.
- `etag` (Boolean) Send If-None-Match on every live read. Schema default: `true`.
- `poll` (Attributes) (see [below for nested schema](#nestedatt--forge--poll))
- `pr_state` (String) Whether the factory opens pull requests as drafts or ready for review. One of `draft`, `ready`. Schema default: `draft`.
- `protected` (Boolean) Never push to protected branches. Schema default: `true`.
- `quota` (Attributes) (see [below for nested schema](#nestedatt--forge--quota))
- `read` (Attributes) (see [below for nested schema](#nestedatt--forge--read))
- `stale` (Attributes) (see [below for nested schema](#nestedatt--forge--stale))

<a id="nestedatt--forge--poll"></a>
### Nested Schema for `forge.poll`

Optional:

- `fallback` (String) What a binding with no webhook does. One of `conditional_poll_10m`, `none`. Schema default: `conditional_poll_10m`.


<a id="nestedatt--forge--quota"></a>
### Nested Schema for `forge.quota`

Optional:

- `factory_share` (Number) Percent of the forge call budget the factory may use. Schema default: `10`.
- `on_exhausted` (String) What the factory does when its share is used up. One of `pause_live_reads`, `defer_to_next_window`, `escalate`. Schema default: `pause_live_reads`.
- `proactive_writes_per_hour` (Number) Most calls an installation makes in a clock hour for default workflows. Schema default: `60`.
- `reserve` (Number) Percent of the forge call budget kept in reserve. At most 90, and with forge.quota.factory_share at most 100. Schema default: `20`.


<a id="nestedatt--forge--read"></a>
### Nested Schema for `forge.read`

Optional:

- `default` (String) Where forge reads come from by default. One of `posture_mirror`, `s3_raw`, `live_api`. Schema default: `posture_mirror`.


<a id="nestedatt--forge--stale"></a>
### Nested Schema for `forge.stale`

Optional:

- `alerts` (String) How old alert data may be before it is read again. A duration from 30s to 7d, such as 30s, 5m, 6h or 2d. Schema default: `1h`.
- `checks` (String) How old check data may be before it is read again. A duration from 30s to 7d, such as 30s, 5m, 6h or 2d. Schema default: `2m`.
- `pulls` (String) How old pull request data may be before it is read again. A duration from 30s to 7d, such as 30s, 5m, 6h or 2d. Schema default: `5m`.
- `repo` (String) How old repository data may be before it is read again. A duration from 30s to 7d, such as 30s, 5m, 6h or 2d. Schema default: `6h`.



<a id="nestedatt--hosts"></a>
### Nested Schema for `hosts`

Optional:

- `run_on` (String) Where a crew node, a default or a scan that names no runOn runs: auto, pix_sandbox, self_hosted, host:<uuid> or label:<name>. It also orders the hosts a node's own runOn admits. Schema default: `auto`.
- `run_on_fallback` (String) The fallback that goes with hosts.runOn. A node that sets its own runOnFallback keeps it. One of `none`, `auto`. Schema default: `auto`.


<a id="nestedatt--jail"></a>
### Nested Schema for `jail`

Optional:

- `enforcement_mode` (String) Whether release jails block (enforce) or only report (warn). One of `enforce`, `warn`. Schema default: `enforce`.


<a id="nestedatt--kanban"></a>
### Nested Schema for `kanban`

Optional:

- `conflict_comment` (Boolean) Post a comment when the kanban sync finds a conflict. Schema default: `true`.
- `mode` (String) The default mode of a kanban binding: which cards are mirrored to forge issues. all: every card. skip_review: a card waits until it leaves Review. off: none. One of `all`, `skip_review`, `off`. Schema default: `all`.
- `opt_out_label` (String) A card with this label is never mirrored to an issue. Schema default: `no-github`.


<a id="nestedatt--notify"></a>
### Nested Schema for `notify`

Optional:

- `connector` (String) Where factory notifications go. One of `slack`, `teams`, `discord`, `google_chat`, `webhook`, `jira`. Schema default: `slack`.
- `events` (String) Which events are sent to notify.connector. One of `escalations_holds_budget`, `everything`, `escalations_only`. Schema default: `escalations_holds_budget`.


<a id="nestedatt--persistence"></a>
### Nested Schema for `persistence`

Optional:

- `export` (String) How factory records are exported. One of `daily_jsonl_s3`, `manual`. Schema default: `manual`.
- `postgres` (String) Keep factory records in Postgres. off is accepted only while no workflow is published. One of `on`, `off`. Schema default: `on`.
- `retention` (String) How long factory records are kept. One of `400d`, `90d`, `30d`, `forever`. Schema default: `400d`.


<a id="nestedatt--pix"></a>
### Nested Schema for `pix`

Optional:

- `classifier` (String) The Pix Sandbox's own classifier. One of `clef`, `jev`, `tev1`. Schema default: `clef`.
- `egress_mode` (String) The Pix Sandbox egress mode. Takes effect at the sandbox's next launch. One of `deny_except_allow`, `allow_all_logged`. Schema default: `deny_except_allow`.
- `extension_mode` (String) Which editor extensions a Pix Sandbox may install: curated packs only, or any indexed extension. One of `packs_only`, `any_indexed`. Schema default: `packs_only`.
- `jev_daily_cap` (Number) Calls per day to the jev classifier, shared with the Pix Sandbox. Schema default: `2000`.
- `maintenance` (Attributes List) Weekly Pix Sandbox maintenance windows. Schema default: `[]`. (see [below for nested schema](#nestedatt--pix--maintenance))
- `reserve_percent` (Number) Percent of the Pix allowance kept for interactive use. Schema default: `20`.
- `size` (String) The Pix Sandbox size. Takes effect at the sandbox's next launch. One of `small`, `medium`, `large`. Schema default: `medium`.
- `version_policy` (String) How a Pix Sandbox picks package versions. safest, stable, latest and rolling are the Safe Harbour picks. Takes effect at the sandbox's next launch. One of `safest`, `stable`, `latest`, `rolling`, `pinned`. Schema default: `safest`.

<a id="nestedatt--pix--maintenance"></a>
### Nested Schema for `pix.maintenance`

Required:

- `days` (Set of String) The days the window starts on.
- `end` (String) End time, HH:MM, in tz. An end before the start runs past midnight.
- `start` (String) Start time, HH:MM, in tz.
- `tz` (String) The IANA time zone of start and end.



<a id="nestedatt--providers"></a>
### Nested Schema for `providers`

Optional:

- `code` (String) The model for code work: pix-smart, pix-fast, or vault:<entry id> for a model in the vault. Schema default: `pix-smart`.


<a id="nestedatt--schedules"></a>
### Nested Schema for `schedules`

Optional:

- `blackout` (Attributes List) Weekly windows in which no scheduled run starts, in schedules.tz. Schema default: `[]`. (see [below for nested schema](#nestedatt--schedules--blackout))
- `missed` (String) What happens to runs a schedule missed while the factory was paused or in a blackout. One of `run_once`, `skip`, `run_all`. Schema default: `run_once`.
- `tz` (String) The IANA time zone of schedules and blackout windows. Schema default: `UTC`.

<a id="nestedatt--schedules--blackout"></a>
### Nested Schema for `schedules.blackout`

Required:

- `from` (Attributes) When the window starts each week. (see [below for nested schema](#nestedatt--schedules--blackout--from))
- `to` (Attributes) When it ends. (see [below for nested schema](#nestedatt--schedules--blackout--to))

Optional:

- `label` (String) A name for the window.

<a id="nestedatt--schedules--blackout--from"></a>
### Nested Schema for `schedules.blackout.from`

Required:

- `day` (String) One of `sun`, `mon`, `tue`, `wed`, `thu`, `fri`, `sat`.
- `time` (String) HH:MM.


<a id="nestedatt--schedules--blackout--to"></a>
### Nested Schema for `schedules.blackout.to`

Required:

- `day` (String) One of `sun`, `mon`, `tue`, `wed`, `thu`, `fri`, `sat`.
- `time` (String) HH:MM.




<a id="nestedatt--scope"></a>
### Nested Schema for `scope`

Optional:

- `archived` (Boolean) Include archived repositories. Schema default: `false`.
- `dependencies` (Boolean) Include public dependency repositories that only appear through SBOM uploads. Schema default: `false`.
- `interest` (String) Which repositories of a managed organisation count as of interest: any_signal (any signal at all), scanned_or_pix (scanned by Vulnetix or cloned in a Pix Sandbox) or watched_only (watched by a member). One of `any_signal`, `scanned_or_pix`, `watched_only`. Schema default: `any_signal`.
- `sync` (String) How often the installation metadata sweep runs: every 6h, every 1h, or webhook_only (no sweep). One of `6h`, `1h`, `webhook_only`. Schema default: `6h`.


<a id="nestedatt--signing"></a>
### Nested Schema for `signing`

Optional:

- `auto_add` (Boolean) Add a new signing key's public key to the person's GitHub account. Schema default: `true`.
- `email` (String) The address on a signed commit: the person's primary address, their noreply address, or ask each time. One of `primary`, `noreply`, `ask`. Schema default: `primary`.
- `mode` (String) Whether commits a person makes from factory pages must be signed. One of `required`, `preferred`, `off`. Schema default: `preferred`.
- `on_missing` (String) What happens when GitHub no longer lists a key: add it again, or ask. One of `readd`, `ask`. Schema default: `readd`.
- `retire_days` (String) Days a replaced key stays on GitHub. 0 removes it right away; -1 leaves it there. One of `30`, `7`, `0`, `-1`. Schema default: `30`.
- `rotate_days` (String) Days before a signing key is rotated. 0 never rotates. One of `365`, `180`, `90`, `0`. Schema default: `365`.


<a id="nestedatt--systemone"></a>
### Nested Schema for `systemone`

Optional:

- `below_threshold` (String) What happens to a decision scored below systemone.min_confidence. One of `escalate_human`, `route_to_fallback_model`, `take_safest_outcome`, `end_run`. Schema default: `escalate_human`.
- `fallback` (String) The model a decision is routed to under route_to_fallback_model. One of `jev`, `pix-fast`, `none`. Schema default: `jev`.
- `min_confidence` (Number) The confidence below which a System One decision is not taken on its own. Schema default: `0.8`.
- `model` (String) The System One adapter that scores decisions. Confidences are not comparable across adapters, so changing it records systemone.min_confidence against the new adapter. One of `clef-flash`, `clef`, `jev`, `tev1`. Schema default: `clef-flash`.
- `on_abstain` (String) What happens when System One abstains. One of `escalate_human`, `retry_once_then_escalate`, `end_run`. Schema default: `escalate_human`.
- `override` (Boolean) Whether a person may override any System One decision. Schema default: `true`.
- `retention` (String) How System One decisions are kept: in full for 400 days, as hashes only, or for 30 days. One of `full_400d`, `hash_only`, `days_30`. Schema default: `full_400d`.


<a id="nestedatt--tea"></a>
### Nested Schema for `tea`

Optional:

- `publish` (String) When releases are published to TEA. One of `on_release_tag`, `on_merge_main`, `never`. Schema default: `on_release_tag`.


<a id="nestedatt--vault"></a>
### Nested Schema for `vault`

Optional:

- `grant` (Attributes) Which vault entries the Pix Sandbox may use. Schema default: `{"mode":"all"}`. (see [below for nested schema](#nestedatt--vault--grant))
- `on_error` (String) What happens when a vault entry cannot be read. One of `pause_affected_workflows`, `fall_back_to_pix`. Schema default: `pause_affected_workflows`.

<a id="nestedatt--vault--grant"></a>
### Nested Schema for `vault.grant`

Required:

- `mode` (String) all: every entry. specific: only those in entries. all_except: every entry except those in entries. One of `specific`, `all`, `all_except`.

Optional:

- `entries` (Set of String) Vault entry ids, for specific and all_except.



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
# Import an existing FactorySettings by its name (metadata.name).
terraform import vulnetix_factory_settings.default default
```
