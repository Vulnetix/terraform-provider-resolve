# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/TeaPublicationTarget/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_tea_publication_target" "globex_tea" {
  name      = "globex-tea"
  root_url  = "https://tea.globex.example.com/tea/v1"
  suspended = false
}
