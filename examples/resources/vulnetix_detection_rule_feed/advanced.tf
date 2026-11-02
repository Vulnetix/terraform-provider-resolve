# Every argument the kind supports. From config-schemas examples/DetectionRuleFeed/advanced.yaml.

resource "vulnetix_detection_rule_feed" "edge_suricata" {
  name         = "edge-suricata"
  description  = "Block-lane Suricata rules for exploited edge appliance CVEs, fetched by the perimeter IPS."
  display_name = "Edge Suricata virtual patches"
  entries = [
    {
      cve        = "CVE-2024-3400"
      rule       = "4d5e6f7a-8b9c-4d0e-9f1a-2b3c4d5e6f7a"
      rule_table = "virtual_patch"
      rule_type  = "SURICATA"
    },
    {
      cve        = "CVE-2023-46805"
      disabled   = true
      rule       = "5e6f7a8b-9c0d-4e1f-8a2b-3c4d5e6f7a8b"
      rule_table = "virtual_patch"
      rule_type  = "SURICATA"
    },
  ]
  format         = "SURICATA"
  publish        = true
  schema_version = 1
  selector = {
    when = "kind == 'SURICATA' && (kev || epss >= 0.5) && cvss >= 7.0 && falsePositiveRisk != 'HIGH' && now - published < duration('8760h')"
  }
  labels = {
    engine = "suricata"
    team   = "netsec"
  }
  annotations = {
    consumer = "Suricata sensors at the edge pull this feed every 15 minutes."
  }
}
