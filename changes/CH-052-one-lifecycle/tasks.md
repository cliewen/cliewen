---
id: CH-052-tasks
type: tasks
status: open
links: [CH-052]
title: CH-052 ordered tasks
---

# CH-052 — Tasks

- [ ] Write ADR-025 — one default lifecycle `draft → active → retired` plus named exceptions
- [ ] Write ADR-026 — unknown/adopter types validate against the default (carry the C-004 re-scoping argument)
- [ ] Restructure `internal/corpus/rules.go`: default lifecycle + exceptions map; `checkStatusVocab` falls back to default for any type
- [ ] Update `internal/corpus/corpus_test.go`: unknown-type positive (valid default status) and negative (bad status against default) cases
- [ ] Flip status `verified → active` on the 2 architecture and 7 analysis records; promote `core.md` `draft → active`
- [ ] Rewrite the status tables in `docs/README.md` and `internal/scaffold/templates/docs/README.md` as default + exceptions
- [ ] Align `guide/corpus.md` status wording if it enumerates vocabularies
- [ ] CHANGELOG entry; confirm no skill drift (`go generate ./internal/skills && git diff --exit-code`)
- [ ] Verify (`go test ./...`, `clue validate`, `npm run guide:build`), digest, mark M-017 done, delete workspace, ready PR
