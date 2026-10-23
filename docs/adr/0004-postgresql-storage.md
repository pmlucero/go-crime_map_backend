# 0004. PostgreSQL as the storage engine

- Status: Accepted
- Date: 2026-10-23

## Context

Crimes have structured fields, are filtered and aggregated, and have coordinates that will be searched by distance (see [vision](../vision.md#5-long-term-direction)).

## Decision

Store data in **PostgreSQL** (14 or later), accessed with `database/sql`, `sqlx` and the `lib/pq` driver. Use the `cube` and `earthdistance` extensions for a GiST index on `ll_to_earth(latitude, longitude)`, ready for radius search.

## Alternatives considered

### PostGIS

The most capable option for geospatial data. It is a larger extension than needed for distance search over points.

### SQLite

No server to run, which is convenient for development. It lacks the geospatial index, and concurrent writes are limited.

### MongoDB

Has geospatial queries built in. The data is relational and aggregations are naturally expressed in SQL.

## Consequences

### Positive

- Reliable, well-known, with good aggregation support.
- The official Docker image includes the required extensions.

### Negative

- Integration tests need a running PostgreSQL. CI starts one as a service.

### Risks and mitigations

- *The schema and the code drift apart.* This happened before: indexes referred to renamed columns. Mitigation: CI applies `schema.sql` to a fresh database and runs the integration tests (spec AC-26).
