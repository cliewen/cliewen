---
id: CH-190
type: change
status: open
links: [P-023, M-097, G-015, PDR-062]
title: The first reusable operational discovery is captured in the home PDR-062 settled, and is reachable from where the work starts
---

# CH-190 — Capture the first reusable discovery

[PDR-062](../../docs/decisions/PDR-062-guidance-capture-has-an-eligibility-bar-and-no-new-type-yet.md) settled where a reusable operational discovery lives, how it is found, and what earns it a place before any capture obligation existed. [P-023](../../docs/plans/P-023-challenge-plans-and-retain-what-work-teaches.md)/M-097 is the first real test of that decision: a discovery this repository actually paid to find, filed into the home PDR-062 named, reachable without a new index.

This change files one discovery: PowerShell's backtick-as-escape-character behavior corrupting multi-line `gh`/git text built as a PowerShell double-quoted string. This repository hit it for real — PR #49's description ran headers together, turned `npm` into `\npm`, and turned a backtick-quoted commit hash into a form-feed character, because a backtick before `n` collapses to a newline and eats the letter, while a literal `\n` typed with intent to mean "newline" is not an escape in that context and survives as text. It meets PDR-062's eligibility bar: it cost a corrupted PR to find, and it is plausibly recurring for any agent on a Windows/PowerShell shell contributing here, since the full loop repeatedly builds PR bodies and commit messages.

**Home.** Not capability-specific and not product meaning — it is about running this repository's own contribution tooling (`gh`, `git`) from its documented shell. PDR-062 names exactly that case as `CONTRIBUTING.md`'s. The full loop already opens a PR in Propose and finalizes its body in Mark the Pull Request Ready, so a subsection there is reached without adding a new read path — PDR-062's "no new index" term, and AGENTS.md's source-repository-conventions table already sends an agent to `CONTRIBUTING.md` before verifying, meaning the file is opened whole rather than by section.

**Shape.** The milestone requires trigger, prerequisites, procedure, expected result, and recovery; PDR-062's evidence-sufficiency term requires stating what was actually observed rather than a universal claim. The entry states both, scoped to Windows PowerShell and `gh`/`git`, and invites a later session to confirm or narrow it under different conditions — which is M-098's job, not this change's.

**The workaround guard.** Before writing this down, the alternative was asked: can the confusing step be removed or automated instead? The corpus does not control which shell an agent runs in, and the discovery only reveals itself when someone chooses PowerShell for text with backticks or intended newlines, so there is no single point in this repository's own tooling to fix; documenting the recovery is the remaining option PDR-062 anticipated.
