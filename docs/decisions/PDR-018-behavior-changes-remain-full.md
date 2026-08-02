---
id: PDR-018
type: decision
status: verified
links: [G-001, P-007, M-029, AN-008, AN-010, PDR-002, PDR-011]
title: Behavior changes remain full until adopter evidence supports a narrower loop
author: agent
accepted-by: Flemming N. Larsen (2026-08-02, conversation)
---

# PDR-018 — Behavior changes remain full until adopter evidence supports a narrower loop

## Context and problem statement

PDR-002 created the light tier because proposal-workspace ceremony costs more than changes that touch meaning without changing it. AN-008 observes that ordinary fixes still appear priced for Cliewen's authoring repositories and asks whether product behavior implemented under an existing criterion can safely use the light tier. AN-010 measures the first adopting product repository's accepted history: the workspace has material cost when used, but the pinned history contains no behavior change under an existing criterion. Should Cliewen widen the tier without an observed example of the category?

## Decision outcome

**Product behavior changes remain full changes, including when an existing acceptance criterion already states the intended behavior.**

An existing criterion removes the need to invent new acceptance meaning; it does not make the implementation semantically inert. Changing behavior changes executable evidence and can expose that the criterion, test boundary, or product reality no longer agrees. The full proposal and digest keep that alignment reviewable before the human merge boundary accepts it. Pure refactors that preserve behavior remain light, as PDR-002 already allows.

The public guide states this defence where it describes the modal product loop: the full path is deliberate because behavior, criterion, and evidence must remain one reviewed delta. The same text names the measured cost and the evidence limit rather than presenting ceremony as intrinsically valuable. A later analysis may reopen the boundary when an adopter history contains enough behavior-under-existing-criterion changes to compare their actual full-loop cost and failure modes.

Tier routing keeps one canonical detailed agent source in `internal/skills/source/shared/change-tiers.md.tmpl`. Generated lifecycle skills receive that source; repository and scaffold routing hubs stay concise enough to classify before corpus loading, and public human-facing carriers explain the rationale without restating every operational obligation.

**Carrier:** the canonical change-tier skill source and generated `clue-delta`/`clue-verify` skills (agent); the repository and scaffold `AGENTS.md` routing hubs (default); and the change-loop and methodology guide pages (human explanation). The tier boundary remains agent-enforced and does not change what `clue validate` judges.

### Rejected: behavior under an existing criterion is light

AN-010 contains no accepted change in this category. Widening it would trade away the proposal and digest on a hypothesis, while CH-002 demonstrates the adjacent failure mode: structurally green criteria can disagree with shipped behavior. Existing criterion text is an input to review, not proof that implementation cannot change meaning.

### Rejected: use a line or file threshold

Size is not semantic risk. A one-line condition can change product behavior and a large mechanical rename can preserve it; PDR-002 already rejects this proxy.

### Rejected: infer adopter cost from Cliewen's own history

Cliewen's changes disproportionately alter methodology and validator meaning, selecting for the full tier by construction. M-029 requires product-repository evidence precisely because this repository cannot supply a representative ratio.
