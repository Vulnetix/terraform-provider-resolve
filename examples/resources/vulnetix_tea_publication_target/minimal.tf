# Only the arguments the kind requires. From config-schemas examples/TeaPublicationTarget/minimal.yaml.

resource "vulnetix_tea_publication_target" "globex_tea" {
  name     = "globex-tea"
  root_url = "https://tea.globex.example.com/tea/v1"
}
