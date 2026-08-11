---
id: G-004
type: goal
status: accepted
links: [G-001]
title: Repository-local verification commands run as documented
---

# G-004 — Repository-local verification commands run as documented

**Who wants it:** repository contributors (2026-08-09), after the documented `-flag=value` commands failed under PowerShell on Windows while the equivalent space-separated forms succeeded. PowerShell 7's default `Windows` native-argument mode splits a single-dash token containing `=` at the first dot in its value, so `go` received `-coverprofile=coverage` and `.out` as two arguments and wrote the profile to a file named `coverage`; the report against `coverage.out` then failed on a name nothing had created. The equals form is valid Go flag syntax and succeeds in Bash and `cmd.exe`, so the defect is in the documented block's dependence on shell argument handling, not in Go.

**Why:** the local verification block is the repository's contract with contributors and hosted CI. A command that fails solely because of invocation syntax makes an otherwise valid change look unverified and encourages improvised substitutions.

**Success looks like:**

- The documented local verification commands run on the supported contributor environments, or their supported boundary and exact alternatives are stated explicitly.
- The coverage report remains an informative check and the verification block stays usable verbatim.
