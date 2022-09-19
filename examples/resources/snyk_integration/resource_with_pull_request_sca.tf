resource "snyk_integration" "github" {
  organization_id = snyk_organization.backend.id

  type  = "github"
  token = "rotated-github-secret-token"

  pull_request_sca = {
    enabled = true

    fail_on_any_issue            = false
    fail_only_on_issues_with_fix = true
  }
}
