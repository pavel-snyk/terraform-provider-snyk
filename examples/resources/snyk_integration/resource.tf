resource "snyk_integration" "github" {
  organization_id = snyk_organization.frontend.id

  type  = "github"
  token = "rotated-github-secret-token"
}
