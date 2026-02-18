data "snyk_app_install" "example" {
  organization_id = data.snyk_organization.frontend.id

  # Search criteria to find the app installation. At least one of
  # the following must be provided. Combining multiple attributes
  # makes the search more specific. For example, providing both `id` and
  # `app_name` will ensure that the app installation with the given
  # `id` also has the correct name, preventing potential mismatches.
  id       = "<app-install-id>"
  app_id   = "<app-id>"
  app_name = "<app-name>"
}

data "snyk_organization" "frontend" {
  name = "Frontend Team"
}
