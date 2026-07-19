---
id: CH-035
type: change
status: open
links: [P-003, M-011]
title: Cut release 0.4.0 — publish the public-readiness release pair
---

# CH-035 — Cut release 0.4.0

Serves P-003/M-011. PDR-009 already fixes v0.4.0 as the goes-public release and fixes its position after readiness and before the visibility flip; this change implements that accepted direction without making a new decision.

## What

- `CHANGELOG.md`: the `[Unreleased]` section becomes `[0.4.0] - 2026-07-19` and gains an install section with the anonymous public install story, free of the historical private-repository authentication caveats.
- The canonical generated-skill frontmatter version bumps from `0.3.0` to `0.4.0`; `go generate ./internal/skills` updates both distributed skill trees so the release binary and skills remain one drift-checked pair.
- The P-003/M-011 evidence records that CH-035 prepared the reviewed release commit while leaving the milestone open until the tag workflow has published and verified the release.
- After merge (human acts, not tasks): tag `v0.4.0` on the merge commit and push the tag. The release workflow builds, stamps, checksums, and publishes the cross-platform assets from the reviewed changelog section.

This is a full change because it touches methodology carriers. The version and release ordering are already accepted in PDR-009, and no acceptance-criterion meaning changes.
