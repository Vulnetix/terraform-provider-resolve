package main

// served lists the kinds the config engine applies: the ones registered in
// vdb-site api/internal/handler/config_engine_kinds.go configKinds() on
// main (and the helpers it appends: the K2 and K3 kinds, the feed kinds and
// the SIEM kinds in config.SiemKindNames). A kind with a released schema but
// no engine adapter is still generated, so its resource is ready the day the
// adapter ships, but it goes into the unserved list and the provider does
// not offer it: a PUT would answer unknown kind.
//
// When an adapter lands, add its kind here and run go generate ./...
//
// Validate-only kinds (x-vulnetix.validateOnly) are registered with the
// engine too, but get no resource at all: the server stores nothing for them.
var served = map[string]bool{
	"AgentProfile":          true,
	"AiFirewallPolicy":      true,
	"AiGovernanceList":      true,
	"AiProvider":            true,
	"AlertTrigger":          true,
	"ApiSpecSource":         true,
	"AwsAccount":            true,
	"AwsConnection":         true,
	"Campaign":              true,
	"CloudSensors":          true,
	"Crew":                  true,
	"DetectionRuleFeed":     true,
	"DetectionRule":         true,
	"FactoryBudget":         true,
	"FactoryManagedOrg":     true,
	"FactorySchedule":       true,
	"FactorySettings":       true,
	"FactoryTrigger":        true,
	"FactoryWorkflow":       true,
	"ForgeConnection":       true,
	"GitHubInstallation":    true,
	"JailExemption":         true,
	"JailPolicy":            true,
	"LibraryItem":           true,
	"LicenseExceptions":     true,
	"LicensePolicy":         true,
	"LogSource":             true,
	"MalwareIoc":            true,
	"Member":                true,
	"MessagingConnection":   true,
	"MessagingRoute":        true,
	"MessagingTemplate":     true,
	"NotificationDefaults":  true,
	"Organization":          true,
	"PackageFirewallPolicy": true,
	"PackageList":           true,
	"QualityGate":           true,
	"RegistryConnection":    true,
	"RemediationPolicy":     true,
	"RepositorySettings":    true,
	"RiskStrategy":          true,
	"RuntimeSensorPolicy":   true,
	"SandboxEgressPolicy":   true,
	"SandboxLaunchProfile":  true,
	"SavedSearch":           true,
	"ScannerProfile":        true,
	"ScheduledReport":       true,
	"ServiceAccount":        true,
	"SsvcPolicy":            true,
	"Suppression":           true,
	"TeaAccessPolicy":       true,
	"TeaCollection":         true,
	"TeaComponent":          true,
	"TeaProduct":            true,
	"TeaPublicationTarget":  true,
	"VaultEntry":            true,
	"VexSource":             true,
	"WizConnection":         true,
}

// ownerOnly are the kinds that decide who can act for the tenant: its
// settings, its people, its machine credentials and its secrets (vdb-site
// handler/config_tokens.go ownerTokenKinds). A token may carry them only
// with the owner role.
var ownerOnly = map[string]bool{
	"Organization": true, "Member": true, "ServiceAccount": true, "VaultEntry": true,
}
