---
id: CH-106
type: change
status: open
links: [P-010, M-049, ADR-019, ADR-035, ADR-039, C-013]
title: Index rows state their record's identity, and a machine keeps them that way
---

# CH-106 — Index rows state their record's identity, and a machine keeps them that way

**Serves P-010's M-049**, with one honest qualification stated up front: this change does **not** decrement the count M-049's exit criterion names. That count is the thirteen constraints (C-001…C-013) the judge reports as agent-enforced, and this change converts none of them. It registers a *new* form rule that is born machine-enforced, so the backlog does not grow while the rule becomes real. M-049 stays `todo` when this merges. The milestone is named because this is the same problem in the same arc — a rule about Cliewen's own form that no machine holds — and because P-010's second arc is explicitly about clearing Cliewen's own meaning and form. If the human would rather this stand plan-less, that is a one-line change to this frontmatter.

## What is wrong

CH-105 spent three commits bringing all 58 rows of `docs/decisions/README.md`'s index block to one rule: every row reads the record's own `id — title`, verbatim from its frontmatter, followed by that record's `status`. The argument for folding in the last exception — the `log.md` row — was stated plainly in that PR: a rule with one exception has to be told to every future reader, and the telling lived only in a PR description, whereas a rule without exceptions is one **a machine can check**.

No machine checks it, and the generator produces violations.

- `checkIndexes` (`internal/corpus/rules.go`) verifies only that index links resolve and that every artifact in the folder is covered by some row. It never reads a row's label or status. Any label at all passes.
- `regenIndex` (`internal/scaffold/scaffold.go`) preserves curated rows that still cover a live target — which is why CH-105's repair survives `clue scaffold` — but appends every *missing* entry as `- [<filename-stem>](<file>)`. No title, no status. That is precisely how the twenty bare rows CH-105 repaired came to exist in the first place.

So the next decision record added to this repository will be auto-appended in the bare stem form, and the invariant CH-105 established across all 58 rows begins decaying immediately. The repair was real; the rule behind it is held by nothing but memory of a merged PR description. This is the exact failure mode C-006 names for methodology contracts — a rule stated nowhere an agent reads is not a rule — applied to a rule that was never written down at all.

The same holds in every adopted repository: `clue init` and `clue scaffold` emit these blocks, so the bare-stem row is shipped behaviour, not a local blemish.

## What changes

1. **An ADR records the index-row contract.** A row states its artifact's `id — title` exactly as the frontmatter spells it, then the artifact's `status`, then any curated suffix. This is corpus architecture with a format that ships to every adopter, and it adds a rule to `clue validate` — a core carrier under [C-013](../../docs/constraints/C-013-core-changes-need-decision.md) — so it takes an explicit decision record and human acceptance, not a log row.
2. **The generator emits the contract.** `regenIndex` reads each missing target's frontmatter and appends `- [<id> — <title>](<file>) · \`<status>\`` instead of the bare stem. Curated rows keep surviving untouched, so no existing index moves.
3. **The judge counts divergence.** `checkIndexes` grows a reading of each row's label and status against the linked artifact's frontmatter, and reports the mismatches. **Whether it reports or fails is the open question below** — it decides the whole shape of the change and it is not mine to pick.
4. **A constraint registers the rule** in the register, born `enforcement: machine` once (2) and (3) land.

## What does not change

**No index content in this repository moves.** CH-105 already brought all 58 rows into line, and `clue scaffold` reports 0 blocks needing regeneration; this change makes the machine agree with what is already true here, which is why it is safe to do now and would have been noise to do before.

**No curated suffix is touched.** The supersession notes on ADR-011, ADR-013, PDR-001, PDR-002, PDR-003, PDR-007, PDR-010, PDR-016, and the log row's ADR-003/ADR-004 note are part of the contract, not exceptions to it: the rule governs the row's opening, not its tail.

**`clue validate` stays offline and deterministic.** The check reads files already in the corpus.

## Reversal cost

High, which is why it is an ADR. The generator's output format is consumed by every repository `clue init` touches, and a judge rule that fails is a rule that can turn a green adopter red on upgrade. Reversing the ADR after adopters have regenerated indexes against it means either leaving their rows in a format nothing states or rewriting them again.
