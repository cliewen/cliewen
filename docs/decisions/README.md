# Decisions

Architecture Decision Records (MADR format) with two-tier provenance: `inferred` (agent-reconstructed, not yet truth) → `verified` (human has accepted — the act that makes provenance auditable). Every ADR carries `author: agent|human` and `accepted-by:`. **Rejected alternatives are half of "why does the system look like this"** — rejected ADRs live here too and are never deleted.

**Promotion rides the PR** ([ADR-014](ADR-014-pr-approval-promotes-adrs.md)): approving a PR accepts the ADRs it introduces; the agent then performs the clerical flip — `status: verified`, `accepted-by:` naming the approver, date and PR — before merge or in the next digest. A reviewer who approves the PR but not an ADR in it says so in review, and that ADR stays `inferred`.

**ADRs are timeless.** Context states the problem, not the episode: a motivating incident earns at most one sentence, and the change history lives in git log and the plans. Concrete mechanisms appear as decision content — the chosen option, the rejected options, the carrier — never as narrative.

**Carrier rule for method decisions:** an ADR that changes the methodology *for adopting projects* must name its carrier — the `clue` rule (machine), the skill text (agent), or the init template (default) that ships it. A method decision without a carrier does not reach new projects and is not yet done. The foundation new projects receive has exactly one authoritative form: the output of `clue init` plus the rules of the `clue` binary — and CAP-001's criteria are what hold that output to account.

<!-- clue:index:start -->
- [ADR-001 — Implementation language: Go](ADR-001-implementation-language.md) · `verified`
- [ADR-002 — The inbox is goals with status: proposed](ADR-002-inbox-is-proposed-goals.md) · `verified`
- [ADR-003 — Parse frontmatter with gopkg.in/yaml.v3](ADR-003-frontmatter-yaml-library.md) · `verified`
- [ADR-004 — Default test-coverage gate at 80% total](ADR-004-coverage-gate-80-percent.md) · `verified`
- [ADR-005 — Tests reference ACs via framework-native tags; names where no tags exist](ADR-005-test-reference-convention.md) · `verified`
- [ADR-006 — Every test declares its purpose from a small taxonomy](ADR-006-test-purpose-taxonomy.md) · `verified`
- [ADR-007 — AC lifecycle: meaning-immutable IDs, retirement by tombstone](ADR-007-ac-lifecycle.md) · `verified`
- [ADR-008 — Brownfield extraction is one generic skill with per-source mappings](ADR-008-extraction-is-a-skill.md) · `verified`
- [ADR-009 — AC IDs are namespaced: criteria declare an ac-prefix](ADR-009-ac-id-namespaces.md) · `verified`
- [ADR-010 — Extracted artifacts carry a provenance field, born inferred](ADR-010-provenance-field.md) · `verified`
- [ADR-011 — clue and the skills are versioned: tag-stamped binary, per-skill markers, drift is a failure](ADR-011-version-stamping.md) · `verified`
- [ADR-012 — Release notes are user-facing and come from CHANGELOG.md: extracted verbatim, missing section fails the release](ADR-012-release-notes-from-changelog.md) · `verified`
- [ADR-013 — What ships to adopters is generic; AGENTS.md is the repo-local layer](ADR-013-ships-generic-vs-repo-local.md) · `verified`
- [ADR-014 — PR approval is ADR acceptance; the agent performs the clerical promotion](ADR-014-pr-approval-promotes-adrs.md) · `verified`
<!-- clue:index:end -->
