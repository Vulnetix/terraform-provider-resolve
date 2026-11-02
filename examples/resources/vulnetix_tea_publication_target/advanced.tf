# Every argument the kind supports. From config-schemas examples/TeaPublicationTarget/advanced.yaml.

variable "globex_tea_credential" {
  type      = string
  sensitive = true
  ephemeral = true
}

resource "vulnetix_tea_publication_target" "globex_tea" {
  name                  = "globex-tea"
  credential_wo         = var.globex_tea_credential
  credential_wo_version = 1
  description           = "Globex's TEA server, for the payments contract."
  root_url              = "https://tea.globex.example.com/tea/v1"
  suspended             = true
  labels = {
    owner = "platform"
  }
}
