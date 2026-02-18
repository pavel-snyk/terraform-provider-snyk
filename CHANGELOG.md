
<a name="v1.0.0-rc5"></a>
## v1.0.0-rc5 (2026-02-18)
### Features
* **resource/snyk_broker_deployment_credential**: add resource to manage credential references
### Maintaining
* **deps**: bump github.com/hashicorp/terraform-plugin-sdk/v2
* **deps**: upgrade snyk-sdk-go to v2 latest dev version

<a name="v1.0.0-rc4"></a>
## v1.0.0-rc4 (2026-02-17)
### Documentation
* **datasource/snyk_app_install**: update documentation with examples
### Features
* **datasource/snyk_app_install**: add datasource for snyk app installation
### Maintaining
* rename files aligning with official naming schema
* show generated changelog during release process

<a name="v1.0.0-rc3"></a>
## v1.0.0-rc3 (2026-02-15)
### Documentation
* **resource/snyk_app_install**: update documentation with examples
### Features
* **resource/snyk_app_install**: add resource to manage app installations
* **resource/snyk_broker_deployment**: add resource to manage broker deployments

<a name="v1.0.0-rc2"></a>
## v1.0.0-rc2 (2026-02-13)
### Bug Fixes
* **tools**: adjust path to git-chglog in release workflow

<a name="v1.0.0-rc1"></a>
## v1.0.0-rc1 (2026-02-13)
### Bug Fixes
* **provider**: catch edge case with empty config token
* **resource/snyk_organization**: wait until tenant is provisioned by creation
### Code Refactoring
* allow to set custom urls with envvars
* re-generate tfdocs
* disable existing resources due terraform-framework migration
* **provider**: simplify resolving of region block
* **provider**: move validation region logic into ValidateConfig
### Documentation
* **datasource/snyk_user**: update documentation
* **resource/snyk_organization**: update documentation
### Features
* refactor usage of new Snyk SDK client
* **datasource/snyk_organization**: migrate datasource to use new Snyk SDK
* **datasource/snyk_user**: migrate datasource to use new Snyk SDK
* **provider**: add support for legacy V1 Snyk API
* **resource/snyk_organization**: migrate resource to use new Snyk SDK
### Maintaining
* **cicd**: fix typo for go pkg directory
* **cicd**: stop running acceptance tests on main branch
* **cicd**: run sweepers after acceptance tests
* **cicd**: pass new envvars for acceptance tests
* **datasource/snyk_organization**: add acceptance tests
* **deps**: bump hashicorp/setup-terraform from 2 to 3
* **deps**: upgrade snyk-sdk-go to v2 dev version
* **deps**: bump actions to latest stable versions
* **deps**: bump all dependencies to latest stable versions
* **lint**: allow usage of retry package from SDKv2
* **resource/snyk_organization**: refactor acceptance tests using new tf framework
* **test**: add sweepers to cleanup leftover infrastructure
* **tools**: migrate to golangci-lint v2
### BREAKING CHANGE

upgrade terraform-framework deps after 4 years shows a huge drift
in API for provider, resources and data sources. Because Snyk SDK will be migrated
from v1 API to REST API, we have to adjust/rewrite all API calls anyway.

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
