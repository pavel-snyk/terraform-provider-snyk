
<a name="v0.7.1"></a>
## v0.7.1 (2023-03-11)
### Maintaining
* bump GitHub action versions to latest
* don't use deprecated GoReleaser options
* use go 1.19 in GitHub action builds
* **deps**: upgrade terraform-plugin-docs to v0.14.1
* **deps**: upgrade git-chglog to v0.15.4
* **deps**: upgrade golangci-lint to v1.51.2
* **deps**: bump github.com/hashicorp/terraform-plugin-sdk/v2
* **deps**: bump github.com/stretchr/testify from 1.8.0 to 1.8.2
* **deps**: bump goreleaser/goreleaser-action from 3 to 4

<a name="v0.7.0"></a>
## v0.7.0 (2022-09-24)
### Bug Fixes
* **resource/snyk_integration**: use oneOf validator for integration type
### Features
* **resource/snyk_integration**: add credentials validator for type attribute
* **resource/snyk_integration**: add configuration options for dependency auto-upgrade
### Maintaining
* upgrade snyk-sdk-go to v0.4.1
* complete migration to terraform-plugin-framework v0.13.0
* don't run acceptance tests in parallel
* update snyk-sdk-go to v0.4.0
* **deps**: bump github.com/hashicorp/terraform-plugin-framework

<a name="v0.6.1"></a>
## v0.6.1 (2022-09-19)
### Code Refactoring
* **snyk_integration**: rename pull_request_testing attribute to pull_request_sca
### Documentation
* fix provider version in example snippet

<a name="v0.6.0"></a>
## v0.6.0 (2022-09-18)
### Features
* **resource/snyk_integration**: add configuration options for pull request testing
### Maintaining
* update snyk-sdk-go to v0.3.1
* fix caching for tools task in test workflow
* **deps**: bump github.com/hashicorp/terraform-plugin-sdk/v2
* **deps**: bump github.com/hashicorp/terraform-plugin-sdk/v2

<a name="v0.5.0"></a>
## v0.5.0 (2022-09-05)
### Features
* **datasource/snyk_project**: add new data source for projects
### Maintaining
* upgrade snyk-sdk-go to v0.1.0

<a name="v0.4.2"></a>
## v0.4.2 (2022-09-05)
### Bug Fixes
* **resource/snyk_integration**: fix plan and state handling by updating integration
### Code Refactoring
* add validator with NotEmptyString() function
* **resource/snyk_organization**: add acceptance tests for organization resource
### Maintaining
* set provider version via ldflags when releasing
* add debug support for provider

<a name="v0.4.1"></a>
## v0.4.1 (2022-09-02)
### Bug Fixes
* consider envvars by provider configuration
### Documentation
* add how-to guide for provider authentication
### Maintaining
* activate unit and acceptance tests workflows
* install golangci-lint via tools-as-dependencies
* **deps**: bump github.com/hashicorp/terraform-plugin-framework
* **golangci-lint**: replace deadcode and varcheck linter with unused

<a name="v0.4.0"></a>
## v0.4.0 (2022-08-22)
### Features
* **resource/snyk_integration**: add new resource for integrations

<a name="v0.3.1"></a>
## v0.3.1 (2022-08-20)
### Maintaining
* use git-chglog for changelog generation
* install tfplugindocs via tool-as-dependencies
* add make targets to generate and validate docs

<a name="v0.3.0"></a>
## v0.3.0 (2022-08-18)
### Features
* **datasource/snyk_organization**: add new data source for organizations

<a name="v0.2.0"></a>
## v0.2.0 (2022-08-16)
### Documentation
* update readme and contribution guidelines
### Features
* **resource/snyk_organization**: add new resource for organizations
### Maintaining
* upgrade snyk-sdk-go to latest dev version

<a name="v0.1.0"></a>
## v0.1.0 (2022-08-14)
### Features
* **datasource/snyk_user**: add new data source for users
