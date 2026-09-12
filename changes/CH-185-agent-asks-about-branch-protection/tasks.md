---
id: CH-185-tasks
type: tasks
status: open
links: [CH-185]
title: Tasks — the agent asks whether the merge boundary is enforced
---

# CH-185 — Tasks

- [ ] Write `docs/decisions/PDR-059-the-agent-asks-whether-the-merge-boundary-is-enforced.md` with `binds: adopter`, citing `internal/skills/source/shared/review-boundary.md.tmpl` as its carrier.
- [ ] Add the rule to `internal/skills/source/shared/review-boundary.md.tmpl`, after the paragraph distinguishing displaying CI from enforcing it.
- [ ] Run `go generate ./internal/skills` and confirm the four generated references and four scaffolded copies carry the rule.
- [ ] Add AC-194 to `docs/capabilities/CAP-006-collaborative-handoffs/criteria.md`: protection adequate, so the agent marks ready without interrupting.
- [ ] Add AC-195: protection missing or inadequate, so the agent stops, states what is missing, and changes no setting unasked.
- [ ] Add AC-196: the human declines, so the pull request still reaches ready and the handoff states the boundary is unenforced with the reason given.
- [ ] Add AC-193: the host cannot be asked, so the agent reports that it does not know and never claims a verified pass.
- [ ] Link PDR-059 and the new criteria from `docs/capabilities/CAP-006-collaborative-handoffs/README.md`, and update its `design.md` where the handoff behaviour is described.
- [ ] Write the evidence for AC-193 through AC-196 in `internal/skills/` or `cmd/clue` as the carrier requires, each test declaring one purpose and its declared type and direction.
- [ ] Add the `[Unreleased]` CHANGELOG entry naming the generated skill text as the adopter-facing surface.
- [ ] Assess documentation impact: `guide/change-loop.md` and `guide/ci-wall.md` describe the enforcement boundary and may need the new step; state in the handoff what changed or why nothing did.
- [ ] Delete the change workspace as the digest.
