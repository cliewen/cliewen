---
id: UC-001
type: use-case
status: active
links: [G-001, CAP-001, CAP-003, CAP-002, CAP-008]
title: A team adopts Cliewen in a repository that already has a specification corpus
---

# UC-001 — Adopting Cliewen in a repository that already has a corpus

A brownfield adoption is the journey this method is most often judged by, and no single capability contains it. Onboarding materializes the convention, extraction rehearses and then transforms the existing corpus, the judge decides whether the result holds, and local verification is what the contributor actually runs. Each is separately correct and separately proven; what none of them states is that the team must be able to stop, or be stopped, before their existing corpus is deleted.

## Actors

The **adopting developer**, who owns the repository and its existing specifications. The **agent**, which performs the reading, the rehearsal, and the transformation. The **maintainer** who accepts the result at the merge boundary — often the same person as the developer, and a distinct role regardless.

## Trigger

The developer points an agent at a repository that already carries specifications in some other format and asks for the repository to be adopted.

## Preconditions

The repository is under Git with a branch the developer may push and a merge boundary a human controls. An existing corpus is present in a format the extraction mappings recognize, or close enough to one that a mapping can be written. `clue` is installed and its version is known.

## Main flow

1. The developer asks for adoption. The agent reads the repository's existing corpus, its documentation, and its tests before proposing anything, and states which source format it found.
2. The agent materializes the convention with `clue init` — the taxonomy, the routing hub, the skills, the CI caller — without overwriting anything the repository already has (CAP-001).
3. The agent rehearses the extraction and reports what it would do: which criteria map, which identities are preserved, which source material has no destination, and which in-flight work would be stranded. Nothing in the target has changed yet (CAP-003).
4. The developer reads the rehearsal and authorizes the mutation, or does not.
5. On authorization, the agent transforms the source corpus into the target corpus, preserving source identities, and produces the parity report that accounts for every criterion.
6. The agent runs the repository's checks and `clue validate`, and repairs what the judge rejects (CAP-002, CAP-008).
7. The agent opens a pull request whose brief states what the adoption binds. The maintainer reads the parity report and the brief, and the merge commit accepts the adoption.

## Alternative and failure flows

- **The rehearsal shows unmapped material.** The developer decides per item — map it, defer it against a named plan door, or drop it deliberately. Nothing is discarded silently, and the deferral is accountable rather than a footnote.
- **The developer does not authorize the mutation.** The repository keeps its `clue init` result and its original corpus. This is a supported end state, not a failure, and the rehearsal is re-runnable later.
- **The source cannot be migrated at all** — a legacy register whose rows a machine cannot classify. Migration blocks and says so rather than guessing; classification returns to a reviewed change with a human in it.
- **The judge rejects the transformed corpus.** The agent repairs it before the pull request is marked ready. A red wall never becomes a reason to weaken the wall.

## Outcome

The repository holds a Cliewen corpus whose identities trace back to the source, a parity report accounting for every source criterion, and a CI wall in front of the judge. The source corpus is removed only after that accounting exists — and the developer had a point at which they could stop, having been shown what would happen, before anything of theirs was changed.

## Open questions

Whether the same journey holds for a repository with no existing specification corpus at all, where extraction has nothing to read and the intent must instead be elicited. That is a different journey, and it does not yet have a use case.
