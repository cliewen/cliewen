---
id: CH-087
type: change
status: open
links: [P-009]
title: Credit JVM acceptance evidence per executable
---

# CH-087 — Credit JVM acceptance evidence per executable

## What

Implement P-009/M-037 by replacing file-level JVM acceptance-evidence harvesting with a conservative per-executable carrier. A Java or Kotlin executable receives classified credit only when one AC identity, one proof type, and one direction attach to that executable through JUnit tags or through a stable named-executable fallback. Unsupported or ambiguous source shapes receive no classified credit and a diagnostic.

The change records the lasting carrier decision in ADR-036, updates CAP-002 and CAP-003 criteria and design, removes the undelivered extraction-installed ArchUnit promise, updates every live guidance and skill carrier that states the JVM evidence contract, and adds focused regression coverage including mixed files, nested and parameterized JUnit cases, ignored proximity comments, and unsupported syntax.

## Why

The current validator collects all JUnit tags from a JVM test file and forms their Cartesian product. In a mixed file, an AC tag on one method can therefore combine with the proof type and direction attached to unrelated methods and make an unproven criterion look covered. This is a false-green acceptance boundary and contradicts G-001's verifiable thread.

ADR-009 tried to compensate for file-level harvesting by delegating per-test purpose enforcement to an ArchUnit rule that extraction promised to install. AN-013 found that the rule was never installed. The deterministic judge must own the carrier it credits and must decline evidence it cannot attribute safely.

## Decision boundary

The implementation statically recognizes Java and Kotlin executable declarations in conventional test source files; it does not compile source, run a JVM framework, or infer semantic test discovery. Method-attached JUnit `@Tag` metadata is the native carrier. A stable executable name carries the same ID, type, and direction for frameworks without native tags. Class-level AC, proof-type, or direction tags do not supply evidence to methods because their scope can spread across unrelated executables; unsupported constructs are reported rather than guessed.

Structured comments remain ignored under ADR-005. The criterion-ID grammar remains ADR-009's current `<PREFIX>-<digits>` contract; P-009/M-038 owns any extension to that grammar.
