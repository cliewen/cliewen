# Use cases

One file per use case: **one actor's end-to-end path across capabilities**, `UC-xxx-<slug>.md`, `type: use-case`.

**This folder is optional and may stay empty.** A use case is worth writing when a journey crosses several capabilities, when ordering or failure recovery carries meaning, when the acceptance criteria are each locally correct and still do not add up to the outcome, when several actors collaborate through the system, or when a brownfield system holds behaviour no single capability explains. It is not worth writing when one capability and its criteria already describe the behaviour, when there is no actor interaction, when it would restate a goal, or when the subject is internal implementation behaviour that design, architecture, a constraint, or an IDR owns ([PDR-054](../decisions/PDR-054-use-cases-are-optional-and-no-requirement-artifact.md)).

A use case never restates an acceptance criterion, and never becomes a design. It says what happens and in what order; the capability says what the system can do, and the criterion says what proves it.

**Links point down.** A use case names the goal it serves and every capability it crosses. Capabilities do not name their use cases: the edge is written once, and `clue context <id>` names the use cases that reach an artifact without following them ([ADR-066](../decisions/ADR-066-intent-links-point-down-and-context-names-what-points-back.md)).

Each use case carries four sections — `## Actors`, `## Trigger`, `## Main flow`, `## Outcome` — and adds preconditions, alternative and failure flows, and open questions when they carry meaning.

<!-- clue:index:start -->
- [UC-001 — A team adopts Cliewen in a repository that already has a specification corpus](UC-001-adopt-cliewen-in-an-existing-repository.md) · `active` — A brownfield adoption is the journey this method is most often judged by, and no single capability contains it.
<!-- clue:index:end -->
