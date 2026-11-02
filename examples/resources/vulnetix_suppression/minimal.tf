# Only the arguments the kind requires. From config-schemas examples/Suppression/minimal.yaml.

resource "vulnetix_suppression" "lodash_prototype_pollution" {
  name = "lodash-prototype-pollution"
  scope = {
    rule = "CVE-2021-23337"
  }
  type = "risk_accepted"
}
