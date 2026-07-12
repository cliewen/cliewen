# Decisions

Architecture Decision Records (MADR format) with two-tier provenance: `inferred` (agent-reconstructed, not yet truth) → `verified` (human has accepted — the act that makes provenance auditable). Every ADR carries `author: agent|human` and `accepted-by:`. **Rejected alternatives are half of "why does the system look like this"** — rejected ADRs live here too and are never deleted.

**Carrier rule for method decisions:** an ADR that changes the methodology *for adopting projects* must name its carrier — the `clue` rule (machine), the skill text (agent), or the init template (default) that ships it. A method decision without a carrier does not reach new projects and is not yet done. The foundation new projects receive has exactly one authoritative form: the output of `clue init` plus the rules of the `clue` binary — and CAP-001's criteria are what hold that output to account.

<!-- clue:index:start -->
- [ADR-001 — Implementation language: Go](ADR-001-implementation-language.md) · `verified`
- [ADR-002 — The inbox is goals with status: proposed](ADR-002-inbox-is-proposed-goals.md) · `verified`
- [ADR-003 — Parse frontmatter with gopkg.in/yaml.v3](ADR-003-frontmatter-yaml-library.md) · `verified`
- [ADR-004 — Default test-coverage gate at 80% total](ADR-004-coverage-gate-80-percent.md) · `verified`
- [ADR-005 — Tests reference ACs via framework-native tags; names where no tags exist](ADR-005-test-reference-convention.md) · `verified`
- [ADR-006 — Every test declares its purpose from a small taxonomy](ADR-006-test-purpose-taxonomy.md) · `inferred` (awaiting human promotion)
- [ADR-007 — AC lifecycle: meaning-immutable IDs, retirement by tombstone](ADR-007-ac-lifecycle.md) · `verified`
<!-- clue:index:end -->
