---
id: G-005
type: goal
status: proposed
links: [G-001, C-016]
title: A constraint's generated index badge shows its enforcement class
---

# G-005 — A constraint's generated index badge shows its enforcement class

**Who wants it:** repository contributors and adopters reading `docs/constraints/README.md` (2026-08-09), found while M-067 (CH-141) minted four new constraints and their freshly generated index rows all read `` `active` `` — the artifact's `status`, not its `enforcement`.

**Why:** the constraints README states its own contract plainly: "Its badge is the enforcement class rather than the artifact status." `regenIndex` (`internal/scaffold/scaffold.go`) appends every new row with `id.Status` regardless of artifact type, so every constraint added since C-017 needed its badge hand-corrected after generation, silently, with nothing reporting the mismatch. A reader trusts the badge to say who holds the rule; `` `active` `` says nothing.

**Success looks like:**

- A newly appended constraint index row shows its `enforcement:` value as the badge, the way every other constraint row already does.
- Existing curated rows are untouched — `clue init`'s regeneration-preserves-existing-rows contract is not part of what this goal reaches.
- `clue validate` reports a constraint row whose badge does not match its own `enforcement:` field, the same way it already counts a filler or undescribed index row, so a future miss is visible rather than silently hand-fixed again.
