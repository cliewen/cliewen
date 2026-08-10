---
id: CH-143-tasks
type: tasks
status: open
links: []
title: Tasks
---

# Tasks

- [x] Confirm the current population of decision index rows carrying a live prose supersession or amendment claim, re-derived at head rather than read off AN-022's now-possibly-stale count.
- [ ] Write ADR-055, declining the widening a second time and stating the three costs the milestone requires: the obligation a superseding change would gain, how `clue validate` would distinguish a live superseded record from a stale one, and how the reverse "what was downstream" question stays unanswered without a reverse walk.
- [ ] State, in ADR-055, the settled prose shape for a decision index row's live supersession or amendment annotation, and apply it to the nine affected rows in `docs/decisions/README.md`.
- [ ] Add ADR-055's own index row and a carrier sentence in `docs/decisions/README.md`'s preamble, next to ADR-034's, naming that a live-but-superseded relationship stays prose rather than a frontmatter field.
- [ ] Close M-069 in `docs/plans/P-014-supersession-edge.md`'s milestone table: status `done`, evidence naming ADR-055.
- [ ] Run local verification (build, test, coverage, `clue validate --forbid-changes`, whitespace check) against the committed candidate.
- [ ] Run the agentic review loop via a context-isolated reviewer and resolve any blocking findings.
- [ ] Delete this change workspace as the digest and push the final commit.
