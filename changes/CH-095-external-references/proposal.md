---
id: CH-095
type: change
status: open
links: [P-009, M-044, AN-013, ADR-005, C-013, CAP-002]
title: External references name their target, and resolving them stays outside the judge
---

# CH-095 — External references name their target, and resolving them stays outside the judge

## What

A Cliewen corpus points outward: at other repositories, at their pull requests and issues, at wikis and published pages. Today it cannot say so precisely, and it cannot say whether what it points at is still there.

This full change closes both halves, in two deliberately separated layers.

**The judge keeps its offline promise and gains a form rule.** The notation follows URI conventions and splits address from identity. Anything with an address — a pull request, an issue, a wiki page, a published document — is written as a full URL in ordinary markdown link syntax, so the host is stated and nothing assumes a particular forge. An identity in another repository's corpus is written with a `clue:` scheme, because a corpus ID is an identity and a URL would be the wrong kind of answer. A bare `#N` becomes a validation failure: it silently means "this repository" and is wrong the moment it does not. The existing unqualified references in this corpus are repaired in the same change. This rule needs no network, returns the same verdict in a year, and works offline.

**Resolving a reference lives in a separate command.** Fetching a target is what tells you it moved, vanished, or still stands, and that cannot happen inside `clue validate` without making a verdict depend on another system's current state. The new command resolves the qualified references the judge has already found and classifies each one: reachable, redirected, gone, or unreachable. A redirect is reported with its new target and offered as a rewrite the human accepts; the command previews by default and writes only when told to, following `clue migrate`. An unreachable target is reported as unknown and never as invalid, because a network outage must not be able to condemn a corpus.

**Foreign acceptance evidence gets its named form on top of that notation.** A criterion whose real proof is a run in another repository points at it as `clue:owner/repo@revision/ID` — the repository, the revision that was proven, and the identifier. The judge treats it as named but locally unproven — never as coverage, never as an imported verdict. `Test-type: Human` stays the proof of record and this is the pointer beside it. The pinned revision is what keeps it honest: it says which state was proven, not what some branch happens to show today.

## Why

AN-013 measured the damage. This corpus carries 50 `#N` references across 25 files, 48 of them unqualified, spanning at least four repository namespaces. One is already confidently wrong: `#5` in AN-003 means another project's issue, but rendered here it is a real merged pull request about something else entirely. Others are dead now and will become wrong live links once this repository's own numbering passes them. The judge sees none of it.

The same analysis found that acceptance evidence already spans repositories in both directions and cannot be written down at all — a foreign identifier in `links:` is a hard failure, leaving `Test-type: Human` and `@draft` as the only honest expressions.

The two are one problem in the right order. A pointer to foreign evidence must name a repository, a revision, and an identifier; while a bare reference cannot even say which repository it means, that pointer has nothing to stand on. AN-013 says so directly: a bare ID becomes ambiguous the moment it crosses a repository, which is precisely what the evidence candidate would have it do.

## Decision boundary

This change extends P-009's M-044 beyond the foreign-evidence form it currently states, so the milestone's exit criterion is revised in the same change, backed by the decision this change records.

The standing rejections hold and are not reopened. Forge state may enforce but never mean: no green check is imported as acceptance, no base branch is read as accepted meaning, no merge status is queried to derive accepted-ness, and no forge registry becomes a second source of truth. Network resolution does not enter `clue validate` — the separated command is the whole point, not a loophole. The judge stays deterministic, offline, and reproducible; anything the resolver reports is advisory until a human accepts it.
