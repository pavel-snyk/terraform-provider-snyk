data "snyk_project" "next-gen-ui" {
  organization_id = data.snyk_organization.frontend.id

  name = "frontend/next-gen-ui:package.json"
}
