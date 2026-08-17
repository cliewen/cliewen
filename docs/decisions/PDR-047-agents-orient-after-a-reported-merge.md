---
id: PDR-047
type: decision
status: inferred
links: [P-016, C-020]
title: Agents orient after a human reports a merge
author: agent
accepted-by: []
---

# PDR-047 — Agents orient after a human reports a merge

## Context

Treating a merge report as permission to start the next change can move the plan without the human knowing what has entered flight.

## Decision

After a human reports a Cliewen change merged, the agent describes the plan's next unfinished step in plain language and asks whether to start it. If no planned step remains, it says so and asks what comes next.
