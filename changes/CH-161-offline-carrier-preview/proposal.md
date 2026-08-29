---
id: CH-161
type: change
status: open
links: [CAP-004, PDR-043]
title: Preview repository carriers when release availability is unknown
---

# CH-161 — Preview repository carriers when release availability is unknown

CH-160 correctly separated published release availability from repository carrier currency, but its workflow still stopped before the migration preview when `clue latest` could not reach the release list. That leaves a local carrier drift that the installed binary can diagnose undiscovered, and a clean local preview must not be mistaken for proof that no newer release exists.

This plan-less follow-up makes the upgrade workflow run the no-apply migration preview regardless of release-list reachability. It reports release freshness as unknown when the availability check has no answer, reports local carriers as matching the installed binary only when the preview is clean, and asks for upgrade authorization only when a newer release or carrier drift is actually known.

The change updates the methodology decision and capability carriers, canonical and generated skill sources, focused positive and negative evidence, public upgrade guidance, and the adopter-facing `[Unreleased]` note.
