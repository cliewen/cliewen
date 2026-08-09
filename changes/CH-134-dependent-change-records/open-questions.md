---
id: CH-134-questions
type: open-questions
status: resolved
links: [CH-134, ADR-038, ADR-040]
title: Resolved questions for dependent-change records
---

# Open questions

## Q-01 — What is the determination for F2's emitted-wall divergence?

AN-013 establishes that the emitted CI wall functions only in an adopter repository and is kept in step manually. M-065 requires a named mechanism, a decline with a recorded cost, or a successor milestone. Should CH-134 design and implement a Cliewen mechanism for naming this foreign evidence, decline it with its cost, or route it to a named successor?

**Answer (Flemming N. Larsen, 2026-08-09, conversation):** ADR-038 is the named mechanism. Its immutable upstream reusable workflow and thin adopter-owned caller remove the copied validation wall whose divergence AN-013 measured. No second foreign-evidence mechanism is warranted.

**Resolved:** Yes. The determination is a re-derivation from the verified ADR-038, not a new contract.

## Q-02 — What is the determination for F3's unqualified external references?

AN-013 establishes that bare forge references lack repository identity and can resolve incorrectly. M-065 requires a named mechanism, a decline with a recorded cost, or a successor milestone. Should CH-134 adopt the proposed qualified-reference convention and lint, decline it with its cost, or route it to a named successor?

**Answer (Flemming N. Larsen, 2026-08-09, conversation):** ADR-040 is the named mechanism. Full URLs and `clue:` identities qualify external targets, while the deterministic judge rejects bare forge numbers and never resolves them. No successor work is warranted.

**Resolved:** Yes. The determination is a re-derivation from the verified ADR-040, not a new contract.
