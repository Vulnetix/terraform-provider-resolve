---
page_title: "vulnetix_vex_source Resource - vulnetix"
subcategory: ""
description: |-
  VEX documents read back as suppression rules. A not_affected statement (CycloneDX not_affected or false_positive) becomes a suppression with the document's justification, timestamp and author; fixed or resolved marks the finding remediated; affected, in_triage and exploitable are recorded as triage context; under_investigation is an annotation. Products are matched to findings by purl. A later statement wins. Suppressions a source creates are owned by it: they go when the statement leaves the document, and cannot be deleted by hand. Name exactly one of tea, s3, git, path or urls.
  An object of this kind can reference VaultEntry objects. Reference those resources' attributes, or add depends_on, so Terraform applies them first.
  Manages the VexSource config kind through /config/objects/VexSource/{name}. Applies are recorded with the provider's manager, so a GitOps vulnetix config apply --prune leaves this object alone.
  Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (such as ApiSpecSource, MalwareIoc and VexSource) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object.
---

# vulnetix_vex_source (Resource)

VEX documents read back as suppression rules. A not_affected statement (CycloneDX not_affected or false_positive) becomes a suppression with the document's justification, timestamp and author; fixed or resolved marks the finding remediated; affected, in_triage and exploitable are recorded as triage context; under_investigation is an annotation. Products are matched to findings by purl. A later statement wins. Suppressions a source creates are owned by it: they go when the statement leaves the document, and cannot be deleted by hand. Name exactly one of tea, s3, git, path or urls.

An object of this kind can reference `VaultEntry` objects. Reference those resources' attributes, or add `depends_on`, so Terraform applies them first.

Manages the `VexSource` config kind through `/config/objects/VexSource/{name}`. Applies are recorded with the provider's `manager`, so a GitOps `vulnetix config apply --prune` leaves this object alone.

Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (such as `ApiSpecSource`, `MalwareIoc` and `VexSource`) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object.

## Example Usage

### Minimal

The arguments the kind requires, and nothing else.

```terraform
# Only the arguments the kind requires. From config-schemas examples/VexSource/minimal.yaml.

resource "vulnetix_vex_source" "local_vex" {
  name = "local-vex"
  path = {
    path = ".vulnetix/vex.openvex.json"
  }
}
```

### Advanced

Every argument the kind supports. A secret is a write-only `*_wo` argument fed from an ephemeral variable, so it never reaches the plan or the state; change `*_wo_version` to send a new value.

