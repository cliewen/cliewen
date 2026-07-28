---
id: CAP-007-design
type: design
status: active
links: [CAP-007]
title: Focused corpus context design
---

# Focused corpus context design

`corpus.Context` builds one identity index per call and resolves the requested identity against it: artifact IDs, milestone IDs declared in plan bodies, and acceptance-criterion tags declared in criteria bodies, including retained tombstones. More than one declaring artifact is ambiguous and fails rather than choosing by scan order. Resolution reuses the same ID shapes the validator already recognizes, so `context` consumes the discovery job once parked as `clue locate`.

A milestone is declared by its row in a plan's milestone table, where the ID is the row's first cell. Plans also name other campaigns' milestones in prose — corpus-global numbering says which range a campaign continues from — and treating a prose mention as a declaration would make ordinary milestone IDs look ambiguous and take every artifact that links them down with the ID.

Traversal follows outgoing `links` only. The starting artifact is depth zero; each breadth-first layer is sorted by repository-relative path, and an artifact already emitted is never enqueued again. A link to a criterion or milestone resolves to that identity's owning artifact before traversal continues. Cycles therefore terminate and output is independent of filesystem walk order.

The command prints a boundary line containing ID and repository-relative path followed by the scanner's complete LF-normalized markdown content for each artifact. It scans but does not run whole-corpus validation: focused reading remains available when an unrelated validation rule is broken. A link the index cannot resolve is reported on stderr as an unfollowed edge and the rest of the slice still prints, because a reader repairing a broken edge is exactly the reader who needs its context; `clue validate` stays the judge of graph health. A scan-level parse failure still stops the command because an incomplete graph cannot produce a trustworthy slice, and an unresolvable requested ID is still an error with no artifact output.
