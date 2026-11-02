data "vulnetix_detection_rule_feed" "suricata" {
  name = "edge-ips"
}

# Point the sensor at the published manifest.
output "suricata_manifest" {
  value = data.vulnetix_detection_rule_feed.suricata.manifest_url
}
