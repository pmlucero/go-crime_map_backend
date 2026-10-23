# CLAUDE.md

Guidance for Claude Code and other coding agents working in this repository.

## Project overview

Crime Map Backend is a REST API in Go for reporting and querying geolocated crime incidents, built with Clean Architecture, Gin and PostgreSQL. It is a public portfolio project, so clarity matters as much as function.

Read these before changing anything non-trivial:

- [docs/vision.md](docs/vision.md): the purpose, principles and non-goals.
- [docs/spec/core.md](docs/spec/core.md): the API contract and acceptance criteria (AC-xx), with which ones are met.
- [docs/adr/](docs/adr/README.md): the architecture decisions and their rationale.
- [docs/plan/v1.md](docs/plan/v1.md): the milestones (M0–M5) that close the remaining gaps, and their order.

## Workflow: Spec-Driven Development

1. **Spec first.** Behavior is defined in `docs/spec/` before it is implemented. If a task requires behavior the spec does not describe, update the spec first, in its own commit, and get approval. An approved change keeps the status `Approved` and bumps `Version` and `Last updated` (see CONTRIBUTING.md).
2. **ADR for significant decisions.** A change to the HTTP contract, the data model or schema, authentication or the layer structure, or a new runtime dependency, needs an ADR (see [ADR 0000](docs/adr/0000-record-architecture-decisions.md)). Code that contradicts an Accepted ADR must come with a superseding ADR in the same PR.
3. **Trace to acceptance criteria.** Tests that verify an acceptance criterion reference it by name, e.g. `// AC-12`. When a PR makes a `Not met` criterion pass, it updates its status in spec §9 in the same PR. Commit messages reference the AC or ADR they implement when applicable.
4. **Scope gate.** Implementation covers the spec only. Anything in spec §10 or beyond it needs an approved spec change first.

## Language

- All repository content is in **English**: code, comments, docs, commit messages, error messages and test names.
- The crime type values (`ROBO`, `HURTO`, …) are API identifiers in Spanish by design (spec §3.2). Do not translate them.
- Conversation with the maintainer happens in Spanish.

## Go conventions

- **Layout:** Clean Architecture ([ADR 0002](docs/adr/0002-clean-architecture-layers.md)). `cmd/api/` is wiring only. `internal/domain` imports nothing from other layers. `internal/usecases` depends only on the domain. Controllers depend on use case interfaces, never on repositories.
- **Dependencies:** prefer the standard library. Each new module dependency is justified in the PR, and in an ADR if it is significant.
- **Errors:** wrap with `fmt.Errorf("doing X: %w", err)`. Error messages are lowercase with no trailing punctuation. Compare errors with `errors.Is` or `errors.As`, never by their message. Never panic outside `main` and impossible states.
- **HTTP errors:** follow spec §6. A `500` response never includes SQL or driver messages.
- **Context:** every I/O-bound function takes `ctx context.Context` as its first parameter.
- **SQL:** PostgreSQL placeholders (`$1`, `$2`, …), never string interpolation of values. Every read of `crimes` filters `deleted_at IS NULL` ([ADR 0006](docs/adr/0006-soft-delete.md)). Schema changes update `schema.sql` and keep it applicable to an empty database (AC-26).
- **Formatting and linting:** `gofmt`/`goimports`, `go vet` and `golangci-lint` clean.
- **Comments:** document exported identifiers and non-obvious *why*. Do not narrate *what* the code does.
- **Swagger:** after changing handler annotations, regenerate `docs/` with `swag init -g cmd/api/main.go`, never edit the generated files by hand.

## Testing

- `go test -race -p 1 ./...` must pass. `-p 1` is required because integration tests share one database.
- Integration tests need PostgreSQL with the `crime_map` and `crime_map_test` databases and `schema.sql` applied (see README, *Testing*).
- Tests live in a `tests` subpackage next to the code they test (e.g. `internal/usecases/tests`).
- Test names follow `TestFunctionName_Scenario`, and test bodies use `// Arrange`, `// Act`, `// Assert` sections. The pre-commit hooks enforce both.
- Prefer table-driven tests. Mocks are generated with mockery.

## Commits and PRs

- **Trunk-based:** branch from `main` as `<type>/<kebab-case-description>`, using the commit types below, and merge back via PR. Never push to `main`. See CONTRIBUTING.md, *Branches and releases*.
- **Small commits:** one logical change per commit. A spec change, an ADR and its implementation are separate commits.
- **Conventional Commits:** `feat:`, `fix:`, `docs:`, `test:`, `refactor:`, `chore:`, with an optional scope, e.g. `fix(api): return 404 for unknown crimes`.
- The subject is imperative and at most 72 characters. The body explains *why* and references `AC-xx` or `ADR NNNN` when relevant.
- Never commit secrets, `.env` files, real API keys or generated reports.

## Code review

Reviews, human or automated, follow [docs/review-guidelines.md](docs/review-guidelines.md): grounded findings, severities, the review dimensions and the output format. When asked to review a change locally, use the same rubric. Every PR also gets an advisory review from `.github/workflows/claude-review.yml`. Comment `@claude review` on a PR to request another pass.

## Things not to do

- Do not add features listed in spec §10 or in the [non-goals](docs/vision.md#4-non-goals) without an approved spec change.
- Do not accept the numeric `id` as input ([ADR 0007](docs/adr/0007-uuid-public-identifiers.md)).
- Do not hard-delete crimes ([ADR 0006](docs/adr/0006-soft-delete.md)).
- Do not add fields that identify people (victims, suspects, witnesses).
