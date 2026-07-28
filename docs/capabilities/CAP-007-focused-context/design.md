---
id: CAP-007-design
type: design
status: active
links: [CAP-007]
title: Focused corpus context design
---

# Focused corpus context design

`corpus.Context` resolves the requested identity against the scanner's artifact map, then against milestone IDs declared in plan bodies and acceptance-criterion tags declared in criteria bodies, including retained tombstones. More than one declaring artifact is ambiguous and fails rather than choosing by scan order. Resolution reuses the same ID shapes the validator already recognizes, so `context` consumes the discovery job once parked as `clue locate`.

Traversal follows outgoing `links` only. The starting artifact is depth zero; each breadth-first layer is sorted by repository-relative path, and an artifact already emitted is never enqueued again. A link to a criterion or milestone resolves to that identity's owning artifact before traversal continues. Cycles therefore terminate and output is independent of filesystem walk order.

The command prints a boundary line containing ID and repository-relative path followed by the scanner's complete LF-normalized markdown content for each artifact. It scans but does not run whole-corpus validation: focused reading remains available when an unrelated validation rule is broken. A scan-level parse failure still stops the command because an incomplete graph cannot produce a trustworthy slice.
