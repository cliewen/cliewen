package corpus

import (
	"strings"
	"testing"
)

// constraintCorpus is validFiles plus a constraints folder carrying one
// constraint file, so a register rule can be exercised on its own.
func constraintCorpus(file string) map[string]string {
	return with(validFiles, map[string]string{
		"docs/README.md":                 "# Corpus\n\n<!-- clue:index:start -->\n- [goals/](goals/README.md)\n- [plans/](plans/README.md)\n- [constraints/](constraints/README.md)\n<!-- clue:index:end -->\n",
		"docs/constraints/README.md":     "# Constraints\n\n<!-- clue:index:start -->\n- [C-001](C-001-rule.md)\n<!-- clue:index:end -->\n",
		"docs/constraints/C-001-rule.md": file,
	})
}

func constraintFile(enforcement, body string) string {
	return "---\nid: C-001\ntype: constraint\nstatus: active\nlinks: []\ntitle: A rule\nsource: AGENTS.md rule 5\nenforcement: " + enforcement + "\n---\n\n# C-001\n" + body
}

func assertClean(t *testing.T, issues []Issue, what string) {
	t.Helper()
	if len(issues) != 0 {
		t.Fatalf("%s: expected no issues, got %v", what, issues)
	}
}

// AC-089: the register carries source, a class from the widened vocabulary,
// and the declarations ADR-045 requires of partial and human.
func TestAC089_UnitPositive_RegisterClassesAndDeclarationsAccepted(t *testing.T) {
	assertClean(t, run(t, constraintCorpus(constraintFile("agent", "\n**Promotion trigger:** a lint.\n")), false), "an agent-enforced constraint with a promotion trigger")
	assertClean(t, run(t, constraintCorpus(constraintFile("machine", "\nHeld by the linter.\n")), false), "a machine-enforced constraint owes no declaration")
	assertClean(t, run(t, constraintCorpus(constraintFile("partial", "\n**Checked by:** `clue validate`.\n\n**Residual:** whether the reason is real.\n")), false), "a partial constraint declaring both")
	assertClean(t, run(t, constraintCorpus(constraintFile("human", "\n**Residual:** meaning, and it costs a missed weakening.\n")), false), "a human constraint declaring its residual")
}

func TestAC089_UnitNegative_RegisterFieldsAndDeclarationsRejected(t *testing.T) {
	noSource := "---\nid: C-001\ntype: constraint\nstatus: active\nlinks: []\ntitle: A rule\nenforcement: agent\n---\n"
	assertIssue(t, run(t, constraintCorpus(noSource), false), "constraint missing or empty source field")

	noEnforcement := "---\nid: C-001\ntype: constraint\nstatus: active\nlinks: []\ntitle: A rule\nsource: AGENTS.md rule 5\n---\n"
	assertIssue(t, run(t, constraintCorpus(noEnforcement), false), "constraint missing or empty enforcement field")

	assertIssue(t, run(t, constraintCorpus(constraintFile("hope", "")), false), "enforcement hope not allowed (allowed: machine, partial, agent, human)")

	// partial owes both declarations: a claimed machine nobody named, and a
	// residual nobody priced, are the two ways this class could be abused.
	assertIssue(t, run(t, constraintCorpus(constraintFile("partial", "\n**Residual:** judgment.\n")), false), "enforcement partial needs a **Checked by:** declaration")
	assertIssue(t, run(t, constraintCorpus(constraintFile("partial", "\n**Checked by:** `clue validate`.\n")), false), "enforcement partial needs a **Residual:** declaration")
	assertIssue(t, run(t, constraintCorpus(constraintFile("human", "\nNo machine can hold this.\n")), false), "enforcement human needs a **Residual:** declaration")
}

// AC-090: a paragraph or list item broken across lines fails; structural line
// breaks do not.
func TestAC090_UnitPositive_StructuralLineBreaksArePermitted(t *testing.T) {
	page := `---
id: G-002
type: goal
status: accepted
links: []
title: Structure only
---

# A heading

One paragraph on one line, however long it runs and whatever punctuation it carries.

- A list item on one line
- Another list item
  - A nested item

> A blockquote line
> A second blockquote line

| Column | Other |
|---|---|
| Cell | Cell |

<!-- an HTML comment -->

` + "```go" + `
a := 1
b := 2
` + "```" + `

[ref]: https://example.com

Final paragraph.
`
	assertClean(t, run(t, with(validFiles, map[string]string{
		"docs/goals/README.md":          "# Goals\n\n<!-- clue:index:start -->\n- [G-001](G-001-first.md)\n- [G-002](G-002-structure.md)\n<!-- clue:index:end -->\n",
		"docs/goals/G-002-structure.md": page,
	}), false), "a page whose only line breaks are structural")
}

