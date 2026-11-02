# Every argument the kind supports. From config-schemas examples/MessagingTemplate/advanced.yaml.

resource "vulnetix_messaging_template" "security_findings" {
  name                 = "security-findings"
  base_system_template = "generic"
  content = {
    accent = "{{.event.severity}}"
    actions = [
      {
        label = "Open in Vulnetix"
        url   = "{{.event.url}}"
      },
    ]
    body = "{{range $k, $v := .data}}- **{{$k}}:** {{truncate 200 $v}}\n{{end}}{{with .mentions}}\n{{mentions .}}{{end}}"
    digest = {
      item  = "{{.event.title}}{{with .subject.label}}: {{.}}{{end}}"
      title = "{{.digest.count}} {{plural .digest.count \"event\" \"events\"}}"
    }
    fields = [
      {
        inline = true
        label  = "Severity"
        value  = "{{severityEmoji .event.severity}} {{severityLabel .event.severity}}"
      },
    ]
    issue = {
      body   = "{{.event.title}}"
      labels = ["vulnetix", "security"]
      title  = "[Vulnetix] {{.event.title}}"
    }
    summary = "{{with .subject.label}}{{.}}{{end}}"
    title   = "{{.event.title}}"
    webhook = {
      include_payload = true
    }
  }
  description  = "Scan results for the security channel, with a link back to the scan."
  display_name = "Security findings"
  event_scope  = ["scan.completed", "scan.build_broken"]
  labels = {
    team = "security"
  }
  annotations = {
    owner = "appsec@acme.example"
  }
}
