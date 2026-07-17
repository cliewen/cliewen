# Capabilities

CAP-xxx: what the system can do — one folder per capability, and the anchor of the requirements thread.

Each capability folder holds three files, updated together in the same change that alters the capability's behavior:

- `README.md` — what the capability is and why it exists (links its goal).
- `criteria.md` — acceptance criteria as Gherkin scenarios, each tagged `@AC-xxx`. Every active AC is verified by tests (a positive and a negative case); `clue validate` enforces the link. When a requirement changes what an AC *means*, the AC is retired (`@AC-xxx @retired` tombstone) and a new ID is minted — IDs are never redefined.
- `design.md` — how the capability works: the design is documented **per capability**, close to the criteria it realizes, not in one distant document.

A capability whose ACs cannot be tested yet stays `status: draft` with the gap stated — honesty over green.

<!-- clue:index:start -->
<!-- clue:index:end -->
