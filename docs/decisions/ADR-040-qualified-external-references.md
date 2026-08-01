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

The notation follows URI conventions rather than inventing a private one, and it separates the two things Cliewen already keeps apart: an **address**, which says where something currently lives, and an **identity**, which says what something is.

*Anything with an address is written as a full URL.* A pull request, an issue, a wiki page, a published document, a site — all of them, in ordinary markdown link syntax, which is already a URI. This replaces the forge shorthand `owner/repo#N`, and the reason is that the shorthand names no host: it silently assumes one forge, so a resolver has to guess a base address, and a self-hosted or non-GitHub target cannot be written at all. A full URL states the host, resolves without assumption, and treats a wiki and a forge issue identically. AN-013 already contains the evidence: the only two references in this corpus that were not wrong are the two carrying a full URL.

The readable label survives — link text is free, so `[cliewen/cliewen#96](https://github.com/cliewen/cliewen/pull/96)` keeps the familiar shorthand where a human reads it while the target stays unambiguous.

*A bare `#N` is a validation failure.* It asserts a namespace it does not name, and it is the one shape already written wrongly here: offline nothing distinguishes a local number from a foreign one, and online a wrong-but-existing number resolves perfectly, so no later check can question it. Requiring an address surfaces the mistake at the moment it is made.

*An identity in another repository is written with a `clue:` scheme.* A corpus ID is an identity, not an address — Cliewen's own rule is that the ID is the identity and the path is only its current address — so answering with a URL would be the wrong kind of answer. The scheme says what the reference is, and the rest says which corpus it belongs to:

```
clue:robocode-dev/tank-royale/CAP-001
```

A foreign acceptance-evidence pointer is the same identity with the revision that was proven:

```
clue:robocode-dev/tank-royale@384d27d5/BR-001
```

**References that stay inside the repository are untouched.** A corpus artifact keeps its bare ID and a file its relative path, exactly as before. Neither gains a scheme, a slug, or an absolute URL. The rule reaches only what points outward.

Fenced code, inline code spans, link targets, and heading anchors such as `#4-resume-game` are not references and are never read as one. A colour literal carrying a hex letter is excluded by the same rule that excludes an anchor. An all-digit colour such as `#777777` is indistinguishable from an issue number by shape alone, and it is read as a reference: excluding runs of that length would silence the numbers adopters actually cite, and a rule that cannot see the common case is not worth the rare false positive it avoids. A colour literal in prose belongs in a code span, which is ordinary markdown practice.

This rule needs no network. It returns the same verdict offline, on a pinned revision, in a year.

*Resolving, separated and advisory.* A distinct command resolves the URLs the judge has already validated, and classifies each into exactly five outcomes. A `clue:` identity is not resolved: it names a corpus artifact in another repository, and following it would require knowing where that repository currently lives, which is the address coupling this notation exists to avoid.

- **reachable** — the target answered and is there.
- **restricted** — the target is not readable from here. Two answers prove it. A `401` or `403` says so outright. A `404` whose owner root still answers says it obliquely, and that case is not exotic: GitHub returns `404` rather than `403` for a private repository, deliberately, so that a refusal cannot confirm the repository exists. Taken at face value that reports a perfectly correct reference as rot, so the owner probe decides between the two. A private wiki, an internal tracker, and a private repository all land here permanently and correctly, and none is a finding. Authentication is a property of where the command runs, never of the corpus, so nothing is recorded to excuse it.
- **redirected** — the target moved. The new location is reported and offered as a rewrite.
- **gone** — `404` or `410`. The reference points at nothing. This is the error.
- **unreachable** — a timeout, a DNS failure, a rate limit. Reported as unknown.

Separating `restricted` from `unreachable` is what keeps the report worth reading: without it, every private target reports unknown forever and a genuine outage hides in that noise. The server states the difference — sometimes plainly, sometimes only by what its owner root admits — so no corpus annotation has to.

The owner probe trades in the safe direction on purpose. A genuinely deleted repository whose owner still exists is reported restricted rather than gone, so real rot can be missed; the reverse error would fail a corpus that is right, and `gone` is the one outcome that fails anything.

**A credential sharpens the answer where it is honoured, and is sent nowhere else.** `GITHUB_TOKEN` or `GH_TOKEN` turns a private target from an oblique `404` into a plain answer — but only against the API host, because a bearer token is inert on the web host, whose pages authenticate by session cookie. The resolver therefore sends a credential only when the address has an API equivalent to ask instead, and only to the service that issued it: attaching a forge token to an arbitrary wiki in the corpus would hand that host a credential it was never meant to see. Whoever holds a credential gets the sharper result; whoever does not still gets a report that never condemns a correct reference.

**Pinned history is reported and never rewritten.** A completed plan records what was observed at the time, and "the guide at this address returned HTTP 200" is a statement about that address on that day. Repointing it at wherever the content lives now would falsify the observation rather than repair it. Those references are reported with the rewrite the human may choose to make by hand; `--apply` leaves them alone.

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

CAP-002 owns the judge's form rule and the foreign-evidence pointer; CAP-004 owns the resolver command's shipped behaviour. The canonical and generated skills, scaffold templates, public and contributor guidance, and implementation explanations state the same two-layer boundary. Focused positive and negative tests cover the form rule, every one of the resolver's five outcomes, and the pointer's separate listing under --coverage. That a pointer never becomes classified evidence is structural rather than tested: evidence is harvested from test carriers in source, never from corpus prose, so no prose can reach that path.

The standing rejections this decision does not reopen: forge state may enforce but never mean, no green check is imported as acceptance, no base branch is read as accepted meaning, and no merge status is queried to derive accepted-ness.
