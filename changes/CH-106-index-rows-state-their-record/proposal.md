---
id: CH-106
type: change
status: open
links: [ADR-017, ADR-019, ADR-035, C-004, C-013]
title: Generated index rows state their record, and the judge counts the filler
---

# CH-106 — Generated index rows state their record, and the judge counts the filler

**This change is plan-less, deliberately.** It was scoped to P-010's M-049 when it was opened, and that declaration was withdrawn once the scope narrowed. M-049's exit criterion is the count of the thirteen constraints (C-001…C-013) the judge reports as agent-enforced reaching zero; this change converts none of them, and claiming the milestone would put a tick against work that leaves its counter exactly where it found it. What remains here is a defect repair in a generator plus a small counted population — worth doing on its own terms, not worth borrowing a milestone for. Reinstating M-049 is a one-line frontmatter change if the human prefers it.

## What is wrong

`regenIndex` (`internal/scaffold/scaffold.go`) appends every missing index entry as `- [<filename-stem>](<file>)` — no title, no status. That is machine-generated filler: it restates the URL as prose and tells a reader nothing the link did not already say.

It is also the sole origin of the twenty bare rows CH-105 spent a commit repairing in `docs/decisions/README.md`. Those rows were not written by anyone; they accumulated, one per new decision record, because the generator emits them and nothing objects. CH-105 repaired the symptom and left the source running, so the next decision record added to this repository reintroduces it — and, since `clue init` and `clue scaffold` ship these blocks, the same filler accumulates in every adopted repository.

Nothing checks index row content either: `checkIndexes` (`internal/corpus/rules.go`) verifies only that links resolve and that every artifact is covered by some row. Any label at all passes.

## What changes

1. **The generator emits the record instead of the filename.** `regenIndex` reads each missing target's frontmatter and appends `- [<id> — <title>](<file>) · \`<status>\``. This is the whole repair: it only ever touches entries that do not yet exist, so not one existing row in this or any adopted repository moves.
2. **The judge counts rows that state nothing.** A new `corpus.IndexRowBacklog` — a sibling of `ProvenanceBacklog`, not an addition to `checkIndexes` — reports rows whose label is exactly the target's filename stem, the generator's own output. `clue validate` names the count on its OK line and lists the rows behind a flag. It never fails on them.
3. **An ADR records the contract and its two rejections**, because this adds a rule to `clue validate` and fixes an output format shipped to every adopter — a core carrier under [C-013](../../docs/constraints/C-013-core-changes-need-decision.md).
4. **A constraint registers the rule**, born `enforcement: machine`.

## What this deliberately does not do

**It does not judge whether a curated label matches its record's title.** CH-105 brought all 58 rows of this repository's decision index to one rule — the label is the record's own `id — title`, verbatim, then its status — and that rule is worth keeping here. It is not worth shipping to adopters through the judge. A bare stem is *provably* uninformative and can be recognized mechanically without an opinion; "the label must equal the frontmatter title character for character" is house style, and some of those titles run past a hundred characters, which is a defensible choice for this corpus and nobody else's to make. The judge counts filler; it does not grade prose.

The honest consequence, stated rather than buried: this repository's 58-row invariant is then held by the generator and by review, not by a machine. That is weaker than CH-105's description implied, and saying so is better than reaching for a check that enforces taste.

**It does not rewrite existing rows.** `scaffold.go` states its own contract in as many words — regeneration "keeps existing entries whose targets still exist (hand-written descriptions survive)". Normalizing a divergent row's opening would reverse that deliberate decision, turn `clue scaffold` into a tool that edits prose inside adopter-owned files it previously only appended to, and force it to guess where an adopter's own note begins on any row that carries no status to split on. A tool that guesses inside someone's README eventually eats a sentence. If that behaviour is ever wanted it needs its own decision against [ADR-019](../../docs/decisions/ADR-019-init-regenerates-indexes.md), not a bullet in this change.

**It ships no corpus migration and breaks no adopter.** Nothing existing is judged; the count is a report, not a failure.

## Reversal cost

High, which is why it is an ADR rather than a log row. The emitted row format is consumed by every repository `clue init` touches, and reversing it after adopters have regenerated indexes against it means either leaving their rows in a shape nothing states or rewriting them again. The counter is cheap to withdraw; the format is not.
