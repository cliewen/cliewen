---
layout: home

hero:
  name: Cliewen
  text: Ship agent-written changes without losing the intent
  tagline: Cliewen keeps requirements, decisions, implementation, and acceptance evidence connected in Git—and catches missing evidence before merge.
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

## Three moves

```mermaid
graph LR
  I["Install clue"] --> P["Prompt your agent"]
  P --> A["It writes the code, the corpus, and the evidence"]
  A --> V["clue validate and CI check the thread"]
  V --> M["You merge"]
```

Install takes one command. Prompting takes ordinary words. The only step Cliewen refuses to do for you is the last one.

## Next

[Start with what Cliewen is and why it exists.](./what-is-cliewen)
