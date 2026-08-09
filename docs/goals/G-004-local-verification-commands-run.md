---
id: G-004
type: goal
status: proposed
links: [G-001]
title: Repository-local verification commands run as documented
---

# G-004 — Repository-local verification commands run as documented

**Who wants it:** repository contributors (2026-08-09), after the documented `go tool cover -func=coverage.out` command failed under Go 1.26.5 on Windows while the equivalent space-separated form succeeded.

**Why:** the local verification block is the repository's contract with contributors and hosted CI. A command that fails solely because of invocation syntax makes an otherwise valid change look unverified and encourages improvised substitutions.

**Success looks like:**

- The documented local verification commands run on the supported contributor environments, or their supported boundary and exact alternatives are stated explicitly.
- The coverage report remains an informative check and the verification block stays usable verbatim.
