---
id: CH-163-tasks
type: tasks
status: open
links: [CH-163]
title: Tasks for the repository role marker and analysis expiry
---

# Tasks

- [x] Record G-011, create P-019 with M-081 and M-082, record ADR-062 amending ADR-013 so a role marker is machine state rather than a second configuration layer, and record PDR-052 for analysis expiry.
- [x] Add AC-153 through AC-156 with focused positive and negative tests for role declaration, the source-only carrier rule, analysis destination, and expiry reporting.
- [x] Implement the `.clue/role.yaml` reader beside the existing ledger reader, have `clue init` write `role: adopter`, and commit this repository's own marker as `role: source` (AC-153).
- [x] Implement the source-repository-only validator rule that an adopter-binding decision or constraint names a shipped carrier (AC-154).
- [x] Add the analysis destination field and its validation (AC-155).
- [ ] Explain the destination field in the scaffolded analysis guidance and the corpus README's field list.
- [x] Add the report-only expiry migration after MIG-011, emitting a notice per spent analysis and never a finding or a change (AC-156).
- [ ] Separate this repository's local rules from the shared routing text, update the canonical skill sources for role and expiry, and regenerate the managed skills.
- [ ] Update the architecture and design overviews and affected capability material, regenerate indexes, and add the adopter-facing release note.
- [ ] Run focused and full verification, complete P-019's milestone digest, and prepare the reviewed pull-request handoff.
