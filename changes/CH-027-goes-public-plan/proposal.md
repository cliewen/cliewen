---
id: CH-027
type: change
status: open
links: [G-001, G-003]
title: P-003 "Cliewen goes public" begins
---

# CH-027 — P-003 "Cliewen goes public" begins

## What

1. **G-003 — "Cliewen is public"** enters the inbox as `proposed` and is promoted to `accepted` by this change's acceptance, which carries the campaign that serves it (the G-002/CH-006 precedent: PR approval is the promotion decision).
2. **PDR-009 — Going public** records the direction decisions the maintainer made on 2026-07-18: the repo flips public only after a readiness campaign; v0.4.0 is cut as the goes-public release, after the readiness changes and before the flip; the public guide is a handwritten VitePress site on GitHub Pages that links to the corpus rather than rendering it; and the corpus never cites private repositories' artifacts as normative references, while historical adopter evidence stays.
3. **P-003 — "Cliewen goes public"** (`active`): readiness cleanup, community front door, the guide site, the v0.4.0 release, and the visibility flip, as milestones M-008…M-012 (corpus-global M-numbering continues from P-002).
4. Housekeeping that belongs to the same bookkeeping: the root README Status section still describes P-002 as under way; it now reflects P-002 `completed` and names P-003.

## Why

P-002 completed on 2026-07-18 with all milestones done and explicitly left "public release / repo visibility decisions" out of scope. Making the repository public is the maintainer's stated goal, and it is exactly the kind of expensive-to-reverse direction decision that C-011 routes to a PDR and clue-plan routes to a new campaign: a public flip cannot be meaningfully undone (clones, caches, and module proxies persist), so the readiness work must precede it deliberately.

This change is the plan-creating change; it serves G-003 directly. The campaign's implementing changes (CH-028…) serve P-003's milestones.

## Decision boundary

This change records the maintainer's already-made direction decisions and creates the campaign structure. It changes no code, no acceptance criteria, and no methodology carriers. Anything beyond that scope becomes an open question and stops the change.
