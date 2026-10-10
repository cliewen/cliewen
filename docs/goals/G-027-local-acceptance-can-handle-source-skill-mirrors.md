---
id: G-027
type: goal
status: accepted
links: [VIS-001]
title: Local acceptance can handle the source repository's managed skill mirrors
---

# Accept source work locally with managed mirrors

The source repository's local acceptance must handle its committed Claude skill mirrors without making the owner perform a manual Git merge. Those internal links must resolve against the same committed revision, and the snapshot must retain all validation content without following live filesystem targets or escaping the private snapshot.

The owner explicitly requested this fix after the readable-reference candidate's native preflight refused a managed mirror. The [local acceptance design](../capabilities/CAP-012-local-acceptance/design.md) explains the committed-graph boundary. Unsafe links remain refused; supporting internal mirrors does not authorize an agent to invent human approval.
