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
