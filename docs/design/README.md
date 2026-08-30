# Design

This is Cliewen's cross-cutting behaviour overview. It explains how the agent workflow, deterministic CLI, durable corpus, CI wall, and human acceptance boundary work together. [Architecture](../architecture/README.md) covers their static boundaries; capability designs hold local implementation detail.

```mermaid
sequenceDiagram
    participant U as User
    participant A as Agent
    participant C as clue
    participant D as /docs
    participant H as Human
    U->>A: Request a change
    A->>C: Read context and validate
    A->>D: Update current truth and evidence
    C-->>A: Deterministic corpus result
    A->>H: Present reviewed pull request
    H->>D: Merge accepts the change
```

The feedback loop is deliberate: a change first asks which durable document a reader will need, updates that document only when the system truth changed, then verifies the corpus against the same revision. Decisions explain why a future choice constrains work and link here or to architecture when they affect this view; they do not duplicate it. Diagrams belong here when they make an interaction or flow easier to review than prose.

<!-- clue:index:start -->
<!-- clue:index:end -->
