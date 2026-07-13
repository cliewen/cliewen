---
id: ADR-011
type: decision
status: inferred
links: [G-002, CAP-004, P-002]
title: clue and the skills are versioned — tag-stamped binary, per-skill markers, drift is a failure
author: agent
accepted-by: —
---

# ADR-011 — Version stamping for clue and the skills

## Context and problem statement

[G-002](../goals/G-002-versioned-clue-and-skills.md) asks that `clue` and the skills carry versions so drift between the judge (the binary), the guidance (the skills), and the corpus conventions is detectable — and lintable. The goal deliberately left three doors for the accepting change ([CH-007](../plans/P-002-leaves-home.md)): how versions are stamped, whether skills version individually or as a set, and whether drift is a warning or a failure.

## Decision outcome

**Versions come from git tags via ldflags; skills carry a per-skill frontmatter `version:` kept consistent as a set; drift is a build failure, with `dev` builds exempt.**

- **Stamping.** Release tags are conventional semver with a `v` prefix (`vX.Y.Z`, as Go module tagging requires). The release workflow strips the `v` and injects the bare semver into `cmd/clue`'s `var version` with `-ldflags "-X main.version=<semver>"`. An unstamped binary falls back to the module version Go embeds in `go install module@vX.Y.Z` builds (`v` trimmed), so tag installs self-stamp; checkout builds and pseudo-version (branch/commit) installs report `dev`. `clue version` / `clue --version` print it.
- **Granularity: per-skill markers, enforced as a set.** Each `.agents/skills/<name>/skill.md` declares its own `version:`; there is no separate set-file. `clue validate` enforces that every skill carries a stamp and that they all agree — so per-skill markers behave as a set without a second source of truth. Bare-semver strings match the binary's, so the drift comparison is string-equal on both sides.
- **Drift is a failure, `dev` exempt.** The wall philosophy (Foundation §2): a machine-checkable rule that only warns gets ignored, so a released `clue` whose skills differ from the binary exits non-zero. `dev`/source builds have no release to drift from, so they skip the binary comparison (they still require the stamps and their mutual agreement).

**Carrier:** `clue version` and `corpus.checkSkillVersions` (machine); the release workflow (default/CI); the `version:` line each skill's frontmatter must carry (agent-maintained, per [CAP-004](../capabilities/CAP-004-ship/README.md)).

### Rejected: a single `.agents/skills/VERSION` set-file

One file is marginally less to edit at release, but it splits "what version is this skill" from the skill itself — an adopter copying a single skill loses its version, and the marker no longer travels with the text it stamps. Per-skill frontmatter keeps the stamp on the artifact; the set-consistency rule recovers the one-number convenience without the second source of truth.

### Rejected: drift as a warning

A warning that does not fail CI is invisible the moment logs scroll — exactly the "machines enforce form" failure mode the methodology exists to prevent. If a release must ship with a known skill lag, that is a decision to record, not a check to soften.

### Deferred: skills outside `.agents/skills/`

The drift rule looks only under `.agents/skills/`; an adopter who relocates skills gets no drift check until the path is made configurable. A door noted in CAP-004's design, not part of this decision. (Stamping bare `go install` was originally deferred here too; the build-info fallback closed it — only checkout and pseudo-version builds still report `dev`.)
