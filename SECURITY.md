# Security policy

## Supported versions

Crime Map Backend has not been released yet. Once it is, only the latest minor release will receive security fixes.

## Reporting a vulnerability

**Do not open a public issue for security problems.**

Report vulnerabilities privately through GitHub: go to the [Security tab](https://github.com/pmlucero/go-crime_map_backend/security) and choose **Report a vulnerability**.

Please include:

- The affected version or commit.
- A description of the issue and its impact.
- Steps to reproduce, ideally with a minimal request. Do not include real API keys or private data.

You can expect an acknowledgement within 7 days. Fixes are coordinated privately and disclosed in a GitHub Security Advisory once a release is available. Reporters are credited unless they prefer otherwise.

## Scope

These areas are especially relevant:

- **Authentication bypass.** For example, reaching a `/crimes` endpoint without a valid, active and unexpired API key (spec §4).
- **SQL injection.** For example, through list parameters or path parameters.
- **Information disclosure.** For example, error responses that expose SQL, driver or stack details (spec §6), or soft-deleted crimes that become visible again ([ADR 0006](docs/adr/0006-soft-delete.md)).
- **API key exposure.** For example, keys written to logs or returned by an endpoint.

Out of scope:

- API keys stored in plain text in the database. This is a known and documented limitation ([ADR 0005](docs/adr/0005-api-key-authentication.md)).
- The default local development credentials in the README and configuration (`postgres`/`postgres`).
- Missing rate limiting (spec §10).
