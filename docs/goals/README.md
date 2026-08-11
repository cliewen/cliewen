# Goals

One file per goal: **who wants it, and why.** Goals are the top of the red thread — every capability traces back to one.

This folder is also the project's **inbox** (ADR-002): an idea or bug report enters the corpus as a goal with `status: proposed`. Promotion to `accepted` is a human decision made through a change/PR.

<!-- clue:index:start -->
- [G-001 — A verifiable thread from goal to acceptance evidence](G-001-verifiable-thread.md) · `accepted`
- [G-002 — clue and the skills carry versions](G-002-versioned-clue-and-skills.md) · `accepted`
- [G-003 — Cliewen is public](G-003-cliewen-is-public.md) · `accepted`
- [G-004 — Repository-local verification commands run as documented](G-004-local-verification-commands-run.md) · `accepted` — **Who wants it:** repository contributors (2026-08-09), after the documented `-flag=value` commands failed under PowerShell on Windows while the equivalent space-separated forms succeeded. PowerShell 7's default `Windows` native-argument mode splits a single-dash token containing `=` at the first dot in its value…
- [G-005 — A constraint's generated index badge shows its enforcement class](G-005-index-badge-uses-enforcement-not-status.md) · `proposed` — **Who wants it:** repository contributors and adopters reading `docs/constraints/README.md` (2026-08-09), found while M-067 (CH-141) minted four new constraints and their freshly generated index rows…
- [G-006 — Milestone IDs are covered by the corpus-wide identity ledger](G-006-milestone-ids-in-the-ledger.md) · `proposed` — **Who wants it:** repository contributors and agents allocating a new milestone number (2026-08-10), found while closing P-013 and opening P-014: `clue id next M` returns `M-001` despite M-001 through M-069 already existing…
<!-- clue:index:end -->
