---
id: CH-199
type: change
status: open
links: [P-024, M-103, G-025, AN-028]
title: One pinned scenario runs from one command, five times on one agent
---

# CH-199 — One pinned scenario runs from one command, five times on one agent

## Proposal

Serve P-024's M-103. Build the smallest harness that runs the routing scenario from AN-028 repeatedly in a container, records the conditions each run needs to be read against, and reports the spread across five runs on Claude Code. Then run it, record the result in a findings document, and write M-103's evidence field with whether the spread allows comparison.

The harness is Go under `tools/scenario-trials`, chosen by the maintainer on 2026-10-05, and the choice is recorded as IDR-009. It builds a `clue` from a named commit, builds a Linux container image with pinned agent CLI versions and an empty home, runs the scenario N times with the subscription token supplied from a file outside the repository, and writes one result directory per invocation. Each run records the `clue` commit, skills version, method variant (always `baseline` here), agent and its version, model, image identity, and operating system. Deterministic checks read the agent's event stream: which route it recommended for each of the two requests, whether it edited anything before recommending, and what the working tree looked like at the end. The summary groups runs by those outcomes so the spread is visible without reading five transcripts.

A run starts only when the maintainer runs the command. Nothing is wired into CI, and nothing is placed under a path shipped to adopters.

## Challenge

**The assumption most likely to undermine this work:** that the deterministic checks on the event stream say anything about whether the agent behaved well. A run can recommend the right routes by luck of phrasing, or the right words while editing first. If the checks only count keywords the harness would report a stable result for an agent that is stable at saying the words and nothing else.

**A credible alternative:** skip the checks and read five transcripts by hand, as AN-028 did. It needs no parsing code and finds what a keyword scan misses, but it cannot be repeated cheaply and it produces no comparable record.

**The cheapest useful test:** the harness records the full event stream beside the checks, and this change reads all five transcripts once against the checks' verdicts. If a transcript contradicts its verdict, the check is revised before M-103's evidence is written. The findings document states how often the checks and the reading agreed.

**What would revise or stop the work:** checks that disagree with the reading on more than an isolated run, in which case the harness reports the stream and a transcript summary for a person and does not give a verdict at all; or a spread so wide that the five runs cannot be grouped, which is itself the result P-024 names as a reason to revise M-105 and M-106.

**How a build could meet every criterion and still fail the maintainer:** a harness that is green and quiet while the container silently carries host configuration, or a `clue` that is not the named commit. The conditions record is therefore read back from inside the container rather than assumed, and the findings document states what it showed.

## Scope boundary

This change builds and runs one scenario on one agent. It does not add a second agent or an adapter interface (M-104), a scenario set or its scrutiny of obligations (M-105), the OpenCode model (M-106), or a removal trial (M-107). The agent adapter is a single function now and becomes an interface in M-104. No new acceptance criterion is added: the harness is repository tooling, not a capability, and its proof is the recorded run. Where the Go tests run on it they follow the coverage floor like any other package.
