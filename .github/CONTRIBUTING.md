# Contributing to the Snyk Terraform Provider

👍🎉 First off, thanks for taking the time to contribute! 🎉👍

The following is a set of guidelines for contributing to the Snyk Terraform
Provider. These are mostly guidelines, not rules. Use your best judgment, and
feel free to propose changes to this document in a pull request.

## Requirements

- Terraform >= 0.13.x
- Go >= 1.18

## Building The Provider

1. Clone the repository
2. Enter the repository directory
3. Build the provider using the `make build`

> Note: you can see all Makefile commands by executing `make help`.

## Using the Provider

[Development overrides for provider developers](https://www.terraform.io/docs/cli/config/config-file.html#development-overrides-for-provider-developers)
can be leveraged in order to use the provider built from source.

To do this, populate a Terraform CLI configuration file (`~/.terraformrc` for
all platforms other than Windows; `terraform.rc` in the `%APPDATA%` directory
when using Windows) with at least the following options:

```
provider_installation {
  dev_overrides {
    "local/snyk" = "[REPLACE WITH CLONED REPOSITORY PATH]/build"
  }
  direct {}
}
```
