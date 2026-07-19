## Summary

<!-- What changes, and why is it needed? -->

## Proposal

- Change ID: `CH-xxx`
- Change tier: <!-- full or light -->
- Plan item served: <!-- P-xxx / M-xxx, or explicitly plan-less -->
- Proposal location: <!-- /changes/CH-xxx-slug/proposal.md for full changes; this PR description for light changes -->

## Traceability and Decisions

- Acceptance criteria or capability meaning changed: <!-- links, or none -->
- Decision records added or changed: <!-- links, or none -->
- Constraints and quality scenarios assessed: <!-- links and effects, or none -->

## Verification

<!-- List the commands run and their results. Include focused positive and negative evidence for every active acceptance criterion affected by the change. -->

## Checklist

- [ ] The branch started from the current tip of `main`, and this is the author's only change in flight.
- [ ] The proposal was committed before implementation, or this is a correctly scoped light change whose proposal is this PR description.
- [ ] The plan item or plan-less declaration is truthful, and all artifact links resolve.
- [ ] Decisions are recorded at the right tier, and active constraints and quality scenarios were assessed.
- [ ] Changed active acceptance criteria have positive and negative executable evidence.
- [ ] User-visible impact is described under `[Unreleased]` in `CHANGELOG.md`, or the change has no user-visible impact.
- [ ] Full-change tasks are complete, plan bookkeeping is current, and no transient `/changes/` workspace remains.
- [ ] Generated artifacts were regenerated from their canonical sources where applicable.
- [ ] `go build ./...`, coverage-gated `go test ./...`, `go run ./cmd/clue validate --forbid-changes`, and `git diff --check` pass.
- [ ] The pull request is ready for human review and merge; no agent will merge it or push to `main`.