// The blocks a page about markdown actually contains. Each one was a false
// positive before CH-110's review: an indented code block is not fence-
// delimited, a table may be written without outer pipes, an HTML block runs
// past its opening line, and a fence documenting a shorter fence must not be
// closed by its own example.
func TestAC090_UnitPositive_VerbatimBlocksAreNotProse(t *testing.T) {
	page := `---
id: G-002
type: goal
status: accepted
links: []
title: Blocks
---

# A heading

An indented code block, which no fence delimits:

    line one of code
    line two of code

A table written without outer pipes:

Column | Other
--- | ---
Cell | Cell

An HTML block that runs past its opening line:

<div class="note">
  <p>first</p>
  <p>second</p>
</div>

A fence holding a shorter fence as an example:

` + "````" + `markdown
This text is inside the outer fence.
` + "```" + `
And so is this.
` + "```" + `
Still inside.
` + "````" + `

Final paragraph.
`
	assertClean(t, run(t, with(validFiles, map[string]string{
		"docs/goals/README.md":       "# Goals\n\n<!-- clue:index:start -->\n- [G-001](G-001-first.md)\n- [G-002](G-002-blocks.md)\n<!-- clue:index:end -->\n",
		"docs/goals/G-002-blocks.md": page,
	}), false), "verbatim and structural blocks")
}

// The block scanner's dangerous direction is silence: a block state that never
// closes stops every check that uses it, and nothing reports a check that did
// not run. Each case here is a shape where that happened.
func TestAC090_UnitNegative_BlockStatesDoNotSwallowProse(t *testing.T) {
	for _, tc := range []struct {
		name string
		body string
	}{
		{
			// The shape every corpus has: an index block whose rows sit
			// directly under the marker, with no blank line to close it.
			name: "under an index-block comment",
			body: "<!-- clue:index:start -->\n- [G-001](G-001-first.md) — a row wrapped\n  onto a second line\n<!-- clue:index:end -->\n",
		},
		{
			name: "after a closed one-line comment",
			body: "<!-- a note -->\nThis paragraph was broken\nacross two lines.\n",
		},
		{
			name: "after an HTML block closed by a blank line",
			body: "<div>\n  <p>markup</p>\n</div>\n\nThis paragraph was broken\nacross two lines.\n",
		},
		{
			// An autolink is prose. Reading it as an HTML block would exempt
			// the paragraph it opens.
			name: "in a paragraph opening with an autolink",
			body: "<https://example.com> is the site\nand this line continues the sentence.\n",
		},
		{
			// Four spaces under a list item is that item's continuation, not a
			// code block — CommonMark needs eight there, and this is prose.
			name: "in a list item's indented continuation",
			body: "- A list item\n\n    A continuation paragraph\n    wrapped across two lines.\n",
		},
		{
			// An indented fence marker is code showing a fence, not a fence
			// opening one. Read as an opener it starts a block nothing closes,
			// and every line-based check goes quiet for the rest of the file.
			name: "after an indented example containing a fence marker",
			body: "How not to write it:\n\n    ```mermaid\n\nThis paragraph was broken\nacross two lines.\n",
		},
		{
			// An indented block ends at the first unindented line, not at the
			// next blank one, so the line right after an example is read.
			name: "on the line immediately after an indented block",
			body: "An example:\n\n    code line\nThis paragraph was broken\nacross two lines.\n",
		},
		{
			// A commented-out table must not leave a table state open past the
			// comment that contains it.
			name: "after a commented-out table",
			body: "<!-- an old table\n| M | Status |\n| --- | --- |\n-->\nThis paragraph was broken\nacross two lines.\n",
		},
		{
			// `<T>` is a type parameter in prose, not markup. Reading it as an
			// HTML block would exempt the paragraph it opens.
			name: "in a paragraph opening with a generic type",
			body: "<T> is the type parameter, and this\nline continues the sentence.\n",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			page := "---\nid: G-002\ntype: goal\nstatus: accepted\nlinks: []\ntitle: Blocks\n---\n\n# A heading\n\n" + tc.body
			assertIssue(t, run(t, with(validFiles, map[string]string{
				"docs/goals/README.md":       "# Goals\n\n<!-- clue:index:start -->\n- [G-001](G-001-first.md)\n- [G-002](G-002-blocks.md)\n<!-- clue:index:end -->\n",
				"docs/goals/G-002-blocks.md": page,
			}), false), "paragraph continues on a new line")
		})
	}
}

