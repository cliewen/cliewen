---
id: CH-210
type: change
status: open
links: [P-025, CAP-013]
title: Reports and handoffs to a person are written in everyday words
---

# CH-210 — Reports and handoffs to a person are written in everyday words

Serves milestone M-115 of plan P-025. It proceeds under draft, inferred VIS-001 without confirming it.

## What and why

The readable-references rule tells an agent to explain what a referenced goal, decision or criterion means, and allows an identifier beside a name. It says nothing about the shape of a status report, a summary of what was done or a handoff, and it does not stop an agent filling a handoff with command lines, file paths and the method's own vocabulary. Observed: a maintainer repeatedly had to ask for explanations without internal names and abbreviations. A baseline run of the repeatable `tracked-brief` scenario on the current skills, with Claude Code, ended with a handoff that names criteria by number (AC-003, AC-004), says "Human-class", "context-isolated pass", "plan-less", "digest", "tracked route", and gives two `clue validate` command lines and a path into `.git`, none of it needed by the reader to act.

The change adds one paragraph to the shared readable-references instruction, which every lifecycle skill inherits. A status report, a summary of what was done and a handoff say what happened, what is needed from the reader, and what comes next, in plain language. An identifier, a file path or a command appears only when the reader asks for it or needs it to act, and then with a few words saying what it is. The existing rule is tightened, not replaced: an identifier beside a name stays allowed in documents.

## Challenge the commitment

The riskiest assumption is that a written rule changes what an agent writes. Text rules in these skills have been followed unevenly, and a rule that merely lists forbidden words is easy to obey on paper and miss in spirit. The credible alternative is a fixed handoff template, which is predictable but turns every handoff into the same form-filling and was already tried for the acceptance brief. The cheapest test is the one the milestone names: run the same scenario before and after and read the two handoffs. The baseline is already in hand. The work stops and returns for a decision if the after-run handoff still carries unexplained internal terms, in which case the rule needs a different shape, not more words.

An implementation could meet every criterion and still fail the reader if the handoff dropped information the reader needs to act, for example never telling them the command to run to accept. So the rule keeps "needs it to act" as the line, and the before/after comparison is read for lost actionable content as well as for jargon.

## Carriers and proof

New criterion AC-241 in CAP-013 (readable references), Unit class, with tests on the generated skill text. Carriers: the shared readable-references source, the regenerated skills, CAP-013 design, the changelog, and plan P-025. Machine-facing data, frontmatter, command arguments and structured evidence keep their identifiers as before. A single run is observation, not a measurement, and is stated so in the milestone note.
