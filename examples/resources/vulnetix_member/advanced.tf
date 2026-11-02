# Every argument the kind supports. From config-schemas examples/Member/advanced.yaml.

resource "vulnetix_member" "alice" {
  name  = "alice"
  email = "alice@acme.com"
  role  = "owner"
  labels = {
    team = "security"
  }
  annotations = {
    ticket = "https://tracker.example.com/ACCESS-42"
  }
}

resource "vulnetix_member" "bob" {
  name  = "bob"
  email = "bob@acme.com"
  role  = "member"
}
