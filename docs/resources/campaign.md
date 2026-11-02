---
page_title: "vulnetix_campaign Resource - vulnetix"
subcategory: ""
description: |-
  A remediation campaign: a named, time-boxed push to drive a chosen set of vulnerabilities, dependencies or repositories to a defined done state. The scope chooses the cohort, the done criteria decide when each member of it is closed, and the campaign records which members were admitted, when, and what later scan closed each one. What a campaign has achieved is runtime state, reported in status and never declared: its cohort, progress, freeze time and whether it has completed. The cohort is evaluated after apply, the first time the campaign is opened or refreshed. The repositories the scope and the criteria name must resolve when the document is applied, as they must when the console saves a campaign. Removing the document archives the campaign rather than deleting it, so the record of what was committed to, and whether it happened, is kept; an archived campaign is no longer part of the tenant's config, and applying the document again starts a new campaign.
  Manages the Campaign config kind through /config/objects/Campaign/{name}. Applies are recorded with the provider's manager, so a GitOps vulnetix config apply --prune leaves this object alone.
  Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (most of them, such as ApiSpecSource, RiskStrategy and VexSource) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object, or lets go of it for a kind the server never deletes (such as Organization). Some kinds wait on a person after they are applied (an invited Member, a Slack connection): status.phase says so, and Terraform does not wait.
---

# vulnetix_campaign (Resource)

A remediation campaign: a named, time-boxed push to drive a chosen set of vulnerabilities, dependencies or repositories to a defined done state. The scope chooses the cohort, the done criteria decide when each member of it is closed, and the campaign records which members were admitted, when, and what later scan closed each one. What a campaign has achieved is runtime state, reported in status and never declared: its cohort, progress, freeze time and whether it has completed. The cohort is evaluated after apply, the first time the campaign is opened or refreshed. The repositories the scope and the criteria name must resolve when the document is applied, as they must when the console saves a campaign. Removing the document archives the campaign rather than deleting it, so the record of what was committed to, and whether it happened, is kept; an archived campaign is no longer part of the tenant's config, and applying the document again starts a new campaign.

Manages the `Campaign` config kind through `/config/objects/Campaign/{name}`. Applies are recorded with the provider's `manager`, so a GitOps `vulnetix config apply --prune` leaves this object alone.

Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (most of them, such as `ApiSpecSource`, `RiskStrategy` and `VexSource`) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object, or lets go of it for a kind the server never deletes (such as `Organization`). Some kinds wait on a person after they are applied (an invited `Member`, a Slack connection): `status.phase` says so, and Terraform does not wait.

## Example Usage

### Minimal

The arguments the kind requires, and nothing else.

```terraform
# Only the arguments the kind requires. From config-schemas examples/Campaign/minimal.yaml.

resource "vulnetix_campaign" "log4shell" {
  name         = "log4shell"
  display_name = "Log4Shell clean-up"
  done = {
    criteria = [
      {
        type = "gone-from-scans"
      },
    ]
  }
  target = "vulnerabilities"
}
```

### Advanced

Every argument the kind supports. A secret is a write-only `*_wo` argument fed from an ephemeral variable, so it never reaches the plan or the state; change `*_wo_version` to send a new value.

