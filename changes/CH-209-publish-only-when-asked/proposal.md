---
id: CH-209
type: change
status: open
links: [P-025, CAP-012]
title: The tracked change loop publishes only when the maintainer has asked it to
---

# CH-209 — The tracked change loop publishes only when the maintainer has asked it to

Serves milestone M-114 of plan P-025. It proceeds under draft, inferred VIS-001 without confirming it.

## What and why

The change-loop skill tells the agent to "commit and push" the proposal and "open the draft PR" as step 2, with no condition. A section further down says local acceptance replaces PR creation, but an agent reading step 1 and 2 first has already branched and published. Observed in the Tank Royale repository: two draft pull requests were opened for a maintainer who wanted none.

The change makes the loop state its publication intent before it branches or opens anything. The intent comes from two things: the acceptance policy and an explicit instruction. Under local acceptance the agent keeps the work on a local branch and neither pushes nor opens a pull request unless the maintainer asks for it in that session. Under pull-request acceptance the policy is the standing instruction, so the loop behaves as it does today. The agent says which of the two applies in one sentence at the start, so the maintainer can correct it before anything leaves the machine. The verification checklist's push-every-turn and draft-PR items are scoped to the pull-request mechanism, since under local acceptance they would contradict the new rule.

## Challenge the commitment

The riskiest assumption is that the policy alone settles intent under pull-request acceptance. A maintainer on that policy may still want one particular change kept local, and the sentence at the start only helps if the agent actually says it before acting. The credible alternative is to ask before every publication regardless of policy, which is safe but makes the pull-request path ask a question whose answer is almost always yes. The cheapest test is a scripted read of the changed skill text under both policies, checking that the intent sentence comes before the first branch or push instruction in the order an agent reads them. The change stops and returns for a decision if a reading under pull-request acceptance shows the loop now asking where it did not, or if any step before the intent sentence still publishes.

An implementation could meet every criterion and still fail the maintainer if the intent sentence were buried after the branching step, so the criterion fixes its position, before the branch is created.

## Carriers and proof

New criterion AC-240 in CAP-012 (local acceptance). Carriers: the change-loop skill source (steps 1 and 2), the shared review-boundary source (a short publication paragraph that the lifecycle skills all inherit), the verification checklist and review-loop wording in the verify skill source, the regenerated skills, CAP-012 design, and a user-facing changelog entry. The acceptance policy, the ledger protocol, identity allocation and the meaning of every other criterion are unchanged; a stated instruction to push or open a pull request under local acceptance is still honoured.
