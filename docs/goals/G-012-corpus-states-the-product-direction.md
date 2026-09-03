---
id: G-012
type: goal
status: accepted
links: [G-001, VIS-001]
title: A corpus states the product direction its goals serve
---

# G-012 — A corpus states the product direction its goals serve

**Who wants it:** anyone orienting in a Cliewen corpus — a contributor arriving at a mature repository, and every agent that starts a change by reading a bounded slice (2026-09-03).

**Why:** the verifiable thread starts at a goal and stops there. A goal answers who wants something and why they want *that*; nothing answers why the product exists at all, whom it serves, or what has been deliberately excluded. The absence is cheap in a small repository and expensive in a grown one: the corpus can hold a dozen accepted goals that are individually justified and collectively unexplained, and a reader who wants the shape of the whole has to reconstruct it from capability folders.

An agent pays a second cost. Asked to weigh a judgement — is this in scope, is this the kind of thing this product does — it has capability-local truth and no statement of direction. It can ask the human every time, which makes routine work expensive, or it can infer direction from implementation structure, which reliably produces a confident answer to a question the code cannot answer: code shows what a system does, never why anyone wanted it.

There is a matching gap one level down. A capability is a unit of what the system can do; an acceptance criterion is a unit of proof. Neither carries an actor's path *across* capabilities, so a journey can be assembled from criteria that are each locally correct and still be one nobody would choose. Most behaviour does not need that artifact, which is why the answer here is an optional one rather than a required one.

**Success looks like:**

- One concise statement of direction exists per corpus, is reachable from the goals that serve it, and is short enough that an agent reads it while orienting rather than instead of working.
- An actor's end-to-end path across capabilities can be recorded when it materially helps, and is not created when it would only repeat a capability.
- Drafted or inferred direction is visibly distinct from direction a human has confirmed.
- A repository that has not yet stated its direction stays valid, and work unrelated to product meaning is never blocked by the gap.
- The judge checks that these artifacts are well-formed and connected, and claims nothing about whether they are right.
