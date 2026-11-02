# Every repository the Vulnetix GitHub App can see in the acme organisation,
# archived ones left out.
data "vulnetix_github_installation_repos" "acme" {
  query = "acme/"
}

# Apply the security team's VEX statements to all of them.
resource "vulnetix_vex_source" "advisories" {
  name = "acme-advisories"
  scope = {
    repos = data.vulnetix_github_installation_repos.acme.full_names
  }
  git = {
    url   = "https://github.com/acme/security-advisories.git"
    ref   = "main"
    paths = ["vex/**/*.openvex.json"]
  }
}
