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
| `vulnetix_ai_firewall_policy` | AiFirewallPolicy, one per tenant | [minimal](examples/resources/vulnetix_ai_firewall_policy/minimal.tf), [advanced](examples/resources/vulnetix_ai_firewall_policy/advanced.tf) | [docs](docs/resources/ai_firewall_policy.md) |
| `vulnetix_ai_governance_list` | AiGovernanceList | [minimal](examples/resources/vulnetix_ai_governance_list/minimal.tf), [advanced](examples/resources/vulnetix_ai_governance_list/advanced.tf) | [docs](docs/resources/ai_governance_list.md) |
| `vulnetix_ai_provider` | AiProvider | [minimal](examples/resources/vulnetix_ai_provider/minimal.tf), [advanced](examples/resources/vulnetix_ai_provider/advanced.tf) | [docs](docs/resources/ai_provider.md) |
| `vulnetix_api_spec_source` | ApiSpecSource | [minimal](examples/resources/vulnetix_api_spec_source/minimal.tf), [advanced](examples/resources/vulnetix_api_spec_source/advanced.tf) | [docs](docs/resources/api_spec_source.md) |
| `vulnetix_jail_exemption` | JailExemption | [minimal](examples/resources/vulnetix_jail_exemption/minimal.tf), [advanced](examples/resources/vulnetix_jail_exemption/advanced.tf) | [docs](docs/resources/jail_exemption.md) |
| `vulnetix_jail_policy` | JailPolicy | [minimal](examples/resources/vulnetix_jail_policy/minimal.tf), [advanced](examples/resources/vulnetix_jail_policy/advanced.tf) | [docs](docs/resources/jail_policy.md) |
| `vulnetix_malware_ioc` | MalwareIoc | [minimal](examples/resources/vulnetix_malware_ioc/minimal.tf), [advanced](examples/resources/vulnetix_malware_ioc/advanced.tf) | [docs](docs/resources/malware_ioc.md) |
| `vulnetix_package_firewall_policy` | PackageFirewallPolicy, one per tenant | [minimal](examples/resources/vulnetix_package_firewall_policy/minimal.tf), [advanced](examples/resources/vulnetix_package_firewall_policy/advanced.tf) | [docs](docs/resources/package_firewall_policy.md) |
| `vulnetix_package_list` | PackageList | [minimal](examples/resources/vulnetix_package_list/minimal.tf), [advanced](examples/resources/vulnetix_package_list/advanced.tf) | [docs](docs/resources/package_list.md) |
| `vulnetix_quality_gate` | QualityGate, one per tenant | [minimal](examples/resources/vulnetix_quality_gate/minimal.tf), [advanced](examples/resources/vulnetix_quality_gate/advanced.tf) | [docs](docs/resources/quality_gate.md) |
| `vulnetix_remediation_policy` | RemediationPolicy, one per tenant | [minimal](examples/resources/vulnetix_remediation_policy/minimal.tf), [advanced](examples/resources/vulnetix_remediation_policy/advanced.tf) | [docs](docs/resources/remediation_policy.md) |
| `vulnetix_risk_strategy` | RiskStrategy | [minimal](examples/resources/vulnetix_risk_strategy/minimal.tf), [advanced](examples/resources/vulnetix_risk_strategy/advanced.tf) | [docs](docs/resources/risk_strategy.md) |
| `vulnetix_suppression` | Suppression | [minimal](examples/resources/vulnetix_suppression/minimal.tf), [advanced](examples/resources/vulnetix_suppression/advanced.tf) | [docs](docs/resources/suppression.md) |
| `vulnetix_vex_source` | VexSource | [minimal](examples/resources/vulnetix_vex_source/minimal.tf), [advanced](examples/resources/vulnetix_vex_source/advanced.tf) | [docs](docs/resources/vex_source.md) |

Each resource calls `GET`, `PUT` and `DELETE /config/objects/{kind}/{name}`. A PUT sends `{"document": <vulnetix.com/v1 document>, "manager": "<manager>"}`. Every request and response shape is in [internal/client/client.go](internal/client/client.go).

Kinds with a released schema but no engine adapter yet (the SIEM, cloud and feed kinds) are generated but not offered; each becomes a resource when the engine serves it. LicensePolicy and LicenseExceptions are validate-only: the server stores nothing for them, so they have no resource.

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

Most kinds merge a document into the stored object, and the tenant's row can hold values the configuration never set. So a refresh only reads the arguments you manage: an argument removed from the configuration stops being managed and keeps its last server value. Wholly declared kinds (ApiSpecSource, MalwareIoc, RemediationPolicy, RiskStrategy, Suppression, VexSource) are replaced instead, so a removed argument goes back to its default. Import reads every field the server holds.

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

Before the first release:

- Add the public half to the `vulnetix` namespace at registry.terraform.io (User Settings, Signing Keys).
- Publish the provider from this repository.
- Submit it to the OpenTofu registry (github.com/opentofu/registry) with the same key.

## License

[MPL-2.0](LICENSE)
