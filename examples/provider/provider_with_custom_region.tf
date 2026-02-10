# Set the variable value in *.tfvars file
# or using -var="snyk_token=..." CLI option
variable "snyk_token" {}

# Configure the Snyk Provider with custom region.
provider "snyk" {
  region = {
    name          = "my-instance"
    app_base_url  = "https://app.my-instance.local/"
    rest_base_url = "https://api.my-instance.local/rest/"
    v1_base_url   = "https://api.my-instance.local/v1/"
  }
  token = var.snyk_token
}
