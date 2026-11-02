# Every model the AI firewall knows whose slug or name mentions "sonnet".
data "vulnetix_ai_models" "sonnet" {
  search = "sonnet"
}

output "sonnet_models" {
  value = [for m in data.vulnetix_ai_models.sonnet.models : "${m.provider_slug}/${m.slug}" if m.is_active]
}
