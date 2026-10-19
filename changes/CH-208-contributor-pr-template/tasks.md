---
id: CH-208-tasks
type: tasks
status: open
links: [CH-208]
title: Contributor-facing pull request template tasks
---

# Tasks

- [ ] Add criterion AC-239 to CAP-006 and amend C-024 to name the tracked template before behavior changes (AC-239).
- [ ] Write the short default template and move the full form to `.github/PULL_REQUEST_TEMPLATE/tracked.md`, in this repository and in the embedded scaffold (AC-239, AC-042, AC-048).
- [ ] Repoint the tests that read the full form at the tracked template and add positive and negative evidence for the short default; regenerate the evidence manifest (AC-239, AC-042, AC-048).
- [ ] Update the change-loop and review-boundary skill sources, regenerate skills, and update CAP-006 and CONTRIBUTING.md (AC-239).
- [ ] Trial in a disposable adopter repository: a pull request without a tracked change shows no maintainer-only sections and a tracked one reaches the brief gate (AC-239).
- [ ] Add the user-facing `[Unreleased]` changelog entry.
