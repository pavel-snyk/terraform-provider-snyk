resource "snyk_broker_connection" "gitlab" {
  app_install_id       = "<app-install-id>"
  tenant_id            = "<tenant-id>"
  broker_deployment_id = "<broker-deployment-id>"

  type = "bitbucket-server"
  name = "dev-bitbucket-connection"
  configuration = {
    broker_client_url                = "https://api.snyk.io"
    bitbucket_hostname               = "<bitbucket-server-hostname>"
    bitbucket_password_credential_id = snyk_broker_deployment_credential.bitbucket.id
    bitbucket_username               = "<bitbucket-username>"
  }
}

resource "snyk_broker_deployment_credential" "bitbucket" {
  app_install_id       = "<app-install-id>"
  tenant_id            = "<tenant-id>"
  broker_deployment_id = "<broker-deplyment-id>"

  environment_variable_name = "BB_USER_PASSWORD_FOR_DEV_CLUSTER"
  broker_connection_type    = "bitbucket-server"
}
