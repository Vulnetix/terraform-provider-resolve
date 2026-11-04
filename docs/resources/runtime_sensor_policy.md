---
page_title: "vulnetix_runtime_sensor_policy Resource - vulnetix"
subcategory: ""
description: |-
  The configuration runtime sensors enroll with and keep pulling: mode, hold, rule feeds, forwarding and tamper protection, and which people or MDM systems it applies to. It is a sensor provisioning profile. A device takes the enabled profile with the highest priority that matches it, and picks up a change on its next config pull. Sensors configured by hand with an ApiKey do not read profiles; a person changes their mode from the console. A profile with enrolled devices cannot be deleted.
  An object of this kind can reference DetectionRuleFeed objects. Reference those resources' attributes, or add depends_on, so Terraform applies them first.
  Manages the RuntimeSensorPolicy config kind through /config/objects/RuntimeSensorPolicy/{name}. Applies are recorded with the provider's manager, so a GitOps vulnetix config apply --prune leaves this object alone.
  Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (most of them, such as ApiSpecSource, RiskStrategy and VexSource) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object, or lets go of it for a kind the server never deletes (such as Organization). Some kinds wait on a person after they are applied (an invited Member, a Slack connection): status.phase says so, and Terraform does not wait.
---

# vulnetix_runtime_sensor_policy (Resource)

The configuration runtime sensors enroll with and keep pulling: mode, hold, rule feeds, forwarding and tamper protection, and which people or MDM systems it applies to. It is a sensor provisioning profile. A device takes the enabled profile with the highest priority that matches it, and picks up a change on its next config pull. Sensors configured by hand with an ApiKey do not read profiles; a person changes their mode from the console. A profile with enrolled devices cannot be deleted.

An object of this kind can reference `DetectionRuleFeed` objects. Reference those resources' attributes, or add `depends_on`, so Terraform applies them first.

Manages the `RuntimeSensorPolicy` config kind through `/config/objects/RuntimeSensorPolicy/{name}`. Applies are recorded with the provider's `manager`, so a GitOps `vulnetix config apply --prune` leaves this object alone.

Most kinds merge: an argument removed from the configuration is no longer managed, and keeps its last value on the server. A kind that is wholly declared (most of them, such as `ApiSpecSource`, `RiskStrategy` and `VexSource`) is replaced instead: a removed argument goes back to its default. Destroying the resource deletes the object, or lets go of it for a kind the server never deletes (such as `Organization`). Some kinds wait on a person after they are applied (an invited `Member`, a Slack connection): `status.phase` says so, and Terraform does not wait.

## Example Usage

### Minimal

The arguments the kind requires, and nothing else.

```terraform
# Only the arguments the kind requires. From config-schemas examples/RuntimeSensorPolicy/minimal.yaml.

resource "vulnetix_runtime_sensor_policy" "engineering_laptops" {
  name = "engineering-laptops"
}
```

### Advanced

Every argument the kind supports. A secret is a write-only `*_wo` argument fed from an ephemeral variable, so it never reaches the plan or the state; change `*_wo_version` to send a new value.

```terraform
# Every argument the kind supports. From config-schemas examples/RuntimeSensorPolicy/advanced.yaml.

resource "vulnetix_runtime_sensor_policy" "engineering_laptops" {
  name               = "engineering-laptops"
  allow_major_update = true
  allow_user_stop    = false
  credential_mode    = "device-key"
  description        = "Every engineer's laptop, in guard mode, plus the security team by address."
  dns_listen         = "127.0.0.1:5353"
  egress_listen      = "127.0.0.1:8181"
  enabled            = true
  feed_key           = "KymtMlKOJvj2IkhUbMpXRl/JBmmOwu/LIj8BekwmQXk="
  feeds = jsonencode([{
    feedRef = "edge-waf"
    }, "https://vulnetix.com/feeds/0b6f2c1e-4d3a-4e5f-8a9b-7c6d5e4f3a2b/9c8b7a6f-5e4d-4c3b-8a2f-1e0d9c8b7a6f", {
    sensorFeed = "7d6c5b4a-3f2e-4d1c-9b8a-7f6e5d4c3b2a"
  }])
  forward                 = ["siem"]
  heartbeat_alert_minutes = 30
  hold                    = "15m"
  inventory_every         = "6h"
  match = {
    domains          = ["acme.example"]
    emails           = ["alice@acme.example", "bob@acme.example"]
    mdm_integrations = ["5e4d3c2b-1a09-4f8e-9d7c-6b5a4f3e2d1c"]
  }
  mode            = "guard"
  priority        = 200
  protection      = "hardened"
  update_base_url = "https://mirror.example.com/vulnetix/cli"
  labels = {
    fleet = "laptops"
  }
  annotations = {
    ticket = "https://jira.example.com/browse/SEC-2210"
  }
}

resource "vulnetix_runtime_sensor_policy" "servers_locked" {
  name            = "servers-locked"
  auto_update     = false
  credential_mode = "org-key"
  description     = "Production servers enrolled by MDM, blocking every match on a block-mode rule."
  feeds = jsonencode([{
    feedRef = "edge-suricata"
  }])
  heartbeat_alert_minutes = 5
  match = {
    mdm_integrations = ["5e4d3c2b-1a09-4f8e-9d7c-6b5a4f3e2d1c"]
  }
  mode       = "enforce"
  priority   = 10
  protection = "locked"
}
```

