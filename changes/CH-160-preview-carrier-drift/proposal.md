---
id: CH-160
type: change
status: open
links: [CAP-004, PDR-043]
title: Preview repository carrier drift during upgrade discovery
---

# CH-160 — Preview repository carrier drift during upgrade discovery

`clue latest` answers whether a newer CLI release is published, but it cannot answer whether the adopted repository's managed carriers match the installed binary. The `clue-upgrade` skill currently treats the release check as sufficient to decide that no upgrade is needed when the binary is current, even though `clue migrate` can preview local carrier drift such as a 0.16.0 repository on a 0.19.0 binary.

This plan-less corrective change makes the upgrade skill run `clue migrate` in preview mode immediately after the availability check and before the human's now-or-later decision. The preview becomes the repository-state check; the existing human authorization, simple-work route, coordinated migration, verification, and human merge boundary remain unchanged.

The change updates the decision and capability carriers, canonical and generated skill sources, focused positive and negative evidence, public upgrade guidance, and the adopter-facing `[Unreleased]` note.
