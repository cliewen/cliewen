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

## How it differs from change-centered specifications

Change-centered frameworks are good at describing what one patch should do. Cliewen keeps that proposal layer, but treats it as temporary. Once accepted, the useful knowledge is digested into the permanent model of the system and Git preserves the historical delta.

That distinction prevents two common failures: a growing archive of stale proposals that must be reconstructed to understand current behavior, and a polished permanent specification whose connection to executable evidence is only assumed.

## What Cliewen is not

Cliewen is not an issue tracker, a project-management service, a replacement for tests, or a way to remove humans from engineering decisions. It is deliberately repo-native: Markdown, Git, the test framework you already use, one small binary, and skills that teach agents the workflow.

Ready to see the pieces? Start with [the verifiable thread](./methodology), or go straight to [installing Cliewen](./getting-started).

