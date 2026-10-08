---
type: guide
title: See the judge work
---

# See the judge work

Cliewen's judge is a command, not a service: `clue validate` reads a repository and either agrees that the thread from intent to acceptance evidence is intact, or names the edge that is broken. This page lets you watch that happen in a disposable Git repository, in about five minutes.

It assumes `clue` is [installed](./install) and that `clue version` prints a release. It does not touch an existing project, and you can remove the whole experiment by deleting one directory.

## 1. Initialize a disposable repository

Create an empty directory instead of experimenting in an existing project:

```sh
mkdir cliewen-demo
cd cliewen-demo
git init
clue init
```

`clue init` creates three marked bootstraps: the architecture and design overviews, and the vision at `docs/vision.md`. An agent completes them from the repository's code and existing documentation, and asks you when a material boundary or intent is unclear. That is why a newly initialized repository is deliberately not green yet. It also creates an empty `docs/use-cases/` folder, which most repositories leave empty.

On a fresh repository, the important final lines look like this:

```text
next: replace the marked vision, architecture, and design bootstraps with repository truth, then run `clue validate`
(docs/use-cases/ is optional and may stay empty — a use case is written only when it explains something the capabilities do not)
```

After the three short bootstraps are replaced, run `clue validate`. It then becomes the normal first green check. The created-file count can change between releases. The top-level tree is:

```text
cliewen-demo/
├── .clue/evidence/      repository-owned exporter examples
├── .agents/skills/       versioned Cliewen skills
├── .claude/skills/       Claude Code mirror of those skills
├── .github/workflows/    the thin CI caller
├── docs/                 the permanent corpus
├── AGENTS.md             routing instructions for coding agents
└── CLAUDE.md             a pointer, because Claude Code does not read AGENTS.md
```

When a newer release is available, ordinary `clue` workflow commands print a one-line notice. It never changes `clue validate` output or its exit status. Run `clue latest` when you want to check deliberately; [Operate Cliewen safely](./operations) explains the update behavior and controls.

`clue init` gives you a starting set of repository files; it does not take ownership of them. You and your agent own the corpus prose and repository-specific instructions. A later `clue init` or `clue scaffold` regenerates only marked README index blocks and skips existing files rather than replacing them. The copied skills and workflow are versioned repository files, not background-managed services.

## 2. See `clue` catch a broken thread

In real work your agent writes these artifacts — that is the whole point of the method, and [what one change produces](./first-change) illustrates one end to end. Here you write them by hand, because three small files are the fastest way to watch the judge change its mind.

Add a tiny goal and capability while keeping its acceptance criterion in `draft`. Create these three files:

:::: code-group

```powershell [Windows PowerShell]
New-Item -ItemType Directory -Force docs/capabilities/CAP-001-greeting | Out-Null
```

```sh [macOS and Linux]
mkdir -p docs/capabilities/CAP-001-greeting
```

::::

:::: code-group

```markdown [docs/goals/G-001-demo.md]
---
id: G-001
type: goal
status: accepted
links: []
title: A greeting is available
---

# G-001 — A greeting is available

A user wants a greeting they can request by name.
```

```markdown [docs/capabilities/CAP-001-greeting/README.md]
---
id: CAP-001
type: capability
status: active
links: [G-001]
title: Return a greeting
goal: G-001
---

# CAP-001 — Return a greeting

The system returns a greeting for a supplied name.
```

````markdown [docs/capabilities/CAP-001-greeting/criteria.md]
---
id: CAP-001-criteria
type: criteria
status: draft
links: [CAP-001]
title: Acceptance criteria for greetings
---

```gherkin
Feature: Return a greeting

  @AC-001
  Scenario: Greet a supplied name
    Given the name "Ada"
    When a greeting is requested
    Then the result is "Hello, Ada"
```
````

::::

Regenerate the two taxonomy indexes and validate:

```sh
clue scaffold
clue validate
```

The draft criterion does not claim verified behavior, so the result is green:

```text
clue scaffold: 2 index block(s) regenerated
clue validate: OK (5 artifacts)
```

