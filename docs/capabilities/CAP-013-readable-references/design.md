---
id: CAP-013-design
type: design
status: active
links: [CAP-013]
title: Keep readable reference presentation separate from corpus identity and judgment
---

# Design

One canonical instruction is included and routed in every independently installable lifecycle skill. Hubs and handoff templates point to that instruction and name the reader's immediate need. Changelog and guide prose describe effects rather than expecting identity lookups. Archived records remain history.

The human CLI uses a local naming index built from artifact titles, the existing criterion declaration parser's scenario names, and declared milestone names. Unknown and ambiguous declarations are not guessed. The index normalizes terminal labels to one safe line and never alters stored metadata or fetches a foreign target.

A producing corpus check registers the local subjects of a diagnostic separately from its raw path and message. The human renderer substitutes only those explicitly registered tokens using offsets in the original text, so a newly inserted name is not rescanned as another reference. It does not name built-in Cliewen policy references from adopter artifacts. Validation still returns the same raw issues, ordering and verdicts; existing proof coverage retains its meaning.

The same instruction also shapes what an agent says to a person. A status report, a summary of what was done and a handoff say what happened, what the reader must do and what comes next in everyday words, and name an identifier, a path or a command only when the reader asks for it or needs it to act ([AC-241](criteria.md)). It is a writing rule, not a template: the check is a reading of a real handoff, and one scenario run before and after is an observation, not a measurement.

Literal ID-allocation output, command arguments, tags, evidence and manifests remain canonical. The [cross-cutting design overview](../../design/README.md) locates this presentation layer relative to the state-only judge and local acceptance boundary.
