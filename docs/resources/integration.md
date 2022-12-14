---
page_title: "snyk_integration Resource - terraform-provider-snyk"
subcategory: ""
description: |-
  The integration resource allows you to manage Snyk integration.
---

# snyk_integration (Resource)

The integration resource allows you to manage Snyk integration.

## Example Usage

### Default settings

```terraform
resource "snyk_integration" "gitlab" {
  organization_id = snyk_organization.frontend.id

  type  = "gitlab"
  token = "gitlab-secret-token"
}
```

### Using custom configuration to update out-of-date dependencies

```terraform
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
```

### Using custom SCA configuration for new opened pull requests

```terraform
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `organization_id` (String) The ID of the organization that the integration belongs to.
- `type` (String) The integration type, e.g. 'github'.

### Optional

- `password` (String, Sensitive) The password used by the integration.
- `pull_request_dependency_upgrade` (Attributes) The pull request configuration for dependency upgrades. Snyk can automatically raise pull requests to update out-of-date dependencies. (see [below for nested schema](#nestedatt--pull_request_dependency_upgrade))
- `pull_request_sca` (Attributes) The pull request testing configuration for SCA (Software Composition Analysis). Snyk will checks projects imported through the SCM integration for security and license issues whenever a new PR is opened. (see [below for nested schema](#nestedatt--pull_request_sca))
- `region` (String) The region used by the integration.
- `registry_url` (String) The URL for container registries used by the integration (e.g. for ECR).
- `role_arn` (String) The role ARN used by the integration (ECR only).
- `token` (String, Sensitive) The token used by the integration.
- `url` (String) The URL used by the integration.
- `username` (String) The username used by the integration.

### Read-Only

- `id` (String) The ID of the integration.

<a id="nestedatt--pull_request_dependency_upgrade"></a>
### Nested Schema for `pull_request_dependency_upgrade`

Optional:

- `enabled` (Boolean) Denotes the pull request automatic dependency upgrade feature should be enabled for this integration
- `ignored_dependencies` (List of String) List of exact names of the dependencies that should not be included in the automatic upgrade operation. You can use only enter lowercase letters.
- `include_major_version` (Boolean) Defines if major version upgrades will be included in the recommendations. By default, only patches and minor versions are included in the upgrade recommendations.
- `limit` (Number) The maximum number of simultaneously opened pull requests with dependency upgrades.


<a id="nestedatt--pull_request_sca"></a>
### Nested Schema for `pull_request_sca`

Optional:

- `enabled` (Boolean) Denotes the pull request SCA feature should be enabled for this integration.
- `fail_on_any_issue` (Boolean) Fails an opened pull request if any vulnerable dependencies have been detected, otherwise the pull request should only fail when a dependency with issues is added.
- `fail_only_for_high_and_critical_severity` (Boolean) Fails an opened pull request if any dependencies are marked as being of high or critical severity.
- `fail_only_on_issues_with_fix` (Boolean) Fails an opened pull request only when issues found have a fix available.
