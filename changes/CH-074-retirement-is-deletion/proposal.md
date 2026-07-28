---
id: CH-074
type: change
status: open
links: [M-027]
title: Retirement is deletion, and decisions die too
---

# CH-074 — Retirement is deletion, and decisions die too

M-027 of [P-007](../../docs/plans/P-007-core-hardening.md). Today an artifact "retires" by sitting on disk with `status: retired` — a state no artifact in the corpus has ever actually reached, because the convention that exists in practice is either a criteria tombstone (`@retired` tag, file stays so the test tag keeps failing loudly) or a decision demoted to a dated log row (file deleted, its content folded into `docs/decisions/log.md`). [ADR-025](../../docs/decisions/ADR-025-one-status-lifecycle.md) still names `retired` as every default-lifecycle type's terminal state, but nothing makes that state reachable once retirement is deletion — the file that would carry `status: retired` is gone.

This change makes the implicit convention explicit and generic: retiring any artifact means deleting its file in the same change; the successor (or, for a decision with no successor record, the log row) carries a machine-visible `supersedes:` frontmatter pointer naming the dead ID; Git history is the archive of the full retired text. The validator rejects a `supersedes:` entry whose ID still exists on disk (the retirement wasn't actually done) and — via the existing `checkLinks` dangling-reference rule — any live artifact still linking to an ID that no longer exists, so cleanup is enforced before merge rather than left to grep. Two exceptions stand, because they're not really "retirement" in this sense: criteria tombstones (a test tag must keep failing loudly, so the criterion's ID must stay findable in its file) and completed plans ([C-008](../../docs/constraints/C-008-completed-plans-immutable.md) keeps them frozen, never deleted).

This is a [C-013](../../docs/constraints/C-013-core-changes-need-decision.md) red-line change: it changes what a green `clue validate` asserts about retirement and revises [PDR-003](../../docs/decisions/PDR-003-decision-log.md)'s demotion mechanic and [ADR-025](../../docs/decisions/ADR-025-one-status-lifecycle.md)'s terminal-state claim. It carries its own decision record, ADR-034.
