resource "snyk_organization" "frontend" {
  name     = "my-awesome-frontend-team"
  // you can specify group_id if you want to create
  // an org inside of the group
  group_id = "c0b44a43-e0a7-4862-b23b-8a581a547081"
}
