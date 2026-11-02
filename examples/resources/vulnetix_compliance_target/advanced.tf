# Every argument the kind supports. From config-schemas examples/ComplianceTarget/advanced.yaml.

resource "vulnetix_compliance_target" "eu_ai_act_payments" {
  name      = "eu-ai-act-payments"
  cadence   = "quarterly"
  enabled   = true
  framework = "eu_ai_act"
  report = {
    context = "The organisation is the provider of a high-risk credit scoring system (Annex III, 5b)."
    title   = "EU AI Act, payments"
  }
  scope = {
    repositories = ["acme/payments-api", "acme/fraud-model"]
  }
  labels = {
    team = "governance"
  }
  annotations = {
    ticket = "https://tracker.example.com/GRC-77"
  }
}

resource "vulnetix_compliance_target" "slsa_build_l2" {
  name      = "slsa-build-l2"
  cadence   = "monthly"
  framework = "slsa"
  slsa = {
    level = 2
  }
}

resource "vulnetix_compliance_target" "ssdf_core" {
  name      = "ssdf-core"
  cadence   = "quarterly"
  framework = "ssdf"
  ssdf = {
    min_met   = 80
    practices = ["PS", "PW", "RV"]
  }
}

resource "vulnetix_compliance_target" "pqc_readiness" {
  name    = "pqc-readiness"
  cadence = "monthly"
  crypto = {
    allow          = "quantumSafe || (primitive in ['block-cipher', 'hash', 'mac'] && classicalSecurityLevel >= 128)"
    applies_to     = "category == 'algorithm' && confidence != 'low'"
    max_violations = 0
  }
  framework = "crypto-policy"
}
