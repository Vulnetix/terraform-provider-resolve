# Every argument the kind supports. From config-schemas examples/LibraryItem/advanced.yaml.

resource "vulnetix_library_item" "install_guard" {
  name = "install-guard"
  body = "{\n  \"name\": \"install-guard\",\n  \"description\": \"Checks package installs before they run\",\n  \"hooks\": {\n    \"PreToolUse\": [\n      {\"matcher\": \"Bash\", \"hooks\": [{\"type\": \"command\", \"command\": \"bash scripts/guard.sh\", \"timeout\": 30}]}\n    ]\n  }\n}\n"
  files = [
    {
      body = "#!/usr/bin/env bash\nexec vulnetix agent hook\n"
      path = "scripts/guard.sh"
    },
    {
      file = "hooks/install-guard/allow.txt"
      path = "scripts/allow.txt"
    },
  ]
  type = "hook"
  labels = {
    team = "platform"
  }
}

resource "vulnetix_library_item" "triage_prompt" {
  name = "triage-prompt"
  file = "prompts/triage.md"
  type = "prompt"
}
