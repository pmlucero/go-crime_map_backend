# Architecture Decision Records

This directory records the significant architectural decisions of Crime Map Backend, following [ADR 0000](0000-record-architecture-decisions.md).

ADRs 0001 to 0007 were recorded on 2026-10-23, after the code. They document decisions already taken during development in 2025, with the alternatives and trade-offs that apply.

## Index

| # | Title | Status |
|---|---|---|
| [0000](0000-record-architecture-decisions.md) | Record architecture decisions | Accepted |
| [0001](0001-go-implementation-language.md) | Go as the implementation language | Accepted |
| [0002](0002-clean-architecture-layers.md) | Clean Architecture layers with inward dependencies | Accepted |
| [0003](0003-gin-http-framework.md) | Gin as the HTTP framework | Accepted |
| [0004](0004-postgresql-storage.md) | PostgreSQL as the storage engine | Accepted |
| [0005](0005-api-key-authentication.md) | API key authentication in the `X-API-Key` header | Accepted |
| [0006](0006-soft-delete.md) | Soft delete for crimes | Accepted |
| [0007](0007-uuid-public-identifiers.md) | UUIDs as public identifiers, numeric IDs internally | Accepted |

## Statuses

- **Proposed**: under discussion. The decision is not yet binding.
- **Accepted**: binding. Code must follow it.
- **Superseded by NNNN**: replaced by a later ADR. Kept for history and never deleted.
- **Rejected**: considered and discarded. Kept so the discussion is not repeated.

## Template

```markdown
# NNNN. Title in imperative or noun form

- Status: Proposed
- Date: YYYY-MM-DD

## Context

What forces are at play: requirements, constraints, unknowns.

## Decision

What we will do, stated plainly.

## Alternatives considered

### Alternative A
Pros / cons, and why it was not chosen.

## Consequences

### Positive
### Negative
### Risks and mitigations
```
