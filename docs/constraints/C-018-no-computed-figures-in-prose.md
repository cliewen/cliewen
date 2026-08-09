---
id: C-018
type: constraint
status: active
links: [P-013]
title: A durable record never states a figure a command computes
source: the shared durable-work fragment
enforcement: agent
---

# C-018 — A durable record never states a figure a command computes

A durable record never states a figure a command computes — an artifact count, a coverage percentage, a reported population size. Name the command instead. A number written into prose becomes a hand-maintained obligation that goes stale on the next change and that every later reviewer re-derives, and repairing one writes new prose carrying new numbers, so the finding regenerates instead of converging. Measurements that are the point of a record — an analysis's own results, a milestone's observed evidence — are stated with what produced them and when.

**Promotion trigger:** a lint that can tell a genuine measurement result apart from a restated computed figure, without false-positiving on every number a record legitimately states (a date, an ID, a count that is itself the finding). No such lint exists today; the register's own experience is the case in point — [AN-018](../analysis/AN-018-skill-statement-register.md)'s first draft stated its populations in prose, which sent the review loop hunting fresh arithmetic on three consecutive passes, and the repair was to move the figures into a script the register cites rather than to write a check.