func TestAC090_UnitNegative_HardWrappedProseAndListItemsRejected(t *testing.T) {
	wrapped := "---\nid: G-002\ntype: goal\nstatus: accepted\nlinks: []\ntitle: Wrapped\n---\n\n# A heading\n\nThis paragraph was broken\nacross two lines.\n"
	issues := run(t, with(validFiles, map[string]string{
		"docs/goals/README.md":        "# Goals\n\n<!-- clue:index:start -->\n- [G-001](G-001-first.md)\n- [G-002](G-002-wrapped.md)\n<!-- clue:index:end -->\n",
		"docs/goals/G-002-wrapped.md": wrapped,
	}), false)
	assertIssue(t, issues, "line 12: paragraph continues on a new line")

	wrappedItem := "---\nid: G-002\ntype: goal\nstatus: accepted\nlinks: []\ntitle: Wrapped\n---\n\n# A heading\n\n- A list item broken\n  across two lines\n"
	assertIssue(t, run(t, with(validFiles, map[string]string{
		"docs/goals/README.md":        "# Goals\n\n<!-- clue:index:start -->\n- [G-001](G-001-first.md)\n- [G-002](G-002-wrapped.md)\n<!-- clue:index:end -->\n",
		"docs/goals/G-002-wrapped.md": wrappedItem,
	}), false), "paragraph continues on a new line")
}

// changeCorpus is validFiles plus one transient workspace.
func changeCorpus(files map[string]string) map[string]string {
	return with(validFiles, files)
}

const validProposal = "---\nid: CH-001-proposal\ntype: change\nstatus: open\nlinks: [P-001]\ntitle: A change\n---\n\n# CH-001\n\nWhat and why.\n"

// AC-091: a `[-]` task carries prose after its checkbox.
func TestAC091_UnitPositive_SkippedTaskWithAReasonPasses(t *testing.T) {
	tasks := "---\nid: CH-001-tasks\ntype: tasks\nstatus: open\nlinks: []\ntitle: Tasks\n---\n\n# Tasks\n\n- [x] Done\n- [ ] Not yet\n- [-] Not feasible: the upstream API never shipped\n"
	assertClean(t, run(t, changeCorpus(map[string]string{
		"changes/CH-001-a/proposal.md": validProposal,
		"changes/CH-001-a/tasks.md":    tasks,
	}), false), "a skipped task carrying its reason")
}

// An indented code block is a fence written without one. Reading its contents
// as live markdown makes an example of what not to write fail as the thing
// itself — a false failure on a rule C-004 forbids weakening, which leaves the
// author no exit but to delete legitimate documentation.
func TestAC091_UnitPositive_IndentedExamplesAreNotTasks(t *testing.T) {
	tasks := "---\nid: CH-001-tasks\ntype: tasks\nstatus: open\nlinks: []\ntitle: Tasks\n---\n\n# Tasks\n\nWhat a skipped task must never look like:\n\n    - [-]\n\n- [x] Done\n"
	assertClean(t, run(t, changeCorpus(map[string]string{
		"changes/CH-001-a/proposal.md": validProposal,
		"changes/CH-001-a/tasks.md":    tasks,
	}), false), "a reasonless task shown as an indented example")
}

func TestAC091_UnitNegative_SkippedTaskWithoutAReasonRejected(t *testing.T) {
	tasks := "---\nid: CH-001-tasks\ntype: tasks\nstatus: open\nlinks: []\ntitle: Tasks\n---\n\n# Tasks\n\n- [x] Done\n- [-]\n"
	assertIssue(t, run(t, changeCorpus(map[string]string{
		"changes/CH-001-a/proposal.md": validProposal,
		"changes/CH-001-a/tasks.md":    tasks,
	}), false), "line 12: skipped task carries no reason")
}

// AC-092: a proposal names its plan item or declares itself plan-less.
func TestAC092_UnitPositive_PlanItemOrPlanLessDeclarationAccepted(t *testing.T) {
	assertClean(t, run(t, changeCorpus(map[string]string{
		"changes/CH-001-a/proposal.md": validProposal,
	}), false), "a proposal linking a plan")

	planLess := "---\nid: CH-001-proposal\ntype: change\nstatus: open\nlinks: []\ntitle: A change\n---\n\n# CH-001\n\nThis change is plan-less: it repairs a broken link.\n"
	assertClean(t, run(t, changeCorpus(map[string]string{
		"changes/CH-001-a/proposal.md": planLess,
	}), false), "a proposal declaring itself plan-less")
}

