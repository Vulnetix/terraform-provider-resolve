---
page_title: "vulnetix_ai_firewall_policy Resource - vulnetix"
subcategory: ""
description: |-
  The tenant's AI Firewall policy: which providers and models its agents may use, the guardrails the gateway applies to every request, the hold for review setting, and capability allow and deny lists. One per tenant. This is the document vulnetix ai-firewall apply reads from .vulnetix/ai-firewall.yaml, and every file that command accepts is valid here; metadata, and metadata.name, may be left out.
  An object of this kind can reference AiProvider objects. Reference those resources' attributes, or add depends_on, so Terraform applies them first.
  Manages the AiFirewallPolicy config kind through /config/objects/AiFirewallPolicy/{name}. Applies are recorded with the provider's manager, so a GitOps vulnetix config apply --prune leaves this object alone.
  Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (such as ApiSpecSource, RiskStrategy and VexSource) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object.
---

# vulnetix_ai_firewall_policy (Resource)

The tenant's AI Firewall policy: which providers and models its agents may use, the guardrails the gateway applies to every request, the hold for review setting, and capability allow and deny lists. One per tenant. This is the document vulnetix ai-firewall apply reads from .vulnetix/ai-firewall.yaml, and every file that command accepts is valid here; metadata, and metadata.name, may be left out.

An object of this kind can reference `AiProvider` objects. Reference those resources' attributes, or add `depends_on`, so Terraform applies them first.

Manages the `AiFirewallPolicy` config kind through `/config/objects/AiFirewallPolicy/{name}`. Applies are recorded with the provider's `manager`, so a GitOps `vulnetix config apply --prune` leaves this object alone.

Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (such as `ApiSpecSource`, `RiskStrategy` and `VexSource`) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object.

## Example Usage

### Minimal

The arguments the kind requires, and nothing else.

```terraform
# Only the arguments the kind requires. From config-schemas examples/AiFirewallPolicy/minimal.yaml.

resource "vulnetix_ai_firewall_policy" "default" {
}
```

### Advanced

Every argument the kind supports. A secret is a write-only `*_wo` argument fed from an ephemeral variable, so it never reaches the plan or the state; change `*_wo_version` to send a new value.

