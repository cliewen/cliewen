# Imported changes

One `imported-change` record per source repository's in-flight pending change that brownfield extraction preserves ([ADR-050](../decisions/ADR-050-imported-change-records.md)). Each record pins the source revision and location it was read from, and carries the source change's intent, design rationale, dependency links, and a task-to-criterion proof-links table — the trace a deleted source `tasks.md` can no longer supply.

A record's lifecycle is `in-progress` → `complete`, not the corpus default: `clue validate` requires every proof-linked criterion on a `complete` record to exist, be undrafted, and not be retired. There is no `retired` value — a record is never deleted once written; it is the permanent evidence of what an extracted, now-gone source change once contained.

<!-- clue:index:start -->
- [IC-001 — Fixture: a proposal, design, dependency, and proof task remain inspectable](IC-001-fixture-inspectable-in-flight-work.md)
- [IC-002 — Fixture: an in-progress record may name an unproven proof link](IC-002-fixture-in-progress-dependent.md)
<!-- clue:index:end -->
