---
id: G-008
type: goal
status: proposed
links: [G-003]
title: The acceptance-brief gate distinguishes a change's tier from its CI scope
---

# G-008 — The acceptance-brief gate distinguishes a change's tier from its CI scope

**Who wants it:** contributors and agents taking a light Cliewen change to a ready pull request, in this repository and in any adopter running the shipped validation workflow (2026-08-11).

**Why:** the gate requiring a completed acceptance brief fires on the CI scope classifier's `full` output, which means "this diff needs the full check suite". Its failure message calls the pull request a "full Cliewen PR" — the Cliewen *tier* sense of the word, and a different thing. The check receives no tier signal at all, so every light change touching `docs/` is asked for a full change's artifact, and the adopter-facing workflow has the same shape. The gate binds only when a pull request is marked ready, so the failure appears after the draft's checks have passed, on a pull request that claims to be a candidate.

Underneath the message is an unanswered question: the change loop says the acceptance brief is what a *full* change fills, while the shipped workflow's own wording treats every Cliewen pull request as owing one. One of the two is wrong, and which one is a question about what the methodology requires rather than a defect in either file. Answering it decides whether the gate needs a tier signal or the loop needs to say that every Cliewen pull request carries a brief.

**Success looks like:**

- The gate and the change loop agree about which pull requests owe an acceptance brief, and the corpus records which reading was chosen and why.
- A contributor who is asked for a brief is told a reason that matches the rule they are being held to.
- A light change that legitimately owes no brief does not fail a required check, in this repository or in an adopter's.
