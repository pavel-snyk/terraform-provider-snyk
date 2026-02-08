# Set the variable value in *.tfvars file
# or using -var="snyk_token=..." CLI option
variable "snyk_token" {}

# Configure the Snyk Provider with default SNYK-US-01 region.
provider "snyk" {
  token = var.snyk_token
}
