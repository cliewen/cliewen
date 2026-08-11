---
id: ADR-041
type: decision
status: verified
links: [CAP-002, CAP-005, ADR-017, ADR-019, ADR-035, ADR-046, C-004, C-013, C-016]
title: Generated index rows state their record, and rows that state only their link are counted
author: agent
accepted-by: Flemming N. Larsen (2026-08-08, conversation)
---

# ADR-041 — Generated index rows state their record, and rows that state only their link are counted

## Context and problem statement

Index generation appended every missing entry as the target's filename with the extension removed: `- [ADR-039-versioned-corpus-migrations](ADR-039-versioned-corpus-migrations.md)`. The label restates the link and adds nothing — no title, so a reader must open the file to learn what the record says; no status, so provenance is invisible in the one place that lists every record in a folder.

These rows are not written by anyone. They accumulate, one per new artifact, because the generator emits them and nothing objects. In this repository twenty of them had collected in the decisions index before a change repaired them by hand, and three more remained in the goals, architecture, and analysis indexes; the same blocks ship through `clue init` and `clue scaffold`, so the same accumulation runs in every adopted repository.

Nothing checked row content either. `checkIndexes` verifies that links resolve and that every artifact is covered, and any label at all passes.

The repair for the generator is obvious. The harder question is what a judge should say about the rows already sitting in repositories — including rows the tool itself wrote there.

## Decision outcome

**A generated index row states the record it links — `id — title` from the artifact's own frontmatter, then its status. A row whose label restates only its own filename is counted by the judge and never failed on.**

- **The generator reads frontmatter.** An appended row is `- [<id> — <title>](<file>) · \`<status>\``, spelling the *parsed* frontmatter values. A title whose value contains a colon is YAML-quoted in the file, and the row carries the value, not the quoting.
- **A subfolder row states a section, not a record.** `docs/README.md`'s block links folder READMEs, which carry no artifact identity of this shape; those rows stay plain.
- **An artifact missing any of id, title, or status degrades to the plain link.** A row is one shape or the other and never a third carrying an empty status badge. Index generation reports nothing and fails on nothing; naming a malformed artifact is the judge's job, and a generator that refused to emit a row would turn one bad file into a missing index entry.
- **The judge counts filler and does not fail on it.** `clue validate` reports rows whose label equals the target's filename stem as a population on its OK line, and `--index-rows` names them. It is a count, not an `Issue`.
- **Nothing existing is rewritten.** Regeneration keeps preserving every row whose target still exists.

[ADR-046](ADR-046-index-rows-say-what-the-artifact-is-about.md) extends the appended row's shape with a description seeded from the artifact's body, and counts the rows that state their record without saying what it is about. The rules above are unchanged: the row still opens by stating its record, an artifact that cannot supply one still degrades to the plain link, and nothing existing is rewritten — which is what makes the seed a first draft rather than an assertion.

**Counting rather than failing follows the shape this corpus already uses twice.** [ADR-035](ADR-035-bounded-provenance-and-reality-edges.md) reports costly unverified meaning as an actionable population instead of a build failure, and [ADR-017](ADR-017-conventions-are-constraints.md) applies the same reading to the constraint backlog — "the backlog is visible, not archival". The inferred-decision counter ran that experiment to completion: it stayed visible for months, drove a campaign milestone, and reached zero.

Failing would also invert where the fault lies. Every filler row in a Cliewen corpus was written by the generator, not by a person. A judge that rejected text the tool emitted in a previous release, in a file the adopter owns, would turn the tool's own defect into the adopter's red build on upgrade. Fixing the generator stops the source; the count shows the remainder.

**The accepted cost:** no command clears this count, because regeneration preserves rows whose targets still exist, so every repair is by hand. That is why the rows are listable rather than only countable. Under [C-004](../constraints/C-004-never-weaken-checks.md) the count is never softened to make the number look better.

**Carrier:** the emitting branch in `regenIndex` (`internal/scaffold/scaffold.go`), `corpus.IndexRowBacklog` and the `--index-rows` surface (`internal/corpus/index.go`, `cmd/clue/main.go`), [C-016](../constraints/C-016-index-rows-state-their-record.md), and CAP-005's and CAP-002's criteria and design.

### Rejected: check that a label matches its record's title

Requiring a row's label to equal the frontmatter title character for character is what this repository chose for its own decisions index, and it is house style rather than correctness. Several titles in this corpus run past a hundred characters, which is a defensible choice here and not one to impose on every adopter's README through the judge. A filename stem can be recognized mechanically and carries no opinion; "this label is not the wording I would have used" cannot.

The consequence is accepted and stated rather than hidden: a repository wanting every row to match its record verbatim holds that rule through the generator and through review, not through a machine.

### Rejected: normalize existing rows during regeneration

Rewriting a divergent row's opening while preserving its curated tail would make the count self-clearing, and it is rejected for this change. Regeneration documents the opposite as deliberate — entries whose targets still exist survive, so hand-written descriptions survive — and reversing that converts `clue scaffold` from a tool that appends into one that edits prose inside adopter-owned files ([ADR-019](ADR-019-init-regenerates-indexes.md)).

It is also ambiguous in practice. The boundary between a generated opening and a curated tail is defined by the status separator, so a row carrying no status offers nothing to split on and the tool must either guess or refuse. A tool that guesses inside someone's README eventually eats a sentence. Worth revisiting only with evidence that the counted backlog does not fall on its own.
