# Pre-contract migration fixture

This fixture is pinned to the corpus shape used before the `reversal-cost` field and the narrowed default lifecycle were enforced. Its artifact shape reproduces the historical blocker measured in AN-012: inferred meaning has no routing field and an analysis uses `status: verified`. The fixture is migration input only; the test supplies the current generated carrier set so the post-apply corpus can be judged by the current validator.
