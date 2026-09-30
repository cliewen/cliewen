# Acceptance evidence across frameworks

Acceptance evidence connects a criterion to the executable that proves it. Since Cliewen 0.27.0, repository-owned exporters write those connections to `.clue/evidence.yaml`. `clue validate`, coverage and parity read the same file. Cliewen does not load test frameworks or run exporters or tests.

## One file, several producers

A producer is a suite or export source. A Java backend, Python scripts, Vitest UI tests and Gatling scenarios can each have a producer in the same manifest. Two suites using the same framework can also have separate producers. Give each producer a unique ID; framework and language names are descriptive values, with no built-in support registry.

Each reference names one canonical acceptance criterion ID, source path, executable subject, test type and direction. The subject identifies the method, test call or scenario within its source file. Parameterized cases belong to one executable. A test proving several criteria needs separate executables, each with its own attribution.

## Choose metadata your tests can carry

Prefer the framework's native executable metadata. If it cannot carry the required values, use a custom annotation, attribute or decorator attached directly to the executable. Keep ordinary test names readable when metadata is available.

Ordinary `# AC:` or `// AC:` comments near a test provide no credit. Metadata on a containing suite, class or assembly also provides no executable proof. An exporter must establish which executable owns the evidence, rather than infer ownership from proximity.

When metadata is unsuitable, use a stable name or literal title. The portable fallback encodes canonical hyphens as underscores, preserves numeric spelling and lowercase suffixes, and combines type and direction:

```text
test_PDO_115_UnitPositive_preserves_pages
[PDO-115 Unit Positive] preserves pages
```

Adapt the prefix to the framework's discovery rules. There is no universal test-name syntax: a producer must recognize the chosen convention and reject ambiguous or dynamic identities. Existing compact Go/JVM names can be retained through the compatibility exporter example.

`clue init` supplies examples under `.clue/evidence/`. Read `frameworks.md` for JUnit, NUnit, pytest and other metadata examples, and `README.md` for the common export contract. These are starting points for repository-owned exporters; test them against your actual discovery configuration before relying on them.

## Establish the export command

For the supplied Python example, put tests in `scripts/test_*.py`, copy `.clue/evidence/producers.example.json` to `.clue/evidence/producers.json`, then run these commands from the repository root:

```sh
python -m unittest discover -s scripts -p 'test_*.py'
python .clue/evidence/aggregate.py
clue validate
```

The example assumes `python` names your Python interpreter. Adjust both the shell command and the producer's `command` field if your environment uses another name. The producer reads source metadata; the separate unittest command establishes whether the tests pass.

For another suite, add a producer to `producers.json` with its own command and include/exclude scopes. The aggregator supplies the discovered input paths to each command and combines its references. Review each scope: it must cover test sources, exporter code and discovery configuration, while excluding dependencies and generated output. A fresh export cannot detect a suite you omitted from its configuration.

Commit the producer configuration, exporter code and `.clue/evidence.yaml` with the tests. The manifest is generated, including normalized source hashes and the complete matching input inventory. Do not maintain its references or hashes by hand. Regenerate after modifying tests, exporters or discovery configuration, including adding or deleting a matching file.

The supplied aggregator runs every declared producer and replaces the manifest only after the complete export succeeds. A producer failure preserves the previous file and fails the command; it must not publish a partial success. Exporter tests and human review establish discovery, attribution and what the assertions prove.

## Validation and execution are separate checks

`clue validate` checks identities, classification, diagnostics and freshness within the declared scopes. A green verdict says the references satisfy the evidence contract. Your normal test runner checks that the tests pass; review checks whether their assertions prove the criterion.

A new or revised machine-proven criterion declares `Unit`, `Integration`, `E2E` or `Performance` and needs positive and negative references, unless it records `(single-direction)`. Genuine `Test-type: Human` proof belongs in the acceptance brief. Neither `Human` nor `@draft` is a repair for missing automated evidence.

Run tests and export locally before validation. CI may also regenerate the manifest and fail if it differs from the committed file, alongside the normal test gates. Cliewen's validation command never invokes that regeneration itself. See [operations](./operations#recover-without-bypassing-the-evidence) for missing exports, stale inputs and producer failures.

## Next

[Learn what each skill does and when your agent uses it.](./skills)
