---
id: ADR-013
type: decision
status: verified
links: [P-002, ADR-008, ADR-011, ADR-012]
title: What ships to adopters is generic; AGENTS.md is the repo-local layer
author: agent
accepted-by: Flemming N. Larsen (2026-07-13, PR #8)
---

# ADR-013 — Ships-generic vs repo-local

## Context and problem statement

The skills and the folder READMEs are the adopter-facing product: they ship to every Cliewen repo, verbatim, under one version ([ADR-011](ADR-011-version-stamping.md)). But each repo also has legitimate local conventions — cliewen itself keeps a changelog ([ADR-012](ADR-012-release-notes-from-changelog.md)); another repo will not. A skill that names one repo's convention stops being the methodology and starts being that repo; a repo with no sanctioned place for local rules will fork the skills instead. Which artifacts are generic, where do local conventions live, and what happens when the two conflict?

## Decision outcome

**Every shipped artifact is generic; AGENTS.md is the one repo-local layer; local rules extend the methodology and never override it.**

- **Skills ship verbatim and stay generic.** A skill never names one repo's convention; where local behavior is legitimate, the skill carries a generic hook — "apply the repo-local conventions AGENTS.md declares" — and nothing more. The drift lint (ADR-011) holds the shipped set to one version, which only makes sense if the set is identical everywhere.
- **Generic also means self-contained: no corpus doc-IDs in shipped skill text.** In the repo a skill ships to, `ADR-008` or `G-002` resolves to nothing — or to that repo's own unrelated documents. A skill (including its mapping files) states each rule's content in its own text; the deciding document stays in this corpus, linked from here, never cited from there. This extends the existing rule that keeps doc-IDs out of the `clue` CLI's user-facing strings to the other exported surface, and the same test file enforces both (`TestSanity_SkillsCarryNoDocIDs`). Placeholder grammar (`CH-xxx`, `@AC-xxx`) is format, not reference, and stays.
- **Folder-README prose is generic; this repo's READMEs are the template sources.** What `clue init` (M-005) scaffolds into a new repo's `/docs` folders is the prose of cliewen's own folder READMEs — the MADR conventions, the carrier rule, the inbox rule already declare this per-file (ADR-002; the coverage-gate entry, since demoted to the decision log); this decision makes it the general rule. Index blocks between the `clue:index` markers are per-repo content: `init` emits them empty and `clue scaffold` (M-006) regenerates them from folder contents.
- **AGENTS.md is repo-local, rewritten per repo.** It routes agents into the corpus and declares the repo's own conventions (cliewen's rule 7 — the changelog — lives there, not in any skill). Extraction rewrites it; greenfield gets a template with routing and empty conventions.
- **Extends, never overrides.** An AGENTS.md rule that contradicts a skill is not resolved silently in either direction: it is a blocking question for `open-questions.md`, whose human answer becomes an ADR — a recorded exception or a fix to AGENTS.md. In brownfield extraction, pre-existing agent instructions are reconciled the same way: absorbed as local conventions when compatible, surfaced as open questions when not (`clue-extract`, contract item 7).
- **Refinement of [ADR-008](ADR-008-extraction-is-a-skill.md):** per-source mapping *sections* become per-source mapping *files* under the one `clue-extract` skill (`mappings/<source>.md`). The one-skill decision stands — a new source format is a new mapping file, not a new skill; the skill text every adopter reads stays lean.

**Carrier:** the skill texts and this repo's folder-README prose themselves (agent — they *are* what ships); the doc-ID sanity test over `.agents/skills/**` (machine); AGENTS.md as the declared local layer; the `clue init` templates once M-005 builds them (default — a door until then).

### Rejected: a `cliewen.yaml` configuration file for agent conventions

AGENTS.md already is the agent-facing configuration layer, and every agent runtime reads it natively. A second file splits routing into two sources of truth, and nothing machine-readable consumes it — §7: a field nobody reads gets removed before it exists. A machine config earns its place only when `clue` itself needs repo-local settings (existing doors: the skills path in ADR-011, changelog linting in ADR-012).

### Rejected: per-repo forked skill text

Letting each repo edit its skills makes every local rule invisible to the drift lint and unmergeable on skill upgrades — the exact failure G-002 exists to detect. Local needs go in AGENTS.md; the skills stay identical everywhere.
