---
layout: home

hero:
  name: Cliewen
  text: Evidence-backed Intent Engineering for coding agents
  tagline: Coding agents produce plausible changes faster than anyone can confidently accept them. Cliewen keeps requirements, decisions, implementation, and acceptance evidence connected in Git—and checks that the connection is still intact before merge.
  image:
    src: /cliewen-logo.svg
    alt: Cliewen logo
  actions:
    - theme: brand
      text: Start reading
      link: /what-is-cliewen

features:
  - title: For agent-driven pull requests
    details: Built for repositories where coding agents implement real product changes and tests through Git branches and pull requests.
  - title: One thread the repository can check
    details: Goals lead to capabilities, acceptance criteria, and declared evidence—classified test references or genuine Human proof. The clue CLI reports broken links and missing evidence locally and in CI without executing tests.
  - title: Methodology, judge, and memory
    details: Cliewen is the methodology, clue is its command-line judge, and the corpus under docs is the permanent system record that agents maintain with the code.
  - title: Deliberately visible overhead
    details: Agents prepare the corpus and verified proposal; humans keep control of intent and merge. Small work that changes no meaning stays outside the full loop.
---

## The thread

```mermaid
graph LR
  G["Goal"] --> C["Capability"]
  C --> A["Acceptance criterion"]
  A --> E["Acceptance evidence"]
  E -. "why do we have this?" .-> G
```

Pick up any artifact and follow it back to why it exists, or forward to what proves it. `clue` checks that no arrow is missing. It does not run your tests, it cannot tell you that a test proves the right thing, and it cannot tell you the goal was worth having—that judgment stays yours, at the merge.

## Three names, three things

- **Cliewen** — the methodology.
- **`clue`** — the deterministic command-line judge.
- **The corpus under `docs/`** — the durable record of the system as it exists.

## Three moves

```mermaid
graph LR
  I["Install clue"] --> P["Prompt your agent"]
  P --> A["It writes the change, the corpus, and the evidence"]
  A --> V["clue validate and CI check the thread"]
  V --> M["You merge"]
```

Install takes one command. Prompting takes ordinary words. The only step Cliewen refuses to do for you is the last one.

## Next

[Start with what Cliewen is and why it exists.](./what-is-cliewen)
