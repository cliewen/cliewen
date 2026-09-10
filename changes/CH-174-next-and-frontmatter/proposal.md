---
id: CH-174
type: change
status: open
links: [AC-186]
title: Make next work and front-matter semantics inspectable
---

# CH-174 — Make next work and front-matter semantics inspectable

This is a plan-less change: the source repository has no active campaign, and the work improves the core's read-only orientation and metadata contract without belonging to a product delivery plan.

Add a deterministic `clue next` command and teach generated agent guidance to use it when a user asks what to do next. Active plans supply actionable milestones; draft plans remain visible as proposed work but do not silently become tasks. The command supports listing alternatives so a human can choose another active plan without reassigning an existing change.

Make `reversal-cost` meaningful only for inferred non-decision artifacts. Inferred artifacts must classify their deferral as `low` or `high`; high remains an activation blocker, while low explicitly permits deferral. Verified artifacts must not retain the now-inert field, and migration removes it without changing surrounding content.
