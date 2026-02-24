resource "snyk_broker_integration" "gitlab" {
  broker_connection_id = snyk_broker_connection.gitlab.id
  organization_id      = "<organization-id>"
  tenant_id            = "<tenant-id>"
  type                 = snyk_broker_connection.gitlab.type
}

resource "snyk_broker_connection" "gitlab" {
  app_install_id       = "<app-install-id>"
  tenant_id            = "<tenant-id>"
  broker_deployment_id = "<broker-deployment-id>"

  type = "gitlab"
  name = "dev-gitlab-connection"
  configuration = {
    broker_client_url          = "https://api.snyk.io"
    gitlab_hostname            = "gitlab.com"
    gitlab_token_credential_id = "<deployment-credential-id>"
  }
}

