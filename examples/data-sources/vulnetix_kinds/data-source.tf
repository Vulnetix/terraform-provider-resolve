data "vulnetix_kinds" "all" {}

output "served_kinds" {
  value = [for k in data.vulnetix_kinds.all.kinds : k.kind]
}
