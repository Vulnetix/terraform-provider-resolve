data "vulnetix_virtual_patches" "waf" {
  feed_name = "edge-waf"
}

output "patched_cves" {
  value = distinct([for p in data.vulnetix_virtual_patches.waf.patches : p.matched_cve if p.matched_cve != ""])
}
