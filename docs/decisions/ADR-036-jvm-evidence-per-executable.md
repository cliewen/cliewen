---
id: ADR-036
type: decision
status: inferred
links: [ADR-005, ADR-006, ADR-009, ADR-032, CAP-002, CAP-003, P-009]
title: JVM evidence is credited only from one statically attributable executable
author: agent
accepted-by: []
---

# ADR-036 — JVM evidence is credited only from one statically attributable executable

## Context and problem statement

ADR-009 made JVM harvesting file-level because `clue` did not parse Java or Kotlin structure, and ADR-032 then classified every AC identity in a file against every proof-type and direction tag found anywhere in that file. A file containing unrelated tests can therefore create evidence triples no executable carries. ADR-009 delegated per-test purpose enforcement to an extraction-installed ArchUnit rule, but AN-013 found that extraction never installed that carrier. Cliewen cannot credit acceptance evidence through a rule neither the judge nor the installed methodology demonstrably supplies.

## Decision outcome

**A JVM acceptance-evidence reference is one statically attributable Java or Kotlin executable carrying exactly one AC identity, one proof type, and one direction.** This replaces ADR-009's file-level JVM harvesting and ADR-032's file-level classified cross-product; their namespace, proof-class, paired-direction, and non-JVM decisions remain unchanged.

For JUnit, the supported native carrier is a contiguous method annotation block attached to a Java method or Kotlin `fun` declaration in a conventional `*Test.java`, `*Tests.java`, `*Test.kt`, or `*Tests.kt` file. The executable must carry `@Test`, `@ParameterizedTest`, `@RepeatedTest`, `@TestFactory`, or `@TestTemplate`; its same block carries literal one-line `@Tag` values for the AC identity, proof type, and direction. Fully qualified annotation names are equivalent. Parameterized tests receive one evidence credit for the executable, not one credit per generated invocation. A method inside any depth of `@Nested` class is treated like any other method; tags on an enclosing class do not become method evidence.

Class-level AC tags receive no credit and produce a diagnostic because they can cover several executables. Class-level proof-type and direction tags remain runner metadata only and do not supply missing parts to a method. A method block containing several AC identities, proof types, or directions is ambiguous: every attached AC reference is still checked for unknown or retired identity, but the block produces no classified credit and a diagnostic. A dynamic, multi-line, or otherwise unsupported `@Tag` form likewise produces no credit and a diagnostic instead of being guessed. Tags outside every declared AC namespace remain ordinary runner metadata.

For JVM frameworks without native tags, the stable fallback is the executable name `test<PREFIX><digits>_<Type><Direction>_<description>`, for example `testMG101_IntegrationPositive_acceptsValidInput`. The name must belong to a Java method or Kotlin `fun` declaration in a conventional JVM test file; its AC prefix must be declared by the corpus. The name itself declares that executable as the test carrier, so no JUnit annotation is required. If native tags and the executable name both carry evidence, their triples must agree or the executable is ambiguous and receives no classified credit.

Structured proximity comments remain ignored under ADR-005. This parser recognizes a deliberately small source convention; it does not compile source, run framework discovery, expand parameterized invocations, or infer that an arbitrary helper is a test. A JVM framework or source shape that cannot use the supported annotation block or named-executable form is reported as unsupported rather than credited.

## Consequences

- Unrelated methods in one source file can no longer satisfy each other's evidence contract.
- `clue validate` owns the complete carrier it credits; extraction no longer promises to install an external ArchUnit substitute.
- Java and Kotlin retain one portable, grep-visible fallback where framework-native tags are unavailable.
- Conservative diagnostics can require a source-only annotation reformat even when the JVM compiler accepts a more dynamic form. That cost is deliberate: false negatives are repairable, while false acceptance evidence crosses the core red line.
