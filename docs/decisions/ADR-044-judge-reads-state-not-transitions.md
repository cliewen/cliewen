---
id: ADR-044
type: decision
status: verified
links: [P-010, CAP-002, ADR-017, C-004, C-008, C-012, C-013]
title: The judge judges a repository state, never a transition
author: agent
accepted-by: Flemming N. Larsen (2026-08-08, conversation)
---

# ADR-044 — The judge judges a repository state, never a transition

## Context and problem statement

Some of the rules that bind a change are about the change itself. Did this one add a changelog entry for the behaviour it altered? Did it loosen a test instead of fixing what the test caught? Did it edit a plan that was finished? Did it touch a core carrier without recording a decision?

Every one of those questions compares two states, and none of them can be answered by looking at one. The obvious place to put the answer was the judge, because the judge is the thing that already fails a change — and so the constraint register accumulated four rules whose stated route to enforcement was "`clue` gains git-diff context". Nothing had decided that it should.

`clue validate` is the deterministic judge ([ARCH-003](../architecture/core.md)): the same corpus yields the same verdict for everyone, everywhere, offline, forever. A verdict computed from a diff is not that verdict. It depends on which branch the caller is standing on, on whether the base is present in the clone at all, and on how deep that clone is — a shallow CI checkout and a developer's worktree would disagree about a corpus whose bytes are identical. The judge would answer a different question depending on where it was asked.

## Decision outcome

**`clue validate` reads the repository as it is. It never reads history, a diff, or a base revision, and no rule about what a change did becomes a validation verdict.**

*The judge answers from what is in front of it.* Its input is the working tree. This is the same boundary [ADR-040](ADR-040-qualified-external-references.md) drew around resolving an external reference and [ADR-042](ADR-042-release-check-outside-the-judge.md) drew around the release list, for the same reason and with the same consequence: a verdict that depends on something outside the files under review is not reproducible against a pinned revision.

*A transition rule is enforced by a machine that is allowed to have a base.* Continuous integration knows what it is merging into; that is its purpose, and comparing against the base is ordinary work for a workflow step. A forge's branch protection knows what reached the default branch and how. A release workflow knows what it is publishing. These are real machines, and a rule enforced by one of them is machine-enforced — the register says which one rather than pretending the rule is unheld because `clue` does not hold it.

*A rule no machine can hold is declared, not queued.* Whether a deleted assertion was a weakening or a refactor, whether a decision record is timeless, whether a change alters the meaning of the core — these are meaning, and meaning is the half of the method that humans verify. Recording them as pending automation misrepresents them: it promises a check that should never be written, and it lets a permanent property masquerade as a backlog item.

## What this retires

The promotion triggers on [C-002](../constraints/C-002-changelog-per-user-visible-change.md), [C-004](../constraints/C-004-never-weaken-checks.md), [C-006](../constraints/C-006-adrs-timeless-with-carrier.md), [C-008](../constraints/C-008-completed-plans-immutable.md), and [C-013](../constraints/C-013-core-changes-need-decision.md) each named diff context inside `clue` as the condition for becoming machine-enforced. That condition is withdrawn. Each of those constraints now names the machine that actually holds its detectable part, or states the judgment that remains and what it costs.

## Rejected: a second, diff-aware `clue` command

Attractive, because it keeps the judge clean while still writing the checks in Go. It fails on what it would create: a second verdict with the authority of the first and none of its guarantees. Two commands that both say a change is unacceptable, one of which answers differently depending on the caller's fetch depth, is a worse contract than one command with a stated edge. The rules that need a base already run somewhere that has one.

## Rejected: read the base only when it happens to be available

The degradation this repository already applies to the network — could not tell, stay quiet, exit 0 — is right for an advisory command and wrong for a judge. A rule that is enforced in CI and silently skipped in a worktree teaches its readers that green means nothing in particular, and the failure mode is invisible: the check that did not run leaves no trace.

## Carrier

CAP-002 owns what a green validate asserts and states this boundary in its design. [ADR-017](ADR-017-conventions-are-constraints.md)'s register carries the consequence per constraint, and [ADR-045](ADR-045-register-names-the-machine.md) defines how a constraint states the machine that holds it. The standing rejections this decision does not reopen: the judge stays offline, deterministic, and local to the files it is given.
