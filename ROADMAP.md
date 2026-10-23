# Roadmap

## Done

- Clean Architecture layout: domain, use cases, interface and infrastructure layers.
- Crime CRUD: create, list (paginated), get by UUID, status change, soft delete.
- Aggregated statistics by type, status and address.
- PostgreSQL persistence with a geospatial index (`earthdistance`).
- API key authentication with expiration and revocation.
- OpenAPI docs with Swagger UI.
- Unit tests with mocks, PostgreSQL integration tests, and CI on GitHub Actions.
- Pre-commit hooks: `gofmt`, `goimports`, `golangci-lint` and test conventions.

## Next

- [ ] Apply the list filters already documented in OpenAPI (`type`, `status`, date range).
- [ ] Return `404` instead of `500` when a crime does not exist.
- [ ] Radius search around a point, using the existing geospatial index.
- [ ] Expose the API key management endpoints (`internal/api`) through the Gin server.
- [ ] Use a single migration tool (golang-migrate) instead of `schema.sql` plus `migrations/`.
- [ ] Load all configuration from environment variables.
- [ ] Run `golangci-lint` in CI.

## Later

- [ ] OAuth2 / JWT and role-based access control.
- [ ] Rate limiting and structured logging.
- [ ] Container image and deployment.
- [ ] Map frontend.
