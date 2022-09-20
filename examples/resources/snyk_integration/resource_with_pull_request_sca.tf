resource "snyk_integration" "github" {
  organization_id = snyk_organization.backend.id

  type  = "github"
  token = "rotated-github-secret-token"

  # with this configuration Snyk will checks imported GitHub projects
  # whenever a new pull request is opened. Fail conditions are:
  #   - the pull request is adding a dependency with issues
  #   - when the issues found have a fix available
  pull_request_sca = {
    enabled = true

    fail_on_any_issue            = false
    fail_only_on_issues_with_fix = true
  }
}
