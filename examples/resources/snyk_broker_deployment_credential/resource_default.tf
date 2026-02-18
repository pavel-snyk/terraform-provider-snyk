resource "snyk_broker_deployment_credential" "gitlab" {
  app_install_id       = "<app-install-id>"
  tenant_id            = "<tenant-id>"
  broker_deployment_id = "<broker-deplyment-id>"

  environment_variable_name = "GITLAB_TOKEN_FOR_DEV_CLUSTER"
  broker_connection_type    = "gitlab"
}
