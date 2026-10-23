# 0001. Go as the implementation language

- Status: Accepted
- Date: 2026-10-23

## Context

The service is a JSON API over a relational database. It has to be easy to read for a wide audience, cheap to run, and simple to deploy. It is also a vehicle to practice Go in a realistic, layered codebase.

## Decision

Implement the service in **Go**, using the standard library where it is enough (`net/http`, `database/sql`, `context`) and a small set of well-known modules for the rest (Gin, `lib/pq`, `sqlx`, `swag`, `testify`).

## Alternatives considered

### Java with Spring Boot

Mature and familiar to the maintainer, with excellent tooling. It brings a heavy runtime and a lot of framework behavior that hides what the code does, which works against the goal of a readable reference.

### Node.js with TypeScript

Fast to iterate and with a large ecosystem. Runtime type safety depends on extra validation layers, and the concurrency model is less explicit than goroutines and `context`.

## Consequences

### Positive

- A single static binary, fast startup and low memory use.
- Explicit error handling and context propagation make request flows easy to follow.

### Negative

- More boilerplate than a full framework, for example in SQL row scanning.

### Risks and mitigations

- *Dependencies creep in.* Mitigation: each new module dependency is justified in its PR, and in an ADR when significant (ADR 0000).
