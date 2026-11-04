# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/AiFirewallPolicy/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_ai_firewall_policy" "default" {
  hold_policy = {
    hold_timeout_action  = "wait"
    hold_timeout_minutes = 0
    hold_unknown         = false
  }
  prune = false
}
