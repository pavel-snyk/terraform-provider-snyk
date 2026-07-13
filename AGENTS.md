# Terraform Provider Snyk - Agent Instructions

Mandatory repository instructions for AI coding agents working on the Snyk Terraform Provider.

## Scope

These instructions apply to the entire repository.

Keep changes focused. Preserve unrelated local changes and work safely with a dirty worktree. Do not perform unrelated refactoring.

## Core Principles

The Terraform schema is a user-facing interface, not a mechanical copy of Snyk SDK types.

Prioritize:

- predictable Terraform behavior and useful diagnostics;
- correct state, refresh, import, and drift detection;
- stable plans and idempotent lifecycle operations;
- safe handling of secrets and credentials;
- consistency across the provider.

Terraform fields may combine, rename, omit, normalize, or derive API fields when that improves usability. Such differences must be deliberate, documented, tested, and compatible
with correct refresh behavior.

Do not encode unverified API assumptions into Terraform state. Verify material behavior through authoritative API documentation, SDK behavior, or tests.

## Resource Behavior

A material change to schema semantics, lifecycle behavior, import, request mapping, response mapping, or state requires reviewing the complete affected resource. Fix nearby defects
reasonably within scope; report defects whose repair would substantially expand the task.

For every changed attribute, define its create, update, read, null, unknown, default, sensitivity, replacement, and drift behavior.

Use `RequiresReplace` for immutable fields. Do not use `UseStateForUnknown` to hide incomplete mapping, suppress drift, or retain stale remote values.

Construct straightforward API requests close to their lifecycle operation. Extract request mapping only when reuse or conversion complexity provides clear implementation value.
Preserve meaningful distinctions between null, unknown, empty, and zero values.

Centralize API-to-state mapping on the Terraform model in a deterministic `fromAPI` method:

- map every remotely backed schema field, including nested values;
- pass contextual identifiers not returned by the API explicitly;
- do not make API calls or log inside `fromAPI`;
- use the same mapper across lifecycle methods when API representations are compatible;
- check conversion diagnostics before writing state.

Capture create-only values before any read-back, map the authoritative object through `fromAPI`, then restore the create-only values. Preserve them deliberately during later reads
when the API no longer returns them.

`Read` must:

- use every identifier required for lookup;
- handle all pages for list-based lookup and match by immutable ID;
- refresh every remotely backed schema field;
- remove absent remote objects from state;
- distinguish not-found from authorization, validation, transport, and server failures;
- produce stable state when configuration and remote state are unchanged.

`Delete` must treat an already absent remote object as successfully deleted and must not remove state after a failed deletion.

Resources that can exist independently of Terraform must support import unless it is technically impossible or semantically meaningless; document any exception.

Import must populate every identifier required by the first `Read`. Use and validate a documented composite import format when the resource ID alone is insufficient. Never use
passthrough import in that case.

## Logging and Sensitive Data

Each implemented resource lifecycle operation must log at `tflog.Debug` when it begins and after it succeeds. Use only safe identifiers and do not defer success logs.

Use `tflog.Trace` for API details:

- `"payload"` is for outgoing request objects;
- `"data"` is for incoming response objects.

Keep the request trace immediately before the API call. Keep the response trace immediately after error handling, with no blank lines separating related statements.

Complete objects may be logged only after the full runtime value has been reviewed recursively and every populated field is known to be non-sensitive. Re-review whole-object
logging whenever the SDK type changes or the SDK is upgraded.

If any part may be sensitive, log an explicitly sanitized map or dedicated safe type and never pass the original object to the logger.

Never log tokens, OAuth client secrets, passwords, credential values, authorization data, cookies, private keys, sensitive headers, Terraform state, or one-time secret values.
Credential identifiers may be logged.

Terraform `Sensitive` limits normal CLI display but does not prevent storage in state. Store sensitive values only when required and never use real secrets in fixtures or tests.

Use `err.Error()` directly only when its contents are known to be non-sensitive. Otherwise return a sanitized, actionable diagnostic. Handle nil API responses before reading status
or request identifiers.

## Tests and Documentation

Add tests proportional to the changed behavior. New resources must cover, where applicable, create and refresh, supported updates, import, an empty plan after refresh and import,
normal destruction, and relevant validation. Add drift, external deletion, pagination, retry, and create-only secret coverage when applicable.

Unit-test non-trivial `fromAPI` and Terraform value-conversion behavior directly.

Do not disable, loosen, skip, or comment out tests merely to make validation pass. A skip requires a documented environment, platform, credential, or API limitation.

Update schema descriptions, templates, examples, generated documentation, and changelog entries when applicable. Include generated documentation with its source changes and do not
edit generated files instead of correcting their source.

Delete obsolete code rather than commenting it out. A TODO must reference a tracked issue and must not leave requested behavior incomplete.

## Commands and Authorization

For Go implementation changes, normally run:

```sh
make build
make lint
make test
```

Run `make docs` for schema, registration, template, example, or documentation changes. Run `go mod tidy` for dependency changes. Install pinned tools with `make tools` when needed.

Do not run acceptance tests without explicit authorization, suitable credentials, and a designated non-production Snyk environment. Do not run sweepers without explicit
authorization; sweepers destroy remote resources.

## Handoff

Report concisely:

- what changed and any material design decisions;
- state, import, drift, sensitivity, or compatibility effects;
- validation commands run and their results;
- acceptance-test status;
- remaining risks, assumptions, or relevant commands not run.

Never claim that all tests passed when only a subset ran.
