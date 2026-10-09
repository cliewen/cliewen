---
id: AN-038
type: analysis
status: active
links: [G-016, CAP-012, G-019]
title: Prepared trial of concrete Human evidence with local acceptance
---

# Human-evidence local-acceptance pilot

## Question and consumer

Do screenshots tied to a specific revision help the owner explain a rejection and acceptance, rather than merely add paperwork? The findings should help a person [make an informed acceptance decision](../goals/G-016-acceptance-is-an-informed-decision.md): explain whether a change is wanted and whether the observed behavior supports it. The ledger setup findings should also help contributors [register and finish a change workspace without editing the identity ledger by hand](../goals/G-019-workspace-identities-need-no-hand-edited-ledger.md). It does not adopt a new evidence format, change the validator, or establish general onboarding usability.

## Evidence boundary

The pilot is a local repository at `C:/Code/cliewen-human-evidence-pilot`, deliberately without a remote, PR or automated application tests. It uses the published [Cliewen 0.30.0 release](https://github.com/cliewen/cliewen/releases/tag/v0.30.0), Windows amd64 asset `clue-0.30.0-windows-amd64.exe`, with SHA-256 `e5ad2fb8537d41f0c80a1edda9fba337e04a8724bb01072f68e5c7469222397f` verified against the published checksums. The environment has Git, PowerShell, Python and an existing Edge installation. It is prepared, not a clean onboarding environment. The owner already knows the purpose of the trial; no unfamiliar-user claim is made. This local pilot has no hosted address: source observations below pin its full path and commit, not an invented forge link.

## Observed setup facts

- `clue init` materialized its existing convention and marked vision, architecture and design bootstraps. Replacing those bootstraps was required before the normal first green validation.
- Initial `clue id next` refused allocation because init had not created the identity ledger. A reviewed migration preview offered ledger backfill; apply created it, after which numeric allocation worked.
- The public CLI still cannot register the opaque workspace identities or retire a digested workspace. A disposable helper under this source checkout's Git metadata calls the existing `ledger.Load`, `MarkLive`, `Retire` and `Save` API against the pilot. No ledger text is hand-edited and no product command is added.
- The owner-authorized setup baseline is `20db106e2ed0aaccce4da6a6452058c6847c72c7`. It contains a draft Human criterion and no working page or claimed Human outcome.
- The proposal was committed before code. Its initial questions artifact used an invalid lifecycle value; validation rejected it. It was corrected before the recorded green diagnostic state. This is authoring friction, not evidence that the Human workflow failed.

## Pinned diagnostic and corrected code

The diagnostic control at `63d39db44191a07f78376ff250a1b535cab22d55` contains a deliberately misleading summary for fixed fictional failed-backup data. The published `clue validate` returned `OK (9 artifacts)` with an active Human declaration and no executable evidence export. This is the judge's stated structural boundary, not a contradicted claim that Cliewen proves visual comprehension. The control is not a clean reviewed or accepted candidate.

The corrected page code is committed at `6ba148fb65bfb8d790057468cf047938df11c276`. It renders the same failed outcome and the need for action. Its existence and green structural validation do not prove the owner understood it.

Committed page bytes are copied to revision-named directories in the pilot's Git metadata and served by Python's standard HTTP server bound only to `127.0.0.1:8765`. Preview content is independent of later branch switches. No real backup is performed.

## Current observation boundary

Actual browser capture is pending. The available computer-use surface lists no browser, and in-app browser creation reports unavailable. The owner has been asked whether to authorize a local headless-browser capture instead. No image, browser interaction or human observation is fabricated.

The owner's control verdict, corrected-page observation, local acceptance, missing-reference probe and post-acceptance history recovery remain pending. No conclusion about the value of artifact-backed Human evidence can yet be drawn. Findings will be updated with source revisions, hashes, exact checks and the owner's actual observations as the trial proceeds.
