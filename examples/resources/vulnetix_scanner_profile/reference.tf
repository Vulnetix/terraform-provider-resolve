# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/ScannerProfile/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_scanner_profile" "default" {
  custom   = []
  evidence = "cdx_slsa"
  scanners = "minimum"
  slsa_min = "L3"
}
