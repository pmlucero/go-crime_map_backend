# Crime Map Backend

[![CI](https://github.com/pmlucero/go-crime_map_backend/actions/workflows/ci.yml/badge.svg)](https://github.com/pmlucero/go-crime_map_backend/actions/workflows/ci.yml)
[![License: Apache 2.0](https://img.shields.io/badge/license-Apache%202.0-blue)](LICENSE)

REST API in Go for reporting and querying geolocated crime incidents, built with
Clean Architecture, Gin and PostgreSQL.

## Features

- Crime incidents with geolocation and full address (street, number, city, province, country, ZIP).
- Paginated listing, lookup by UUID, status changes (`ACTIVE`, `INACTIVE`, `RESOLVED`) and soft delete.
- Aggregated statistics by type, status and address.
- API key authentication (`X-API-Key` header) with expiration and revocation.
- OpenAPI docs served with Swagger UI.
- Unit tests with mocks, plus integration tests against a real PostgreSQL database, run in CI.

## Architecture

```text
cmd/api/                    # Entry point: wiring and HTTP server startup
internal/
├── domain/                 # Entities, repository and use case contracts
├── usecases/               # Application business rules
├── interface/
│   ├── controllers/        # HTTP handlers (Gin)
│   ├── middleware/         # API key authentication
│   └── routes/             # Route registration
├── infrastructure/
│   ├── database/           # PostgreSQL connection and schema
│   ├── repositories/       # PostgreSQL repository implementations
│   └── server/             # HTTP server setup
└── security/               # API key storage and validation
docs/                       # Vision, spec, ADRs, plan, and the generated OpenAPI spec
migrations/                 # SQL migrations
```

Dependencies point inward: use cases depend only on domain interfaces, and
infrastructure provides the implementations.

## Getting started

Requirements: Go 1.23+ and PostgreSQL 14+ (with the `cube` and `earthdistance`
extensions, included in the official image).

```bash
# 1. Start PostgreSQL
docker run -d --name crime-map-db -p 5432:5432 \
  -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=crime_map postgres:16

# 2. Create the schema
psql postgres://postgres:postgres@localhost:5432/crime_map \
  -f internal/infrastructure/database/schema.sql

# 3. Create a user and an API key
psql postgres://postgres:postgres@localhost:5432/crime_map -c "
  INSERT INTO users (id, username, email) VALUES ('u1', 'demo', 'demo@example.com');
  INSERT INTO api_keys (key, user_id, expires_at) VALUES ('demo-key', 'u1', NOW() + INTERVAL '30 days');"

# 4. Run the API
go run ./cmd/api
```

The server listens on `http://localhost:8080`. Set `DATABASE_URL` to use another
database (default: `postgres://postgres:postgres@localhost:5432/crime_map?sslmode=disable`).

## API

All `/crimes` endpoints require the `X-API-Key` header. Interactive docs:
`http://localhost:8080/swagger/index.html`.

| Method | Path | Description |
| --- | --- | --- |
| `GET` | `/health` | Health check (public) |
| `POST` | `/crimes` | Report a crime |
| `GET` | `/crimes?page=1&limit=10` | List crimes (paginated) |
| `GET` | `/crimes/{uuid}` | Get a crime |
| `PATCH` | `/crimes/{uuid}/status` | Change its status |
| `DELETE` | `/crimes/{uuid}` | Soft delete |
| `GET` | `/crimes/stats` | Statistics by type, status and address |

```bash
curl -X POST http://localhost:8080/crimes \
  -H 'X-API-Key: demo-key' -H 'Content-Type: application/json' \
  -d '{
    "title": "Robbery at a store",
    "description": "Armed robbery",
    "type": "ROBO",
    "latitude": -33.123, "longitude": -64.349,
    "address": "San Martin", "address_number": "100",
    "city": "Rio Cuarto", "province": "Cordoba",
    "country": "Argentina", "zip_code": "5800"
  }'
```

## Testing

```bash
# Unit tests only need Go. Integration tests expect PostgreSQL on localhost:5432
# with the crime_map and crime_map_test databases and the schema applied.
createdb -h localhost -U postgres crime_map_test
psql postgres://postgres:postgres@localhost:5432/crime_map_test \
  -f internal/infrastructure/database/schema.sql

go test ./...
./coverage.sh   # coverage summary
```

Run a single test: `go test ./internal/interface/controllers/tests -run TestGetCrimeStats -v`.

Regenerate the OpenAPI docs after changing handler annotations:
`go install github.com/swaggo/swag/cmd/swag@latest && swag init -g cmd/api/main.go`.

## Design and status

The project follows Spec-Driven Development:

- [Vision](docs/vision.md): purpose, principles and non-goals.
- [Core spec](docs/spec/core.md): the API contract and its acceptance criteria, with which ones the code meets today.
- [ADRs](docs/adr/README.md): the architecture decisions and their rationale.
- [v1 plan](docs/plan/v1.md): the order in which the remaining gaps are closed.

## License

[Apache License 2.0](LICENSE)
