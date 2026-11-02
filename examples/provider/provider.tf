terraform {
  required_providers {
    vulnetix = {
      source = "vulnetix/vulnetix"
    }
  }
}

# The token is read from VULNETIX_CONFIG_TOKEN when not set here. Keep it out
# of the configuration: pass it from your CI secret store.
provider "vulnetix" {
  # endpoint = "https://www.vulnetix.com/api/site/v1"
  # manager  = "terraform"
}