Now change only `status: draft` to `status: active` in `criteria.md` and run `clue validate` again. The command exits with status 1 and names the broken edge:

```text
docs/capabilities/CAP-001-greeting/criteria.md: AC-001 has no test (ADR-071: export an executable reference in .clue/evidence.yaml)
```

Additional diagnostics can report the missing `.clue/evidence.yaml`.

That is the product's job: an active machine-proven promise cannot silently lose its acceptance evidence. This deliberately small example is an unannotated legacy criterion, so one supported reference would satisfy it. To return the demo to green, set the whole criteria file back to `draft`.

::: details Evidence rules for production criteria

Before validation, have the agent establish the repository-owned export command using the examples in `.clue/evidence/`. Run it, commit `.clue/evidence.yaml`, and regenerate after changing tests, exporter code or discovery configuration. The same file can combine several frameworks and suites.

A new or revised machine-proven criterion declares `Test-type: Unit`, `Integration`, `E2E`, or `Performance` and requires classified positive and negative executable references in `.clue/evidence.yaml`, unless it explicitly records `(single-direction)`. Each executable carries one canonical AC identity, type and direction in native metadata, executable-bound custom metadata, or a stable naming fallback. Repository-owned exporters attribute and aggregate those references across every framework. The judge checks their identities and input fingerprints without running exporters or tests. A genuine `Test-type: Human` uses the acceptance brief; `@draft` exempts one not-yet-proven criterion; an unannotated legacy criterion retains one reference.

For a criterion whose proof is inherently human, declare `Test-type: Human`; naming it in the acceptance brief is its proof, and no code test is invented. `clue validate` recognizes the Human declaration and waives code evidence, but it cannot check that the brief supplies the proof — the selected acceptance workflow and human merge gate do that. The judge also checks classified pairs, single-direction declarations, and `@draft`; it does not run tests, so your normal test runner remains responsible for whether executable evidence passes.

:::

## 3. Return to green with executable proof

To keep the criterion active, add `Test-type: Unit` to the scenario and provide two tests: one returns `Hello, Ada`, and one rejects an empty name. Use your framework's native metadata, or custom executable metadata, to attach `AC-001`, `Unit` and the appropriate direction to each test.

For a Python trial, create `scripts/` and copy `.clue/evidence/test_example.py` to `scripts/test_greeting.py`. Keep its imports and `cliewen` decorator. Replace the example class with this greeting implementation and its tests:

```python
def greet(name):
    if not name:
        raise ValueError("name is required")
    return f"Hello, {name}"


class GreetingTest(unittest.TestCase):
    @cliewen(ac="AC-001", type="Unit", direction="positive")
    def test_supplied_name(self):
        self.assertEqual(greet("Ada"), "Hello, Ada")

    @cliewen(ac="AC-001", type="Unit", direction="negative")
    def test_empty_name(self):
        with self.assertRaises(ValueError):
            greet("")
```

In a real repository, import the product implementation into the test. The supplied, unchanged sample tests integer parsing and does not prove the greeting criterion.

Copy `.clue/evidence/producers.example.json` to `.clue/evidence/producers.json`. With Python available as `python`, run from the repository root:

```sh
python -m unittest discover -s scripts -p 'test_*.py'
python .clue/evidence/aggregate.py
clue validate
```

With passing greeting tests and a fresh export containing both directions, the criterion has executable proof and validation can be green. Commit the tests, configuration and generated `.clue/evidence.yaml` together. [Acceptance evidence](./acceptance-evidence) explains how to add other frameworks to the same export.

## 4. Remove the experiment or continue

The experiment changed only `cliewen-demo`. Leave that directory, then delete it with your file manager or normal directory-removal command to undo the entire trial. Removing the separately installed `clue` binary is not required.

If the failure made sense and you want to use Cliewen, start again in a new project or read the greenfield and brownfield guide before initializing an existing repository. Before real work lands, configure the protected default branch for human-controlled merge commits only, choose and arm the thin CI caller, and run its disposable probe so a broken thread or provenance-losing merge cannot land.

## Next

[Choose a greenfield or brownfield adoption path.](./adoption)
