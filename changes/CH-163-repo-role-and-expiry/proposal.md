---
id: CH-163
type: change
status: open
links: [G-011, P-019, M-081, M-082, PDR-052, ADR-013, ADR-034, CAP-001, CAP-002]
title: A repository declares its Cliewen role, and a spent analysis is reported for retirement
---

# CH-163 — A repository declares its Cliewen role, and a spent analysis is reported for retirement

## Proposal

Serve P-019's M-081 and M-082 by making a repository's Cliewen role an observable fact, and by giving a spent analysis a reported route out of the corpus.

Nothing in a repository states whether it is Cliewen's own source repository or an adopter of it. Both carry `docs/`, `.clue/`, and `.agents/skills/`, so the distinction survives only as the incidental presence of `cmd/clue` and as prose in a routing hub that mixes shipped rules with repository-local ones. [ADR-013](../../docs/decisions/ADR-013-ships-generic-vs-repo-local.md) draws the boundary but leaves it entirely to human discipline: a rule can be recorded in this repository's corpus and never reach the shipped surface, or the reverse, and no check fails. An agent reading a hub cannot tell which of its rules bind where.

The same absence of a forgetting mechanism affects analysis. [ADR-034](../../docs/decisions/ADR-034-retirement-is-deletion.md) already establishes that retirement means deleting the file and that `status: retired` is not a resting state, but nothing applies that to a spike. An analysis exists to retire a risk; when its findings reach a durable artifact it is spent, and no signal says so.

This change adds a role marker under `.clue/`, teaches `clue init`, `clue migrate`, and the generated skills to read it, and adds a source-repository-only validator rule that a decision or constraint binding adopter behaviour must name a shipped carrier. It adds a declared destination to analysis frontmatter and a report-only migration that names a spent analysis together with the evidence that it is spent. It separates this repository's local rules from the shared routing text so the boundary is explicit rather than incidental.

## Scope boundary

Migration reports and never deletes: retirement remains a reviewed change under ADR-034, and `internal/migrate` gains no deletion code. `clue validate` gains no expiry rule, because an unexpired analysis is not an invalid corpus and the judge reads state rather than judgment ([ADR-044](../../docs/decisions/ADR-044-judge-reads-state-not-transitions.md)). The role marker records role alone: the Cliewen version already has two canonical carriers and a third copy would be the drift ADR-013 warns about. A repository with no marker is an adopter and is never blocked for lacking a file it had no way to write.

This change does not prune any corpus. Pruning this repository's own analysis, indexes, and overviews is separable work whose review is not served by being bundled with the mechanism that enables it, and the three adopter repositories are separate repositories with their own review boundaries.
