---
id: C-023
type: constraint
status: active
links: [ADR-067, PDR-054, VIS-001, CAP-009]
title: A full change's acceptance brief states the vision it proceeds under
source: clue-delta skill step 5 (Verify, review, and propose for acceptance), and ADR-067
enforcement: partial
binds: adopter
---

# C-023 — Full work discloses the vision it proceeds under

A full change's acceptance brief names the active vision it serves, or states that the repository has none and that the change proceeds without one. A simple change makes no such statement, because it has no brief and changes no accepted meaning.

The point is not the sentence. It is that a repository which has decided it does not need a stated direction says so where a human reads it, once per meaning-changing change, instead of being asked repeatedly by tooling that cannot tell a deliberate absence from an oversight.

**Checked by:** the acceptance-brief line in `internal/scaffold/templates/github/pull_request_template.md`, whose `<!-- REQUIRED` placeholder the shipped validation workflow rejects on a pull request that is no longer a draft, and `clue validate --intent`, which states the corpus's vision and its provenance without being asked to interpret them.

**Residual:** whether the line is true — whether the named vision is the one the change actually serves, and whether "proceeding without one" was a decision rather than a habit. A brief can name an active vision the change contradicts and pass every check. The merge review is where that is caught.
