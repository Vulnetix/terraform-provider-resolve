variable "vulnetix_token" {
  type        = string
  sensitive   = true
  description = "An Authentik OAuth2 access token or a Vulnetix session JWT, for an owner or admin."
}

# Bearer alternative to the ApiKey: an Authentik access token, for example
# from a client-credentials grant, sent as "Authorization: Bearer <token>".
# It names its own organisation, so there is no org_id. Set either token or
# api_key with org_id, never both. Leave it out and set VULNETIX_API_TOKEN to
# keep it out of the configuration.
provider "vulnetix" {
  token = var.vulnetix_token
}
