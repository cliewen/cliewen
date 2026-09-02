---
id: PDR-052
type: decision
status: verified
links: [G-011, P-019, M-082, ADR-034, ADR-044, ADR-062, PDR-046, CAP-001, CAP-002]
title: An analysis declares where its findings landed, and a spent one is reported rather than retained
binds: adopter
author: agent
accepted-by: Flemming N. Larsen (2026-09-02)
---

# PDR-052 — A spent analysis is reported, not retained

## Context

An analysis exists to retire a risk. [ADR-034](ADR-034-retirement-is-deletion.md) already settled that retirement means deleting the file and that `status: retired` is not a state a file rests in, but nothing applies that to a spike, so every analysis stays `active` indefinitely. The cost is not only volume: a spike's value is a measurement pinned to the revision it names, and a corpus that keeps every such measurement forever asks each later reader to work out which ones still describe the system.

The obvious test — the plan it served is complete — does not work alone. Campaign completion says the campaign ended; it says nothing about whether the findings ever reached durable form. In a corpus whose plans are all completed, that test sweeps up every analysis at once.

## Decision

**An analysis declares the durable artifacts now carrying its findings in a `carried-by:` field, and an analysis is reported as spent only when every plan it serves is completed *and* every carrier it declares resolves.** The second condition is the one no tool can infer, so a human declares it; the field reuses the corpus's existing carrier vocabulary rather than introducing a second word for the same idea.

**`clue migrate` reports a spent analysis and never acts on it.** The report is a non-blocking notice: it becomes neither a finding nor a write, migration gains no deletion code, and retirement remains a reviewed change under ADR-034 with Git history as the archive. This follows the asymmetry the existing migrations already keep — anything touching authored prose is reported, and only mechanically provable bytes are rewritten.

**`clue validate` gains no expiry rule.** An unexpired analysis is not an invalid corpus, and the judge reads state rather than judgment ([ADR-044](ADR-044-judge-reads-state-not-transitions.md)). Validation checks only that a declared carrier is honest: the field belongs to analysis, names at least one artifact, and names something that resolves and is not itself another spike. A carrier that does not resolve is worse than an absent field, because it claims the findings survived somewhere they did not.

The field is optional. Requiring it on every existing analysis would turn a corpus green yesterday red today for a field its authors could not have written.

## Rejected: hollow out the document instead of deleting it

Trimming a spent analysis to a stub keeps the file, its index row, and its link targets while removing the only part with value. The navigation cost — which is what a large corpus actually charges a reader — survives intact, and ADR-034 already chose deletion for exactly this reason.

## Carrier

The `carried-by:` field and its validation, the spent-analysis derivation, the report-only migration, the spike-retirement instruction in `internal/skills/source/skills/clue-analysis.md.tmpl` and the `MIG-013` walk in `internal/skills/source/skills/clue-upgrade.md.tmpl`, the scaffolded guidance in `internal/scaffold/templates/docs/analysis/README.md`, and the adopter-facing changelog entry.
