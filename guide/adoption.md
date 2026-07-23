# Greenfield and brownfield

Cliewen works for a new system and for one that already has years of history. The difference is the first move: a greenfield project can state its intended outcomes directly, while a brownfield project must first find and reconcile the intent that already exists.

## Who keeps the documentation current?

Tell the agent what outcome you want. You should not have to mirror every code change into `/docs` by hand. The agent reads `AGENTS.md`, loads the relevant Cliewen skill, and updates the implementation and durable corpus together on the same branch.

`clue validate` checks the parts a machine can judge: artifact structure, links, generated indexes, and traceability from active acceptance criteria to tests. A human still reviews whether the documentation and implementation say the right thing.

This is agent-maintained documentation, not background synchronization. `clue` does not watch a wiki or ticket system, and it does not invent missing intent from code. The change loop requires local validation before a Cliewen pull request is ready for review. Once the [generated CI wall](./ci-wall) is armed and its validation job is a required check, broken traceability blocks merge. Plain changes keep the same required job but do not invoke the corpus validator.

## Start with the minimum

Do not fill every corpus folder because the scaffold created it. A useful first thread needs four things:

| Record | Question it answers |
|---|---|
| Goal | Who needs an outcome, and why? |
| Capability | What must the system be able to do? |
| Acceptance criterion | What observable example makes the behavior specific enough to accept? |
| Positive and negative test evidence | Does the behavior work, and does it reject or survive the important counter-case? |

The criterion carries a stable ID such as `AC-001`. The focused positive and negative tests both reference that ID using the convention for their test framework. Keep the criterion draft until the implementation and both tests are ready, then activate it in the same change.

`clue validate` checks that an active criterion has at least one supported test reference. It does not count or classify the pair yet. The change loop and review hold the stronger positive-and-negative rule.

## Add the wider corpus when it earns its keep

The rest of the taxonomy solves problems that appear as work grows:

| Add | When |
|---|---|
| A plan | The outcome needs several ordered changes or milestones |
| A decision record | Reversing a choice later would be expensive, or the repository needs to remember why one path won |
| A constraint | Law, policy, compatibility, licensing, another non-negotiable boundary, or a system quality (performance, reliability, usability) needs a concrete, checked threshold |
| Architecture | Several capabilities depend on the same boundary or an expensive-to-change structure |
| Analysis findings | Important unknowns need investigation before anyone can plan honestly |

Leave unused categories empty. Cliewen is supposed to expose necessary reasoning, not reward document volume.

## When Cliewen is a poor fit

Do not adopt Cliewen when the repository cannot own both the intent and its acceptance evidence. The current method is also a poor fit when:

- Work does not go through Git branches and a human-controlled merge boundary.
- The project cannot run reliable tests or enforce a stable CI check before integration.
- The code is a disposable prototype, generated output, or vendored source whose behavior is accepted somewhere else.
- One corpus would need to claim test evidence spread across several repositories. Current validation is repository-local.
- The team will not let agents update the corpus with the implementation or will not review the meaning before merge.

In those cases, use the project's existing lightweight notes and tests instead of creating a corpus nobody will maintain.

## Prompts that get useful work started

You do not need to speak Cliewen's internal language. Describe what you want in ordinary terms; the repository's `AGENTS.md` tells the agent which workflow to follow. For example:

### Start a greenfield project

After `clue init`, give the agent the first outcome rather than a proposed file layout:

```text
I'm starting a new system that should <outcome>. Help me work out a small first version.
```

The agent should establish the goal, make uncertainty visible, and propose the smallest verifiable plan before implementation.

### Make a routine change

Once the corpus exists, name the behavior and ask for the complete change:

```text
Please add <behavior> and get it ready for review.
```

The agent follows the change loop and leaves the merge decision to a human.

### Adopt one existing repository

Use `clue-extract` once when the repository already contains specifications, decision notes, tagged tests, or other durable intent:

```text
Bring this repository into Cliewen. Keep the links between its existing specifications and tests, and flag anything that disagrees.
```

Extraction is a meaning-level conversion, not a file copy. Existing evidence is mapped into one Cliewen corpus, and every extracted artifact begins inferred: non-decision artifacts use `provenance: inferred`, while decisions use `status: inferred` and `author: agent`. Human review promotes only the meaning it verifies. The old parallel specification corpus is removed in the same pull request.

One extraction mapping ships today: OpenSpec as extended in [Intent Engineering for Coding Agents](https://intent-engineering-for-coding-agents.github.io/book/). Stock OpenSpec does not tag scenarios with stable IDs, so the mapping expects the book's conventions. Specs become capabilities, scenario IDs survive as acceptance-criterion IDs, so existing test tags keep working, and archived changes remain Git history. If the source format has no extraction mapping yet, writing that mapping is the first extraction task.

## When the system spans repositories, wikis, and tickets

A `clue-analysis` discovery pass is useful when evidence and ownership are distributed across several systems. It can establish what the sources are, how fresh they are, and where they disagree before you choose which repository-local extractions to propose.

```text
Before we adopt Cliewen, investigate the risks and unknowns around where our intent lives across <repositories, wiki, tickets>. Find what is still current and what conflicts, then recommend what should live in each repository.
```

Wiki pages and tickets can be evidence, preferably through revision-pinned links or stable exports. They do not become a second system of record, and Cliewen does not live-sync them after adoption.

The current tooling has a deliberate repository boundary:

- `clue-extract` adopts one repository at a time.
- `clue validate` discovers acceptance evidence only inside the repository being validated.
- Several repositories can be adopted separately, each with its own corpus and local test evidence.
- One unified corpus that claims acceptance evidence from tests spread across several repositories is not supported yet. Supporting it would require a future capability rather than a broader reading of the current tools.

## Next

[Follow the change loop for the first adopted change.](./change-loop)
