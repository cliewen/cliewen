---
id: G-010
type: goal
status: accepted
links: [G-005]
title: A generated index row stops disagreeing with the artifact it names
---

# G-010 — A generated index row stops disagreeing with the artifact it names

**Who wants it:** anyone choosing an entry point from a generated `/docs` index, human or agent (2026-08-11), found while CH-150 promoted a goal from `proposed` to `accepted`.

**Why:** `clue scaffold` builds an index row only for an artifact the block does not already reference. A row that already covers its target is kept exactly as it stands, which is deliberate for the description — the author writes what the artifact is about and regeneration must not overwrite it — but the status badge on the same row is generated content that no longer tracks its source. Change an artifact's status and the badge keeps the old value until someone edits the row by hand, and nothing reports the disagreement: the judge counts rows that say nothing about their artifact, not rows that say something untrue about it.

The cost is that the index is least reliable exactly where it is most used. A reader scanning for accepted work sees a `proposed` badge on an accepted goal, or the reverse, and the natural repair — running the regenerator — changes nothing, which teaches the reader that the row is authoritative when it is stale. The same mechanism keeps a description that has drifted from prose the artifact no longer contains.

**Success looks like:**

- A generated badge reflects the artifact's current status, or a disagreement between the two is reported rather than silent.
- Author-written descriptions still survive regeneration, so the fix does not trade one silent staleness for a silent overwrite.
