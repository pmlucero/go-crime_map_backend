# 0003. Gin as the HTTP framework

- Status: Accepted
- Date: 2026-10-23

## Context

The API needs routing with path parameters, JSON binding with validation, middleware for authentication, and generated OpenAPI documentation.

## Decision

Use **Gin** for routing, binding (`binding:"required"` tags), middleware and responses, and **swag** with `gin-swagger` to generate and serve the OpenAPI documentation from handler annotations.

## Alternatives considered

### Standard library `net/http`

Since Go 1.22 it supports method and path patterns, which covers routing. Binding, validation and Swagger integration would have to be assembled by hand.

### Echo or Chi

Comparable options. Gin was chosen for its popularity, its documentation and its `swag` integration.

## Consequences

### Positive

- Little code for binding, validation and middleware.
- Swagger UI served at `/swagger/index.html`.

### Negative

- Controllers depend on `gin.Context`, so they are tied to Gin. Use cases are not (ADR 0002).
- `binding:"required"` rejects zero values, so a latitude or longitude of `0` is rejected (spec AC-07).

### Risks and mitigations

- *Generated docs drift from the routes.* Mitigation: regenerate them with `swag init` whenever handler annotations change, and review `docs/` in PRs.
