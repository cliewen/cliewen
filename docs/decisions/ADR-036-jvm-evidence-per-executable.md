---
id: ADR-036
type: decision
status: verified
links: [ADR-005, ADR-006, ADR-009, ADR-032, CAP-002, CAP-003, P-009]
title: JVM evidence is credited only from one statically attributable executable
author: agent
accepted-by: Flemming N. Larsen (2026-08-02, conversation)
---

# ADR-036 — JVM evidence per executable

## Context and problem statement

File-level JVM harvesting and tag cross-products can credit one executable from evidence attached to another. Classified evidence must instead be attributable to exactly one Java or Kotlin test executable.

## Decision outcome

**A JVM acceptance-evidence reference is one statically attributable Java or Kotlin executable carrying exactly one AC identity, one proof type, and one direction.** In conventional `*Test` or `*Tests` files, supported JUnit methods carry a contiguous annotation block with literal one-line tags; parameterized invocations count once, nested methods count normally, and enclosing-class tags do not supply method evidence.

Class-level AC tags, ambiguous multiple identities/types/directions, and dynamic or multiline tag forms receive diagnostics and no classified credit; unrelated tags remain ordinary runner metadata. Frameworks without native tags use the stable executable name `test<PREFIX><digits><lowercase-suffix>_<Type><Direction>_<description>`, with ADR-037 prefix normalization. Native tags and the fallback name must agree when both are present. Proximity comments, source compilation, framework discovery, and parameterized expansion are not part of the parser.

**Carrier:** the JVM evidence harvester and its diagnostics in `checkACTests`, the canonical ID grammar, and the acceptance-evidence criteria.
