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

3. **AC↔test contract and test-purpose taxonomy** (`corpus.checkACTests`, CH-003): `@AC-xxx` tags harvested from `criteria.md` bodies; test purposes harvested per [ADR-005](../../decisions/ADR-005-test-reference-convention.md)/[ADR-006](../../decisions/ADR-006-test-purpose-taxonomy.md) — framework-native tags are the general convention, and Go's `testing` has none, so the Go harvester reads function-name prefixes (`TestAC<digits>_…`, `TestUnit_…`, `TestSanity_…`, `TestArch_…`) from `*_test.go`, scanning the code tree while skipping `docs/`, `changes/` and hidden directories. Three checks: an AC in **active** criteria without a test fails (draft criteria are exempt — that is what draft means), a test referencing an undeclared AC fails, and a test declaring no purpose at all fails.

## Deliberate limits (doors, not gaps)

- **Positive+negative pair enforcement per AC** — needs a labeling convention on top of ADR-005; deferred until real use shows which labeling survives table-driven tests. Today the contract is ≥1 test per active AC.
- **Non-Go test harvesting** — post-baseline language profiles; `checkACTests`'s harvest step is the extension point.
- **Completed-plan immutability** — needs git diff context, not just a tree snapshot.
- **Layout linting beyond READMEs/indexes** — default-on profile rule per §4, later.

Exit codes: 0 valid · 1 issues · 2 usage error. Output is one `path: message` line per issue, sorted, so CI logs diff cleanly.
