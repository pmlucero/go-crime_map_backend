# Review guidelines

How pull requests to Crime Map Backend are reviewed, by humans and by the automated reviewer (`.github/workflows/claude-review.yml`). The goal is the review an experienced architect would give: it explains *why* something matters, it is grounded in the project's own documents, and it focuses on what will hurt later.

## Principles

1. **Grounded, not generic.** Every finding cites its basis: a spec section, an ADR, a rule in `CLAUDE.md`, a line in the diff, or a concrete failure scenario. If no basis can be stated, the finding is not reported.
2. **Proportional.** Review depth follows the change. A docs-only PR is not reviewed for SQL. Dimensions that do not apply are marked `N/A` in one line.
3. **Signal over volume.** Important issues come first. Nits are capped (at most 5) and never drown out real problems. Praise is allowed when it is specific.
4. **Actionable.** Each finding says what is wrong, why it matters here, and what to do instead, preferably with a concrete suggestion.
5. **Advisory.** The automated review never blocks a merge. The maintainer is the approver.

## Severity

| Level | Meaning | Examples |
|---|---|---|
| 🔴 **Blocker** | Must be fixed before merging | Contradicts an Accepted ADR; bypasses authentication; SQL injection; hard-deletes or exposes soft-deleted crimes; breaks the HTTP contract of the spec |
| 🟠 **Major** | Should be fixed before merging, or tracked explicitly | Wrong status code; an internal error leaked in a `500`; errors compared by message; a schema change that does not apply to an empty database; a test that does not verify what it claims; behavior not in the spec |
| 🟡 **Minor** | Worth fixing, not urgent | An unclear name; a function doing two things; a missing edge-case test |
| ⚪ **Nit** | Optional polish | Wording, small style points that linters do not catch |

## Dimensions

### 1. Architecture and product

- Does the change serve the [vision](vision.md), and stay out of the [non-goals](vision.md#4-non-goals)?
- Is the behavior described in the [spec](spec/core.md)? Behavior that is not in the spec is a finding (CLAUDE.md, *Workflow*).
- Is it consistent with the Accepted [ADRs](adr/README.md)? If it contradicts one, is there a superseding ADR?
- Is the *why* of the change clear from the PR, the commits and the code?
- Is it globally consistent: same concepts, same names and same patterns as the rest of the codebase?

### 2. Design: layers, coupling and dependencies

- Do dependencies point inward ([ADR 0002](adr/0002-clean-architecture-layers.md))? The domain imports no other layer, use cases do not import Gin or SQL, controllers do not call repositories.
- Are interfaces small and justified? Is a new abstraction used by more than one implementation or by a test?
- Is every new module dependency justified? Does it need an ADR (ADR 0000)?
- **Overengineering:** abstractions without a second implementation, configuration nobody asked for, generic frameworks for a single use. Simpler is preferred unless the spec requires otherwise.
- **Duplication:** a second implementation of something that already exists (the codebase has had duplicate use cases and repositories).

### 3. HTTP contract

- Status codes and error bodies follow spec §5 and §6: `400` for invalid input, `404` for unknown or soft-deleted crimes, `401` for authentication, `500` only for unexpected failures.
- Errors are mapped with `errors.Is` or `errors.As`, never by comparing messages.
- Request and response field names match the spec. Breaking changes need a spec change.
- The OpenAPI docs in `docs/` were regenerated with `swag` when handler annotations changed, and they match the routes.

### 4. Data and persistence

- Every read of `crimes` filters `deleted_at IS NULL` ([ADR 0006](adr/0006-soft-delete.md)).
- Queries use placeholders, never interpolated values. List queries are bounded by `limit`.
- Schema changes keep `schema.sql` applicable to an empty database (AC-26), and are consistent with the queries that use them.
- Indexes exist for new filters and sort orders. No N+1 queries.
- The numeric `id` is never accepted as input ([ADR 0007](adr/0007-uuid-public-identifiers.md)).

### 5. Resilience and error handling

- **Failures:** database unavailable, timeouts, constraint violations, rows not found.
- **Edge cases:** empty lists, boundary coordinates (`0`, `±90`, `±180`), very long strings, unknown enum values, pagination beyond the last page.
- **Errors are wrapped with context, never swallowed.** No panics outside `main` or impossible states.
- `context.Context` is propagated to every query.

### 6. Security

- Authentication cannot be bypassed: every `/crimes` route is behind the middleware (spec §4).
- API keys never appear in logs, responses or error messages.
- `500` responses expose no SQL, driver or stack details (AC-25).
- No SQL or shell injection, including in CI workflows.
- No fields that identify people are added (non-goals).

### 7. Readability and maintainability

- **Names explain themselves.** A reader should not need the implementation to understand what a function, type or variable is.
- **Single responsibility.** Functions and types do one thing. Long functions with several levels of nesting are a finding.
- **Junior-friendly.** Could a junior developer understand, modify and debug this code? If clever code is needed, its *why* is commented.
- Comments explain *why*, not *what*. Exported identifiers are documented.
- Everything is in English, except the crime type values (CLAUDE.md, *Language*).

### 8. Tests

- **Quality over coverage.** Tests assert behavior, not implementation details. A test that cannot fail is a finding.
- Acceptance criteria are traced: tests reference the `AC-xx` they verify, claimed ACs are actually covered, and criteria made to pass are marked `Met` in spec §9.
- Edge and failure cases from dimension 5 are tested, not just the happy path.
- Integration tests run against PostgreSQL, not only mocks, for anything that depends on SQL (the authentication path once broke because it was only tested with SQLite).
- Tests follow the repository conventions: `tests` subpackage, `TestFunctionName_Scenario`, Arrange-Act-Assert.

### 9. Git and PR hygiene

- Commits are small, atomic and in Conventional Commits format. The body explains *why*. Spec, ADR and implementation are in separate commits.
- The PR description fills in the traceability section (spec sections, AC-xx, ADRs) truthfully.
- The PR is focused. Unrelated changes are a finding.
- Nothing is committed that should not be: secrets, `.env` files, editor files, coverage reports, build artifacts.

### 10. Documentation

- If behavior or the HTTP contract changed, the spec, the README and the OpenAPI docs are updated in the same PR.
- ADR statuses and the ADR index are consistent (also enforced by CI).

## Output format (automated reviewer)

The automated reviewer posts **one summary comment** and **inline comments** for line-specific findings.

The summary comment uses this structure:

```markdown
## 🔍 Review: <one-line verdict>

**Verdict:** ✅ Ready to merge | 💬 Merge after addressing majors | ⛔ Blockers found
**Scope:** <PR type> · <files changed> · <dimensions that apply>

### Summary
<2–4 sentences: what the PR does, whether it achieves it, the most important finding>

### Findings
| # | Severity | Dimension | Location | Finding |
|---|---|---|---|---|
| 1 | 🟠 | HTTP contract | `internal/interface/controllers/crime_controller.go:42` | … |

### By dimension
<One short paragraph per applicable dimension. Use one line "N/A: <reason>" for the rest.>

### Questions for the author
<Optional: things that need context, not necessarily problems>
```

Inline comments start with the severity and the dimension, e.g. `🟠 HTTP contract: ...`, and include the fix when possible.
