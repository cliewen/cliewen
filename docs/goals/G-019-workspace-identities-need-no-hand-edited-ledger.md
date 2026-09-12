---
id: G-019
type: goal
status: proposed
links: [G-009, G-013, VIS-001]
title: A change workspace's own identities are recorded by a command, not written into the ledger by hand
---

# G-019 — A change workspace's own identities are recorded by a command, not written into the ledger by hand

**Who wants it:** every agent and contributor running a full Cliewen change, in this repository and in any adopter (2026-09-12, found while digesting CH-185).

**Why:** a change workspace holds `tasks.md` and `open-questions.md`, and each carries its own ID — `CH-xxx-tasks`, `CH-xxx-open-questions`. `clue validate` requires every artifact ID in the corpus to appear in the identity ledger. No command puts them there. `clue id next` rejects the prefix because it is not numeric, and `clue id live` refuses an ID that was never reserved, so both doors are closed to the only two artifacts the full loop creates automatically.

The ledger library can express these identities — it has an opaque kind and a reservation path for them — but nothing on the command line reaches it. So the practice that has grown up instead is to open `.clue/id-ledger.yaml` and type the two entries in, which is what CH-181 did and what CH-185 did after it.

That practice is the problem, not the inconvenience. The ledger is an append-only event log whose own guidance says never to hand-edit it, because Git's union merge combines independent additions and a hand-written entry can duplicate, conflict, or land in the wrong place. Cliewen's own full loop cannot be completed without breaking that rule, which makes the rule unenforceable and teaches every adopter that hand-editing is normal. [G-009](G-009-workspace-identities-have-a-command.md) asked for workspace identities to have a command rather than a hand-edited ledger and was accepted; this is the part of it that remains true.

There is a second question underneath, worth settling rather than inheriting. These two artifacts are transient — the digest deletes them — yet their entries persist, and the repository is inconsistent about what state they should end in: several changes left them `live` after the workspace was gone, two retired them. Whatever the answer, a command should apply it rather than each change deciding afresh.

The same gap reaches one step further than the workspace. When a digest marks a milestone `done`, `clue validate` requires the milestone's ledger entry to be `retired`, and no command writes that event either: the only code that retires a milestone is the migration that seeds a ledger which does not yet exist. CH-186 and CH-187 both appended the line by hand. It is the same fault in a second place — a state transition the full loop requires, which only an editor can make.

**Success looks like:**

- Opening a change workspace records its artifacts' identities through a command, and a contributor never opens the ledger in an editor to complete the loop.
- The digest leaves those entries in one defined state, the same way every time, without a human choosing.
- Marking a milestone `done` or `dropped` in a digest records its ledger retirement through a command.
- An adopter following the shipped skills never needs to know the ledger's file format.
- The fix does not require the transient artifacts to carry numeric identities, and does not make the workspace heavier to open.
