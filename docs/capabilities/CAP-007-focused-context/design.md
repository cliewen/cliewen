---
id: CAP-007-design
type: design
status: active
links: [CAP-007]
title: Focused corpus context design
---

# Focused corpus context design

`corpus.Context` builds one identity index per call and resolves the requested identity against it: artifact IDs, milestone IDs declared in plan bodies, and acceptance-criterion tags declared in criteria bodies, including retained tombstones. Acceptance-criterion identities use the validator's canonical segmented-prefix, decimal-number, and lowercase-suffix grammar, so `SNAP-SQS-001` and `ADP-045b` resolve without renaming. More than one declaring artifact is ambiguous and fails rather than choosing by scan order. Resolution reuses the same ID shapes the validator already recognizes, so `context` consumes the discovery job once parked as `clue locate`.

A milestone is declared by its row in a plan's milestone table, where the ID is the row's first cell. Plans also name other campaigns' milestones in prose — corpus-global numbering says which range a campaign continues from — and treating a prose mention as a declaration would make ordinary milestone IDs look ambiguous and take every artifact that links them down with the ID.

Traversal follows outgoing `links` only. The starting artifact is depth zero; each breadth-first layer is sorted by repository-relative path, and an artifact already emitted is never enqueued again. A link to a criterion or milestone resolves to that identity's owning artifact before traversal continues. Cycles therefore terminate and output is independent of filesystem walk order.

Traversal stops at a stated depth, defaulting to one hop. Beyond it, breadth-first search continues without emitting content, so the artifacts the bound held back are known rather than merely absent: those one hop past the bound are named with their ID and title, and the remainder is counted. Naming the whole remainder would grow with the corpus and reproduce the cost the bound removes, while the next hop is bounded by the slice's own out-degree and is the only part a reader reaches by widening once. The frontier is ordered like the slice itself, by hop count and then repository-relative path, so the whole output stays deterministic.

The library's depth option has a natural zero — the root alone — so unlike this package's other options its zero value is not "use the default"; the command supplies the default a human gets. This keeps an unset bound safe rather than unbounded.

The command prints a boundary line containing ID and repository-relative path followed by the scanner's complete LF-normalized markdown content for each artifact. An artifact inside the slice is never truncated: the bound governs how many artifacts are read, never how much of one. It scans but does not run whole-corpus validation: focused reading remains available when an unrelated validation rule is broken. A link the index cannot resolve is reported on stderr as an unfollowed edge and the rest of the slice still prints, because a reader repairing a broken edge is exactly the reader who needs its context; `clue validate` stays the judge of graph health. That report covers only edges leaving artifacts the slice included — an edge leaving an artifact the reader never sees describes a part of the graph this slice does not cover, and reporting it would turn a bounded read into someone else's backlog. A scan-level parse failure still stops the command because an incomplete graph cannot produce a trustworthy slice, and an unresolvable requested ID is still an error with no artifact output.
