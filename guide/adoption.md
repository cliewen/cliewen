# Greenfield and brownfield

Cliewen works for a new system and for one that already has years of history. The difference is the first move: a greenfield project can state its intended outcomes directly, while a brownfield project must first find and reconcile the intent that already exists.

## Who keeps the documentation current?

Tell the agent what outcome you want. You should not have to mirror every code change into `/docs` by hand. The agent reads `AGENTS.md`, loads the relevant Cliewen skill, and updates the implementation and durable corpus together on the same branch.

`clue validate` checks the parts a machine can judge: artifact structure, links, generated indexes, and traceability from active acceptance criteria to tests. A human still reviews whether the documentation and implementation say the right thing.

This is agent-maintained documentation, not background synchronization. `clue` does not watch a wiki or ticket system, and it does not invent missing intent from code. The change loop requires local validation before a Cliewen pull request is ready for review. Once the generated CI wall is armed and its validation job is a required check, broken traceability blocks merge. Plain changes keep the same required job but do not invoke the corpus validator.

## Prompts that get useful work started

These are examples, not magic phrases. The repository's `AGENTS.md` handles skill routing, so a prompt can stay focused on the outcome.

### Start a greenfield project

After `clue init`, give the agent the first outcome rather than a proposed file layout:

```text
This is a new system. Use Cliewen to turn this outcome into the first goal and a small plan: <outcome>.
```

The agent should establish the goal, make uncertainty visible, and propose the smallest verifiable plan before implementation.

### Make a routine change

Once the corpus exists, name the behavior and ask for the complete change:

```text
Using Cliewen, add <behavior>. Keep the acceptance criteria, tests, code, and durable documentation aligned, then prepare the pull request.
```

The agent follows the change loop and leaves the merge decision to a human.

### Adopt one existing repository

Use `clue-extract` once when the repository already contains specifications, decision notes, tagged tests, or other durable intent:

```text
Use clue-extract to adopt this repository. Preserve existing IDs and test traceability. Mark every extracted non-decision artifact with provenance: inferred, and every extracted decision with status: inferred and author: agent. Stop when sources conflict.
```

Extraction is a meaning-level conversion, not a file copy. Existing evidence is mapped into one Cliewen corpus, and every extracted artifact begins inferred: non-decision artifacts use `provenance: inferred`, while decisions use `status: inferred` and `author: agent`. Human review promotes only the meaning it verifies. The old parallel specification corpus is removed in the same pull request.

One extraction mapping ships today: OpenSpec as extended in [Intent Engineering for Coding Agents](https://intent-engineering-for-coding-agents.github.io/book/). Stock OpenSpec does not tag scenarios with stable IDs, so the mapping expects the book's conventions. Specs become capabilities, scenario IDs survive as acceptance-criterion IDs, so existing test tags keep working, and archived changes remain Git history. If the source format has no extraction mapping yet, writing that mapping is the first extraction task.

## When the system spans repositories, wikis, and tickets

A `clue-analysis` discovery pass is useful when evidence and ownership are distributed across several systems. It can establish what the sources are, how fresh they are, and where they disagree before you choose which repository-local extractions to propose.

```text
Use clue-analysis to inventory evidence across <repositories, wiki, tickets>. Record source revisions and freshness, distinguish observation from inference, surface conflicts, and recommend corpus boundaries for repository-local extraction.
```

Wiki pages and tickets can be evidence, preferably through revision-pinned links or stable exports. They do not become a second system of record, and Cliewen does not live-sync them after adoption.

The current tooling has a deliberate repository boundary:

- `clue-extract` adopts one repository at a time.
- `clue validate` discovers acceptance evidence only inside the repository being validated.
- Several repositories can be adopted separately, each with its own corpus and local test evidence.
- One unified corpus that claims acceptance evidence from tests spread across several repositories is not supported yet. Supporting it would require a future capability rather than a broader reading of the current tools.

Once the boundary is clear, return to [the change loop](./change-loop) for the first adopted change.
