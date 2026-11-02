# Only the arguments the kind requires. From config-schemas examples/MessagingTemplate/minimal.yaml.

resource "vulnetix_messaging_template" "plain_title" {
  name = "plain-title"
  content = {
    title = "{{.event.title}}"
  }
}
