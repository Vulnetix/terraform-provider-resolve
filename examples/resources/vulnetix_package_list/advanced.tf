# Every argument the kind supports. From config-schemas examples/PackageList/advanced.yaml.

resource "vulnetix_package_list" "first_party" {
  name = "first-party"
  entries = [
    {
      purl   = "pkg:npm/@acme/*"
      reason = "Our own npm scope; releases must not wait out the cooldown."
      scope  = "full"
      when   = "!malware && !badActor"
    },
    {
      purl   = "pkg:golang/github.com/acme/*"
      reason = "Our own Go modules."
    },
    {
      purl   = "pkg:pypi/requests@2.32.3"
      reason = "CVSS on this release is a false positive for our use."
      when   = "cvss < 9.0 && !kev"
    },
    {
      purl   = "pkg:npm/left-pad"
      reason = "Maintainer attribution reviewed and judged a false positive."
      scope  = "bad_actor"
    },
  ]
  list = "allow"
  labels = {
    team = "platform-security"
  }
  annotations = {
    ticket = "https://jira.example.com/browse/SEC-2211"
  }
}

resource "vulnetix_package_list" "banned" {
  name = "banned"
  entries = [
    {
      purl   = "pkg:npm/event-stream@3.3.6"
      reason = "Compromised release."
      scope  = "full"
    },
    {
      purl   = "pkg:npm/ua-parser-js"
      reason = "Hold back fresh releases after the 2021 hijack."
      when   = "ageDays < 3 || versionLag > 5"
    },
  ]
  list = "deny"
}
