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
