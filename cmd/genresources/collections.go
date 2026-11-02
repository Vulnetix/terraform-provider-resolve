package main

// collection says how one array property is carried in Terraform.
//
// A list keeps the order the configuration gives, and a reorder is a diff. A
// set has no order, so it only fits where the server ignores order: there a
// list would show a diff on every plan whenever the server reads the array
// back in its own order (sorted by key, say), which is most of them.
type collection struct {
	ordered bool
	// why is the evidence for the choice: the schema description or the
	// vdb-site engine adapter (api/internal/config/<kind>.go).
	why string
}

// collections decides every array in the released kinds, keyed by its path:
// Kind.spec.<property>, with [] for a list element and {} for a map value.
// TestEveryArrayDecided fails on an array of a served kind until it is
// listed here, so the choice is made from the adapter, not left to a
// default. An unserved kind falls back to uniqueItems: a set when the
// schema says the items are unique, else a list.
//
// An array the generator carries as a JSON string (VexSource urls, whose
// items are a string or an object) needs no entry: its JSON keeps the order.
var collections = map[string]collection{
	"JailPolicy.spec.rules": {true, "schema: \"The rules, in order\"; order decides which rule gives the " +
		"headline reason. jail.go stores each rule's position (RuleInput.OrderIndex) and plans rules.order."},
	"JailPolicy.spec.rules[].flags": {false, "schema: every flag listed must hold, uniqueItems. " +
		"A conjunction has no order."},
	"RiskStrategy.spec.rules": {true, "schema: \"The rules in priority order: the first is tier 0\". " +
		"riskstrategy.go writes each rule at its priority (RiskPrioritizationRule.priority)."},

	"PackageList.spec.entries": {false, "packagelist.go keys entries by purl and reads them back sorted by " +
		"purl; an entry's effect does not depend on its position."},
	"PackageFirewallPolicy.spec.mirrors": {false, "pkgfirewall.go matches mirrors by (ecosystem, url) and " +
		"reads them back sorted by ecosystem then priority; the priority field, not the position, is the order."},

	"AiProvider.spec.models": {false, "aiprovider.go makes the AiModel rows the declared set and reads " +
		"them back sorted by id."},
	"AiFirewallPolicy.spec.baseline.exclude": {false, "baseline guardrail ids to leave out: membership only."},
	"AiFirewallPolicy.spec.providers": {false, "aifirewallpolicy.go keeps one pin per provider slug; " +
		"read back in the store's order."},
	"AiFirewallPolicy.spec.models": {false, "aifirewallpolicy.go expands entries to provider/slug keys; " +
		"read back in the store's order."},
	"AiFirewallPolicy.spec.guardrails": {false, "schema: reconciled by name. aifirewallpolicy.go reads " +
		"them back sorted by name; the priority field orders evaluation."},
	"AiFirewallPolicy.spec.capabilities{}.patterns": {false, "uniqueItems; aigovernance.go stores one " +
		"guardrail per pattern, trimmed and sorted."},
	"AiGovernanceList.spec.patterns": {false, "uniqueItems; aigovernance.go stores one guardrail per " +
		"pattern, trimmed and sorted."},

	"ApiSpecSource.spec.tags": {true, "apispecsource.go stores the tags as given and reads them back in " +
		"the same order; free labels, not deduplicated."},

	"VexSource.spec.scope.repos":          {false, "uniqueItems: the repositories in scope."},
	"VexSource.spec.scope.products":       {false, "uniqueItems: the TEA products in scope."},
	"VexSource.spec.trust.allowedAuthors": {false, "uniqueItems: an allow list of authors."},
	"VexSource.spec.git.paths": {false, "globs selecting documents; every match is read, so the order " +
		"of the globs changes nothing."},
}
