---
page_title: "vulnetix_library_item Resource - vulnetix"
subcategory: ""
description: |-
  One item of the organisation's Belai library that is not an agent or a crew: a skill (an Agent Skills SKILL.md), prompt, command, process, MCP server, repo, budget set, provider set, bash rewrite table, hook bundle or launch configuration. The body is the item exactly as the library stores it, and the item's name is the one the body carries; the library holds one live item of a name per type, and neither the type nor the name can change under a metadata.name. The server checks the body with the library's own validator, the one the console and the Belai CLI use. Every apply that changes it saves a new version; deleting the document hides the item and keeps its versions. An item nobody declared is exported under its name.
  Manages the LibraryItem config kind through /config/objects/LibraryItem/{name}. Applies are recorded with the provider's manager, so a GitOps vulnetix config apply --prune leaves this object alone.
  Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (most of them, such as ApiSpecSource, RiskStrategy and VexSource) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object, or lets go of it for a kind the server never deletes (such as Organization). Some kinds wait on a person after they are applied (an invited Member, a Slack connection): status.phase says so, and Terraform does not wait.
---

# vulnetix_library_item (Resource)

One item of the organisation's Belai library that is not an agent or a crew: a skill (an Agent Skills SKILL.md), prompt, command, process, MCP server, repo, budget set, provider set, bash rewrite table, hook bundle or launch configuration. The body is the item exactly as the library stores it, and the item's name is the one the body carries; the library holds one live item of a name per type, and neither the type nor the name can change under a metadata.name. The server checks the body with the library's own validator, the one the console and the Belai CLI use. Every apply that changes it saves a new version; deleting the document hides the item and keeps its versions. An item nobody declared is exported under its name.

Manages the `LibraryItem` config kind through `/config/objects/LibraryItem/{name}`. Applies are recorded with the provider's `manager`, so a GitOps `vulnetix config apply --prune` leaves this object alone.

Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (most of them, such as `ApiSpecSource`, `RiskStrategy` and `VexSource`) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object, or lets go of it for a kind the server never deletes (such as `Organization`). Some kinds wait on a person after they are applied (an invited `Member`, a Slack connection): `status.phase` says so, and Terraform does not wait.

## Example Usage

### Minimal

The arguments the kind requires, and nothing else.

```terraform
# Only the arguments the kind requires. From config-schemas examples/LibraryItem/minimal.yaml.

resource "vulnetix_library_item" "pdf_tools" {
  name = "pdf-tools"
  body = "---\nname: pdf-tools\ndescription: Read, split and merge PDF files.\n---\nUse the pdf tools on the path the user gives.\n"
  type = "skill"
}
```

### Advanced

Every argument the kind supports. A secret is a write-only `*_wo` argument fed from an ephemeral variable, so it never reaches the plan or the state; change `*_wo_version` to send a new value.

```terraform
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
```

### Reference

Every argument that has a schema default, set to that default, plus the required arguments. Planned against a tenant it only sets values the tenant left unset. Every kind's reference file is also in [Vulnetix/platform-iac-examples](https://github.com/Vulnetix/platform-iac-examples), with the same document as YAML.

```terraform
# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/LibraryItem/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_library_item" "pdf_tools" {
  name = "pdf-tools"
  body = "---\nname: pdf-tools\ndescription: Read, split and merge PDF files.\n---\nUse the pdf tools on the path the user gives.\n"
  type = "skill"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The object's name (`metadata.name`), unique per kind within the tenant. Changing it replaces the object.
- `type` (String) skill, prompt and command are markdown with front matter; the others are JSON documents in the library's format for that type. One of `skill`, `prompt`, `command`, `process`, `mcp`, `repo`, `budget`, `provider`, `rewrite`, `hook`, `launch`.

### Optional

- `annotations` (Map of String) Free key/value notes for people and tools (`metadata.annotations`), such as a ticket link.
- `body` (String, Sensitive) The item: markdown for skill, prompt and command, JSON for the rest. The library's size limit for the type applies.
- `file` (String) A file holding the item, the same as body.
- `files` (Attributes Set) The scripts a hook bundle runs, written beside its hooks.json. Only a hook carries files: at most 32, each a relative path of safe segments (not hooks.json) and UTF-8 text of up to 256 KiB. (see [below for nested schema](#nestedatt--files))
- `labels` (Map of String) Free key/value labels (`metadata.labels`). Stored with the object, never interpreted by the server.

### Read-Only

- `id` (String) The object's name, the same as `name`.
- `status` (Attributes) What the server says about the object beyond its spec. (see [below for nested schema](#nestedatt--status))

<a id="nestedatt--files"></a>
### Nested Schema for `files`

Required:

- `path` (String) Where the file is written, relative to the agent: one line, no .. segment, naming a file.

Optional:

- `body` (String, Sensitive) The file's contents.
- `file` (String) A file holding the contents, the same as body.


<a id="nestedatt--status"></a>
### Nested Schema for `status`

Read-Only:

- `authorization_url` (String) Where a person authorizes the object while the phase is `AwaitingAuthorization`.
- `extra` (Map of String) Other fields the kind reports, such as an invitation's expiry or a service account's principal. Structured values are JSON.
- `message` (String) Why the object is in its phase, when the server says.
- `object_uuid` (String) The object's id on the server.
- `once` (Map of String, Sensitive) Credentials the server shows only when it creates the object: a service account's `apiKey` and `secret`, a messaging connection's `signingSecret`, a log source's ingest token. They are kept in the state from the create, so keep the state encrypted. They cannot be read again: replace the object to get new ones.
- `phase` (String) `Applied` for most objects. `Pending` when the object waits on someone (an invited member who has not accepted), `AwaitingAuthorization` when a person has to authorize it in a browser at `authorization_url`. Terraform does not wait for either. `Unmanaged` means the object exists but no manager owns it.

## Import

Import an existing object by its name (`metadata.name`). The first read after an import fills in every field the server holds.

```shell
# Import an existing LibraryItem by its name (metadata.name).
terraform import vulnetix_library_item.pdf_tools pdf-tools
```
