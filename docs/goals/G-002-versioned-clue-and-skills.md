---
id: G-002
type: goal
status: accepted
links: [G-001]
title: clue and the skills carry versions
---

# G-002 — clue and the skills carry versions

> Accepted 2026-07-13 with [P-002](../plans/P-002-leaves-home.md), which carries it as milestone M-004.

**Who wants it:** the maintainer (2026-07-12), prompted by the first real distribution friction — `go install` builds whatever the checkout has, and nothing tells an adopted repo whether its installed skills or binary have drifted behind cliewen's main.

**Why:** the carrier rule ships method decisions into adopting repos as binary rules and skill text. Without versions there is no way to ask "is model2diagram running current Cliewen?" — drift between the judge (`clue`), the guidance (skills), and the corpus conventions would be invisible until something breaks. A version on the binary (`clue --version`, stamped at release) and a version marker in each skill would make drift detectable — and lintable, if `clue validate` learns to compare.

**Door, not design:** how versions are stamped (release tags/ldflags), whether skills version individually or as a set, and whether drift is a warning or a failure are decisions for the change that accepts this goal — natural P-002 material alongside the release pipeline.
