---
id: CH-195
type: change
status: open
links: [PDR-042, CAP-006, CAP-004, AC-139, AC-142]
title: Change routes are named direct and tracked
---

# CH-195 — Change routes are named direct and tracked

## Proposal

This change is plan-less: no active plan has an unfinished milestone, and the change was requested directly by the maintainer.

Rename the two change routes from `simple` and `full` to `direct` and `tracked` in every live carrier, without changing what either route requires or when it applies.

The current names describe how big the work is, but the rule that picks a route is about something else: whether the accepted contract changes. A large refactor is "simple" and a one-line edit to an acceptance criterion is "full". Readers new to Cliewen take both names the wrong way. "Simple" sounds like it should be easy, and "full" sounds like the thorough choice a careful person should make. `direct` says how the work proceeds (do it and run the relevant checks) and is the default. `tracked` says what the work gets: it is followed from proposal through a CH workspace, digest, verification, and human acceptance.

`planned` was considered and rejected because "plan" is already a Cliewen artifact (`P-xxx`, `clue-plan`). A tracked change may be plan-less (C-005), and "a plan-less planned change" would mislead exactly the newcomers this change is for.

The change will:

- record a PDR that names the routes, states that the rename changes vocabulary and not route meaning, and fixes the compatibility rule for override trailers;
- retire AC-139 and AC-142 and mint successors that use the new names, because the recommendation an agent says aloud changes (`Recommended route: direct` / `tracked`);
- add a criterion that the shipped validation workflow treats an override as complete in either spelling (`Cliewen-Route: direct` with `Cliewen-Recommendation: tracked`, or the legacy `simple` with `full`), so adopter history and adopters still on older skills keep passing;
- add a `clue migrate` notice, in the MIG-006 pattern, reporting an adopter `AGENTS.md` hub that still names the old routes; migration never rewrites the hub (PDR-023);
- update the complete live-carrier inventory: canonical skill sources and generated skills, the `AGENTS.md` template and this repository's hub, the shipped PR template and workflow, the local CI scope script, live decisions, constraints, criteria, architecture and design overviews, `/guide`, `CONTRIBUTING.md`, and CLI text. Completed plans, analyses, and changelog entries stay as pinned history.

## Scope boundary

The routes keep their meaning: the same work goes to the same route, with the same obligations. Reference file names such as `change-scope-and-tiers.md` and decision file names such as `PDR-018-behavior-changes-remain-full.md` keep their paths, because renaming them would break links and adopter references for no gain to readers. Prose inside them is updated. Historical records are not rewritten.

## Challenge

**Riskiest assumption:** that new names help newcomers more than two vocabularies in flight hurt existing adopters. An adopter's own `AGENTS.md` keeps saying `simple`/`full` after an upgrade, because Cliewen never rewrites the hub. Their agents then read one pair of names in the hub and another in the skills.

**Credible alternative:** keep `simple`/`full` and only explain them better ("simple, the default"). That is cheaper and avoids the split, but it leaves names that describe size when the rule is about the contract, and every new reader would still have to unlearn the first impression.

**Cheapest test:** for each place the old names will survive (adopter hubs, adopter PR templates, past commit trailers), check that an agent or CI step reading the old spelling still reaches the right behavior. The trailer criterion and the migrate notice make this testable. If a surviving old spelling would make CI fail or route work differently, stop and revise before implementation goes further.

**Meets every criterion and still fails the reader:** every string is renamed but the guide and hub still lead with the mechanism, so a newcomer learns two new words without learning that the choice depends on whether the contract changes. The guide's first mention of each route therefore states the deciding question, not just the name.
