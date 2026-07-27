---
id: CH-070
type: change
status: open
links: []
title: Make multi-agent PR handoffs durable and SHA-bound
---

# CH-070 — Make multi-agent PR handoffs durable and SHA-bound

This plan-less change closes the recurring gap where one agent reviews or repairs another agent's pull request but leaves findings or fixes only in its private session. Independent changes remain parallel; collaboration on one pull request is coordinated through its hosted head commit and unresolved review conversations.

The change records the repeated failure, decides the shared handoff contract, adds a collaboration capability with executable evidence, strengthens the generated lifecycle skills and repository routing rules, and teaches adopters to require conversation resolution alongside the existing protected validation check.
