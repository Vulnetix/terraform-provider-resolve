# Every argument the kind supports. From config-schemas examples/LogSource/advanced.yaml.

variable "cloudflare_edge_cloudflare_api_token" {
  type      = string
  sensitive = true
  ephemeral = true
}

resource "vulnetix_log_source" "cloudflare_edge" {
  name = "cloudflare-edge"
  cloudflare = {
    account_id           = "0123456789abcdef0123456789abcdef"
    api_token_wo         = var.cloudflare_edge_cloudflare_api_token
    api_token_wo_version = 1
    datasets             = ["firewall_events", "http_requests", "gateway_dns"]
    zone_ids             = ["fedcba9876543210fedcba9876543210"]
  }
  label = "Cloudflare edge and Gateway"
  type  = "cloudflare"
  labels = {
    team = "netsec"
  }
  annotations = {
    note = "The Cloudflare token is used for this apply only; deleting the source needs it again."
  }
}

resource "vulnetix_log_source" "aws_prod" {
  name  = "aws-prod"
  label = "Production VPCs"
  type  = "aws"
}

resource "vulnetix_log_source" "sensor_events" {
  name  = "sensor-events"
  label = "Runtime sensors"
  type  = "events"
}
