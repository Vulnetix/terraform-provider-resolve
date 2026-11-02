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

	// K2 and K3 kinds. Arrays the schema marks uniqueItems need no entry (a set
	// compares without order); these may repeat, so each is decided here.
	//
	// Lists keep the order the document gives, which the store reads back
	// unchanged because it holds the document as JSON.
	"SsvcPolicy.spec.mapping":                 {true, "schema: when several mapping rules match a cell the last one wins, so order is meaning; ssvcpolicy.go keeps the rules as written."},
	"SsvcPolicy.spec.overrides":               {true, "schema: the first override that holds sets the outcome, so order is meaning; ssvcpolicy.go keeps the rules as written."},
	"Campaign.spec.done.criteria":             {true, "campaign.go stores the done criteria as versioned JSON and reads them back in the order written."},
	"Crew.spec.members":                       {true, "crew.go keeps the members as the library item stores them, in the order written; a crew is composed in that order."},
	"FactoryWorkflow.spec.graph.nodes":        {true, "factoryworkflow.go stores the graph as the version's canonical JSON; the console normaliser keeps node and edge order."},
	"FactoryWorkflow.spec.graph.edges":        {true, "factoryworkflow.go stores the graph as the version's canonical JSON; the console normaliser keeps node and edge order."},
	"FactorySettings.spec.budget.alerts":      {true, "factorysettings.go stores the settings document as written."},
	"FactorySettings.spec.pix.maintenance":    {true, "factorysettings.go stores the settings document as written; windows are listed in the order given."},
	"FactorySettings.spec.schedules.blackout": {true, "factorysettings.go stores the settings document as written; windows are listed in the order given."},
	"MessagingTemplate.spec.content.actions":  {true, "messagingtemplate.go keeps the buttons in the order written; they render in that order."},
	"MessagingTemplate.spec.content.fields":   {true, "messagingtemplate.go keeps the fields in the order written; they render in that order."},
	"SandboxLaunchProfile.spec.extensions":    {true, "sandboxlaunchprofile.go lists each extension once, in the order written."},
	"ScannerProfile.spec.custom":              {true, "the scanner profile keeps its custom scanners in the order written."},

	// Sets: membership, or the server reads them back in its own order.
	"AgentProfile.spec.files":                      {false, "agentprofile.go reads the files back in path order."},
	"LibraryItem.spec.files":                       {false, "libraryitem.go reads the files back in path order, each path once."},
	"Campaign.spec.contributors":                   {false, "campaign.go sorts contributors by key (sortContributors)."},
	"DetectionRule.spec.suppressions":              {false, "membership: each suppression applies on its own."},
	"DetectionRule.spec.threshold.by":              {false, "membership: the fields a threshold groups by."},
	"DetectionRuleFeed.spec.entries":               {false, "detectionrulefeed.go says their order in a document does not matter, and sorts them by key."},
	"FactoryBudget.spec.alertAt":                   {false, "membership: the percentages at which to alert."},
	"FactorySettings.spec.approve.sensitive_paths": {false, "membership: globs of paths that need approval."},
	"FactorySettings.spec.pix.maintenance[].days":  {false, "membership: the days of the week."},
	"FactorySettings.spec.vault.grant.entries":     {false, "membership: the vault entries granted."},
	"MessagingRoute.spec.targets":                  {false, "messagingroute.go sorts the targets by connection: their order means nothing to delivery."},
	"MessagingRoute.spec.targets[].options.labels": {false, "membership: labels put on the issue."},
	"MessagingTemplate.spec.content.issue.labels":  {false, "membership: labels put on the issue."},
	"SandboxLaunchProfile.spec.packs":              {false, "sandboxlaunchprofile.go trims the pack ids and lists each once: membership."},
	"TeaAccessPolicy.spec.sharedWith":              {false, "membership: the organisations a product is shared with."},
	"VaultEntry.spec.grants":                       {false, "vaultentry.go stores grants with their ids lower case, sorted and unique; an entry is granted to a set."},
	"VaultEntry.spec.grants[].ids":                 {false, "vaultentry.go: grant ids are lower case, sorted and unique."},
	"FactoryWorkflow.spec.graph.scope.orgs":        {false, "membership: the organisations a workflow runs for."},
	"FactoryWorkflow.spec.graph.scope.repos":       {false, "membership: the repositories a workflow runs for."},

	// FactoryTrigger filters are membership tests: a trigger fires when the
	// event matches any (or all, for labels.all) of the values, so their
	// order means nothing.
	"FactoryTrigger.spec.categories":               {false, "filter: membership."},
	"FactoryTrigger.spec.conditions.actions":       {false, "filter: membership."},
	"FactoryTrigger.spec.conditions.authors.allow": {false, "filter: membership."},
	"FactoryTrigger.spec.conditions.authors.deny":  {false, "filter: membership."},
	"FactoryTrigger.spec.conditions.baseBranch":    {false, "filter: membership."},
	"FactoryTrigger.spec.conditions.conclusion":    {false, "filter: membership."},
	"FactoryTrigger.spec.conditions.environment":   {false, "filter: membership."},
	"FactoryTrigger.spec.conditions.headBranch":    {false, "filter: membership."},
	"FactoryTrigger.spec.conditions.labels.all":    {false, "filter: every label must be present, in any order."},
	"FactoryTrigger.spec.conditions.labels.any":    {false, "filter: membership."},
	"FactoryTrigger.spec.conditions.packageType":   {false, "filter: membership."},
	"FactoryTrigger.spec.conditions.paths":         {false, "filter: a path matching any glob; order means nothing."},
	"FactoryTrigger.spec.conditions.pathsIgnore":   {false, "filter: a path matching any glob; order means nothing."},
	"FactoryTrigger.spec.conditions.reviewState":   {false, "filter: membership."},
	"FactoryTrigger.spec.conditions.severity":      {false, "filter: membership."},
	"FactoryTrigger.spec.decision":                 {false, "filter: membership."},
	"FactoryTrigger.spec.ecosystems":               {false, "filter: membership."},
	"FactoryTrigger.spec.events":                   {false, "filter: membership."},
	"FactoryTrigger.spec.kinds":                    {false, "filter: membership."},
	"FactoryTrigger.spec.providers":                {false, "filter: membership."},
	"FactoryTrigger.spec.sandboxes":                {false, "filter: membership."},
	"FactoryTrigger.spec.tables":                   {false, "filter: membership."},
	"FactoryTrigger.spec.triggers":                 {false, "filter: membership."},

	"VexSource.spec.scope.repos":          {false, "uniqueItems: the repositories in scope."},
	"VexSource.spec.scope.products":       {false, "uniqueItems: the TEA products in scope."},
	"VexSource.spec.trust.allowedAuthors": {false, "uniqueItems: an allow list of authors."},
	"VexSource.spec.git.paths": {false, "globs selecting documents; every match is read, so the order " +
		"of the globs changes nothing."},
}

// renames moves a top-level spec property whose snake case name the
// resource reserves (name, labels, annotations, status, id and Terraform's
// meta-arguments) to another attribute name, keyed by Kind.spec.<property>.
// The document key is unchanged. A collision with no entry here fails
// generation.
var renames = map[string]string{
	// The resource's own name is metadata.name; the campaign's title, as the
	// console shows it, is spec.name.
	"Campaign.spec.name":             "display_name",
	"FactoryBudget.spec.provider":    "model_providers",
	"RepositorySettings.spec.labels": "asset_labels",
	"TeaComponent.spec.name":         "component_name",
	"TeaProduct.spec.name":           "product_name",
	"VaultEntry.spec.name":           "entry_name",
	"VaultEntry.spec.provider":       "provider_slug",
}
