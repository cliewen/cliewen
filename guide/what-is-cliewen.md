# What is Cliewen?

Cliewen is a methodology and a command-line tool for building software with coding agents while keeping intent, implementation, and evidence connected. Its name comes from the Old English word for a ball of thread: the same word that became *clue*.

The central idea is simple: the durable documentation describes the system as it exists, not a pile of past change requests. A goal leads to a capability, a capability owns acceptance criteria, and each active criterion has focused test evidence. The `clue` command checks that this thread is intact.

## Why another workflow?

An agent can produce a convincing patch without understanding why the system exists. It can also update a specification without updating the tests, leave decisions buried in a chat, or quietly reinterpret an old acceptance criterion. More agent autonomy makes those gaps more expensive, not less.

Cliewen gives each kind of truth one durable home and makes the hand-offs explicit:

- The corpus under `/docs` is the system of record.
- A branch is a proposal and a pull request is the human review gate.
- A transient `/changes/CH-xxx-*` workspace holds the delta while it is being built, then disappears before merge.
- The `clue` CLI enforces structure and traceability.
- The human reviewer decides whether the change means the right thing.

## Born from Intent Engineering and spec-driven development

Cliewen builds on the ideas in [Intent Engineering for Coding Agents](https://intent-engineering-for-coding-agents.github.io/book/), written by Cliewen's author, Flemming N. Larsen: human intent is written down before an agent implements it, and the shared ground between human and agent lives in the repository under version control. Cliewen carries that approach one step further. The durable documentation is where that intent lives, and the `clue` binary enforces what the book otherwise leaves to discipline.

The book's working example of spec-driven development is [OpenSpec](https://github.com/Fission-AI/OpenSpec), where a change-sized spec is proposed, applied, and then moved to an archive folder to keep the workspace clean. Cliewen keeps that proposal layer but needs no archive step: by the time a pull request merges, the transient `/changes` workspace has been digested into the durable documentation under `/docs` and deleted, and the merged pull request itself is the historical record. Instead of a spec that goes stale after implementation, the documentation is the spec, and every change is required to leave it true. A repository already using the book's extended OpenSpec format can be adopted with its IDs and test traceability intact; the [greenfield and brownfield guide](./adoption) shows how.

Decisions follow the same rhythm. A decision an agent records during a change is born `inferred`; merging the pull request makes it binding, and a later explicit human approval (from a solo developer or any team member, at whatever pace review actually happens) promotes it to `verified`. Shipping never blocks on an approval ritual, and the count of decisions no human has yet endorsed stays honestly visible.

That combination prevents two common failures of change-centered specifications: a growing archive of stale proposals that must be reconstructed to understand current behavior, and a polished permanent specification whose connection to executable evidence is only assumed.

## What Cliewen is not

Cliewen is not an issue tracker, a project-management service, or a way to remove humans from engineering decisions. It is also not a replacement for tests. It depends on them: every active acceptance criterion must have focused positive and negative tests that show the stated intent is met, and `clue validate` fails when that evidence is missing. It is deliberately repo-native: Markdown, Git, the test framework you already use, one small binary, and skills that teach agents the workflow.

Ready to see the pieces? Start with [the verifiable thread](./methodology), or go straight to [installing Cliewen](./getting-started).

