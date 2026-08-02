# CH-104 tasks

- [ ] Write `docs/plans/P-010-adopters-keep-current.md`: frontmatter, campaign prose naming P-009 as predecessor, the M-045…M-050 milestone table, an explicit out-of-campaign section, and the standard mutation rules
- [ ] Record the campaign-opening decision as a dated row in `docs/decisions/log.md`, following the P-007/P-008/P-009 precedent that opening a plan is a log row and each milestone routes its own record
- [ ] Regenerate the plans index with `clue scaffold` so `docs/plans/README.md` lists P-010
- [ ] Verify: `clue validate --forbid-changes`, `go test ./...`, `go build ./...`, `go vet ./...`, `gofmt -l .`
