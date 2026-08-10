---
id: AN-012
type: analysis
status: active
provenance: verified
reversal-cost: low
links: [P-008, M-035, ADR-013, ADR-011, ADR-022, ADR-030, ADR-035, AN-010, CAP-001, CAP-004, C-011, P-013, PDR-026, ADR-038, ADR-039, C-001]
title: The adopter needs an upgrade path and a wall it need not fork, not a configuration file
---

# AN-012 — The adopter needs an upgrade path and a wall it need not fork, not a configuration file

## Risk

Cliewen may be about to build a configuration interface — a `cliewen.yaml`, a mirror option, or an extension point — for adopter needs that no adopting repository has demonstrated, paying the cost [ADR-013](../decisions/ADR-013-ships-generic-vs-repo-local.md) named (a second source of truth for agent routing) to solve a problem that may not exist.

## Evidence boundary

The spike inspected the public Robocode Tank Royale repository, the first product repository to adopt Cliewen, pinned at fetched `origin/main` commit `86a4bd58514bcdc4d36f1dd374900e6eae3b29f3` (2026-07-23) — the same revision [AN-010](AN-010-adopter-change-overhead.md) pinned, unchanged since. Its adoption merge is `b7fb320ccec1f4742ef923cb315e7dd84f7e824f` (PR [#218](https://github.com/robocode-dev/tank-royale/pull/218)). Cliewen was pinned at `c4549ef18b8c3ae13f221e348448cdeb96fe969e`, the accepted `main` this change branched from.

**Every reproduced result in this analysis is from a prepared environment, and none of it is onboarding-reproducibility evidence.** The prerequisites were: an existing local clone of the adopter, a Go 1.26.5 toolchain used to build `clue` from source, and a pre-existing released `clue 0.6.0` already on `PATH`. Reproductions ran on Microsoft Windows NT 10.0.26200 under Git Bash. Two binaries were used deliberately: a source build of `c4549ef`, which reports version `dev` and therefore **skips the skill-drift comparison by design** ([ADR-011](../decisions/ADR-011-version-stamping.md), AC-033), and a released `clue 0.6.0`, used only to observe drift behaviour a released binary would produce. No result here was produced in a clean disposable environment, so nothing in it establishes what a new adopter would experience on first contact.

Nothing in the adopter repository was modified. `clue validate` is read-only, so it was run in place; no reproduction required a disposable clone.

Counts below are of validator issues over one pinned corpus, not a sample: there is no population, no sampling, and no uncertainty interval to report. They are deterministic outputs of one binary over one revision pair, and they say nothing about how any other adopter's corpus would fare. Repository activity is evidence of activity; where maintainer intent matters it is marked as an open question rather than inferred, because the human who selected this adopter also maintains it.

## What was measured

### The runner and action assumptions cost this repository nothing

The vendored wall hardcodes `runs-on: ubuntu-latest` and `actions/checkout@v4`. Across the adopter's six workflows, ten jobs already declare `ubuntu-latest` and one uses an OS matrix; all twenty-two of its action references are floating major tags, and **no workflow SHA-pins any action**. The wall therefore matches the repository's own established CI conventions exactly. There are no self-hosted runners and no pinning policy for it to violate. Measured cost: zero edits, zero failures.

The install path `/usr/local/bin/clue` was likewise never touched. On GitHub-hosted runners the default user can write there, so the assumption never surfaced. This is evidence about GitHub-hosted runners only; it establishes nothing about a self-hosted or container runner without root.

### `.claude` placement is already solved, and the adopter uses the solution

`.claude/commands` and `.claude/skills` are **committed symbolic links** into `../.agents/commands` and `../.agents/skills`. Cliewen already anticipates exactly this: [AC-038](../capabilities/CAP-001-onboarding/criteria.md) makes `clue init` detect a symlinked `.claude/skills` — or a symlinked `.claude` ancestor — and skip emitting the mirror rather than writing through the link. Measured cost: zero edits, zero failures.

One residual hazard is the adopter's own, not Cliewen's: committed symlinks require `core.symlinks` support on the checkout, which is not the Windows default.

### Local skills beside managed skills is already solved, measured at zero

The adopter's `.agents/skills/` tree holds twelve skill directories: the five managed `clue-*` skills and **seven repository-local ones** (`deploy-sample-bots`, `dot-audit`, `dot-prime`, `dot-scout`, `release`, `structurizr`, `update-deps`). The local skills use the Claude Code `SKILL.md` spelling and native `name`/`description`/`argument-hint`/`allowed-tools` frontmatter; the managed ones use lowercase `skill.md` with `cliewen-skill: true` and a version stamp.

Head `clue validate` reports **no issue about any of the seven**. [ADR-022](../decisions/ADR-022-skill-ownership-marker.md)'s ownership marker holds in practice, and `checkSkillVersions` case-folds the manifest name by design (AC-037), so the two spellings coexist without a host-dependent verdict. Measured cost: zero edits, zero failures. The blessed-extension-point question P-007 deferred has no demonstrated need behind it.

### The adopter's one wall edit forked it

Change CH-004 (`e0fb8faaecd283a150af341368027b4a83e9c789`, 2026-07-22) rewrote the wall's install step: instead of verifying a `clue` binary vendored under `.github/tools/`, it downloads the release asset directly and verifies it against `SHA256SUMS`. `.github/tools/` no longer exists in the repository.

The significant part is not the edit but its permanence. The wall ships as a **copied file**, so an adopter's change to it is a fork, and upstream improvements never arrive. The adopter's copy is missing three things head's template now has:

- **change-scope detection** — without it, every pull request runs corpus validation, including ones that touch no corpus file at all;
- the **armed-check warning** step;
- the **acceptance-brief gate** introduced by P-007/M-024, which fails a Cliewen pull request whose body carries no completed brief.

**Maintainer classification: declining the vendored binary was a preference, not a mandatory adopter need.** The Tank Royale maintainer states that the only requirement was for the repository to use Cliewen; Cliewen imposed no requirement about its binary's delivery. The direct verified-release download is therefore Tank Royale's implementation choice. It does not change the finding: a configuration key does not express *how the binary arrives* better than an edit does, and the candidate remedy below addresses the edit without treating the preference as a generic configuration need.

The third divergence is the consequential one: the human merge gate that M-024 exists to fill with content is **not enforced on this adopter's pull requests**, and nothing in the repository or in Cliewen reports that absence. The wall's own runs are green and fast — the six most recent all succeeded, in 11 to 14 seconds, through 2026-07-23 — which is precisely why the gap is invisible.

### Version upgrade, not configuration, is the dominant measured cost

The adopter pins `CLUE_VERSION: 0.5.1` and its five managed skills are stamped `0.5.1`; Cliewen's head is four minor versions ahead. Running head `clue validate` against the pinned adopter corpus produces **69 blocking issues**: 68 artifacts whose `provenance: inferred` now requires the `reversal-cost` field [ADR-035](../decisions/ADR-035-bounded-provenance-and-reality-edges.md) introduced in P-007/M-028, and one analysis record whose `status: verified` is no longer in the allowed vocabulary. A **released** head binary would add five skill-drift issues on top; the source build reports `dev` and skips that comparison.

So a version bump this adopter has not yet attempted starts from 69 blocking corpus issues — 68 of them the same mechanical field addition, one a vocabulary fix — plus re-vendoring the five managed skills, which is a single act that clears all five drift issues at once. Against that, `clue init` is documented as deliberately non-destructive and explicitly "not an updater", and [`guide/operations.md`](../../guide/operations.md)'s manual upgrade procedure still instructs the reader to replace `.github/tools/clue-<version>-linux-amd64` and its `SHA256SUMS` — files this adopter deleted in CH-004. **The documented upgrade path does not match the repository it was written for.**

## Finding

**[ADR-013](../decisions/ADR-013-ships-generic-vs-repo-local.md)'s condition for a machine configuration file is still not met, and this adopter supplies no evidence for one.** ADR-013 rejected `cliewen.yaml` and left the door open only for when "`clue` itself needs repo-local settings". Two of the three assumption families P-007 deferred behind a pilot — `.claude` placement and local skills beside managed skills — turn out to be solved already, and measured at zero cost. The third, the CI wall, cost the adopter exactly one edit, and configurability would not have prevented it: what the adopter changed was *how the binary arrives*, a choice a configuration key does not express better than an edit does.

The real costs sit next to the ones M-035 anticipated: **a wall that must be forked to be adapted**, and **an upgrade with no path**. Neither is a configuration problem.

## AGENTS.md local-layer assessment

[ADR-013](../decisions/ADR-013-ships-generic-vs-repo-local.md) keeps `AGENTS.md` as the repository-local layer for agent routing and instructions. None of the measured needs belongs there. The binary delivery choice and the copied CI wall are executable workflow concerns: an `AGENTS.md` instruction cannot prevent a wall fork or select how GitHub Actions installs `clue`. The agent-directory symlink and local skills already coexist without an added routing convention, so the adopter has no demonstrated agent-local setting to add to `AGENTS.md` either. The boundary therefore remains appropriate and supplies no reason to add a machine configuration file.

## Rejected

- **A `cliewen.yaml` configuration file.** Nothing measured needs one. ADR-013's reasoning stands unchanged and is now backed by evidence rather than argument: the two placement questions are solved, and the wall's single real edit is not expressible as configuration.
- **Making the CI wall's runner label, action refs, and install path configurable.** This was the milestone's leading hypothesis. Zero measured edits and zero measured failures across all three. Building keys for them would be a field nobody reads.
- **A configurable skills-mirror location.** The adopter's symlink already achieves the relocation, and AC-038 already handles it. [ADR-011](../decisions/ADR-011-version-stamping.md)'s standing door — a relocated skills tree gets no drift check until the path is configurable — remains open but is **not** what this adopter needed: its skills are in the default location and the symlink points the other way.
- **A blessed local-extension point beside managed skills.** Seven local skills already coexist with the managed set at zero cost. The marker-based ownership boundary is the extension point.

## Candidates for a successor plan

Routed by reversal cost, as candidates only — this change decides none of them:

1. **ADR candidate — the wall stops being a file adopters fork.** Ship the corpus wall as a reusable workflow, or reduce the vendored copy to a thin caller whose version and install strategy are inputs, so adaptation stops severing the upgrade channel. High reversal cost: like [ADR-030](../decisions/ADR-030-verified-install-scripts.md)'s release-asset names, whatever shape this takes becomes a public contract with every already-adopted repository as a dependent.
2. **ADR candidate — corpus upgrades have a mechanical path.** A required-field addition like ADR-035's `reversal-cost` should not leave an adopter with dozens of manual edits and no tooling. Corpus format and lint rules, so an ADR under [C-011](../constraints/C-011-decision-records-typed.md).
3. **PDR candidate — a release that changes corpus obligations states its migration.** Process, not architecture: what a release must tell an adopter about newly required fields and vocabulary narrowing.
4. **Log-row candidate — `guide/operations.md`'s upgrade procedure is stale.** It names `.github/tools/` files this adopter no longer has. Cheap and local to reverse.
5. **For the adopter, not a Cliewen decision — the acceptance-brief gate is absent from its wall.** A repository-local fix in Tank Royale, recorded here because Cliewen's own evidence found it; candidate 1 would stop the divergence recurring but does not repair this instance.

**Named consumer:** a successor campaign **P-009**, to be proposed after P-008 closes. This analysis has no consumer inside P-008: M-035 forbids implementing an interface, and M-036 is a separate investigation. Candidates 1 and 2 are the intended spine of that campaign; candidate 4 may ride any change touching the guide.

## What this analysis does not establish

One adopter, one pinned revision, one prepared environment, no clean-environment result. Every "zero measured cost" above means *this repository did not have to change this thing*, never *no adopter will*. The specific blind spots are a runner without root, a non-GitHub forge, a case-sensitive filesystem (untested here — the host was case-insensitive), a corpus large enough for validation time to matter, and any adopter who did not have Cliewen's author as its maintainer.

## Re-derived at head (P-013/M-066, 2026-08-10)

[PDR-026](../decisions/PDR-026-campaigns-close-on-re-derived-evidence.md) requires this cost surface to be measured again rather than read off a milestone table. The adopter has moved to `origin/main` tip `2005b1d3a77f564a4714a7b0e64ed5545a1a396f` (2026-08-03); Cliewen was built from its own current `main` (a source build, so it again reports `dev` and again skips skill-drift by design). Every reproduction below is again read-only, again from one prepared environment rather than a clean one, and this time under native PowerShell 5.1.26100.8875 rather than the original's Git Bash — a boundary difference on top of the boundary the original section already states.

**Candidate 1 (the wall stops being a file adopters fork) is resolved and observably adopted.** [ADR-038](../decisions/ADR-038-upstream-validation-workflow.md) shipped this candidate under P-009/M-040, and this adopter's `.github/workflows/clue.yml` now confirms it in practice, not just in Cliewen's own tests: the file is a four-input thin caller to `cliewen/cliewen/.github/workflows/clue-validation.yml@16dc3a5...`, and its own comments state the checkout, change-scope detection, and acceptance-brief gate all now live upstream. `.github/tools/` still does not exist — CH-004's binary-delivery choice stands, exactly as this analysis said a configuration key would not have improved it. The fork this analysis measured is gone; the adopter did not need to re-fork to get there, it only needed to update one reference.

**Candidate 2 (a mechanical corpus upgrade path) is resolved as a capability, not yet as an event.** [ADR-039](../decisions/ADR-039-versioned-corpus-migrations.md) shipped `clue migrate` under P-009/M-041, and the adopter's own recorded activity shows the *specific* 69-issue gap this analysis measured closing by ordinary means before the tool would have been needed: CH-006 (`origin/main` PR #223, part of the accepted-unit range [AN-010](AN-010-adopter-change-overhead.md) now covers) added the required `reversal-cost` field to 67 existing artifacts and corrected the one bad `status` value, in a full-tier change, months before `clue migrate` existed. Running the current source build's `clue validate` against the adopter's current corpus finds **zero** of the original 69 issues — no `reversal-cost`, `provenance`, `status`, or `ADR-035` hits at all.

**What replaced them is a different, larger population this campaign did not anticipate.** The same command reports **536 blocking issues** against the current corpus: 474 are [C-001](../constraints/C-001-no-hard-wrapped-markdown.md) (hard-wrapped Markdown, a rule this analysis's original pin predates), 58 are C-009 (a capability README's required `goal` field, also newer than the pin), 2 are ADR-045's residual-declaration rule, and 2 are a newer constraint-source-resolution check pointing at a source file the adopter's corpus does not have. The adopter is pinned at `clue-version: 0.11.2` against Cliewen's current release `0.14.1` — three minor releases behind, narrower than the four this analysis originally measured, but not closed, and `clue migrate` (candidate 2's shipped remedy) has not been run: its existence lowers the cost of closing this gap, it does not by itself close it.

**What changed:** the forked wall (candidate 1) and the un-migratable field addition (candidate 2, for the specific instance measured) are both gone, one by an upstream reference update and one by ordinary full-tier work that predated the tool. **What did not change:** a version-pinned adopter still accumulates blocking issues every time Cliewen narrows or adds a corpus obligation, and closing that gap is still a decision the adopter has not made, exactly as [`guide/operations.md`](../../guide/operations.md)'s upgrade procedure already assumed it would be. **What the campaign declined:** auditing whether [ADR-039](../decisions/ADR-039-versioned-corpus-migrations.md)'s `MIG-001`..`MIG-003` registry mechanically covers the new C-001 and C-009 obligations, or would report them for semantic review instead — that is a question about the migration tool's own coverage, not about this adopter's measured cost, and re-deriving it belongs to whichever change next touches `clue migrate`'s registry rather than to a closing milestone that only re-measures.
