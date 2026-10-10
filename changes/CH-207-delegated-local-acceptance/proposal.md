---
id: CH-207
type: change
status: open
links: [G-016, CAP-012, G-027]
plan: none
title: Materialize internal links safely and execute recorded human approval
---

# Proposal

The owner requested a fix so integration can be executed by the agent after conversational approval, without manual Git steps. This change permits safe internal tracked links in the isolated Git-object snapshot and adds an explicitly recorded human-approval path to local acceptance. The human still decides whether an exact candidate is accepted; the agent may mechanically execute that recorded decision and an explicitly authorized push. No agent fabricates an approval or decides acceptance for itself.

## Authorized dependency

This branch starts from the reviewed but unintegrated readable-reference candidate `2416c0dff2bdd658268fe68b8177cac2bdcfd967`, rather than accepted main `5be6182506f1d2f58da3ea87ee359fe259d0578e`. The owner approved that candidate with “accepteret”, then explicitly requested “fix det, så du selv kan lave commits. Jeg vil ikke lave manuelle skridt” after the source symlink blocker was explained. The fix serves that requested integration and carries the readable-reference meaning as an authorized unaccepted dependency. The combined new candidate needs its own exact human approval after verification and review; the earlier approval does not authorize changed bytes.

## Challenge

A record written by an agent cannot cryptographically prove that a human authored the approval. The existing terminal prompt and Git identity also cannot authenticate human presence. Record the actual approver, source, statement and time, bind candidate, base and complete brief hash, and state that this is procedural trust. A missing, stale, malformed or mismatched record must fail without integrating. An unattended yes flag or silently replayed approval is not an alternative.

Symlinks must be resolved only against the committed revision. A live checkout target could escape or change during validation even when its apparent path looks safe. Read blobs directly, resolve an in-revision graph, copy regular target content into a private snapshot, and reject external or missing targets, cycles, metadata paths and submodules. The actual acceptance tree retains the original candidate objects and modes. Test hostile links as well as source skill mirrors before claiming support.

The repository's vision is draft and inferred. This plan-less change supports informed human acceptance while delegating execution; it does not claim automatic proof of observations, tests or approval authenticity.
