---
id: CH-123
type: change
status: open
links: [P-012, M-058, CAP-003, ADR-049]
title: Deferred migration criteria carry inspectable accountability
---

# CH-123 — Deferred migration criteria carry inspectable accountability

## Proposal

Serve P-012's M-058 by making every migration-parity disposition for `@draft`, `Human`, or retired criteria carry structured source-location and plan-door accountability, and by reporting the deferred population as derived parity output.

The existing parity manifest accepts an opaque `justification` string. It can therefore state the same prose for every disposition without proving that the original source was found or that a real plan owns the outstanding work. This change will retain a readable explanation but make the accountable location and plan door independently checkable against the reconciled corpus.

The change will extend the parity contract and its CLI report, add AC-125 and focused positive and negative Unit evidence, update the extraction guidance that creates manifests, and record the durable manifest decision. Its digest will mark M-058 done only after the exit criterion is completely evidenced.

## Scope boundary

This change concerns migration parity only. It does not make ordinary greenfield criteria carry migration provenance, and it does not decide M-059 through M-061.
