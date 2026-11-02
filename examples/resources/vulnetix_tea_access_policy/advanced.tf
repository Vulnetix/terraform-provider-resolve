# Every argument the kind supports. From config-schemas examples/TeaAccessPolicy/advanced.yaml.

resource "vulnetix_tea_access_policy" "acme_payments" {
  name           = "acme-payments"
  confirm_public = true
  product        = "acme-payments"
  visibility     = "public"
  annotations = {
    approved-by = "https://tracker.example.com/SEC-210"
  }
}

resource "vulnetix_tea_access_policy" "payments_api" {
  name      = "payments-api"
  component = "payments-api"
  shared_with = [
    {
      expires           = "2027-06-30T00:00:00Z"
      organization      = "6f1c2a8e-3b4d-4e5f-8a9b-0c1d2e3f4a5b"
      organization_name = "Globex Corporation"
    },
  ]
  visibility = "shared"
}

resource "vulnetix_tea_access_policy" "acme_payments_2_4_0" {
  name       = "acme-payments-2-4-0"
  collection = "acme-payments-2-4-0"
  visibility = "private"
}
