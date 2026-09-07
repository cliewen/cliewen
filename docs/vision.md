---
id: VIS-001
type: vision
status: draft
links: []
title: Cliewen — durable intent that a machine can check before a human accepts it
provenance: inferred
reversal-cost: low
---

# VIS-001 — Cliewen

**What this is.** A methodology and one small command-line judge that keep a repository's durable intent — what is wanted, what the system can do, what proves it — in the repository itself, wired as a graph that `clue validate` checks and a human accepts by merging.

**Whom it serves.** Developers working with coding agents, and the team leads, architects, and reviewers who have to trust what those agents produced. Secondarily the auditor or newcomer who arrives later and needs to know why the system is the way it is.

**The problem.** Agents can now produce a great deal of plausible work quickly, and the reviewing human is the bottleneck. The two ordinary answers both fail: reading every diff does not scale, and trusting a confident summary is not evidence. Meanwhile the intent behind the work lives in issue trackers, chat, and people's heads, so a year later nobody can say which parts of the system were wanted and which merely happened.

**What it is for.** A reviewer should be able to see, mechanically, that a change's claims are connected to declared intent and declared proof — and then spend their judgement on whether the intent is right, which is the only part a machine cannot do. Intent stays current truth in the repository, versioned with the code that satisfies it.

**In scope.** The durable corpus and its identity graph; the deterministic judge and the CI wall in front of it; the agent skills that carry the method; greenfield adoption and brownfield extraction from an existing corpus; the human acceptance boundary at a merge commit.

**Out of scope.** Executing tests or judging whether they are good ones. Project management, estimation, and scheduling. Being an agent, orchestrating agents, or depending on any particular vendor's agent. Storing state outside the repository. Deciding, in any form, whether the intent recorded is the right intent.

**What constrains the direction.** Ceremony stays proportional to what changes — most work is simple work and pays for none of the loop. The judge checks form and never claims to have checked meaning; anything it cannot honestly verify is named as a residual rather than implied. Nothing the tool cannot verify is presented as verified. The corpus holds current truth, and Git holds the history. Agents prepare; humans accept.

**Succeeding looks like.** A reviewer trusting a merge because of what the corpus and the wall show, not because an agent sounded certain. An agent orienting in an unfamiliar repository from a bounded read rather than a full scan. An adopter's corpus still being maintained, and still true, a year after adoption. Someone choosing not to adopt Cliewen because its own stated limits told them it was the wrong fit.

**Still uncertain.** Whether the method holds at a scale larger than the repositories that have adopted it so far. Whether teams without a coding agent get enough from it to keep it current. Whether the shipped skills stay proportionate as the method's surface grows, or eventually need to be split by audience.
