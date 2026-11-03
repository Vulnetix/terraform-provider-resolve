# terraform-provider-vulnetix

The `vulnetix/vulnetix` provider for Terraform and OpenTofu. Each resource is one `vulnetix.com/v1` config document, applied through the same server-side engine as `vulnetix config apply` and the console's Download as code. Resources are generated from the kind schemas in [config-schemas](https://github.com/Vulnetix/config-schemas).

Reference docs for every resource and data source are in [docs/](docs/), and are what the Terraform and OpenTofu registries publish.

## Install

Terraform 1.11 or later, or OpenTofu 1.11 or later (secrets use write-only arguments).

```hcl
terraform {
  required_providers {
    vulnetix = {
      source = "vulnetix/vulnetix"
    }
  }
}

provider "vulnetix" {} # token from VULNETIX_CONFIG_TOKEN

resource "vulnetix_quality_gate" "default" {
  severity      = "high"
  block_malware = true
}
```

Then `terraform init && terraform apply`, or `tofu init && tofu apply`. OpenTofu resolves the same `vulnetix/vulnetix` source from its own registry.

## Authentication

| Argument   | Environment variable    | Default                                | Notes |
|------------|-------------------------|----------------------------------------|-------|
| `token`    | `VULNETIX_CONFIG_TOKEN` | none                                   | A `vmt_` config management token, sent as `Authorization: Bearer`. Sensitive. |
| `endpoint` | `VULNETIX_ENDPOINT`     | `https://www.vulnetix.com/api/site/v1` | The console API. Config routes are under `/config/...`. |
| `manager`  | none                    | `terraform`                            | Recorded as the owner of each object, so `vulnetix config apply --prune` leaves Terraform's objects alone. |

A `vmt_` token acts as a service account of the tenant with a role and an allow list of kinds. An owner or admin, signed in, issues one with `POST /v1/config/tokens` (`name`, `principalUuid`, `role`, `kinds`, `expiresInDays`); the token is in that answer only. Rotate it with `POST /v1/config/tokens/{uuid}/rotate` and revoke it with `DELETE /v1/config/tokens/{uuid}`. Applying needs the `owner` or `admin` role and every kind the configuration manages.

## Resources

Every resource has a minimal example (only the required arguments) and an advanced one (every argument), generated from config-schemas `examples/<Kind>/minimal.yaml` and `advanced.yaml`:

