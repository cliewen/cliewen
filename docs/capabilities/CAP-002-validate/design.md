---
id: CAP-002-design
type: design
status: active
links: [CAP-002, ADR-003]
title: Design for CAP-002 clue validate
---

# Design — CAP-002 `clue validate`

## Shape

`cmd/clue` (entry, stdlib `flag` only — no CLI framework until a second command earns one) and `internal/corpus` (scan + rules). Two phases:

1. **Scan** (`corpus.Scan`): walk `docs/` and `changes/`, split YAML frontmatter from body (CRLF-normalized first — Windows checkouts), parse via yaml.v3 ([ADR-003](../../decisions/ADR-003-frontmatter-yaml-library.md)), build the artifact graph keyed by `id`. Discovery is by frontmatter, never by path (Foundation §4).
2. **Rules** (`corpus.Validate`): core fields, duplicate IDs, status vocabulary per type (the `statusVocab` map in `rules.go` is the single source; the docs/README.md table mirrors it), link resolution (M-xxx harvested from plan bodies), README.md in every `/docs` folder, index-block integrity on taxonomy READMEs (`<!-- clue:index:start/end -->`: block links must exist, siblings and artifact-bearing subfolders must be referenced), and the `--forbid-changes` gate.

## Deliberate limits (doors, not gaps)

- **AC↔test traceability** — the rest of M-002; needs the test-tag convention. Test names already carry AC-IDs (`TestAC004_…`) as groundwork.
- **Completed-plan immutability** — needs git diff context, not just a tree snapshot.
- **Layout linting beyond READMEs/indexes** — default-on profile rule per §4, later.

Exit codes: 0 valid · 1 issues · 2 usage error. Output is one `path: message` line per issue, sorted, so CI logs diff cleanly.
