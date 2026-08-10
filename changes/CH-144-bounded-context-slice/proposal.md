---
id: CH-144
type: change
status: open
links: [P-015, M-070, CAP-007, PDR-034]
title: A context slice is bounded, and states what the bound held back
---

# CH-144 — A context slice is bounded, and states what the bound held back

## Proposal

Serve P-015's M-070 by giving `clue context` a bounded default, a way to widen it, and a frontier report naming what the bound held back.

`clue context` followed every outgoing edge to exhaustion and printed each artifact in full. PDR-034 obliges an agent to read the corpus narrowly and names this command as the mechanism, describing it as emitting one artifact and its link slice; in a corpus dense enough for a goal or a plan to sit on most paths, the closure from an ordinary artifact is most of the repository, so the obligation could not be honoured with the tool it names.

The change will make the slice follow a stated number of link hops, default to the root and what it links to directly, and offer `--depth` to widen — including `--depth=all` for the previous behaviour. Beyond the printed depth it will name the artifacts one hop out and count the rest, so a reader always sees that an edge exists and can ask for it. `--stats` will report what the slice cost. An unfollowed edge will be reported only when it leaves an artifact the slice included, because an edge leaving something the reader never sees describes a part of the graph this slice does not cover.

AC-053 states that the command emits each transitively linked artifact. The bound contradicts that meaning, so AC-053 will be retired with a tombstone and a new criterion minted, per ADR-007. A decision record will engage PDR-034's rejection of a read cap in its own terms and make the frontier report a condition of the bound rather than a feature of it.

Because P-014 closed without designating a successor, this change also opens P-015 from the campaign it serves.

## Scope boundary

This change bounds what the command prints and records why. It does not add the read-cost measurements P-015's M-071 defines, does not split any artifact, does not trim the generated skills, and does not repair index rows. It changes no rule about what an agent may read: nothing is truncated, every artifact remains reachable, and the widest read costs exactly what it cost before.
