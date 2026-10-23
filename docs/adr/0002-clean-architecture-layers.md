# 0002. Clean Architecture layers with inward dependencies

- Status: Accepted
- Date: 2026-10-23

## Context

A CRUD API can be written as handlers that talk to the database directly. That is fast at first, but business rules end up spread across handlers and SQL, and they cannot be tested without HTTP and a database. The project also aims to be a reference for a layered Go service.

## Decision

Organize the code in Clean Architecture layers, with dependencies pointing inward:

- `internal/domain`: entities, and the repository and use case interfaces. It imports nothing from other layers.
- `internal/usecases`: business rules. It depends only on the domain.
- `internal/interface`: HTTP controllers, middleware and routes. It depends on use case interfaces.
- `internal/infrastructure`: PostgreSQL repositories, database setup and the HTTP server. It implements the domain interfaces.
- `cmd/api`: wiring only.

## Alternatives considered

### Flat package per feature

Simpler, with fewer files. Business rules would depend on Gin and SQL, and testing them would need HTTP and a database.

### Hexagonal architecture with explicit ports and adapters packages

Equivalent in spirit. Clean Architecture's layer names were chosen because they are widely recognized.

## Consequences

### Positive

- Use cases are unit-tested with mocked repositories.
- Storage or HTTP framework can change without touching business rules.

### Negative

- More packages and interfaces than the size of the API strictly needs.

### Risks and mitigations

- *Duplicate or dead layers accumulate.* The codebase currently has two use case implementations (`internal/usecases` and `internal/infrastructure/usecases`) and unused HTTP handlers. Mitigation: cleanup is scheduled in the [plan](../plan/v1.md), and reviews check that new code respects the layer rules.
