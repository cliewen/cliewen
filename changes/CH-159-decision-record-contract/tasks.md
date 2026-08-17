---
id: CH-159-tasks
type: tasks
status: open
links: [CH-159]
title: Tasks for CH-159
---

# Tasks

- [x] Add the M-075 acceptance scenarios before implementation: validator taxonomy and legacy rejection (AC-143), reviewed legacy-log migration (AC-144), initialized and scaffolded decision corpus (AC-145), brownfield extraction routing (AC-146), and generated-skill parity (AC-147).
- [x] Write the decision record that settles the subject-typed taxonomy, one routing test, compact record shape, legacy-row disposition, and complete live-carrier inventory.
- [x] Implement validation for the settled decision types and filenames, including positive coverage and rejection of legacy logs and out-of-taxonomy decision names (AC-143).
- [x] Implement an idempotent versioned migration that inventories every legacy row, refuses to guess classifications, blocks mutation until reviewed destinations are complete, and permits removal only when no durable row is lost (AC-144).
- [x] Update `clue init`, index scaffolding, and their canonical templates so a new corpus materializes only the settled decision-record contract (AC-145).
- [x] Update the canonical extraction contract and source mappings so brownfield decisions route to the settled types and a legacy log cannot be silently carried or discarded (AC-146).
- [x] Update the canonical lifecycle-skill sources and regenerate both distributed skill trees with focused parity evidence (AC-147).
- [ ] Classify every row in this repository's decision log, create or amend the reviewed destinations, repair all live references, retire LOG-001, and remove the legacy log without changing immutable completed-plan prose beyond link targets.
- [ ] Update current corpus explanations, public guidance, constraints and capabilities affected by the new contract, plus the adopter-facing `[Unreleased]` note.
- [ ] Run the repository-local verification required for every changed surface and confirm the migration is a no-op on the converted repository.
