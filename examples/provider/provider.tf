terraform {
  required_providers {
    snyk = {
      source  = "pavel-snyk/snyk"
      version = "~> 0.1.0"
    }
  }
}

# Set the variable value in *.tfvars file
# or using -var="snyk_token=..." CLI option
variable "snyk_token" {}

# Configure the Snyk Provider
provider "snyk" {
  token = var.snyk_token
}
