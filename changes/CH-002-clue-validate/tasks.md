---
id: CH-002-tasks
type: tasks
status: open
links: [CH-002]
title: Task breakdown for CH-002
---

# Tasks — CH-002

- [ ] Go module `github.com/cliewen/cliewen`, `cmd/clue`, `internal/corpus`
- [ ] Frontmatter scan + artifact graph (yaml.v3, CRLF-safe)
- [ ] Rules: core fields, unique IDs, link resolution (incl. M-xxx),
      status vocab, README presence, index integrity, --forbid-changes
- [ ] Unit tests per rule (fixture corpora) + self-test on this repo
- [ ] CAP-002-validate: README / criteria (AC-004..AC-008) / design
- [ ] ADR-003: frontmatter parsed with yaml.v3 (author: agent, inferred)
- [ ] CI workflow: job `validate` on pull_request + push main
- [ ] Digest: docs indexes updated, status vocab table extended,
      M-002 → in-progress, delete /changes/CH-002-clue-validate/
- [ ] Push branch, open PR, CI green — human reviews and merges
