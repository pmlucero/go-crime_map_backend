# 0005. API key authentication in the `X-API-Key` header

- Status: Accepted
- Date: 2026-10-23

## Context

The clients are applications acting for an organization, not end users (see [non-goals](../vision.md#4-non-goals)). Access must be revocable and time-limited.

## Decision

Protect every `/crimes` endpoint with an API key sent in the `X-API-Key` header. Keys are stored in the `api_keys` table with an owner (`user_id`), an `is_active` flag and an `expires_at` time. A Gin middleware rejects missing, unknown, inactive and expired keys with `401` (spec §4).

## Alternatives considered

### OAuth2 / JWT

The standard for end-user and delegated access. It adds an identity provider, token issuance and validation, with no benefit for application clients today.

### Bearer token in the `Authorization` header

Equivalent security. A dedicated header makes the intent explicit and avoids confusion with future OAuth2 support.

## Consequences

### Positive

- Simple to implement, to use with `curl`, and to revoke.

### Negative

- Keys are stored in plain text. A database leak exposes them.
- There is no API to create or revoke keys yet. Operators do it in SQL.

### Risks and mitigations

- *Plain-text keys.* Mitigation, planned: store a hash of the key instead (requires a superseding ADR).
- *The middleware is only tested with mocks.* A SQL dialect bug once broke all authentication against PostgreSQL. Mitigation: integration tests for the authentication path.
