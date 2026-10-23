# Vision

- Status: Approved
- Last updated: 2026-10-23

> A small, well-built REST API to report geolocated crime incidents and query them, as a reference for Clean Architecture in Go.

## 1. Problem statement

Neighbors, local organizations and journalists often want to answer simple questions about safety in their area: what happened, where, and how often. The information is usually scattered across social media posts and news articles, and it cannot be queried or aggregated.

Crime Map Backend provides the data layer for that: a service where incidents are reported with a location and a type, and then listed, inspected and summarized.

It is also a public portfolio project. How it is built matters as much as what it does: the code is meant to be read as an example of a layered Go service with tests, documentation and a disciplined workflow.

## 2. Target users and core use cases

**Users:** client applications (a map frontend, scripts, integrations) that act on behalf of an organization holding an API key.

| Use case | What the API provides |
|---|---|
| **Report an incident** | Create a crime with a title, description, type, coordinates and full address. |
| **Browse incidents** | A paginated list, newest first, and the detail of a single incident. |
| **Maintain incidents** | Mark an incident as active or inactive, and remove it without losing its history (soft delete). |
| **Summarize** | Aggregated counts by type, status and address. |

## 3. Product principles

1. **The spec is the contract.** Behavior is defined in [spec/core.md](spec/core.md) before it is implemented, and the code is verified against its acceptance criteria.
2. **Simple and readable over clever.** A junior developer should be able to follow a request from the route to the database.
3. **Dependencies point inward.** The domain does not know about HTTP or PostgreSQL ([ADR 0002](adr/0002-clean-architecture-layers.md)).
4. **Nothing is lost by accident.** Deleting an incident hides it; it does not erase it ([ADR 0006](adr/0006-soft-delete.md)).
5. **Honest status.** Known gaps between the spec and the code are listed, not hidden.

## 4. Non-goals

- **Being an official crime record.** Reports are not verified, and the API does not claim they are.
- **Identifying people.** The data model has no fields for victims, suspects or witnesses, and none will be added.
- **End-user accounts.** Access is per organization, through API keys. User registration and login are out of scope.
- **Real-time features.** No push notifications, websockets or streaming.
- **The map frontend.** It is a separate project that consumes this API.

## 5. Long-term direction

After the current spec is fully met (see [plan/v1.md](plan/v1.md)), the most likely next steps are:

- Search by radius around a point, using the existing geospatial index.
- Self-service management of API keys through the API.
- A map frontend in a separate repository.

Each of them needs a spec change first.

## 6. Success criteria

- Every acceptance criterion in the spec is met and covered by a test that references it.
- CI runs vet, lint and the full test suite, including PostgreSQL integration tests, on every PR.
- A new reader can run the API locally and make an authenticated request in under 10 minutes by following the README.
