---
id: ADR-056
type: decision
status: inferred
links: [CAP-007, PDR-034, ADR-007]
title: A context slice is bounded by default, and states what the bound held back
author: agent
accepted-by: []
---

# ADR-056 — A context slice is bounded by default, and states what the bound held back

## Context and problem statement

`clue context` exists so an agent can read what one identity depends on instead of reading a folder, and [PDR-034](PDR-034-the-corpus-is-read-narrowly-by-default.md) makes that the rule: reading starts at the narrowest point that answers the question. That decision names this command as the mechanism and describes it as emitting one artifact and its link slice.

The description was true of a young corpus and stops being true of a mature one. The slice followed every outgoing edge to exhaustion, and a corpus densifies as campaigns close: a goal accumulates inbound edges from everything that serves it, a plan points at every artifact in its campaign and is pointed back at, and a core constraint sits on most paths by construction. Once any such artifact is reachable, breadth-first traversal passes through it and out again into everything it touches. The closure from an ordinary artifact is then most of the repository, and `clue context --stats` will report as much.

Nothing announced the transition. The command behaved identically before and after; only the graph changed. So an adopter's corpus crosses this line silently, some campaigns in, and the first symptom is an agent that reads for a long time and then works badly — which looks like a model problem, not a corpus problem.

This leaves the obligation and its mechanism in contradiction: narrow reading is the rule, and the tool named to make it possible cannot.

## Decision outcome

**A slice follows a stated number of link hops, defaults to the root and what it links to directly, and names what the bound held back.**

*The bound is a default, not a limit.* `--depth` widens it to any number of hops, and `--depth=all` follows every edge to exhaustion — the behaviour that was previously the only behaviour. Nothing is unreachable, nothing is truncated mid-artifact, and the widest read costs exactly what it cost before. An artifact included in a slice is printed whole, as it always was.

*The frontier report is a condition of the bound, not a feature of it.* Beyond the printed depth, the command names the artifacts one hop out and counts the remainder. This is what separates this decision from the one PDR-034 rejected, and it is not a convenience: a slice that quietly stops at an unstated boundary makes a reader believe they have seen the dependencies, which is worse than an expensive read and worse than no read. A bounded slice is honest only when it says where it stopped and how to go further. Widening therefore stays a judgement the reader makes on evidence, which is precisely the shape PDR-034 asks for.

*Only the next hop is named; the rest is counted.* Naming the whole remainder would grow with the corpus and reproduce the cost the bound removes. The artifacts one hop out are the ones a reader can act on, because one widening reaches them.

*An unfollowed edge is reported only when it leaves an artifact the slice included.* An edge leaving an artifact the reader never sees describes a part of the graph this slice does not cover, and reporting it turns a bounded read into someone else's backlog. `clue validate` remains the judge of graph health and sees the whole graph regardless of how anyone reads it.

*The default is one hop.* Zero is the artifact alone, which a reader can obtain by opening the file. One hop is the smallest slice that answers the question the command exists for — what does this identity depend on — and every further hop answers a question the reader has not yet asked.

Bounding what a command prints does not change what the corpus means, what a thread connects, what a merge accepts, or what a green `clue validate` asserts.

## The criterion this retires

AC-053 states that the command emits each transitively linked artifact. The bound contradicts that meaning rather than rewording it, so the criterion is retired with a tombstone and a new one minted, per [ADR-007](ADR-007-ac-lifecycle.md). The retired criterion's assertions survive as the widened behaviour and are what `--depth=all` continues to guarantee.

## Rejected: leave the slice unbounded and rely on the obligation

The state this decision ends. PDR-034 tells an agent to read narrowly and points at a command that reads everything, so following the rule and using the tool are the same act with opposite results. An obligation whose only mechanism defeats it is not an obligation; it is an aspiration with a citation.

## Rejected: bound the slice silently

Cheaper to implement and dishonest in exactly the way PDR-034 warned about. A reader who cannot see that neighbours were withheld cannot decide to widen, so the bound stops being a default and becomes the cap that decision refused. This is the whole reason the frontier report is a condition rather than an enhancement.

## Rejected: name every artifact beyond the bound

Complete, and self-defeating: the list grows with the corpus, so the report inherits the cost the bound was added to remove. The next hop is bounded by the slice's own out-degree and is the only part a reader can act on without widening twice.

## Rejected: follow reverse links to decide what matters

Proposed periodically as a way to keep a bound smart — rank neighbours by how central they are and follow the important ones. It fails on the same ground the capability already refuses reverse traversal: centrality is computed from the whole graph, so the command would read everything to decide what not to print. A bound that must load the corpus to apply itself has no bound.

## Carrier

CAP-007's README and design state what a slice contains and where it stops; its criteria state the bound, the widening, and the frontier report. The CLI usage text states the flags. The routing hub, and the scaffolded hub an adopter receives, state how an agent enters the corpus — the instruction this decision changes for every reader who never opens this record. The generated skills carry no statement of what a slice contains; they defer corpus entry to the hub, which is why the hub is the carrier that moves.
