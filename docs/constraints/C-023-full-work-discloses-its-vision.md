---
id: C-023
type: constraint
status: active
links: [ADR-067, PDR-054, PDR-073, VIS-001, CAP-009]
title: A tracked change's proposal states the vision it proceeds under
source: clue-delta and clue-verify skills (proposal and verification checklist), ADR-067 and PDR-073
enforcement: partial
binds: adopter
---

# C-023 — Tracked work discloses the vision it proceeds under

A tracked change's proposal names the active vision it serves, or states that the repository has none and that the change proceeds without one. A direct change makes no such statement, because it has no proposal and changes no accepted meaning.

The point is not the sentence. It is that a repository which has decided it does not need a stated direction says so where a human reads it, once per meaning-changing change, instead of being asked repeatedly by tooling that cannot tell a deliberate absence from an oversight.

**Checked by:** the `clue-verify` checklist item in `internal/skills/source/skills/clue-verify.md.tmpl` that asks whether the proposal names the vision, the intent rule in `internal/skills/source/shared/intent-model.md.tmpl`, and `clue validate --intent`, which states the corpus's vision and its provenance without being asked to interpret them.

**Residual:** whether the statement is there and true — whether the named vision is the one the change actually serves, and whether "proceeding without one" was a decision rather than a habit. No machine check reads the sentence; a proposal can name an active vision the change contradicts. The human's acceptance decision is where that is caught.
