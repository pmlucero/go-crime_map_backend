# 0006. Soft delete for crimes

- Status: Accepted
- Date: 2026-10-23

## Context

Crime reports are a record of what was reported. Deleting one by mistake, or on purpose, should not erase it irrecoverably, and statistics over time should be explainable.

## Decision

Deleting a crime sets its `status` to `DELETED` and `deleted_at` to the current time. The row is kept. Every read filters by `deleted_at IS NULL` (spec §7), and an index on `deleted_at` supports that filter.

## Alternatives considered

### Hard delete

Simpler queries and no hidden rows. Data is lost, and there is no way to audit or restore it.

### Separate archive table

Keeps the main table small. It needs a move operation in a transaction, and the archive has to follow every schema change.

## Consequences

### Positive

- Deletions can be audited and reverted by an operator.

### Negative

- Every query must remember the `deleted_at IS NULL` filter.
- `status` and `deleted_at` carry overlapping information.

### Risks and mitigations

- *A query forgets the filter and leaks deleted crimes.* Mitigation: spec AC-12, AC-13 and AC-21 are tested.
