resource "snyk_broker_deployment" "dev" {
  app_install_id  = "<app-install-id>"
  organization_id = "<organization-id>"
  tenant_id       = "<tenant-id>"

  metadata = {
    env     = "dev"
    name    = "deployment for frontend team on staging cluster"
    comment = "managed by terraform"
  }
}

output "broker_deployment_id" {
  value = snyk_broker_deployment.dev.id
}
