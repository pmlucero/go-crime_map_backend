## Summary

<!-- What does this PR change, and why? Link the issue if there is one: "Closes #123". -->

## Type

<!-- Keep the ones that apply. -->

- Spec change (`docs/spec/`)
- ADR (`docs/adr/`)
- Implementation
- Tests only
- Docs / chore

## Traceability

<!-- Spec-Driven Development: implementation follows an approved spec. -->

- Spec sections: <!-- e.g. core.md §5.4 -->
- Acceptance criteria: <!-- e.g. AC-11, AC-19 -->
- ADRs: <!-- e.g. ADR 0006 -->

## Checklist

- [ ] Behavior is described in the spec, or this PR updates the spec in a separate commit
- [ ] Significant decisions (HTTP contract, data model, authentication, layers, dependencies) have an ADR
- [ ] Tests reference the acceptance criteria they verify (`// AC-xx`)
- [ ] `go test -race -p 1 ./...`, `go vet ./...` and `golangci-lint run` pass
- [ ] Schema changes update `schema.sql`, which still applies to an empty database (AC-26)
- [ ] Criteria made to pass are marked `Met` in spec §9
- [ ] No secrets, `.env` files or real API keys are committed
- [ ] Commits are small and follow Conventional Commits

## Notes for the reviewer

<!-- Trade-offs, open questions, what to look at first. -->
