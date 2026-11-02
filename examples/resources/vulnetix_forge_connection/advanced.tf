# Every argument the kind supports. From config-schemas examples/ForgeConnection/advanced.yaml.

variable "gitlab_platform_oauth_client_secret" {
  type      = string
  sensitive = true
  ephemeral = true
}
variable "gitlab_platform_token" {
  type      = string
  sensitive = true
  ephemeral = true
}
variable "ghes_acme_app_private_key" {
  type      = string
  sensitive = true
  ephemeral = true
}

resource "vulnetix_forge_connection" "gitlab_platform" {
  name     = "gitlab-platform"
  account  = "platform"
  api_base = "https://gitlab.example.com/api/v4"
  host     = "gitlab.example.com"
  oauth = {
    client_id                = "6f1e2d3c4b5a69788796a5b4c3d2e1f0"
    client_secret_wo         = var.gitlab_platform_oauth_client_secret
    client_secret_wo_version = 1
  }
  token_wo         = var.gitlab_platform_token
  token_wo_version = 1
  type             = "gitlab"
  labels = {
    team = "platform"
  }
  annotations = {
    note = "The bot account is vulnetix-bot, a Developer in the platform group."
  }
}

resource "vulnetix_forge_connection" "ghes_acme" {
  name    = "ghes-acme"
  account = "acme"
  app = {
    app_id                 = 12
    app_slug               = "vulnetix"
    installation_id        = 3
    private_key_wo         = var.ghes_acme_app_private_key
    private_key_wo_version = 1
  }
  host = "github.acme.example"
  type = "ghes"
}