| Resource | Kind | Examples | Docs |
|----------|------|----------|------|
| `vulnetix_agent_profile` | AgentProfile | [minimal](examples/resources/vulnetix_agent_profile/minimal.tf), [advanced](examples/resources/vulnetix_agent_profile/advanced.tf) | [docs](docs/resources/agent_profile.md) |
| `vulnetix_ai_firewall_policy` | AiFirewallPolicy, one per tenant | [minimal](examples/resources/vulnetix_ai_firewall_policy/minimal.tf), [advanced](examples/resources/vulnetix_ai_firewall_policy/advanced.tf) | [docs](docs/resources/ai_firewall_policy.md) |
| `vulnetix_ai_governance_list` | AiGovernanceList | [minimal](examples/resources/vulnetix_ai_governance_list/minimal.tf), [advanced](examples/resources/vulnetix_ai_governance_list/advanced.tf) | [docs](docs/resources/ai_governance_list.md) |
| `vulnetix_ai_provider` | AiProvider | [minimal](examples/resources/vulnetix_ai_provider/minimal.tf), [advanced](examples/resources/vulnetix_ai_provider/advanced.tf) | [docs](docs/resources/ai_provider.md) |
| `vulnetix_alert_trigger` | AlertTrigger | [minimal](examples/resources/vulnetix_alert_trigger/minimal.tf), [advanced](examples/resources/vulnetix_alert_trigger/advanced.tf) | [docs](docs/resources/alert_trigger.md) |
| `vulnetix_api_spec_source` | ApiSpecSource | [minimal](examples/resources/vulnetix_api_spec_source/minimal.tf), [advanced](examples/resources/vulnetix_api_spec_source/advanced.tf) | [docs](docs/resources/api_spec_source.md) |
| `vulnetix_aws_account` | AwsAccount | [minimal](examples/resources/vulnetix_aws_account/minimal.tf), [advanced](examples/resources/vulnetix_aws_account/advanced.tf) | [docs](docs/resources/aws_account.md) |
| `vulnetix_aws_connection` | AwsConnection, one per tenant | [minimal](examples/resources/vulnetix_aws_connection/minimal.tf), [advanced](examples/resources/vulnetix_aws_connection/advanced.tf) | [docs](docs/resources/aws_connection.md) |
| `vulnetix_campaign` | Campaign | [minimal](examples/resources/vulnetix_campaign/minimal.tf), [advanced](examples/resources/vulnetix_campaign/advanced.tf) | [docs](docs/resources/campaign.md) |
| `vulnetix_cloud_sensors` | CloudSensors | [minimal](examples/resources/vulnetix_cloud_sensors/minimal.tf), [advanced](examples/resources/vulnetix_cloud_sensors/advanced.tf) | [docs](docs/resources/cloud_sensors.md) |
| `vulnetix_compliance_target` | ComplianceTarget | [minimal](examples/resources/vulnetix_compliance_target/minimal.tf), [advanced](examples/resources/vulnetix_compliance_target/advanced.tf) | [docs](docs/resources/compliance_target.md) |
| `vulnetix_crew` | Crew | [minimal](examples/resources/vulnetix_crew/minimal.tf), [advanced](examples/resources/vulnetix_crew/advanced.tf) | [docs](docs/resources/crew.md) |
| `vulnetix_detection_rule` | DetectionRule | [minimal](examples/resources/vulnetix_detection_rule/minimal.tf), [advanced](examples/resources/vulnetix_detection_rule/advanced.tf) | [docs](docs/resources/detection_rule.md) |
| `vulnetix_detection_rule_feed` | DetectionRuleFeed | [minimal](examples/resources/vulnetix_detection_rule_feed/minimal.tf), [advanced](examples/resources/vulnetix_detection_rule_feed/advanced.tf) | [docs](docs/resources/detection_rule_feed.md) |
| `vulnetix_factory_budget` | FactoryBudget | [minimal](examples/resources/vulnetix_factory_budget/minimal.tf), [advanced](examples/resources/vulnetix_factory_budget/advanced.tf) | [docs](docs/resources/factory_budget.md) |
| `vulnetix_factory_managed_org` | FactoryManagedOrg | [minimal](examples/resources/vulnetix_factory_managed_org/minimal.tf), [advanced](examples/resources/vulnetix_factory_managed_org/advanced.tf) | [docs](docs/resources/factory_managed_org.md) |
| `vulnetix_factory_schedule` | FactorySchedule | [minimal](examples/resources/vulnetix_factory_schedule/minimal.tf), [advanced](examples/resources/vulnetix_factory_schedule/advanced.tf) | [docs](docs/resources/factory_schedule.md) |
| `vulnetix_factory_settings` | FactorySettings, one per tenant | [minimal](examples/resources/vulnetix_factory_settings/minimal.tf), [advanced](examples/resources/vulnetix_factory_settings/advanced.tf) | [docs](docs/resources/factory_settings.md) |
| `vulnetix_factory_trigger` | FactoryTrigger | [minimal](examples/resources/vulnetix_factory_trigger/minimal.tf), [advanced](examples/resources/vulnetix_factory_trigger/advanced.tf) | [docs](docs/resources/factory_trigger.md) |
| `vulnetix_factory_workflow` | FactoryWorkflow | [minimal](examples/resources/vulnetix_factory_workflow/minimal.tf), [advanced](examples/resources/vulnetix_factory_workflow/advanced.tf) | [docs](docs/resources/factory_workflow.md) |
| `vulnetix_forge_connection` | ForgeConnection | [minimal](examples/resources/vulnetix_forge_connection/minimal.tf), [advanced](examples/resources/vulnetix_forge_connection/advanced.tf) | [docs](docs/resources/forge_connection.md) |
| `vulnetix_git_hub_installation` | GitHubInstallation | [minimal](examples/resources/vulnetix_git_hub_installation/minimal.tf), [advanced](examples/resources/vulnetix_git_hub_installation/advanced.tf) | [docs](docs/resources/git_hub_installation.md) |
| `vulnetix_jail_exemption` | JailExemption | [minimal](examples/resources/vulnetix_jail_exemption/minimal.tf), [advanced](examples/resources/vulnetix_jail_exemption/advanced.tf) | [docs](docs/resources/jail_exemption.md) |
| `vulnetix_jail_policy` | JailPolicy | [minimal](examples/resources/vulnetix_jail_policy/minimal.tf), [advanced](examples/resources/vulnetix_jail_policy/advanced.tf) | [docs](docs/resources/jail_policy.md) |
| `vulnetix_library_item` | LibraryItem | [minimal](examples/resources/vulnetix_library_item/minimal.tf), [advanced](examples/resources/vulnetix_library_item/advanced.tf) | [docs](docs/resources/library_item.md) |
| `vulnetix_log_source` | LogSource | [minimal](examples/resources/vulnetix_log_source/minimal.tf), [advanced](examples/resources/vulnetix_log_source/advanced.tf) | [docs](docs/resources/log_source.md) |
| `vulnetix_malware_ioc` | MalwareIoc | [minimal](examples/resources/vulnetix_malware_ioc/minimal.tf), [advanced](examples/resources/vulnetix_malware_ioc/advanced.tf) | [docs](docs/resources/malware_ioc.md) |
| `vulnetix_member` | Member, owner token | [minimal](examples/resources/vulnetix_member/minimal.tf), [advanced](examples/resources/vulnetix_member/advanced.tf) | [docs](docs/resources/member.md) |
| `vulnetix_messaging_connection` | MessagingConnection | [minimal](examples/resources/vulnetix_messaging_connection/minimal.tf), [advanced](examples/resources/vulnetix_messaging_connection/advanced.tf) | [docs](docs/resources/messaging_connection.md) |
| `vulnetix_messaging_route` | MessagingRoute | [minimal](examples/resources/vulnetix_messaging_route/minimal.tf), [advanced](examples/resources/vulnetix_messaging_route/advanced.tf) | [docs](docs/resources/messaging_route.md) |
| `vulnetix_messaging_template` | MessagingTemplate | [minimal](examples/resources/vulnetix_messaging_template/minimal.tf), [advanced](examples/resources/vulnetix_messaging_template/advanced.tf) | [docs](docs/resources/messaging_template.md) |
| `vulnetix_notification_defaults` | NotificationDefaults, one per tenant | [minimal](examples/resources/vulnetix_notification_defaults/minimal.tf), [advanced](examples/resources/vulnetix_notification_defaults/advanced.tf) | [docs](docs/resources/notification_defaults.md) |
| `vulnetix_organization` | Organization, one per tenant, owner token | [minimal](examples/resources/vulnetix_organization/minimal.tf), [advanced](examples/resources/vulnetix_organization/advanced.tf) | [docs](docs/resources/organization.md) |
| `vulnetix_package_firewall_policy` | PackageFirewallPolicy, one per tenant | [minimal](examples/resources/vulnetix_package_firewall_policy/minimal.tf), [advanced](examples/resources/vulnetix_package_firewall_policy/advanced.tf) | [docs](docs/resources/package_firewall_policy.md) |
| `vulnetix_package_list` | PackageList | [minimal](examples/resources/vulnetix_package_list/minimal.tf), [advanced](examples/resources/vulnetix_package_list/advanced.tf) | [docs](docs/resources/package_list.md) |
| `vulnetix_quality_gate` | QualityGate, one per tenant | [minimal](examples/resources/vulnetix_quality_gate/minimal.tf), [advanced](examples/resources/vulnetix_quality_gate/advanced.tf) | [docs](docs/resources/quality_gate.md) |
| `vulnetix_registry_connection` | RegistryConnection | [minimal](examples/resources/vulnetix_registry_connection/minimal.tf), [advanced](examples/resources/vulnetix_registry_connection/advanced.tf) | [docs](docs/resources/registry_connection.md) |
| `vulnetix_remediation_policy` | RemediationPolicy, one per tenant | [minimal](examples/resources/vulnetix_remediation_policy/minimal.tf), [advanced](examples/resources/vulnetix_remediation_policy/advanced.tf) | [docs](docs/resources/remediation_policy.md) |
| `vulnetix_repository_settings` | RepositorySettings | [minimal](examples/resources/vulnetix_repository_settings/minimal.tf), [advanced](examples/resources/vulnetix_repository_settings/advanced.tf) | [docs](docs/resources/repository_settings.md) |
| `vulnetix_risk_strategy` | RiskStrategy | [minimal](examples/resources/vulnetix_risk_strategy/minimal.tf), [advanced](examples/resources/vulnetix_risk_strategy/advanced.tf) | [docs](docs/resources/risk_strategy.md) |
| `vulnetix_runtime_sensor_policy` | RuntimeSensorPolicy | [minimal](examples/resources/vulnetix_runtime_sensor_policy/minimal.tf), [advanced](examples/resources/vulnetix_runtime_sensor_policy/advanced.tf) | [docs](docs/resources/runtime_sensor_policy.md) |
| `vulnetix_sandbox_egress_policy` | SandboxEgressPolicy | [minimal](examples/resources/vulnetix_sandbox_egress_policy/minimal.tf), [advanced](examples/resources/vulnetix_sandbox_egress_policy/advanced.tf) | [docs](docs/resources/sandbox_egress_policy.md) |
| `vulnetix_sandbox_launch_profile` | SandboxLaunchProfile | [minimal](examples/resources/vulnetix_sandbox_launch_profile/minimal.tf), [advanced](examples/resources/vulnetix_sandbox_launch_profile/advanced.tf) | [docs](docs/resources/sandbox_launch_profile.md) |
| `vulnetix_saved_search` | SavedSearch | [minimal](examples/resources/vulnetix_saved_search/minimal.tf), [advanced](examples/resources/vulnetix_saved_search/advanced.tf) | [docs](docs/resources/saved_search.md) |
| `vulnetix_scanner_profile` | ScannerProfile, one per tenant | [minimal](examples/resources/vulnetix_scanner_profile/minimal.tf), [advanced](examples/resources/vulnetix_scanner_profile/advanced.tf) | [docs](docs/resources/scanner_profile.md) |
| `vulnetix_scheduled_report` | ScheduledReport | [minimal](examples/resources/vulnetix_scheduled_report/minimal.tf), [advanced](examples/resources/vulnetix_scheduled_report/advanced.tf) | [docs](docs/resources/scheduled_report.md) |
| `vulnetix_service_account` | ServiceAccount, owner token | [minimal](examples/resources/vulnetix_service_account/minimal.tf), [advanced](examples/resources/vulnetix_service_account/advanced.tf) | [docs](docs/resources/service_account.md) |
| `vulnetix_ssvc_policy` | SsvcPolicy, one per tenant | [minimal](examples/resources/vulnetix_ssvc_policy/minimal.tf), [advanced](examples/resources/vulnetix_ssvc_policy/advanced.tf) | [docs](docs/resources/ssvc_policy.md) |
| `vulnetix_suppression` | Suppression | [minimal](examples/resources/vulnetix_suppression/minimal.tf), [advanced](examples/resources/vulnetix_suppression/advanced.tf) | [docs](docs/resources/suppression.md) |
| `vulnetix_tea_access_policy` | TeaAccessPolicy | [minimal](examples/resources/vulnetix_tea_access_policy/minimal.tf), [advanced](examples/resources/vulnetix_tea_access_policy/advanced.tf) | [docs](docs/resources/tea_access_policy.md) |
| `vulnetix_tea_collection` | TeaCollection | [minimal](examples/resources/vulnetix_tea_collection/minimal.tf), [advanced](examples/resources/vulnetix_tea_collection/advanced.tf) | [docs](docs/resources/tea_collection.md) |
| `vulnetix_tea_component` | TeaComponent | [minimal](examples/resources/vulnetix_tea_component/minimal.tf), [advanced](examples/resources/vulnetix_tea_component/advanced.tf) | [docs](docs/resources/tea_component.md) |
| `vulnetix_tea_product` | TeaProduct | [minimal](examples/resources/vulnetix_tea_product/minimal.tf), [advanced](examples/resources/vulnetix_tea_product/advanced.tf) | [docs](docs/resources/tea_product.md) |
| `vulnetix_tea_publication_target` | TeaPublicationTarget | [minimal](examples/resources/vulnetix_tea_publication_target/minimal.tf), [advanced](examples/resources/vulnetix_tea_publication_target/advanced.tf) | [docs](docs/resources/tea_publication_target.md) |
| `vulnetix_vault_entry` | VaultEntry, owner token | [minimal](examples/resources/vulnetix_vault_entry/minimal.tf), [advanced](examples/resources/vulnetix_vault_entry/advanced.tf) | [docs](docs/resources/vault_entry.md) |
| `vulnetix_vex_source` | VexSource | [minimal](examples/resources/vulnetix_vex_source/minimal.tf), [advanced](examples/resources/vulnetix_vex_source/advanced.tf) | [docs](docs/resources/vex_source.md) |
| `vulnetix_wiz_connection` | WizConnection, one per tenant | [minimal](examples/resources/vulnetix_wiz_connection/minimal.tf), [advanced](examples/resources/vulnetix_wiz_connection/advanced.tf) | [docs](docs/resources/wiz_connection.md) |

