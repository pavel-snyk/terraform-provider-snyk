---
page_title: "Authenticate with Snyk using service accounts"
description: |-
  Authenticate with Snyk using service accounts
---

# Authenticate with Snyk using service accounts

~> **Note:** Due to Snyk API v1 restrictions, it is strongly recommended
to create a group with a service account and give this group API entitlements.
This feature is available for Business and Enterprise plans.

## Before you begin

- To use the provider, you must first get your token from Snyk. We recommended
  to create an [service account token](https://docs.snyk.io/features/user-and-group-management/structure-account-for-high-application-performance/service-accounts)
  for group with the **Group Admin** role.
- [Install Terraform](https://learn.hashicorp.com/tutorials/terraform/install-cli)
  and read the Terraform getting started guide that follows. This guide will
  assume basic proficiency with Terraform - it is an introduction to the Snyk
  provider.

## How to set up a service account

- Log in to your account and navigate to the relevant group and organization that you want to manage.
- Click on settings **⚙️ > Service accounts** to view existing service accounts and their details.
- Click **Create a service account** to create a new one.
- Give it the role **Group Admin** in combobox, since it will be writing resources.

## Two options to configure the provider
Save the token as the environment variables `SNYK TOKEN`. Or, configure
the provider with the token by copy-pasting the values directly into provider
config.

```terraform
// Token can be set explicitly or via the environment variable SNYK_TOKEN
provider "snyk" {
  token = "service-account-token"
}
```