```terraform
# Every argument the kind supports. From config-schemas examples/AiFirewallPolicy/advanced.yaml.

resource "vulnetix_ai_firewall_policy" "default" {
  baseline = {
    enabled = true
    exclude = ["pii-phone"]
    ref     = "recommended"
  }
  capabilities = {
    approved-mcp = {
      action   = "block"
      family   = "mcp"
      list     = "allow"
      patterns = ["github", "linear*"]
    }
    banned-clients = {
      action   = "block"
      family   = "client"
      list     = "deny"
      patterns = ["unknown"]
      priority = 10
    }
  }
  guardrails = [
    {
      action      = "block"
      baseline_id = "secrets-aws-access-key"
      enabled     = true
      name        = "block-aws-keys"
      pattern     = "(?i)AKIA[0-9A-Z]{16}"
      priority    = 10
      rule_type   = "blocked_pattern"
    },
    {
      action    = "block"
      name      = "cap-conversation"
      pattern   = "200"
      priority  = 50
      rule_type = "max_messages"
    },
    {
      action    = "redact"
      name      = "redact-pii"
      priority  = 20
      rule_type = "pii_redact"
    },
    {
      action    = "block"
      name      = "sealed-blocks"
      priority  = 5
      rule_type = "delimiter_integrity"
    },
    {
      action    = "strip"
      name      = "strip-shell-for-codex"
      pattern   = "Bash"
      priority  = 30
      rule_type = "tool_deny"
      when      = "client == 'codex'"
    },
    {
      action    = "flag"
      enabled   = false
      name      = "flag-filesystem-writes"
      priority  = 40
      rule_type = "expression"
      when      = "'filesystem' in mcpServers && tools.exists(t, t.startsWith('mcp__filesystem__write'))"
    },
  ]
  hold_policy = {
    hold_timeout_action  = "deny"
    hold_timeout_minutes = 240
    hold_unknown         = true
  }
  models = [
    {
      action       = "allow"
      any_provider = false
      provider     = "anthropic"
      slug         = "claude-sonnet-4-5"
    },
    {
      action   = "allow"
      provider = "acme-vllm"
      slug     = "qwen3-coder-30b"
    },
    {
      action       = "deny"
      any_provider = true
      slug         = "gpt-4o-mini"
    },
  ]
  providers = [
    {
      action = "allow"
      key = {
        from_env = "ANTHROPIC_PROVIDER_KEY"
      }
      slug = "anthropic"
    },
    {
      action = "allow"
      key = {
        from_file = "~/.config/acme/vllm.key"
      }
      slug = "acme-vllm"
    },
    {
      action = "deny"
      slug   = "openrouter"
    },
    {
      action = "default"
      slug   = "groq"
    },
  ]
  prune = true
  settings = {
    logs_enabled = true
  }
  labels = {
    team = "platform-security"
  }
  annotations = {
    ticket = "https://tracker.example.com/SEC-2210"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `annotations` (Map of String) Free key/value notes for people and tools (`metadata.annotations`), such as a ticket link.
- `baseline` (Attributes) The server's recommended guardrails. Left out, they are added to the desired set; set enabled to false to decline them. (see [below for nested schema](#nestedatt--baseline))
- `capabilities` (Attributes Map) Capability allow and deny lists, keyed by a name for the list. Each pattern becomes one guardrail, the same rows the guardrails list writes; this is the shorter way to write a list of them. An AiGovernanceList is the same thing as its own object. (see [below for nested schema](#nestedatt--capabilities))
- `guardrails` (Attributes Set) Rules the gateway applies to each request, reconciled by name, so names must be unique. Content rules match what a request says; capability rules match the tools, MCP servers and skills it offers and the client sending it. (see [below for nested schema](#nestedatt--guardrails))
- `hold_policy` (Attributes) Hold for review: what the gateway does with a tool, MCP server, skill or client it has not seen before. An approved capability is recorded as reviewed rather than written as an allow rule, which would put its whole family into allowlist mode. (see [below for nested schema](#nestedatt--hold_policy))
- `labels` (Map of String) Free key/value labels (`metadata.labels`). Stored with the object, never interpreted by the server.
- `models` (Attributes Set) Model allow and deny entries. Once a provider has any allowed model, its other models are denied. (see [below for nested schema](#nestedatt--models))
- `name` (String) The object's name. AiFirewallPolicy is one per tenant, so it is always `default`.
- `providers` (Attributes Set) Provider allow and deny entries, and where each provider's API key comes from. A provider with no entry is usable. (see [below for nested schema](#nestedatt--providers))
- `prune` (Boolean) Delete guardrails and model entries on the server that this document does not mention. Off by default, so an apply never removes rules someone made in the console. A provider entry is never pruned: clearing a deny opens access. Schema default: `false`.
- `settings` (Attributes) (see [below for nested schema](#nestedatt--settings))

### Read-Only

- `id` (String) The object's name, the same as `name`.

<a id="nestedatt--baseline"></a>
### Nested Schema for `baseline`

Optional:

- `enabled` (Boolean)
- `exclude` (Set of String) Baseline guardrail ids to leave out.
- `ref` (String) A named baseline set to pin.


<a id="nestedatt--capabilities"></a>
### Nested Schema for `capabilities`

Required:

- `action` (String) block refuses the request. flag lets it through and records the match. strip removes the matched tool or MCP server and sends the rest of the request; it applies to the tool and mcp families only, because a skill is a value inside another tool and a client is a header. One of `block`, `flag`, `strip`.
- `family` (String) tool: a tool offered to the model. mcp: an MCP server wired into the session. skill: an agent skill. client: the calling client, named by its headers (claude-code, codex-cli). One of `tool`, `mcp`, `skill`, `client`.
- `list` (String) One of `allow`, `deny`.
- `patterns` (Set of String)

Optional:

- `enabled` (Boolean) Schema default: `true`.
- `priority` (Number) Evaluation order among guardrails, lowest first. Schema default: `100`.


<a id="nestedatt--guardrails"></a>
### Nested Schema for `guardrails`

Required:

- `action` (String) block refuses the request. redact rewrites the matched content (content rules only). flag lets it through and records the match. strip removes the matched tool or MCP server (tool_* and mcp_* rules only). One of `block`, `redact`, `flag`, `strip`.
- `name` (String)
- `rule_type` (String) Content rules: blocked_pattern (an RE2 regex over the request text), max_messages (a cap on messages per request), pii_redact (email, card, SSN and phone detectors, or your own regex), delimiter_integrity (sealed harness blocks left intact). Capability rules: <family>_allow and <family>_deny for tool, mcp, skill and client. expression: the request matches when its when expression is true, over the request variables. One of `blocked_pattern`, `max_messages`, `pii_redact`, `delimiter_integrity`, `tool_allow`, `tool_deny`, `mcp_allow`, `mcp_deny`, `skill_allow`, `skill_deny`, `client_allow`, `client_deny`, `expression`.

Optional:

- `baseline_id` (String) The baseline guardrail this one came from.
- `enabled` (Boolean) Schema default: `true`.
- `pattern` (String) What the rule matches. blocked_pattern: an RE2 regex, required (RE2 has no lookahead or lookbehind). pii_redact: an RE2 regex, or empty for the built-in detectors. max_messages: a positive integer, required. delimiter_integrity: none. Capability rules: a name or * glob, required, anchors implied.
- `priority` (Number) Evaluation order, lowest first.
- `when` (String) A CEL condition over the request variables (provider, model, client, tools, mcpServers, skills, now). The rule applies only to requests for which it is true; left out, it applies to every request. For an expression rule it is the match itself, and is required. A CEL expression, evaluated in the `request` environment.


<a id="nestedatt--hold_policy"></a>
### Nested Schema for `hold_policy`

Optional:

- `hold_timeout_action` (String) What happens when the wait runs out: keep waiting, deny the capability, or allow it. One of `wait`, `deny`, `allow`. Schema default: `wait`.
- `hold_timeout_minutes` (Number) How long a hold waits for a decision, up to 7 days. 0 waits indefinitely. Schema default: `0`.
- `hold_unknown` (Boolean) Hold a request that carries an unseen capability until someone approves or denies it. Schema default: `false`.


<a id="nestedatt--models"></a>
### Nested Schema for `models`

Required:

- `action` (String) One of `allow`, `deny`.
- `slug` (String) The model's id at its provider.

Optional:

- `any_provider` (Boolean) Apply the entry to every provider listing this model.
- `provider` (String) The provider slug the entry applies to.


<a id="nestedatt--providers"></a>
### Nested Schema for `providers`

Required:

- `action` (String) allow pins the provider as allowed, deny blocks it, default removes any pin. One of `allow`, `deny`, `default`.
- `slug` (String) A built-in provider slug, or the metadata.name of an AiProvider.

Optional:

- `key` (Attributes) Where the client reads the provider API key it uploads. Never the key itself. fromEnv wins when both are set. (see [below for nested schema](#nestedatt--providers--key))

<a id="nestedatt--providers--key"></a>
### Nested Schema for `providers.key`

Optional:

- `from_env` (String)
- `from_file` (String)



<a id="nestedatt--settings"></a>
### Nested Schema for `settings`

Optional:

- `logs_enabled` (Boolean) Record per-inference metadata (model, decision, tokens, latency, never the prompt or completion). Paid plans only.

## Import

Import an existing object by its name (`metadata.name`). The first read after an import fills in every field the server holds.

```shell
# Import an existing AiFirewallPolicy by its name (metadata.name).
terraform import vulnetix_ai_firewall_policy.default default
```
