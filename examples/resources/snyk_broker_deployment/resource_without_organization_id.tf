# The "organization_id" attribute is optional but recommended
# for faster performance. If omitted, the provider searches all
# accessible organizations to locate the Universal Broker
# App installation, which will increase creation time.
resource "snyk_broker_deployment" "dev" {
  app_install_id = "<app-install-id>"
  tenant_id      = "<tenant-id>"
}