```terraform
# Every argument the kind supports. From config-schemas examples/VexSource/advanced.yaml.

resource "vulnetix_vex_source" "vendor_tea_product" {
  name     = "vendor-tea-product"
  format   = "auto"
  on_error = "keep-last"
  refresh = {
    interval = "6h"
  }
  scope = {
    org = true
  }
  tea = {
    product = "7f3c1a52-1b8e-4d0e-9c3f-2a6b8e5d4c10"
    server  = "https://tea.acme.example"
  }
  trust = {
    allowed_authors   = ["security@acme.example"]
    require_signature = true
  }
  labels = {
    team = "product-security"
  }
  annotations = {
    ticket = "https://jira.example.com/browse/SEC-2217"
  }
}

resource "vulnetix_vex_source" "vendor_tea_release" {
  name   = "vendor-tea-release"
  format = "cyclonedx"
  refresh = {
    interval = "1d"
  }
  scope = {
    products = ["urn:tei:uuid:products.acme.example:7f3c1a52-1b8e-4d0e-9c3f-2a6b8e5d4c10"]
  }
  tea = {
    product_release = "urn:tei:uuid:products.acme.example:4b2e9d17-6a3c-4f81-b5d0-8e1c7a9f2d36"
    server          = "https://tea.acme.example"
  }
}

resource "vulnetix_vex_source" "vendor_tea_collection" {
  name = "vendor-tea-collection"
  refresh = {
    webhook = true
  }
  scope = {
    repos = ["acme/payments-api"]
  }
  tea = {
    collection = "2c8f4e61-9d7a-4b35-a0e2-5f6c3b8d1e94"
    server     = "https://tea.vendor.example.com"
  }
}

resource "vulnetix_vex_source" "vex_bucket" {
  name     = "vex-bucket"
  format   = "openvex"
  on_error = "keep-last"
  refresh = {
    webhook = true
  }
  s_3 = {
    bucket      = "acme-vex-statements"
    external_id = "acme-vex-read-7f3c1a52"
    prefix      = "openvex/"
    region      = "us-east-1"
    role_arn    = "arn:aws:iam::123456789012:role/VulnetixVexRead"
  }
}

resource "vulnetix_vex_source" "vex_r2_object" {
  name   = "vex-r2-object"
  format = "openvex"
  refresh = {
    interval = "30m"
  }
  s_3 = {
    bucket         = "acme-vex-public"
    credential_ref = "r2-vex-reader"
    endpoint       = "https://0123456789abcdef0123456789abcdef.r2.cloudflarestorage.com"
    key            = "products/ledger/ledger.openvex.json"
  }
}

resource "vulnetix_vex_source" "platform_repo" {
  name = "platform-repo"
  git = {
    credential_ref = "github-advisories-token"
    paths          = ["vex/**/*.openvex.json", "vex/**/*.cdx.json"]
    ref            = "main"
    url            = "https://github.com/acme/security-advisories.git"
  }
  scope = {
    products = ["7f3c1a52-1b8e-4d0e-9c3f-2a6b8e5d4c10"]
    repos    = ["acme/payments-api", "acme/ledger"]
  }
}

resource "vulnetix_vex_source" "local_vex" {
  name   = "local-vex"
  format = "cyclonedx"
  path = {
    path = ".vulnetix/vex.cdx.json"
  }
}

resource "vulnetix_vex_source" "upstream_urls" {
  name     = "upstream-urls"
  on_error = "clear"
  refresh = {
    interval = "1d"
  }
  urls = jsonencode(["https://downloads.example.org/vex/product.openvex.json", {
    auth = {
      credentialRef = "vendor-vex-token"
      header        = "Authorization"
    }
    sha256 = "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"
    url    = "https://vendor.example.com/vex/all.cdx.json"
  }])
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The object's name (`metadata.name`), unique per kind within the tenant. Changing it replaces the object.

### Optional

- `annotations` (Map of String) Free key/value notes for people and tools (`metadata.annotations`), such as a ticket link.
- `format` (String) openvex: OpenVEX 0.2.0. cyclonedx: CycloneDX VEX 1.4 to 1.7 (vulnerabilities[].analysis). auto: decided per document. One of `openvex`, `cyclonedx`, `auto`. Schema default: `auto`.
- `git` (Attributes) Documents in a git repository over HTTPS (GitHub, GitHub Enterprise Server, GitLab, Gitea, Forgejo or any git server). A matching ForgeConnection is used when there is one. (see [below for nested schema](#nestedatt--git))
- `labels` (Map of String) Free key/value labels (`metadata.labels`). Stored with the object, never interpreted by the server.
- `on_error` (String) When a fetch fails: keep-last keeps the suppressions from the last good fetch; clear removes them. One of `keep-last`, `clear`. Schema default: `keep-last`.
- `path` (Attributes) A file or directory, absolute or relative to the config or repository root. Read by the CLI only: at scan time it applies locally, and config apply uploads a snapshot so server views match. (see [below for nested schema](#nestedatt--path))
- `refresh` (Attributes) When the documents are fetched again: every interval, or when a webhook says they changed. (see [below for nested schema](#nestedatt--refresh))
- `s_3` (Attributes) Documents in an S3 bucket, read through a cross-account role whose external id is the tenant uuid (the cloud onboarding trust pattern), or from an S3-compatible endpoint (R2, MinIO). Name exactly one of prefix or key. (see [below for nested schema](#nestedatt--s_3))
- `scope` (Attributes) Which findings the statements may suppress. Left out, the whole tenant. (see [below for nested schema](#nestedatt--scope))
- `tea` (Attributes) VEX artifacts from a TEA collection, on the Vulnetix TEA publisher or any TEA server. Name exactly one of product, productRelease or collection. (see [below for nested schema](#nestedatt--tea))
- `trust` (Attributes) (see [below for nested schema](#nestedatt--trust))
- `urls` (String) Plain HTTPS URLs, each fetched with ETag and If-Modified-Since. A JSON-encoded value; use jsonencode().

### Read-Only

- `id` (String) The object's name, the same as `name`.

<a id="nestedatt--git"></a>
### Nested Schema for `git`

Required:

- `paths` (Set of String) Globs, relative to the repository root, selecting the documents.
- `url` (String)

Optional:

- `credential_ref` (String) The metadata.name of a VaultEntry holding a credential. The value is read from the tenant's vault when it is needed and is never copied into the referring object.
- `ref` (String) A branch, tag or commit. Defaults to the default branch.


<a id="nestedatt--path"></a>
### Nested Schema for `path`

Required:

- `path` (String)


<a id="nestedatt--refresh"></a>
### Nested Schema for `refresh`

Optional:

- `interval` (String) A whole number of minutes, hours or days: 30m, 6h, 1d.
- `webhook` (Boolean)


<a id="nestedatt--s_3"></a>
### Nested Schema for `s_3`

Required:

- `bucket` (String)

Optional:

- `credential_ref` (String) With endpoint: the VaultEntry holding the access key for it.
- `endpoint` (String) An S3-compatible endpoint instead of AWS.
- `external_id` (String) The external id the role's trust policy requires. Defaults to the tenant uuid.
- `key` (String) One object.
- `prefix` (String) Every object under this prefix.
- `region` (String)
- `role_arn` (String) The role to assume. Required unless endpoint is set.


<a id="nestedatt--scope"></a>
### Nested Schema for `scope`

Optional:

- `org` (Boolean) The whole tenant.
- `products` (Set of String) TEA products, by TEI or uuid.
- `repos` (Set of String)


<a id="nestedatt--tea"></a>
### Nested Schema for `tea`

Required:

- `server` (String) The TEA API root, or the host whose /.well-known/tea discovery names it.

Optional:

- `collection` (String) A TEA identifier (urn:tei:...) or a uuid.
- `product` (String) A TEA identifier (urn:tei:...) or a uuid.
- `product_release` (String) A TEA identifier (urn:tei:...) or a uuid.


<a id="nestedatt--trust"></a>
### Nested Schema for `trust`

Optional:

- `allowed_authors` (Set of String) Only statements by these authors (the VEX author field) are used.
- `require_signature` (Boolean) Refuse a document that is not signed. Schema default: `false`.

## Import

Import an existing object by its name (`metadata.name`). The first read after an import fills in every field the server holds.

```shell
# Import an existing VexSource by its name (metadata.name).
terraform import vulnetix_vex_source.local_vex local-vex
```
