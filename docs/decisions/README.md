# Decisions

Three records live here ([PDR-006](PDR-006-decision-records-are-typed.md), [PDR-003](PDR-003-decision-log.md)), routed by two questions. Is reversing the decision cheap and local? Then it is a row in the **[decision log](log.md)** — one dated line per decision. Otherwise, is it about architecture or about how the project works? **ADRs** (Architectural Decision Records, in Nygard's strict sense) record decisions about the structure of the software and the corpus format; **PDRs** (Project/Process Decision Records) record decisions about how the project works — change tiers, acceptance mechanics, validation strategy. A decision that adopts a well-established practice cites it by name and records only the local why and deviations.

ADRs and PDRs share the MADR format and two-tier provenance: `inferred` (binding once merged, but no human has signed it) → `verified` (at least one human has explicitly approved — the act that makes provenance auditable). Every record carries `author: agent|human` and `accepted-by:`. **Rejected alternatives are half of "why does the system look like this"** — rejected records live here too. **Decisions are never deleted** — retention applies to the decision, not the file format: a record demoted under the litmus test survives as a dated log row, a record filed under the wrong type is renamed into the right series, and git history keeps the full provenance.

**Merge binds, approval signs** ([PDR-004](PDR-004-merge-binds-approval-signs.md), superseding [PDR-001](PDR-001-pr-approval-promotes-adrs.md)): merging a PR makes its `inferred` decisions binding — in force immediately, no signature required — but does not touch their status. Only an explicit human approval (review approval, review comment, or a stated "approved") flips a decision to `verified`: each approver signs `accepted-by:`, approvals accumulate, and the acceptance date is the first approval. An explicit objection keeps a decision `inferred` regardless of other approvals and becomes an open question; unresolved reviewer disagreement is an objection. The agent performs the clerical signing, citing approver, date, and venue.

**Records are timeless.** Context states the problem, not the episode: a motivating incident earns at most one sentence, and the change history lives in git log and the plans. Concrete mechanisms appear as decision content — the chosen option, the rejected options, the carrier — never as narrative.

**Carrier rule for method decisions:** a decision that changes the methodology *for adopting projects* — usually a PDR — must name its carrier: the `clue` rule (machine), the skill text (agent), or the init template (default) that ships it. A method decision without a carrier does not reach new projects and is not yet done. The foundation new projects receive has exactly one authoritative form: the output of `clue init` plus the rules of the `clue` binary — and CAP-001's criteria are what hold that output to account.

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
- [ADR-013 — What ships to adopters is generic; AGENTS.md is the repo-local layer](ADR-013-ships-generic-vs-repo-local.md) · `verified` · index-block clause superseded by ADR-019
- [ADR-017 — Prose conventions register as constraint artifacts with enforcement classes](ADR-017-conventions-are-constraints.md) · `verified`
- [ADR-018 — The init scaffolding is embedded in the clue binary](ADR-018-init-templates-embedded.md) · `verified`
- [ADR-019 — Index regeneration runs in clue init; ADR-013's emits-empty clause is superseded](ADR-019-init-regenerates-indexes.md) · `verified`
- [ADR-020 — The scaffolded register seeds only conventions without a versioned carrier](ADR-020-scaffolded-register-scope.md) · `verified`
- [ADR-021 — Skills are generated as standalone artifacts from shared canonical sources](ADR-021-generated-standalone-skills.md) · `verified`
- [PDR-001 — PR approval is decision acceptance; the agent performs the clerical promotion](PDR-001-pr-approval-promotes-adrs.md) · `verified` · superseded by PDR-004
- [PDR-002 — A light change tier: the PR description is the proposal](PDR-002-light-change-tier.md) · `verified`
- [PDR-003 — Records for the expensive-to-reverse; a decision log for the rest](PDR-003-decision-log.md) · `verified` · superseded by PDR-006
- [PDR-004 — Merge makes a decision binding; approval verifies it](PDR-004-merge-binds-approval-signs.md) · `verified`
- [PDR-005 — Validation requires foreign soil: trials on external repos, as findings not adoptions](PDR-005-foreign-soil-trials.md) · `verified`
- [PDR-006 — Decision records are typed: ADRs for architecture, PDRs for project/process](PDR-006-decision-records-are-typed.md) · `verified`
- [PDR-007 — The review boundary is real: changes root at main, one in flight per author, humans merge](PDR-007-review-boundary.md) · `verified`
- [PDR-008 — A declared plan revision may ride with its implementing change](PDR-008-plan-revisions-may-ride.md) · `verified`
- [Decision log](log.md) — dated rows for the cheap-to-reverse (ADR-003 and ADR-004 demoted here)
- [PDR-009-going-public](PDR-009-going-public.md)
<!-- clue:index:end -->
