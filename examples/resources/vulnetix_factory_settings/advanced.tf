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
