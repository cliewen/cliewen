# Metadata patterns to adapt

These examples use AC-001 only as a placeholder. Replace it with a real Gherkin criterion and retain its assertions when translating an existing test. Each example illustrates the positive and negative directions of the same behavior: a numeric string parses, and an invalid string is rejected. The test runner still executes the tests; export records their references.

## JUnit

```java
@Test @Tag("AC-001") @Tag("Unit") @Tag("positive")
void validInput() { assertEquals(42, Integer.parseInt("42")); }

@Test @Tag("AC-001") @Tag("Unit") @Tag("negative")
void invalidInput() {
    assertThrows(NumberFormatException.class, () -> Integer.parseInt("bad"));
}
```

The source repository's compatibility producer in `tools/export-evidence` shows conservative literal JUnit attribution outside the judge. For a native discovery producer, use the [JUnit Platform Launcher API](https://docs.junit.org/current/advanced-topics/launcher-api.html), inspect each MethodSource and direct method metadata, and emit source-qualified method references. Reject class-level AC metadata before interpreting inherited discovery tags. A method-targeted custom annotation can carry structured fields; read those fields explicitly rather than assuming annotation arguments become native JUnit tags. [JUnit tags](https://docs.junit.org/current/writing-tests/tagging-and-filtering.html) and composed annotations are the native starting points.

## NUnit

```csharp
[Test, Property("Cliewen.AC", "AC-001"), Property("Cliewen.Type", "Unit"), Property("Cliewen.Direction", "positive")]
public void ValidInput() => Assert.That(int.Parse("42"), Is.EqualTo(42));

[Test, Property("Cliewen.AC", "AC-001"), Property("Cliewen.Type", "Unit"), Property("Cliewen.Direction", "negative")]
public void InvalidInput() => Assert.Throws<FormatException>(() => int.Parse("bad"));
```

[PropertyAttribute](https://docs.nunit.org/articles/nunit/writing-tests/attributes/property.html) supports named metadata and custom derived attributes. Use direct method metadata, not properties inherited from a fixture or assembly. An exporter can use [NUnit exploration](https://docs.nunit.org/articles/nunit/running-tests/Console-Command-Line.html) and inspect the compiled method's direct attributes to check origin. Resolve the source from the project's actual source/PDB layout, fingerprint those sources and exporter/configuration, and emit one reference per executable. Parameterized results must collapse to their declaration, not increase proof by their case count.

## Pytest and unittest

```python
@pytest.mark.cliewen(ac="AC-001", type="Unit", direction="positive")
def test_valid_input():
    assert int("42") == 42

@pytest.mark.cliewen(ac="AC-001", type="Unit", direction="negative")
def test_invalid_input():
    with pytest.raises(ValueError):
        int("bad")
```

Register one `cliewen(ac, type, direction)` marker in pytest configuration; there is no per-AC marker registry. [Pytest custom markers](https://docs.pytest.org/en/stable/example/markers.html) expose their arguments and originating nodes during collection. Reject markers originating on a class/module. `python_producer.py` gives a bounded AST alternative for literal direct markers, and `test_example.py` shows a custom unittest decorator preserving the original method. Adapt discovery for custom loaders, imported classes, or nonstandard patterns rather than silently ignoring them.

## Vitest

```typescript
test('valid input', { tags: ['AC-001', 'Unit', 'positive'] }, () => {
  expect(Number('42')).toBe(42)
})
test('invalid input', { tags: ['AC-001', 'Unit', 'negative'] }, () => {
  expect(Number.isNaN(Number('bad'))).toBe(true)
})
```

Use [Vitest tags](https://vitest.dev/guide/test-tags.html) and establish the repository's tag configuration. Where the pinned version permits it, `strictTags: false` allows AC identities without per-ID registration; classification values are still checked by export. `vitest_producer.mjs` demonstrates [collection without executing test callbacks](https://vitest.dev/api/advanced/vitest.html#collect). Its explicit collection loads test modules, so top-level side effects remain the repository's responsibility. The producer rejects suite AC metadata and checks parameterized declaration locations. Validate the pinned API before adopting the example. Use a literal bracketed title when tags are unsuitable.

## Go and Cucumber

Go uses `TestAC_001_UnitPositive_validInput` and `TestAC_001_UnitNegative_invalidInput`, with the required `Test` discovery prefix. The compatibility producer also accepts existing compact names such as `TestAC001_UnitPositive_validInput`. Source extraction uses the Go AST so strings and comments cannot supply test declarations. Cucumber keeps `@AC-001 @unit @positive` and `@AC-001 @unit @negative` directly on their scenarios; feature-level AC metadata is not executable evidence.

## Gatling

```java
ScenarioBuilder scn = scenario("[AC-001 Performance Positive] valid input")
    .exec(http("parse").get("/parse?value=42").check(status().is(200)));
setUp(scn.injectOpen(atOnceUsers(1)))
    .assertions(global().failedRequests().count().is(0L));
```

The negative executable uses its own literal `[AC-001 Performance Negative] invalid input` title and an assertion/check for the intended failure response. `gatling_producer.py` illustrates a bounded Java/Kotlin source shape with variable attribution and `setUp(...).assertions(...)`. It diagnoses tagged scenarios missing registration/assertions. Inspect and test a native or richer source producer when the simulation uses helper classes, multiple constructors, chained setup variants, or shared scenario builders. Review the connection between each criterion and its [Gatling assertions](https://docs.gatling.io/concepts/assertions/); an assertion's presence does not prove its relevance.
