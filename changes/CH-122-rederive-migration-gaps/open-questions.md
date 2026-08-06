---
id: CH-122-questions
type: open-questions
status: open
links: [CH-122, P-012, M-057, CAP-003, ADR-038]
title: Blocking questions for CH-122
---

# CH-122 — open questions

## Q1 — How should this repository seed the ledger when an unrelated legacy-carrier finding blocks `clue migrate --apply`?

`go run ./cmd/clue migrate` plans `MIG-008 .clue/id-ledger.yaml: seed the identity ledger with 149 live id(s) from the current corpus scan`. The required `go run ./cmd/clue migrate --apply` refuses to apply that safe change because the repository lacks `.github/workflows/clue.yml`, which it reports as `MIG-003 ... thin CI caller is missing; run clue init before migrating`.

The repository instead carries its current CI files under `.github/workflows/`, and M-057 does not say whether to add the legacy caller, change migration atomicity so unrelated findings do not prevent MIG-008, or defer the repository's ledger seed. Adding the caller or changing migration failure semantics could alter the CI or migration contract, so neither is assumed here.

**Decision needed:** authorize one of those routes, including whether this is M-057 scope or requires a separately scoped core change.
