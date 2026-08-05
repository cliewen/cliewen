---
id: AN-017
type: analysis
status: active
provenance: inferred
reversal-cost: low
reality: contradicted
links: [CAP-003, ADR-036, PDR-020]
title: Heimdall extraction silently discarded multi-criterion JVM evidence
---

# AN-017 — Heimdall extraction silently discarded multi-criterion JVM evidence

## Finding

At [Heimdall extraction commit `b443613c`](https://github.com/elsevier-research/iip-heimdall/commit/b443613c), `UploadDeletionServiceTest.deleteUpload_openSearchCleanupFails_throwsAndPreservesDbRow` retained `@Tag("UA-071")` while the extraction diff removed `@Tag("UA-069")` and `@Tag("UA-070")`. The extraction report says 31 Java test files were normalized by removing duplicate AC tags and keeping the last one.

The source method's assertions throw an OpenSearch-cleanup exception and verify that the database row is preserved. The tag order does not establish that UA-071 is primary or that UA-069 and UA-070 are not proven. ADR-036 correctly treats several identities on one executable as ambiguous and gives no classified credit; the extraction contradicted that boundary by resolving ambiguity through silent tag loss.

## Evidence boundary

This is a read-only comparison of `b443613c` with its first parent in the local Heimdall checkout on 2026-08-05. It proves the committed tag deletion and the extraction report's recorded strategy. It does not establish the correct source-level resolution for every normalized method; that requires the per-method human review PDR-024 introduces.

## Rejected response

Marking every discarded identity `@draft` after deleting its tag was rejected. A draft may be the reviewed resolution, but it cannot stand in for deciding what the original assertions prove.

## Consumer

PDR-024 requires a rehearsal inventory and human resolution for every multi-identity executable. A separate Heimdall recovery change must review the affected methods and restore or replace their evidence.
