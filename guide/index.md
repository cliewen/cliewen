---
layout: home

hero:
  name: Cliewen
  text: Evidence-backed Intent Engineering for coding agents
  tagline: Coding agents can produce plausible changes faster than a person can confidently accept them. Cliewen keeps the requirement, implementation, and acceptance evidence connected in Git, then checks that the connection is intact before merge.
  image:
    src: /cliewen-logo.svg
    alt: Cliewen logo
  actions:
    - theme: brand
      text: Start reading
      link: /what-is-cliewen
    - theme: alt
      text: The skills
      link: /skills
    - theme: alt
      text: The corpus
      link: /corpus

features:
  - title: Your agent follows the skills
    details: clue init installs six skills and an AGENTS.md routing hub. Any coding agent that reads them knows how to plan, make a change, verify it, and upgrade, so your prompts can stay in ordinary words.
    link: /skills
  - title: Your docs/ folder is the memory
    details: The corpus under docs/ says what the product is for, what it can do, and what proves it. Agents update it with the code, so the next session starts from a record instead of a lost chat.
    link: /corpus
  - title: One thread the repository can check
    details: A vision and goals lead to capabilities, acceptance criteria, and declared evidence, either classified test references or genuine Human proof. The clue CLI reports broken links and missing evidence locally and in CI without executing tests.
    link: /methodology
  - title: You decide what merges
    details: Agents prepare the change, the corpus, and the evidence. Humans keep control of intent and of the merge. Small work that changes no promise stays outside the full loop.
    link: /change-loop
---

## The thread

```mermaid
graph LR
  V["Vision"] --> G["Goal"]
  G --> C["Capability"]
  C --> A["Acceptance criterion"]
  A --> E["Acceptance evidence"]
  E -. "why do we have this?" .-> G
```

Pick up any artifact and follow it back to why it exists, or forward to what proves it. `clue` checks that no arrow is missing. It does not run your tests, it cannot tell you that a test proves the right thing, and it cannot tell you the goal was worth having. That judgment stays yours, at the merge.

## What you work with

- **The skills** in `.agents/skills/`: the process your coding agent follows.
- **The corpus** under `docs/`: the lasting record of the system as it exists.
- **`clue`**: the command-line judge that checks the corpus and its evidence.
- **Cliewen**: the name of the method that ties them together.

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
