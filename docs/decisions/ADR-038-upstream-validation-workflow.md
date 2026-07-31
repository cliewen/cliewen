---
id: ADR-038
type: decision
status: inferred
links: [P-009, M-040, AN-012, AN-014, CAP-001, CAP-006, CAP-004, ADR-030, PDR-021]
title: The CI wall is an upstream reusable workflow with a thin caller
author: agent
accepted-by: []
---

# ADR-038 — The CI wall is an upstream reusable workflow with a thin caller

## Context and problem statement

The scaffolded CI wall is copied into every adopter repository. AN-012 measured that an adopter's legitimate binary-delivery change forked the wall and allowed upstream validation and acceptance-brief fixes to disappear. AN-014 reports self-hosted/no-root runner, immutable action-reference, and repository-approved binary-install policies, but those target policies were not inspectable or reproducible. M-040 requires a public boundary based on a synthetic fixture rather than a general configuration file or an unverified enterprise report.

## Decision outcome

**Cliewen owns one reusable validation workflow, and an adopter owns only a thin caller.** The upstream workflow is `.github/workflows/clue-validation.yml`; the generated caller is `.github/workflows/clue.yml`. The caller invokes the upstream file at a full source commit when the emitting binary carries VCS metadata, or at the corresponding protected release tag when a module build has no VCS metadata. Release tags used as fallback references are immutable release records: the release process never reuses an existing tag, and adopters must protect release tags from force updates.

Both reference forms are immutable, so the fallback is a safe default rather than a degraded one, and emitting a commit the binary cannot vouch for is worse than emitting the tag: an unresolvable reference fails the adopter's first CI run with nothing in their repository to explain it. A revision is therefore emitted only when the emitting tree is identified as this project by its own contents, its tracked files still match that commit, and Go did not mark the build dirty. Publication cannot be checked offline and is not claimed; the release path, where the commit is always a published tag's commit, is the one that carries the guarantee.

The caller carries the event triggers, the stable `validate` job name, the upstream reference, and four `workflow_call` inputs: `clue-version` from the generated skill pair, a JSON runner-label list, `clue-source` (`vendored` or `release`), and a writable `clue-install-directory`. The runner and acquisition inputs are the only local policy choices demonstrated by the synthetic fixture; action references and validation steps are not caller inputs and are never copied into the caller.

The reusable workflow owns checkout, changed-surface detection, the plain-change boundary, the unarmed warning, checksum verification, release acquisition, executable staging, `clue validate --forbid-changes`, and the completed acceptance-brief gate. It executes the verified binary from the caller's checkout or the caller-selected writable directory and never assumes a root-only path such as `/usr/local/bin`. Its action references are full commit SHAs. A vendored source requires the matching Linux binary and `SHA256SUMS`; a release source downloads both from the pinned Cliewen release and verifies the checksum before execution.

The reusable job and caller job are both named `validate`. After the first run, branch protection selects the exact status check GitHub presents for that caller, including any reusable-workflow qualification such as `validate / validate`; the probe records that exact name rather than relying on a substring. The required check, conversation-resolution rule, and supported merge mode remain the protected human merge boundary from PDR-021.

Updating the one upstream workflow reference imports scope, warning, acceptance-brief, and validation fixes while leaving the caller's runner and acquisition choices intact. New adopters receive the thin caller. Existing adopters transition from the copied wall by replacing it with the caller, retaining their vendored assets or choosing the release source; no automatic updater or general `cliewen.yaml` is introduced.

## Consequences

- Upstream validation has one maintained implementation and adopters no longer fork its logic to adapt runner or binary-delivery policy.
- A self-hosted/no-root adopter can select runner labels and a writable install directory without granting the workflow a root-only path.
- The caller remains a public workflow contract: its reference, input names, release asset names, checksum file, and stable check are dependents that must move with the release contract.
- GitHub Actions availability, POSIX shell tools, public release access for `clue-source: release`, and the exact status-check name remain adoption prerequisites; the synthetic fixture does not claim support for other forges or arbitrary runners.
- `clue init` remains non-destructive, so upgrading an existing repository is reviewed repository work and remains outside this change's automatic behavior.
- The checksum and the binary share the release's origin, for both sources. Verification detects a truncated or corrupted download, not a compromised release; this is the same boundary [ADR-030](ADR-030-verified-install-scripts.md) draws for the install scripts, and it is not widened here.
- A `workflow_call`-only file is not executed by the repository that owns it, so the unit adopters depend on could ship having never run. An advisory in-repository caller exercises it — the unarmed path on every change to either file, and acquisition plus validation on manual dispatch against a published release. It is deliberately not named `validate`, so it cannot be confused with this repository's own required check.

## Rejected alternatives

- **Keep the copied all-in-one wall.** It preserves the measured fork and makes every upstream repair a manual comparison.
- **Add a general configuration file.** M-040 demonstrates only workflow inputs for runner and binary acquisition; a repository-wide second source of truth is not justified.
- **Let adopters pin or edit upstream action references.** Action pins are upstream-owned security and maintenance choices; caller-owned action inputs would recreate the fork at a lower level.
- **Install into `/usr/local/bin` unconditionally.** That fails on no-root runners and is unnecessary when the verified executable can run from a writable staging directory.
- **Use a branch reference for the reusable workflow.** A moving branch could change validation without a caller diff; a full commit or protected release tag makes the imported unit reviewable and stable.