func TestAC092_UnitNegative_ProposalWithoutADeclarationRejected(t *testing.T) {
	silent := "---\nid: CH-001-proposal\ntype: change\nstatus: open\nlinks: [G-001]\ntitle: A change\n---\n\n# CH-001\n\nWhat and why, but nothing about a plan.\n"
	assertIssue(t, run(t, changeCorpus(map[string]string{
		"changes/CH-001-a/proposal.md": silent,
	}), false), "proposal names no plan item")
}

// AC-093: diagrams are inline Mermaid, never images.
func TestAC093_UnitPositive_MermaidAndOrdinaryLinksAccepted(t *testing.T) {
	page := "---\nid: G-002\ntype: goal\nstatus: accepted\nlinks: []\ntitle: Diagrams\n---\n\n# G-002\n\nA link to [the first goal](G-001-first.md) and a code span showing the form `![alt](x.png)`.\n\n```mermaid\ngraph TD\n  A --> B\n```\n\nA fenced example:\n\n```markdown\n![alt](diagram.png)\n```\n\nAnd the same example indented rather than fenced:\n\n    ![alt](diagram.png)\n"
	assertClean(t, run(t, with(validFiles, map[string]string{
		"docs/goals/README.md":         "# Goals\n\n<!-- clue:index:start -->\n- [G-001](G-001-first.md)\n- [G-002](G-002-diagrams.md)\n<!-- clue:index:end -->\n",
		"docs/goals/G-002-diagrams.md": page,
	}), false), "inline Mermaid, an ordinary link, and images shown as examples")
}

func TestAC093_UnitNegative_ImageLinksAndImageFilesRejected(t *testing.T) {
	page := "---\nid: G-002\ntype: goal\nstatus: accepted\nlinks: []\ntitle: Diagrams\n---\n\n# G-002\n\n![the architecture](architecture.png)\n"
	issues := run(t, with(validFiles, map[string]string{
		"docs/goals/README.md":         "# Goals\n\n<!-- clue:index:start -->\n- [G-001](G-001-first.md)\n- [G-002](G-002-diagrams.md)\n<!-- clue:index:end -->\n",
		"docs/goals/G-002-diagrams.md": page,
		"docs/goals/architecture.png":  "\x89PNG not really",
	}), false)
	assertIssue(t, issues, "line 11: image link — diagrams in the corpus are inline Mermaid")
	assertIssue(t, issues, "docs/goals/architecture.png: image file under docs/")
}

// AC-094: each type carries its own frontmatter extensions.
func TestAC094_UnitPositive_TypeExtensionsPresentAndEmptySignatureAccepted(t *testing.T) {
	decision := "---\nid: ADR-001\ntype: decision\nstatus: inferred\nlinks: []\ntitle: A decision\nauthor: agent\naccepted-by: []\n---\n\n# ADR-001\n"
	assertClean(t, run(t, with(validFiles, map[string]string{
		"docs/README.md":              "# Corpus\n\n<!-- clue:index:start -->\n- [goals/](goals/README.md)\n- [plans/](plans/README.md)\n- [decisions/](decisions/README.md)\n<!-- clue:index:end -->\n",
		"docs/decisions/README.md":    "# Decisions\n\n<!-- clue:index:start -->\n- [ADR-001](ADR-001-a.md)\n<!-- clue:index:end -->\n",
		"docs/decisions/ADR-001-a.md": decision,
	}), false), "a decision carrying author and an empty accepted-by")
}

func TestAC094_UnitNegative_MissingTypeExtensionsRejected(t *testing.T) {
	noAuthor := "---\nid: ADR-001\ntype: decision\nstatus: inferred\nlinks: []\ntitle: A decision\n---\n\n# ADR-001\n"
	assertIssue(t, run(t, with(validFiles, map[string]string{
		"docs/README.md":              "# Corpus\n\n<!-- clue:index:start -->\n- [goals/](goals/README.md)\n- [plans/](plans/README.md)\n- [decisions/](decisions/README.md)\n<!-- clue:index:end -->\n",
		"docs/decisions/README.md":    "# Decisions\n\n<!-- clue:index:start -->\n- [ADR-001](ADR-001-a.md)\n<!-- clue:index:end -->\n",
		"docs/decisions/ADR-001-a.md": noAuthor,
	}), false), "decision missing or empty field(s): accepted-by, author")

	noGoal := "---\nid: CAP-001\ntype: capability\nstatus: draft\nlinks: [G-001]\ntitle: A capability\n---\n\n# CAP-001\n"
	assertIssue(t, run(t, with(validFiles, map[string]string{
		"docs/README.md":                        "# Corpus\n\n<!-- clue:index:start -->\n- [goals/](goals/README.md)\n- [plans/](plans/README.md)\n- [capabilities/](capabilities/README.md)\n<!-- clue:index:end -->\n",
		"docs/capabilities/README.md":           "# Capabilities\n\n<!-- clue:index:start -->\n- [CAP-001](CAP-001-a/README.md)\n<!-- clue:index:end -->\n",
		"docs/capabilities/CAP-001-a/README.md": noGoal,
	}), false), "capability missing or empty field(s): goal")
}

