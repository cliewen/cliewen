---
id: ADR-040
type: decision
status: inferred
links: [P-009, AN-013, ADR-005, ADR-009, ADR-039, C-013, CAP-002, CAP-004]
title: External references name their target, and resolving them stays outside the judge
author: agent
accepted-by: []
---

# ADR-040 — External references name their target, and resolving them stays outside the judge

## Context and problem statement

A Cliewen corpus points outward. It cites another repository's pull requests and issues, links published pages and wikis, and sometimes its real acceptance evidence is a run that happened somewhere else. None of that can be stated precisely today, and none of it can be checked.

A bare `#N` silently means "this repository". AN-013 measured 50 such references across 25 files, 48 unqualified, spanning at least four repository namespaces — and one that is not broken but confidently wrong: written to mean another project's issue, it resolves in this repository to a real merged pull request about something else. Corpus IDs collide across repositories for the same reason, so a bare ID stops meaning one thing the moment it crosses a boundary.

The obvious repair is to have the judge look references up. That cannot be done: a verdict that depends on another system's present state is no longer deterministic, no longer reproducible against a pinned revision, and no longer available offline. It would also make an outage in someone else's infrastructure able to condemn a corpus that has not changed.

So the question is not whether to check external references. It is where each half of the check belongs.

## Decision outcome

**An external reference names its target, and the judge enforces only that. Resolving the target is a separate command that never gates a merge.**

*Naming, enforced offline.* A forge reference is written `owner/repo#N`. Anything else external — a wiki page, a published document, a site — is written as a full URL. A corpus ID belonging to another repository is qualified by that repository. A bare `#N` in corpus prose is a validation failure, because it asserts a namespace it does not name.

**References that stay inside the repository are untouched.** A corpus artifact is cited by its bare ID and a file by a relative path, exactly as before: identity is the ID and the path is only its current address, and neither gains a slug, a prefix, or an absolute URL. The rule reaches only what points outward. One exception is deliberate: a forge number names its repository even when it means this one, because `#N` is the single shape already written wrongly in this corpus — offline nothing distinguishes a local number from a foreign one, and online a wrong-but-existing number resolves perfectly. Requiring the author to state the repository is what surfaces the mistake at the moment it is made.

The qualification is itself the notation: a reference that names its repository or carries a scheme is external by construction, and nothing further marks it. Fenced code, inline code spans, link targets, heading anchors such as `#4-resume-game`, and colour literals such as `#777777` are not references and are never read as one.

This rule needs no network. It returns the same verdict offline, on a pinned revision, in a year.

*Resolving, separated and advisory.* A distinct command resolves the qualified references the judge has already validated, and classifies each into exactly five outcomes:

- **reachable** — the target answered and is there.
- **restricted** — the target answered `401` or `403`. It exists; this runner may not read it. A private wiki or an internal tracker lands here permanently and correctly, and it is not a finding. Authentication is a property of where the command runs, never of the corpus, so nothing is recorded to excuse it.
- **redirected** — the target moved. The new location is reported and offered as a rewrite.
- **gone** — `404` or `410`. The reference points at nothing. This is the error.
- **unreachable** — a timeout, a DNS failure, a rate limit. Reported as unknown.

Separating `restricted` from `unreachable` is what keeps the report worth reading: without it, every private target reports unknown forever and a genuine outage hides in that noise. The server already states the difference, so no corpus annotation has to.

Like `clue migrate`, the command previews by default and writes only when explicitly told to, so a moved target becomes a reviewed corpus edit rather than a silent one. **Neither `restricted` nor `unreachable` is a failure.** The corpus did not change and must not be condemned by weather elsewhere, or by where the command happened to run. Only **gone** is an error, and it is an error the command reports — never a verdict `clue validate` reaches.

The resolver may inform and may open findings where humans read them. It may not be a required status check. Making it one would put another organisation's uptime between a change and its merge, which is exactly the coupling this separation exists to prevent.

*Foreign acceptance evidence.* A criterion whose real proof is a run in another repository names the repository, the pinned revision, and the identifier. The judge treats it as **named but locally unproven**: never coverage, never an imported verdict. `Test-type: Human` remains the proof of record, and this pointer sits beside it. The pinned revision is what keeps the claim honest — it states which state was proven, not what a branch shows today.

## Adoption

This narrows a corpus obligation, so the release carrying it ships a migration under ADR-039 — and that migration reports, it does not repair. A bare `#N` cannot be qualified mechanically: nothing in the file says which repository was meant, and defaulting to the adopter's own slug would cement precisely the mistake this rule exists to catch, converting a confidently wrong reference into a confidently wrong *qualified* reference that no later check can question. The migration therefore lists every bare reference with its file and line and stops, in the same shape ADR-039 already uses for semantic cases; the adopter resolves them in a reviewed change.

The resolver may assist that work by reporting what each bare number resolves to in the repository being migrated. That is information for the human deciding, not a decision the tool takes.

An adopter that has not upgraded is unaffected. The obligation arrives with the release, on the adopter's schedule, which is what the versioned-migration contract is for.

## Rejected: resolve references inside `clue validate`

The convenient design, and the one AN-013 already ruled out. It puts network access in the deterministic judge, makes a verdict depend on another repository's current state, and stops validation working offline. The two-layer split is not a compromise around that rejection; it is what the rejection implies once the underlying need is taken seriously.

## Rejected: an explicit marker syntax for externality

A dedicated sigil or attribute declaring "this reference is external" would be a second thing to keep in step with the reference itself, and a corpus could carry a marked reference that is local and an unmarked one that is not. Qualification already carries the fact unambiguously: a reference naming a repository or a scheme is external, and one that names neither is local. One rule, nothing to desynchronise.

## Rejected: a registry of external targets

Recording forge identity as corpus truth was rejected in AN-013 and stays rejected. The defect is that references are unqualified, not that the corpus lacks an inventory. A registry would be a second source of truth that ages independently of the references it describes.

## Rejected: treating a redirect as authoritative and rewriting silently

A redirect is forge state. It may inform a proposed edit; it may not change corpus meaning on its own. Renames and transfers leave redirects that a mirror or an offline reader never sees, so the durable repair is the rewritten reference a human accepted, not a hop that happens to work today.

## Carrier

CAP-002 owns the judge's form rule and the foreign-evidence pointer; CAP-004 owns the resolver command's shipped behaviour. The canonical and generated skills, scaffold templates, public and contributor guidance, and implementation explanations state the same two-layer boundary. Focused positive and negative tests cover the form rule, each of the resolver's four outcomes, and the pointer's exclusion from coverage.

The standing rejections this decision does not reopen: forge state may enforce but never mean, no green check is imported as acceptance, no base branch is read as accepted meaning, and no merge status is queried to derive accepted-ness.
