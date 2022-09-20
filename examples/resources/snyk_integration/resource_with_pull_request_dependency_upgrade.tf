resource "snyk_integration" "github" {
  organization_id = snyk_organization.backend.id

  type  = "github"
  token = "rotated-github-secret-token"

  # with this configuration Snyk will automatically raise pull requests
  # to update out-of-date dependencies:
  #   - updates for "jquery" and "lodash" dependencies will be ignored
  #   - maximum number of opened pull requests is 5
  pull_request_dependency_upgrade = {
    enabled = true

    ignored_dependencies = ["jquery", "lodash"]
    limit                = 5
  }
}
