# terraform-provider-vulnetix

The `vulnetix/vulnetix` provider for Terraform and OpenTofu. Each resource is one `vulnetix.com/v1` config document, applied through the same server-side engine as `vulnetix config apply`.

```hcl
terraform {
  required_providers {
    vulnetix = { source = "vulnetix/vulnetix" }
  }
}

provider "vulnetix" {} # token from VULNETIX_CONFIG_TOKEN

resource "vulnetix_quality_gate" "default" {
  severity      = "high"
  block_malware = true
}
```

Reference docs for every resource and data source are in [docs/](docs/). Examples are in [examples/](examples/).

## Provider configuration

| Argument   | Environment variable    | Default                                | Notes |
|------------|-------------------------|----------------------------------------|-------|
| `endpoint` | `VULNETIX_ENDPOINT`     | `https://www.vulnetix.com/api/site/v1` | The console API. Config routes are under `/config/...`. |
| `token`    | `VULNETIX_CONFIG_TOKEN` | none                                   | A `vmt_` config management token, sent as `Authorization: Bearer`. Sensitive. A legacy org API key `<orgId>:<hmac>` is sent as `Authorization: ApiKey`. |
| `manager`  | none                    | `terraform`                            | Recorded as the owner of each object, so `vulnetix config apply --prune` leaves Terraform's objects alone. |

## Resources and data sources

| Name | API |
|------|-----|
| `vulnetix_quality_gate` | `GET/PUT/DELETE /config/objects/QualityGate/default` |
| `data.vulnetix_kinds` | `GET /config/kinds` |
| `data.vulnetix_detection_rule_feed` | `GET /detection-rule-feeds`, `GET /detection-rule-feeds/{uuid}` |
| `data.vulnetix_virtual_patches` | The same feed routes. Returns the entries whose `rule_table` is `virtual_patch`. |

PUT sends `{"document": <vulnetix.com/v1 document>, "manager": "<manager>"}`. Every request and response shape is in [internal/client/client.go](internal/client/client.go).

### How state is read

The engine merges a document into the stored object. The tenant's row can also hold values the configuration never set. So `Read` only refreshes the arguments you manage. An argument removed from the configuration stops being managed and keeps its last server value. Import is the exception: `terraform import` reads every field the server holds.

Import an object by its name:

```sh
terraform import vulnetix_quality_gate.default default
```

### Secrets

A kind property defined as a `valueFrom` reference (`common.schema.json#/$defs/valueFrom`) is generated as two arguments:

- `<name>_wo`: a write-only argument (Terraform 1.11 or later). It is sent in the document as a plain value and never stored in state.
- `<name>_wo_version`: change this number to send the value again. Terraform cannot see a change to a write-only value.

QualityGate has no secrets. The generator tests cover this path with the test kind in `cmd/genresources/testdata/schemas/Widget.schema.json`.

## Generated resources

Resources are generated from the JSON Schemas in [config-schemas](https://github.com/Vulnetix/config-schemas), at the version `go.mod` requires:

```sh
go get github.com/vulnetix/config-schemas@vX.Y.Z
go generate ./...   # cmd/genresources, then tfplugindocs
```

`cmd/genresources` writes `internal/provider/zz_generated_<kind>.go` for each kind. Each file is a `kindSpec`: the framework schema attributes plus the field map between attributes and spec keys. One generic resource (`internal/provider/resource_kind.go`) serves every kind. Property names become snake case (`within30Days` becomes `within_30_days`). Shapes that HCL cannot type exactly, such as a `oneOf` of differing types, become a JSON string compared semantically.

CI fails if the generated files or docs are not what `go generate` writes.

## Development

```sh
go build ./... && go vet ./... && go test ./...
```

The unit tests run a real `terraform` binary (from `PATH`, or downloaded) against an `httptest` fake of the config API.

Acceptance tests run against a real server and change its configuration. Destroy deletes the tenant's QualityGate, so point them at a test tenant only:

```sh
TF_ACC=1 VULNETIX_ENDPOINT=https://www.vulnetix.com/api/site/v1 \
VULNETIX_CONFIG_TOKEN=vmt_... go test ./internal/provider -run TestAcc
```

They skip unless `TF_ACC`, `VULNETIX_ENDPOINT` and `VULNETIX_CONFIG_TOKEN` are all set.

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
