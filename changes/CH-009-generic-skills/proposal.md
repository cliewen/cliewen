---
id: CH-009
type: change
status: open
links: [P-002, G-001, ADR-008]
title: Skills ship generic — AGENTS.md is the repo-local layer; PR approval promotes ADRs
---

# CH-009 — Skills ship generic; AGENTS.md is the repo-local layer

## What

Serves **[P-002](../../docs/plans/P-002-leaves-home.md)** at campaign level — adoptability is P-002's theme, but no single milestone's exit criterion covers the adopter-facing surface being repo-agnostic, so this change binds to the plan, not a milestone. Review of CH-008 surfaced the principle; this change records and applies it:

1. **ADR-013 — what ships to adopters is generic; AGENTS.md is the repo-local layer.** Skills ship verbatim and never carry one repo's conventions (generic hooks reference AGENTS.md; AGENTS.md extends the methodology, never overrides it — conflicts are open questions). Folder-README prose is generic methodology and this repo's own READMEs are the `clue init` template sources; index blocks are per-repo and emit empty at init. Rejected: a `cliewen.yaml` for agent conventions; per-repo forked skill text. Refines [ADR-008](../../docs/decisions/ADR-008-extraction-is-a-skill.md): per-source mapping *sections* become per-source mapping *files* under the one skill.
2. **`clue-extract` applies it**: the OpenSpec mapping moves verbatim to `.agents/skills/clue-extract/mappings/openspec.md`; `skill.md` keeps the generic target contract plus a pointer. Contract item 7 (routing) gains the reconciliation rule: pre-existing agent instructions are absorbed as repo-local conventions when compatible, surfaced as open questions when they conflict — never left to contradict the skills.
3. **ADR-014 — PR approval promotes the PR's ADRs.** Manually editing `status: inferred → verified` after the fact is friction that delays provenance. The human's approval of the PR that introduces an ADR *is* the acceptance event; the agent then performs the clerical flip (`status: verified`, `accepted-by:` naming the approver, date, and PR) before merge or in the next change's digest. Applied retroactively in this change: [ADR-011](../../docs/decisions/ADR-011-version-stamping.md) (PR #6) and [ADR-012](../../docs/decisions/ADR-012-release-notes-from-changelog.md) (PR #7) flip to `verified` citing today's approvals. ADR-014 itself ships `inferred` and is promoted by this PR's own approval — the rule demonstrates itself.
4. **Decisions README**: the ADR style rule from CH-008's review (ADRs are timeless — Context states the problem, not the episode) and the promotion mechanics per ADR-014.
5. **ADR-011 timelessness trim**: the "originally deferred, later closed" narrative parenthetical becomes a plain statement of current behavior (meaning unchanged).

## Why

The skills are the adopter-facing product; CH-008 briefly leaked a cliewen-only convention (`CHANGELOG.md`) into them, and `clue-extract` carries a source mapping that crowds the generic contract. One recorded principle prevents the next leak, and the reviewer's question — "what if a brownfield AGENTS.md works against Cliewen?" — needs a corpus answer, not a review-thread answer. The ADR-promotion friction is the same theme: the PR is already the human gate, so acceptance should ride on it instead of demanding a second manual act.

## Out of scope

`clue init` itself (M-005); a `clue` lint for changelog or AGENTS.md content; retiring ADR-008 (ADR-013 refines its carrier detail, the one-skill decision stands); the AGENTS.md skills table gaining a `clue-extract` row (worth a look at M-005, when the routing template is built).
