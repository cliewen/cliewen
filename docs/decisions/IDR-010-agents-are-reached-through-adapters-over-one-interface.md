---
id: IDR-010
type: decision
status: verified
links: [IDR-009, P-024, G-025]
title: Agents are reached through adapters over one interface, not through a terminal-driving orchestrator
author: agent
accepted-by: Flemming N. Larsen (2026-10-04, conversation)
---

# IDR-010 — Agents are reached through adapters over one interface, not through a terminal-driving orchestrator

## Context

The scenario-trial harness ([IDR-009](IDR-009-scenario-trial-harness-is-go-repository-tooling.md)) has to run the same scenario on different agents and models and compare what they did, without the method depending on any one vendor's agent.

## Decision

Each agent is reached through an adapter in `tools/scenario-trials` that supplies four things: how the login is given from outside the repository, the command that runs the agent headless in the container, a probe for its version and model, and a parser that turns its own event stream into a vendor-neutral transcript. The runner, the checks, the conditions record and the summary read only that interface and that transcript, so a vendor's name appears only in its own adapter file and in the container image that installs the agent.

A terminal-driving orchestrator such as Orca is not used. It controls an agent through its terminal interface, so "idle" stands in for "finished" and the record is screen content, not the agent's own event stream. It needs its application running, which makes a trial depend on one machine's session, and the documentation pages read, two, do not show how to start an agent on a chosen model and prompt from a clean configuration. A second copy of the first agent's path with the second agent's parsing was also rejected: it keeps vendor names in the checks and makes every further agent a further copy.
