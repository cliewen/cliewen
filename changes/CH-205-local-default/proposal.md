---
id: CH-205
type: change
status: open
links: [G-016, CAP-001, CAP-012, PDR-067]
title: Local acceptance by default for new adopters
---

# CH-205 — Local acceptance by default

This change is plan-less and implements the conversational plan accepted on 2026-10-09. It proceeds under draft, inferred VIS-001 without confirming it.

## Accepted direction

Fresh adoption through `clue init` materializes local acceptance on `main`. `clue init --acceptance=pr` actively selects PR acceptance; explicit local selection is also supported. Existing Cliewen adoption without policy retains PR and receives an explicit PR policy through init or the reviewed migration. Existing valid policy is preserved; contradictory init flags and malformed policies fail before writes. This source repository explicitly selects PR and remains PR-only. The ledger, allocation modes, evidence requirements, exact-candidate merge and procedural human boundary remain unchanged.

## Challenge the commitment

The riskiest assumption is that absence of a policy means fresh adoption. Older corpora can lack role markers, so treating them as new would silently change accepted workflow. The credible alternative is keeping PR as the universal default or asking every init caller. The cheapest useful test is a fixture matrix for empty/new brownfield projects, canonical skills and old markerless corpora, explicit local/PR policies, malformed configuration, symlink boundaries, and source role. Any implicit switch of an existing workflow stops the work until classification is corrected. A green implementation could still fail its user if an init option appeared to change existing policy but merely skipped it, so explicit conflicting flags must fail visibly before any write.

## Carriers and proof

AC-223 covers default/explicit init selection and idempotency. AC-224 covers migration and legacy compatibility. AC-225 succeeds retired AC-221 while preserving base/candidate agreement and identity independence. AC-226 covers strict shared policy parsing and source-role behavior. PDR-068 records the default and compatibility decision.

Live carriers: shared policy reader and detection, init CLI/help and scaffold emission, migration registry/preview/apply, local acceptance loader/tests, canonical review boundary and lifecycle references, generated skills, hub and onboarding templates, core/design/capability truth, relevant decision refinements, public guide/navigation, and changelog including migration guidance. Existing user-authored hubs are never rewritten; notices point to needed policy-wording repair. No branch protection or release publication is changed.
