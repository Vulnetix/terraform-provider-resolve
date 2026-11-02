package main

// served lists the kinds the config engine applies: the ones registered in
// vdb-site api/internal/handler/config_engine_kinds.go configKinds() on
// main. A kind with a released schema but no engine adapter is still
// generated, so its resource is ready the day the adapter ships, but it goes
// into the unserved list and the provider does not offer it: a PUT would
// answer unknown kind.
//
// When an adapter lands, add its kind here and run go generate ./...
//
// Validate-only kinds (x-vulnetix.validateOnly) are registered with the
// engine too, but get no resource at all; see validateOnly.
var served = map[string]bool{
	"AiFirewallPolicy":      true,
	"AiGovernanceList":      true,
	"AiProvider":            true,
	"ApiSpecSource":         true,
	"JailExemption":         true,
	"JailPolicy":            true,
	"LicensePolicy":         true,
	"MalwareIoc":            true,
	"PackageFirewallPolicy": true,
	"PackageList":           true,
	"QualityGate":           true,
	"VexSource":             true,
}
