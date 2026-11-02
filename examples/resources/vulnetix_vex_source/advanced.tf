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
