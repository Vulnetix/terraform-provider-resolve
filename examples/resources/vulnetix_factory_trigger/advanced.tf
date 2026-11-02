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