```terraform
# Every argument the kind supports. From config-schemas examples/Campaign/advanced.yaml.

resource "vulnetix_campaign" "kev_burndown_q1" {
  name = "kev-burndown-q1"
  contributors = [
    {
      display_name = "Dana Reyes"
      principal    = "2b7e1516-28ae-4d2a-a6ab-f7158809cf4f"
      role         = "owner"
    },
    {
      login = "payments-oncall"
      role  = "contributor"
    },
    {
      display_name = "External auditor"
      email        = "auditor@example.com"
      role         = "contact"
    },
  ]
  deadline     = "2027-03-31"
  description  = "Close every internet-exposed finding in the CISA KEV catalogue before the audit."
  display_name = "KEV burndown, Q1"
  done = {
    criteria = [
      {
        category = "sca"
        repos    = ["acme/payments-api"]
        type     = "gone-from-scans"
      },
      {
        type = "safe-version"
      },
      {
        mode = "supported"
        type = "eol-supported"
      },
    ]
    join = "any"
  }
  owner_email = "appsec-lead@example.com"
  scope = {
    category = "sca"
    exclusions = {
      dev_deps        = true
      min_confidence  = 0.5
      suppressed      = true
      test_code       = true
      transitive_deps = true
    }
    filters = {
      attack_chain = {
        attack_path_ids = ["ap-ingress-payments"]
        tactics         = ["TA0001"]
        techniques      = ["T1190"]
      }
      exploit_asserted = {
        kev            = true
        msrc_exploited = true
        ransomware     = true
        vulncheck_kev  = true
      }
      exploit_observed = {
        any     = true
        sources = ["shadowserver"]
      }
      exposure = {
        audiences        = ["public"]
        internet_exposed = true
      }
      identity = {
        branch_name       = ["main"]
        category          = ["sca"]
        commit_sha        = ["4f2a9c1"]
        cpe               = ["cpe:2.3:a:apache:log4j:2.14.1:*:*:*:*:*:*:*"]
        finding_id        = ["CVE-2021-44228"]
        package_ecosystem = ["maven"]
        package_name      = ["log4j-core"]
        package_version   = ["2.14.1"]
        product           = ["log4j"]
        purl              = ["pkg:maven/org.apache.logging.log4j/log4j-core@2.14.1"]
        repo_name         = ["acme/payments-api"]
        source            = ["vulnetix"]
        vendor            = ["apache"]
      }
      impact = {
        min_cess            = 5
        min_cvss            = 7
        min_cwss            = 50
        min_epss            = 0.1
        min_epss_percentile = 0.9
      }
      licence = {
        expressions     = ["GPL-3.0-only"]
        violates_policy = true
      }
      lifecycle = {
        eol          = true
        unmaintained = true
      }
      reachability = {
        assessed_only = true
        verdicts      = ["REACHABLE"]
      }
      supply_chain = {
        abandoned       = true
        bad_actor       = true
        malicious       = true
        scorecard_below = 4
        squatted        = true
      }
      weakness = {
        aliases               = ["GHSA-jfh8-c2jp-5v3q"]
        clean_code_categories = ["intentionality"]
        cwes                  = ["CWE-502"]
        rule_keys             = ["java.lang.security.audit.jndi-injection"]
      }
    }
    repos = ["acme/payments-api", "acme/web"]
  }
  target = "vulnerabilities"
  timing = "frozen"
  labels = {
    team = "platform-security"
  }
  annotations = {
    ticket = "https://tracker.example.com/SEC-118"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `display_name` (String) The campaign's name, as the console shows it.
- `done` (Attributes) When a member of the cohort is done. (see [below for nested schema](#nestedatt--done))
- `name` (String) The object's name (`metadata.name`), unique per kind within the tenant. Changing it replaces the object.
- `target` (String) What the cohort is made of. vulnerabilities: findings. dependencies: dependencies in use. repositories: both the findings and the dependencies of the repositories in scope. Fixed once the campaign exists: a changed target is refused, because the cohort it already admitted would no longer fit it. One of `vulnerabilities`, `dependencies`, `repositories`.

### Optional

- `annotations` (Map of String) Free key/value notes for people and tools (`metadata.annotations`), such as a ticket link.
- `contributors` (Attributes Set) Who is on the hook. Each contributor is identified by exactly one of principal, login or email, because a console member, a code owner and an external contact are different identities. Plan lists each contributor added or removed, and apply changes only those. (see [below for nested schema](#nestedatt--contributors))
- `deadline` (String) When the campaign is due. Recorded as a deadline set by hand.
- `description` (String) What the campaign is for, in a sentence or two.
- `labels` (Map of String) Free key/value labels (`metadata.labels`). Stored with the object, never interpreted by the server.
- `owner_email` (String) The email address of the person accountable for the campaign.
- `scope` (Attributes) The filter that chooses the cohort. Left out, or empty, the campaign covers the whole estate. Every filter set narrows the cohort further. Some dimensions are recorded and not yet evaluated (attackChain, licence.violatesPolicy, supplyChain.scorecardBelow, exploitObserved.sources, weakness.cleanCodeCategories and weakness.ruleKeys); the console reports them as deferred. (see [below for nested schema](#nestedatt--scope))
- `timing` (String) frozen: the cohort is locked the first time the campaign is evaluated, so entities found later never join and the burndown's denominator cannot move. scopeCreep: every evaluation admits whatever newly matches the scope. A frozen campaign's scope and criteria cannot be changed; set scopeCreep first, which unfreezes it. One of `frozen`, `scopeCreep`. Schema default: `frozen`.

### Read-Only

- `id` (String) The object's name, the same as `name`.
- `status` (Attributes) What the server says about the object beyond its spec. (see [below for nested schema](#nestedatt--status))

<a id="nestedatt--done"></a>
### Nested Schema for `done`

Required:

- `criteria` (Attributes List) The done criteria. For a vulnerabilities campaign: gone-from-scans, safe-version and eol-supported. For a dependencies campaign: package-removed, latest-version and safe-version-used. A repositories campaign takes either family. A finding the scanner marks fixed always closes, whatever the criteria. (see [below for nested schema](#nestedatt--done--criteria))

Optional:

- `join` (String) all: a member is done when every criterion holds. any: when one does. One of `all`, `any`. Schema default: `all`.

<a id="nestedatt--done--criteria"></a>
### Nested Schema for `done.criteria`

Required:

- `type` (String) gone-from-scans: the finding is no longer reported. safe-version: the package moved to a version the advisory says is fixed. eol-supported: the package is on a supported release line. package-removed: the dependency is no longer used. latest-version: the dependency is on its latest version. safe-version-used: the dependency is on a version with no known vulnerability. One of `gone-from-scans`, `safe-version`, `eol-supported`, `package-removed`, `latest-version`, `safe-version-used`.

Optional:

- `category` (String) gone-from-scans only: the finding category whose scans count. One of `sca`, `sast`, `secrets`, `iac`, `license`.
- `mode` (String) eol-supported only. latest: the latest release line. supported: any release line still supported. One of `latest`, `supported`.
- `repos` (Set of String) gone-from-scans only: the repositories whose scans must no longer report the finding, which may differ from the scope. Each must resolve as the scope's repositories do.



<a id="nestedatt--contributors"></a>
### Nested Schema for `contributors`

Optional:

- `display_name` (String) The name to show.
- `email` (String) An email address, for a contact outside the console.
- `login` (String) A GitHub login, such as a code owner.
- `principal` (String) A console member, by their Vulnetix principal UUID.
- `role` (String) owner: accountable for the campaign. contributor: does the work. contact: kept informed. One of `owner`, `contributor`, `contact`. Schema default: `contributor`.


<a id="nestedatt--scope"></a>
### Nested Schema for `scope`

Optional:

- `category` (String) One finding category. Left out, every category. One of `sca`, `sast`, `secrets`, `iac`, `license`.
- `exclusions` (Attributes) Noise left out of the cohort. (see [below for nested schema](#nestedatt--scope--exclusions))
- `filters` (Attributes) (see [below for nested schema](#nestedatt--scope--filters))
- `repos` (Set of String) The repositories the campaign covers. For a repositories campaign they are the target itself; for the others they narrow it. Each must resolve to a repository Vulnetix holds or can fetch from its forge.

<a id="nestedatt--scope--exclusions"></a>
### Nested Schema for `scope.exclusions`

Optional:

- `dev_deps` (Boolean) Leave out development dependencies.
- `min_confidence` (Number) Leave out findings whose scanner confidence is below this value. Findings with no confidence score are kept.
- `suppressed` (Boolean) Leave out suppressed findings.
- `test_code` (Boolean) Leave out findings in test code.
- `transitive_deps` (Boolean) Leave out transitive dependencies.


<a id="nestedatt--scope--filters"></a>
### Nested Schema for `scope.filters`

Optional:

- `attack_chain` (Attributes) Findings on an attack path. Recorded, not yet evaluated. (see [below for nested schema](#nestedatt--scope--filters--attack_chain))
- `exploit_asserted` (Attributes) Findings someone with standing says are exploited. Each flag set must hold. (see [below for nested schema](#nestedatt--scope--filters--exploit_asserted))
- `exploit_observed` (Attributes) Findings with exploitation observed in the wild. (see [below for nested schema](#nestedatt--scope--filters--exploit_observed))
- `exposure` (Attributes) Findings by exposure. (see [below for nested schema](#nestedatt--scope--filters--exposure))
- `identity` (Attributes) Match findings by what they are and where they were found. Each list matches any of its values. (see [below for nested schema](#nestedatt--scope--filters--identity))
- `impact` (Attributes) Minimum scores. A finding must reach every minimum set. (see [below for nested schema](#nestedatt--scope--filters--impact))
- `licence` (Attributes) Findings by licence. (see [below for nested schema](#nestedatt--scope--filters--licence))
- `lifecycle` (Attributes) Findings by package lifecycle. (see [below for nested schema](#nestedatt--scope--filters--lifecycle))
- `reachability` (Attributes) Findings by reachability analysis. (see [below for nested schema](#nestedatt--scope--filters--reachability))
- `supply_chain` (Attributes) Dependencies with supply chain risk. Each flag set must hold. (see [below for nested schema](#nestedatt--scope--filters--supply_chain))
- `weakness` (Attributes) Match findings by weakness. (see [below for nested schema](#nestedatt--scope--filters--weakness))

<a id="nestedatt--scope--filters--attack_chain"></a>
### Nested Schema for `scope.filters.attack_chain`

Optional:

- `attack_path_ids` (Set of String) Attack path ids.
- `tactics` (Set of String) MITRE ATT&CK tactic ids, such as TA0001.
- `techniques` (Set of String) MITRE ATT&CK technique ids, such as T1190.


<a id="nestedatt--scope--filters--exploit_asserted"></a>
### Nested Schema for `scope.filters.exploit_asserted`

Optional:

- `kev` (Boolean) In the CISA Known Exploited Vulnerabilities catalogue.
- `msrc_exploited` (Boolean) Marked exploited by the Microsoft Security Response Center.
- `ransomware` (Boolean) Known to be used in ransomware campaigns.
- `vulncheck_kev` (Boolean) In the VulnCheck KEV catalogue.


<a id="nestedatt--scope--filters--exploit_observed"></a>
### Nested Schema for `scope.filters.exploit_observed`

Optional:

- `any` (Boolean) Observed by any source.
- `sources` (Set of String) Observed by one of these sources. Recorded, not yet evaluated.


<a id="nestedatt--scope--filters--exposure"></a>
### Nested Schema for `scope.filters.exposure`

Optional:

- `audiences` (Set of String) Exposure audiences, such as public or internal.
- `internet_exposed` (Boolean) Only findings in services exposed to the internet.


<a id="nestedatt--scope--filters--identity"></a>
### Nested Schema for `scope.filters.identity`

Optional:

- `branch_name` (Set of String) Branch names.
- `category` (Set of String) Finding categories, as the findings list shows them.
- `commit_sha` (Set of String) Commit SHAs.
- `cpe` (Set of String) CPE names.
- `finding_id` (Set of String) Finding ids, such as CVE-2021-44228 or GHSA ids.
- `package_ecosystem` (Set of String) Package ecosystems, such as npm or pypi.
- `package_name` (Set of String) Package names.
- `package_version` (Set of String) Package versions.
- `product` (Set of String) Products, as the advisory names them.
- `purl` (Set of String) Package URLs.
- `repo_name` (Set of String) Repository full names.
- `source` (Set of String) The scanner or tool that reported the finding.
- `vendor` (Set of String) Vendors, as the advisory names them.


<a id="nestedatt--scope--filters--impact"></a>
### Nested Schema for `scope.filters.impact`

Optional:

- `min_cess` (Number) Minimum effective CESS score.
- `min_cvss` (Number) Minimum effective CVSS base score.
- `min_cwss` (Number) Minimum effective CWSS score.
- `min_epss` (Number) Minimum EPSS probability, from 0 to 1.
- `min_epss_percentile` (Number) Minimum EPSS percentile, from 0 to 1.


<a id="nestedatt--scope--filters--licence"></a>
### Nested Schema for `scope.filters.licence`

Optional:

- `expressions` (Set of String) SPDX licence expressions.
- `violates_policy` (Boolean) Licences the licence policy refuses. Recorded, not yet evaluated.


<a id="nestedatt--scope--filters--lifecycle"></a>
### Nested Schema for `scope.filters.lifecycle`

Optional:

- `eol` (Boolean) Packages past end of life.
- `unmaintained` (Boolean) Unmaintained packages.


<a id="nestedatt--scope--filters--reachability"></a>
### Nested Schema for `scope.filters.reachability`

Optional:

- `assessed_only` (Boolean) Only findings reachability analysis looked at.
- `verdicts` (Set of String) Reachability verdicts, such as REACHABLE or UNREACHABLE. Assessed and unreachable is a different state from never assessed, and either may be targeted.


<a id="nestedatt--scope--filters--supply_chain"></a>
### Nested Schema for `scope.filters.supply_chain`

Optional:

- `abandoned` (Boolean) Abandoned packages.
- `bad_actor` (Boolean) Packages published by a known bad actor.
- `malicious` (Boolean) Known malicious packages.
- `scorecard_below` (Number) OpenSSF Scorecard score below this value. Recorded, not yet evaluated.
- `squatted` (Boolean) Typosquatted or namespace-squatted packages.


<a id="nestedatt--scope--filters--weakness"></a>
### Nested Schema for `scope.filters.weakness`

Optional:

- `aliases` (Set of String) Advisory aliases.
- `clean_code_categories` (Set of String) Clean code categories. Recorded, not yet evaluated.
- `cwes` (Set of String) CWE ids, such as CWE-79.
- `rule_keys` (Set of String) Static analysis rule keys. Recorded, not yet evaluated.




<a id="nestedatt--status"></a>
### Nested Schema for `status`

Read-Only:

- `authorization_url` (String) Where a person authorizes the object while the phase is `AwaitingAuthorization`.
- `extra` (Map of String) Other fields the kind reports, such as an invitation's expiry or a service account's principal. Structured values are JSON.
- `message` (String) Why the object is in its phase, when the server says.
- `object_uuid` (String) The object's id on the server.
- `once` (Map of String, Sensitive) Credentials the server shows only when it creates the object: a service account's `apiKey` and `secret`, a messaging connection's `signingSecret`, a log source's ingest token. They are kept in the state from the create, so keep the state encrypted. They cannot be read again: replace the object to get new ones.
- `phase` (String) `Applied` for most objects. `Pending` when the object waits on someone (an invited member who has not accepted), `AwaitingAuthorization` when a person has to authorize it in a browser at `authorization_url`. Terraform does not wait for either. `Unmanaged` means the object exists but no manager owns it.

## Import

Import an existing object by its name (`metadata.name`). The first read after an import fills in every field the server holds.

```shell
# Import an existing Campaign by its name (metadata.name).
terraform import vulnetix_campaign.log4shell log4shell
```
