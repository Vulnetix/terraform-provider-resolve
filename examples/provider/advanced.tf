variable "vulnetix_org_id" {
  type        = string
  description = "The organisation UUID the ApiKey belongs to."
}

variable "vulnetix_api_key" {
  type        = string
  sensitive   = true
  description = "The ApiKey of an admin or owner service account."
}

# Advanced: every argument, with the ApiKey passed as variables. Set them with
# TF_VAR_vulnetix_org_id and TF_VAR_vulnetix_api_key, or leave api_key and
# org_id out and set VULNETIX_API_KEY and VULNETIX_ORG_ID instead.
provider "vulnetix" {
  # The console API. Defaults to $VULNETIX_ENDPOINT, then
  # https://www.vulnetix.com/api/site/v1.
  endpoint = "https://www.vulnetix.com/api/site/v1"

  # Sent as "Authorization: ApiKey <org_id>:<api_key>". The principal must be
  # an owner or admin; a member is refused with HTTP 403 on apply.
  org_id  = var.vulnetix_org_id
  api_key = var.vulnetix_api_key

  # Recorded as the owner of every object this provider applies, so
  # "vulnetix config apply --prune" leaves them alone.
  manager = "terraform"
}