### Reference

Every argument that has a schema default, set to that default, plus the required arguments. Planned against a tenant it only sets values the tenant left unset. Every kind's reference file is also in [Vulnetix/platform-iac-examples](https://github.com/Vulnetix/platform-iac-examples), with the same document as YAML.

```terraform
# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/RuntimeSensorPolicy/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_runtime_sensor_policy" "engineering_laptops" {
  name                    = "engineering-laptops"
  allow_major_update      = false
  allow_user_stop         = false
  auto_update             = true
  credential_mode         = "device-key"
  enabled                 = true
  heartbeat_alert_minutes = 15
  hold                    = "15m"
  inventory_every         = "6h"
  mode                    = "monitor"
  priority                = 100
  protection              = "hardened"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The object's name (`metadata.name`), unique per kind within the tenant. Changing it replaces the object.

### Optional

- `allow_major_update` (Boolean) Let enrolled sensors take a release with a higher major version. Off by default: a major release can change the profile format, the units and the protocol, so a sensor stays on its major version until the profile says otherwise. Has no effect when autoUpdate is false. Schema default: `false`.
- `allow_user_stop` (Boolean) Let the device's user stop the sensor. Schema default: `false`.
- `annotations` (Map of String) Free key/value notes for people and tools (`metadata.annotations`), such as a ticket link.
- `auto_update` (Boolean) Let enrolled sensors install new CLI releases themselves. The guard checks the latest release about every 6 hours, verifies its checksum, keeps the previous binary as .prev and restarts onto the new one. false holds a sensor on its installed version, for hosts that take releases on a change window. Schema default: `true`.
- `credential_mode` (String) The credential an enrolled device gets. device-key: a key that works on the sensor routes only, of which the server keeps a hash. user-key: the enrolling person's. org-key: the organization's. One of `device-key`, `user-key`, `org-key`. Schema default: `device-key`.
- `description` (String)
- `dns_listen` (String) Where the DNS guard listens.
- `egress_listen` (String) Where the egress proxy listens.
- `enabled` (Boolean) Schema default: `true`.
- `feed_key` (String) The base64 Ed25519 public key feed manifests must be signed with. Left out, the sensor's built-in key.
- `feeds` (String) The feeds the sensor pulls: a DetectionRuleFeed by name, which the server resolves to its published URL; one of the organization's own indicator feeds, by uuid; or an https URL, which must be a published Vulnetix rule feed (https://vulnetix.com/feeds/<org>/<feed>) or a feed in the sensor feed catalogue. Left out, the profile pulls the catalogue's default feeds; an empty list pulls none. A JSON-encoded value; use jsonencode().
- `forward` (Set of String) siem: send every DNS and egress decision to Guardian as dns_events and net_events rows.
- `heartbeat_alert_minutes` (Number) Alert when an enrolled sensor has been silent this long. Schema default: `15`.
- `hold` (String) How long guard mode holds a match, as a Go duration from 1s to 24h. Schema default: `15m`.
- `inventory_every` (String) How often the sensor reports the packages on its host, installed and loaded by a running process, as a Go duration from 1m to 168h. Schema default: `6h`.
- `labels` (Map of String) Free key/value labels (`metadata.labels`). Stored with the object, never interpreted by the server.
- `match` (Attributes) Who enrolls with this profile: a person by exact email or email domain (* for everyone), or a device an MDM integration enrolls. (see [below for nested schema](#nestedatt--match))
- `mode` (String) monitor: report matches, block nothing. guard: hold a match on a rule safe to block (a request denied, an exec refused, a connection dropped) until a person confirms a block or marks a false positive, or the hold lapses; an action that cannot be undone is never held. enforce: block every match on a block-mode rule, a running process included. In every mode a rule from a feed whose signature does not verify only detects. One of `monitor`, `guard`, `enforce`. Schema default: `monitor`.
- `priority` (Number) Highest wins when several profiles match a device; ties go to the oldest. Schema default: `100`.
- `protection` (String) How hard the install resists being stopped or removed on the host. standard: a plain service. hardened: a guard restores a sensor that was stopped, disabled or masked, and the user cannot stop it unless allowUserStop. locked: as hardened, and uninstalling needs a one-time code from the console. One of `standard`, `hardened`, `locked`. Schema default: `hardened`.
- `update_base_url` (String) An https mirror of the sensor's release layout (<base>/latest redirecting to <base>/tag/<tag>, and <base>/download/<tag>/checksums.txt with the binaries) for hosts that cannot reach GitHub. Absent: GitHub. Credentials in the URL are refused; a pinned release signing key still applies to a mirror.

### Read-Only

- `id` (String) The object's name, the same as `name`.
- `status` (Attributes) What the server says about the object beyond its spec. (see [below for nested schema](#nestedatt--status))

<a id="nestedatt--match"></a>
### Nested Schema for `match`

Optional:

- `domains` (Set of String)
- `emails` (Set of String)
- `mdm_integrations` (Set of String) MDM integrations, by uuid.


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
# Import an existing RuntimeSensorPolicy by its name (metadata.name).
terraform import vulnetix_runtime_sensor_policy.engineering_laptops engineering-laptops
```
