---
id: CH-208
type: change
status: open
links: [P-025, CAP-006, C-024]
title: The default pull request template does not mislead an outside contributor
---

# CH-208 — The default pull request template does not mislead an outside contributor

Serves milestone M-113 of plan P-025. It proceeds under draft, inferred VIS-001 without confirming it.

## What and why

`clue init` installs one pull request template, and it is the full maintainer template: an acceptance brief, a proposal block and about twenty checklist items about reviews, ledgers and evidence. Observed in the Tank Royale repository: a first-time contributor opening a registration pull request was handed all of it, none of which they could answer.

The change makes the default template a short contributor-facing description (what changes and why, how it was checked) and moves the full tracked-change form to a named template, `.github/PULL_REQUEST_TEMPLATE/tracked.md`. The change-loop skill tells the agent to open a tracked pull request from that file. The text of the tracked form is unchanged, so the acceptance-brief gate and every criterion about the brief's wording still hold, now against the new file.

## Challenge the commitment

The riskiest assumption was that the hosted gate keys on the template, so a plain pull request without a brief would fail it. Reading the shipped workflow shows otherwise: "Require a completed acceptance brief" runs only when scope detection finds a change proposal in the pull request's history, so a pull request without a tracked change is never asked for a brief. The split is therefore safe for the gate.

The credible alternative is to keep one template and fence the maintainer sections behind a "delete everything below" instruction or an HTML comment. It is cheaper but leaves the contributor looking at the checklist, which is the failure observed. A second risk is that an agent forgets to use the tracked file and a tracked pull request goes out with the short body; the gate then fails with its existing message, so the miss is visible and cheap to repair, and the skill names the file where the pull request is opened.

An implementation could meet every criterion and still fail the contributor if the short template were itself jargon-heavy, so its wording is checked against the readable-references rule and by opening a pull request in a disposable adopter repository.

The work stops and returns for a decision if the disposable-repository trial shows the forge ignoring the named template or the short default still tripping the gate.

## Carriers and proof

New criterion AC-239 in CAP-006 states that the default template is contributor-facing and the full form lives in the named tracked template. AC-042 and AC-048 keep their meaning; their tests now read the tracked template. The shipped copy of each template equals this repository's, as before. Constraint C-024 names the header-less pull request template files, so it also names the new one.

Carriers: the two embedded templates and this repository's two copies, the scaffold materializer if it lists files, the change-loop and review-boundary references in the skill source and the regenerated skills, CAP-006 criteria and design, C-024, CONTRIBUTING.md where it names the template, and a user-facing changelog entry. Adopters who already have a template keep it: it is adopter-owned after arrival and no migration rewrites it, which the changelog says. The acceptance policy, the gate and the meaning of every other criterion are unchanged.
