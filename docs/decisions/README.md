# Decisions

Two records live here ([ADR-016](ADR-016-decision-log.md)): **ADRs** for decisions that are expensive to reverse — architecture, methodology mechanics, public contracts, AC semantics — and the **[decision log](log.md)** for everything else, one dated row per decision. Litmus test: if reversing it later is cheap and local, it's a log row; if it constrains future changes, it's an ADR.

ADRs are MADR format with two-tier provenance: `inferred` (agent-reconstructed, not yet truth) → `verified` (human has accepted — the act that makes provenance auditable). Every ADR carries `author: agent|human` and `accepted-by:`. **Rejected alternatives are half of "why does the system look like this"** — rejected ADRs live here too. **Decisions are never deleted** — retention applies to the decision, not the file format: an ADR demoted under the litmus test survives as a dated log row, with git history keeping its full text.

**The merge is the acceptance** ([ADR-018](ADR-018-merge-is-acceptance.md), superseding [ADR-014](ADR-014-pr-approval-promotes-adrs.md)): merging a PR accepts the `inferred` decisions it introduces — whoever merges accepts on behalf of the reviewing team, solo or not. An explicit review objection to a specific decision keeps it `inferred` through the merge and becomes an open question; unresolved reviewer disagreement is an objection. The agent performs the clerical flip — `status: verified`, `accepted-by:` naming the merger, date and PR — in the next digest, or in a final pre-merge commit when the merge was explicitly instructed.

**ADRs are timeless.** Context states the problem, not the episode: a motivating incident earns at most one sentence, and the change history lives in git log and the plans. Concrete mechanisms appear as decision content — the chosen option, the rejected options, the carrier — never as narrative.

**Carrier rule for method decisions:** an ADR that changes the methodology *for adopting projects* must name its carrier — the `clue` rule (machine), the skill text (agent), or the init template (default) that ships it. A method decision without a carrier does not reach new projects and is not yet done. The foundation new projects receive has exactly one authoritative form: the output of `clue init` plus the rules of the `clue` binary — and CAP-001's criteria are what hold that output to account.

<!-- clue:index:start -->
- [ADR-001 — Implementation language: Go](ADR-001-implementation-language.md) · `verified`
- [ADR-002 — The inbox is goals with status: proposed](ADR-002-inbox-is-proposed-goals.md) · `verified`
- [ADR-005 — Tests reference ACs via framework-native tags; names where no tags exist](ADR-005-test-reference-convention.md) · `verified`
- [ADR-006 — Every test declares its purpose from a small taxonomy](ADR-006-test-purpose-taxonomy.md) · `verified`
- [ADR-007 — AC lifecycle: meaning-immutable IDs, retirement by tombstone](ADR-007-ac-lifecycle.md) · `verified`
- [ADR-008 — Brownfield extraction is one generic skill with per-source mappings](ADR-008-extraction-is-a-skill.md) · `verified`
- [ADR-009 — AC IDs are namespaced: criteria declare an ac-prefix](ADR-009-ac-id-namespaces.md) · `verified`
- [ADR-010 — Extracted artifacts carry a provenance field, born inferred](ADR-010-provenance-field.md) · `verified`
- [ADR-011 — clue and the skills are versioned: tag-stamped binary, per-skill markers, drift is a failure](ADR-011-version-stamping.md) · `verified`
- [ADR-012 — Release notes are user-facing and come from CHANGELOG.md: extracted verbatim, missing section fails the release](ADR-012-release-notes-from-changelog.md) · `verified`
- [ADR-013 — What ships to adopters is generic; AGENTS.md is the repo-local layer](ADR-013-ships-generic-vs-repo-local.md) · `verified`
- [ADR-014 — PR approval is ADR acceptance; the agent performs the clerical promotion](ADR-014-pr-approval-promotes-adrs.md) · `verified` · superseded by ADR-018
- [ADR-015 — A light change tier: the PR description is the proposal](ADR-015-light-change-tier.md) · `verified`
- [ADR-016 — ADRs for the expensive-to-reverse; a decision log for the rest](ADR-016-decision-log.md) · `verified`
- [ADR-017 — Prose conventions register as constraint artifacts with enforcement classes](ADR-017-conventions-are-constraints.md) · `verified`
- [ADR-018 — The merge is ADR acceptance; objection blocks, teams included](ADR-018-merge-is-acceptance.md) · `inferred`
- [Decision log](log.md) — dated rows for the cheap-to-reverse (ADR-003 and ADR-004 demoted here)
<!-- clue:index:end -->
