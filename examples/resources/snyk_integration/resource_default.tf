resource "snyk_integration" "gitlab" {
  organization_id = snyk_organization.frontend.id

  type  = "gitlab"
  token = "gitlab-secret-token"
}
