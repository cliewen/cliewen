---
id: ADR-014
type: decision
status: verified
links: [ADR-010, ADR-011, ADR-012]
title: PR approval is ADR acceptance — the agent performs the clerical promotion
author: agent
accepted-by: Flemming N. Larsen (2026-07-13, PR #8)
---

# ADR-014 — PR approval promotes the PR's ADRs

## Context and problem statement

Agent-authored decisions are born `status: inferred` and become truth only when a human accepts them ([ADR-010](ADR-010-provenance-field.md)'s two-tier provenance). But the acceptance *mechanics* were a second manual act: after reviewing and approving the PR that introduces an ADR, the human still had to edit `status:` and `accepted-by:` by hand in some later change. The result is friction and lag — decisions the human has plainly accepted (the PR merged on their approval) sit `inferred` indefinitely, and the inferred count stops meaning "unverified" and starts meaning "unbookkept".

## Decision outcome

**Approving a PR is accepting the ADRs it introduces. The human's act is the approval; the agent performs the clerical flip.**

- **The acceptance event is the PR approval** — the same human gate that already verifies the change's meaning (the merge is the acceptance of the change; the approval covers the decisions the change records). An approval given in review comments or an explicit "approved, merge" instruction counts the same as the forge's approve button: what matters is the recorded human judgment, not the widget.
- **The agent does the bookkeeping**: on approval, flip `status: inferred → verified` and set `accepted-by:` to the approver, date, and PR (e.g. `Flemming N. Larsen (2026-07-13, PR #7)`) — in a final commit before merge, or in the next change's digest when the approval and merge were simultaneous. The provenance chain stays auditable: the named PR *is* the acceptance record.
- **Partial acceptance stays possible**: a reviewer who approves the PR but not an ADR in it says so in review — that ADR stays `inferred` (or the change is revised). Silence plus approval is acceptance; objection is objection, exactly as with any other reviewed line.
- **Nothing else self-promotes.** This covers `status` on decisions and the `provenance` field ([ADR-010](ADR-010-provenance-field.md)) alike: the flip always cites a specific human approval event. An agent still never promotes its own work absent one.

**Carrier:** the promotion paragraph in `docs/decisions/README.md` (ships as `clue init` template prose per [ADR-013](ADR-013-ships-generic-vs-repo-local.md)); the `clue-verify` checklist item already binds the two-tier rule and stays as written — only the mechanics of *who edits the file* changed.

Applied retroactively in the change that introduces this ADR: [ADR-011](ADR-011-version-stamping.md) (accepted with PR #6) and [ADR-012](ADR-012-release-notes-from-changelog.md) (accepted with PR #7) flip to `verified`. This ADR itself ships `inferred` and is promoted by its own PR's approval — the rule demonstrating itself.

### Rejected: humans edit the frontmatter themselves

The status quo. It doubles the human's acts (approve, then bookkeep) for zero added judgment — the judgment already happened in review — and in practice the second act lags or never comes.

### Rejected: ADRs born `verified` when a human is in the loop anyway

Destroys the distinction ADR-010 exists for: birth provenance (`author: agent`) and acceptance are different facts, and an unreviewed draft must never carry the accepted marker, even briefly.
