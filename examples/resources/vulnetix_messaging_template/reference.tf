# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/MessagingTemplate/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_messaging_template" "plain_title" {
  name = "plain-title"
  content = {
    title = "{{.event.title}}"
  }
}
