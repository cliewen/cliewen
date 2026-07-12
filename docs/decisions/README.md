# Decisions

Architecture Decision Records (MADR format) with two-tier provenance: `inferred` (agent-reconstructed, not yet truth) → `verified` (human has accepted — the act that makes provenance auditable). Every ADR carries `author: agent|human` and `accepted-by:`. **Rejected alternatives are half of "why does the system look like this"** — rejected ADRs live here too and are never deleted.

<!-- clue:index:start -->
- [ADR-001 — Implementation language: Go](ADR-001-implementation-language.md) · `verified`
- [ADR-002 — The inbox is goals with status: proposed](ADR-002-inbox-is-proposed-goals.md) · `verified`
- [ADR-003 — Parse frontmatter with gopkg.in/yaml.v3](ADR-003-frontmatter-yaml-library.md) · `inferred` (awaiting human promotion)
- [ADR-004 — Default test-coverage gate at 80% total](ADR-004-coverage-gate-80-percent.md) · `verified`
<!-- clue:index:end -->