// AC-095: milestone status cells follow one vocabulary.
func TestAC095_UnitPositive_DeclaredVocabularyAndUnstatusedTablesAccepted(t *testing.T) {
	plan := "---\nid: P-002\ntype: plan\nstatus: active\nlinks: [G-001]\ntitle: A plan\n---\n\n# P-002\n\n| ID | Milestone | Status | Evidence |\n|---|---|---|---|\n"
	for _, v := range MilestoneStatuses {
		body := plan + "| M-002 | do it | `" + v + "` | |\n"
		assertClean(t, run(t, planCorpus(body), false), "milestone status "+v)
	}

	// A table written without outer pipes is the same table. The prose lint
	// already reads one, and a status column the two checks disagreed about
	// would be a hole in exactly the shape the other check taught a writer to use.
	bare := "---\nid: P-002\ntype: plan\nstatus: active\nlinks: [G-001]\ntitle: A plan\n---\n\n# P-002\n\nID | Milestone | Status\n--- | --- | ---\nM-002 | do it | `done`\n"
	assertClean(t, run(t, planCorpus(bare), false), "a milestone table written without outer pipes")

	// A prose cell carrying an escaped pipe and a code span with a pipe in it.
	// Splitting on every pipe would shift every later column, reading a
	// neighbouring cell as the status — a false failure, and a real one missed.
	pipes := "---\nid: P-002\ntype: plan\nstatus: active\nlinks: [G-001]\ntitle: A plan\n---\n\n# P-002\n\n| ID | Milestone | Status | Evidence |\n|---|---|---|---|\n| M-002 | supports `a \\| b` and ``x | y`` forms \\| and more | `done` | none |\n"
	assertClean(t, run(t, planCorpus(pipes), false), "a milestone row whose prose carries escaped and code-span pipes")

	// A plan may show the wrong vocabulary as an example, the way every other
	// check allows one — indented, fenced, or commented out.
	example := `---
id: P-002
type: plan
status: active
links: [G-001]
title: A plan
---

# P-002

Do not write it like this:

    | ID | Status |
    |---|---|
    | M-009 | in progress |

Or like this:

` + "```markdown" + `
| ID | Status |
|---|---|
| M-009 | wip |
` + "```" + `
`
	assertClean(t, run(t, planCorpus(example), false), "a milestone table shown as an indented or fenced example")

	// A table with no status column is not a milestone table, whatever it
	// happens to contain.
	other := "---\nid: P-002\ntype: plan\nstatus: active\nlinks: [G-001]\ntitle: A plan\n---\n\n# P-002\n\n| ID | Note |\n|---|---|\n| M-002 | shipped |\n"
	assertClean(t, run(t, planCorpus(other), false), "a plan table declaring no status column")
}

func TestAC095_UnitNegative_StatusOutsideTheVocabularyRejected(t *testing.T) {
	plan := "---\nid: P-002\ntype: plan\nstatus: active\nlinks: [G-001]\ntitle: A plan\n---\n\n# P-002\n\n| ID | Milestone | Status | Evidence |\n|---|---|---|---|\n| M-002 | do it | `wip` | |\n"
	issues := run(t, planCorpus(plan), false)
	assertIssue(t, issues, "M-002: milestone status wip is not one of todo, doing, done, dropped")

	// The same bad value in a table written without outer pipes: read as a
	// table by the prose lint, so it is read as one here too.
	bare := "---\nid: P-002\ntype: plan\nstatus: active\nlinks: [G-001]\ntitle: A plan\n---\n\n# P-002\n\nID | Milestone | Status\n--- | --- | ---\nM-002 | do it | `wip`\n"
	assertIssue(t, run(t, planCorpus(bare), false), "M-002: milestone status wip is not one of")
	for _, i := range issues {
		if strings.Contains(i.Msg, "milestone status Status") {
			t.Fatalf("the header row is not a value: %v", issues)
		}
	}
}

func planCorpus(plan string) map[string]string {
	return with(validFiles, map[string]string{
		"docs/plans/README.md":       "# Plans\n\n<!-- clue:index:start -->\n- [P-001](P-001-baseline.md)\n- [P-002](P-002-second.md)\n<!-- clue:index:end -->\n",
		"docs/plans/P-002-second.md": plan,
	})
}
