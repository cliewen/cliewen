---
id: C-022
type: constraint
status: active
links: [ADR-057, P-015]
title: Every durable artifact serves one primary consumer
source: ADR-057 and P-015/M-071
enforcement: partial
---

# C-022 — Every durable artifact serves one primary consumer

Each durable artifact serves one primary consumer. A file that contains several rendered documents is split by consumer and connected to its parent by durable links when that is the truthful structure; a file one reader must consume as a whole remains together with a stated reason.

**Checked by:** `clue validate` reports, without failing, every non-completed artifact under `docs/` whose body carries more than one rendered H1 outside fenced examples, and every durable identity whose default one-hop context slice prints more than the budget [ADR-057](../decisions/ADR-057-read-cost-measurements.md) states. `clue validate --read-cost` names both populations.

**Residual:** only a reader can decide whether two rendered documents truly have different primary consumers, whether splitting preserves a historical finding verbatim, or whether an exception is coherent. The cost of that judgment is a corpus whose reported backlog must be reviewed rather than silently treated as green.
