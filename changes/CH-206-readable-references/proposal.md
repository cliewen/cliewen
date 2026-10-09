---
id: CH-206
type: change
status: open
links: [G-001, CAP-002, CAP-004]
plan: none
title: Make references understandable without requiring identity lookups
---

# Proposal

The owner selected one method for all human-facing surfaces: explain the relevant meaning inline, use descriptive links for optional detail, and show an identity only as secondary information when it helps precise identification. Metadata, literal command arguments and identity-allocation results retain their canonical identities. Changelog entries describe user-visible changes and upgrade actions without internal rule identifiers.

This plan-less change adds a readable-reference goal and capability, a project decision and agent-enforced convention, one canonical instruction available in all six standalone skills, aligned scaffold and handoff templates, guide examples, and named human CLI reports and diagnostics. It preserves structural validation, proof semantics, states, ordering, exit codes and machine-facing data. Published notes and completed plans are not rewritten. The pilot's editorial correction remains on its independent analysis branch, whose unaccepted work is not a dependency.

## Challenge

The riskiest assumption is that adding a title makes a reference understandable. A document could follow the formatting rule and still hide the consequence the reader needs, or make the reader navigate a long reference table. The instruction therefore requires the relevant effect in the sentence itself; a link and name alone are insufficient when they leave the argument implicit. The cheapest check is to review the pilot paragraph, a changelog entry, a handoff and representative CLI outputs without opening links. If their purpose or required action is still unclear, revise the text before handoff. The alternative is to repair one document only, which would leave the same failure in the other selected surfaces.

A technically correct implementation can still mislead: a formatter could label a Cliewen source rule using an adopter's unrelated artifact with the same identifier. CLI naming must use explicit local-reference origin, not blindly substitute every identity token in a message. Unknown, ambiguous and foreign references remain honest about what can be named without network access.

The repository states a draft inferred vision for evidence-backed intent and human acceptance; this change supports its human-readable review boundary without claiming that the vision has been confirmed. Agents prepare and preflight; the human performs local acceptance.