Each resource calls `GET`, `PUT` and `DELETE /config/objects/{kind}/{name}`. A PUT sends `{"document": <vulnetix.com/v1 document>, "manager": "<manager>"}`. Every request and response shape is in [internal/client/client.go](internal/client/client.go).

Every kind with a released schema that the engine applies has a resource: 57 as of config-schemas v1.11.1, including the SIEM kinds (written in Guardian), the cloud and feed kinds, the organisation, vault and messaging kinds, and the Software Factory, Belai library, Pix sandbox and TEA kinds. A kind with a schema but no engine adapter is generated but not offered until the engine serves it. LicensePolicy and LicenseExceptions are validate-only: the server stores nothing for them, so they have no resource.

A spec property whose name the resource reserves (`name`, `labels`, `annotations`, `status`, `id` and Terraform's meta-arguments) is renamed: Campaign `name` is `display_name`, TeaProduct `name` is `product_name`, TeaComponent `name` is `component_name`, VaultEntry `name` is `entry_name` and `provider` is `provider_slug`, and FactoryBudget `provider` is `model_providers`, and RepositorySettings `labels` (the repository as an asset) is `asset_labels`. The document key is unchanged.

### Status, phases and values shown once

Every resource has a computed `status` block:

| Attribute | Meaning |
|-----------|---------|
| `status.phase` | `Applied` for most objects. `Pending` for a Member whose invitation is not accepted. `AwaitingAuthorization` for a messaging connection, such as Slack, that a person must authorize in a browser. `Unmanaged` when no manager owns the object. Terraform does not wait for a person: Pending and AwaitingAuthorization are successful applies, and a refresh shows the phase when it changes. |
| `status.authorization_url` | Where to authorize, while the phase is `AwaitingAuthorization`. |
| `status.extra` | Other fields the kind reports, such as `invitationExpires`. |
| `status.once` | Sensitive. Credentials the server shows only when it creates the object: a ServiceAccount's `apiKey` and `secret`, a messaging connection's `signingSecret`, a log source's ingest token. Terraform has no write-only computed attribute, so they are kept in the state from the create. Keep the state encrypted. They cannot be read again; replace the object to get new ones. |

A SIEM kind (SavedSearch, DetectionRule, AlertTrigger, ScheduledReport, LogSource) is written in Guardian after the config database. When Guardian refuses, the server answers 502 and records nothing, and the apply fails with Guardian's reason. Apply again once it is fixed.

### Owner-only kinds

Organization, Member, ServiceAccount and VaultEntry decide who can act for the tenant. A token can carry them only with the owner role, and the token must list them among its kinds. An apply with any other token fails with a 403 that says so. Organization, Member and ServiceAccount also need an owner or admin to apply them.

## Data sources

| Data source | API |
|-------------|-----|
| `vulnetix_kinds` | `GET /config/kinds`: the kinds the server applies, with their schemas. |
| `vulnetix_ai_models` | `GET /ai-firewall/models`: the AI firewall model catalogue, filtered by provider or search. |
| `vulnetix_vex_statements` | `GET /vex/statements`: every VEX statement the tenant holds, in every status, with its products and purls, filtered by status, vulnerability, purl or repository. A token needs the `VexSource` or `Suppression` kind. |
| `vulnetix_github_installation_repos` | `GET /integrations/github/installations/repos`: the repositories the tenant's GitHub App installations reach, from Vulnetix's cache. A token needs the `github-installations` scope. |
| `vulnetix_detection_rule_feed` | `GET /detection-rule-feeds`, `GET /detection-rule-feeds/{uuid}`: one feed and its entries (Pro plan). |
| `vulnetix_virtual_patches` | The same feed routes: the entries whose `rule_table` is `virtual_patch`. |

## Import

Import an object by its name (`metadata.name`). One-per-tenant kinds are named `default`:

```sh
terraform import vulnetix_quality_gate.default default
terraform import vulnetix_jail_policy.payments_api payments-api
```

Each resource's `import.sh` beside its examples has the command for its example object.

## How state is read

Most kinds merge a document into the stored object, and the tenant's row can hold values the configuration never set. So a refresh only reads the arguments you manage: an argument removed from the configuration stops being managed and keeps its last server value. Only these kinds merge: AiFirewallPolicy, AiProvider, AwsAccount, AwsConnection, CloudSensors, JailExemption, JailPolicy, Organization, PackageFirewallPolicy, PackageList and QualityGate. Every other kind is wholly declared and is replaced instead, so a removed argument goes back to its default. Import reads every field the server holds.

## Secrets

A `valueFrom` property (listed in the schema's `x-vulnetix.secrets`) is generated as two arguments:

- `<name>_wo`: write-only (Terraform or OpenTofu 1.11 or later). It is sent in the document and never stored in the plan or the state. The server never returns it either.
- `<name>_wo_version`: change this number to send the value again. Terraform cannot see a change to a write-only value.

```hcl
variable "vllm_api_key" {
  type      = string
  sensitive = true
  ephemeral = true
}

resource "vulnetix_ai_provider" "vllm" {
  name               = "acme-vllm"
  auth               = "bearer"
  base_url           = "https://llm.acme.example"
  protocol           = "openai-chat"
  api_key_wo         = var.vllm_api_key
  api_key_wo_version = 1
}
```

## OpenTofu

The provider is the same binary for both tools and is published to both registries. Everything above applies to OpenTofu 1.11 or later: write-only arguments, ephemeral variables and import. CI runs the acceptance suite under both.

## Generated code

Resources, examples and docs are generated from the config-schemas version `go.mod` requires:

```sh
go get github.com/vulnetix/config-schemas@vX.Y.Z
go generate ./...   # cmd/genresources, then tfplugindocs
```

`cmd/genresources` writes `internal/provider/zz_generated_<kind>.go` for each kind (a `kindSpec`: the framework schema attributes plus the field map between attributes and spec keys) and `examples/resources/vulnetix_<kind>/{minimal,advanced}.tf` and `import.sh`. One generic resource (`internal/provider/resource_kind.go`) serves every kind. It reads each schema's `x-vulnetix` block for `validateOnly`, `singleton`, `dependsOn`, `secrets` and `celFields`. Two lists stay by hand, each checked by a test:

- [`cmd/genresources/served.go`](cmd/genresources/served.go): the kinds the engine serves (vdb-site `handler/config_engine_kinds.go`). Add a kind there when its adapter ships.
- [`cmd/genresources/collections.go`](cmd/genresources/collections.go): whether each array of a served kind is a list or a set, decided from how its engine adapter stores it.

`templates/resources.md.tmpl` renders the Minimal, Advanced and Import sections on every resource page. CI fails when the generated code, examples or docs are not what `go generate` writes.

## Development

```sh
go build ./... && go vet ./... && go test ./...
```

The unit tests run a real `terraform` binary (1.11 or later, from `PATH`) against an `httptest` fake of the config API. For every served kind they create from the minimal example, update to the advanced one, import each object, destroy, and check that no secret reaches the state.

Acceptance tests (`TestAcc<Kind>`, one per served kind) run the same lifecycle against a real server and change its configuration, singletons included. Point them at a test tenant only:

```sh
TF_ACC=1 VULNETIX_ENDPOINT=https://www.vulnetix.com/api/site/v1 \
VULNETIX_CONFIG_TOKEN=vmt_... go test ./internal/provider -run TestAcc
```

Under OpenTofu, add `TF_ACC_TERRAFORM_PATH=$(command -v tofu) TF_ACC_PROVIDER_NAMESPACE=vulnetix`. They skip unless `TF_ACC`, `VULNETIX_ENDPOINT` and `VULNETIX_CONFIG_TOKEN` are all set. In CI they run when the repository has the `VULNETIX_ENDPOINT` and `VULNETIX_CONFIG_TOKEN` secrets.

## Releasing

Push a `vX.Y.Z` tag. [`.github/workflows/release.yml`](.github/workflows/release.yml) runs goreleaser on the `vulnetix-local` runners. It publishes a GitHub release in the registry layout: zips per platform, `_manifest.json`, `_SHA256SUMS` and a detached GPG `_SHA256SUMS.sig`. registry.terraform.io and the OpenTofu registry both read that release.

Repository secrets the workflow needs:

| Secret | What it is |
|--------|------------|
| `GPG_PRIVATE_KEY` | ASCII-armoured GPG private key (`gpg --armor --export-secret-keys <id>`). RSA or DSA; the Terraform registry does not accept ECC keys. |
| `PASSPHRASE` | That key's passphrase. |

The `justfile` does this without prompts. The key ID and passphrase are kept in the gitignored `.env.production` (mode 600):

| Recipe | What it does |
|--------|--------------|
| `just gpg-key` | Creates an RSA 4096 signing key with a random passphrase (non-interactive `gpg --batch`) and writes `GPG_FINGERPRINT` and `PASSPHRASE` to `.env.production`. |
| `just secrets` | Sets the `GPG_PRIVATE_KEY` and `PASSPHRASE` repository secrets with `gh secret set`, piping the values so they never appear in argv or the terminal. |
| `just public-key` | Prints the armoured public key, for the registry signing-key settings. |
| `just release X.Y.Z` | Tags `origin/main` as `vX.Y.Z` and pushes the tag, which starts the release workflow. |

Before the first release:

- Add the public half to the `vulnetix` namespace at registry.terraform.io (User Settings, Signing Keys).
- Publish the provider from this repository.
- Submit it to the OpenTofu registry (github.com/opentofu/registry) with the same key.

## License

[MPL-2.0](LICENSE)
