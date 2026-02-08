# Set the variable value in *.tfvars file
# or using -var="snyk_token=..." CLI option
variable "snyk_token" {}

# Configure the Snyk Provider with predefined SNYK-EU-01 region.
provider "snyk" {
  region = {
    name = "SNYK-EU-01"
  }
  token = var.snyk_token
}
