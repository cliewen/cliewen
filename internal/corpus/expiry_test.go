package corpus

import "testing"

// analysisFiles extends the baseline corpus with an analysis folder and a
// completed plan, the shape every carried-by and expiry case builds on.
var analysisFiles = map[string]string{
	"docs/README.md":                      "# Corpus\n\n<!-- clue:index:start -->\n- [goals/](goals/README.md)\n- [plans/](plans/README.md)\n- [analysis/](analysis/README.md)\n- [architecture/](architecture/README.md)\n- [design/](design/README.md)\n<!-- clue:index:end -->\n",
	"docs/design/README.md":               "# Design\n\nCross-cutting behaviour.\n\n<!-- clue:index:start -->\n<!-- clue:index:end -->\n",
	"docs/plans/README.md":                "# Plans\n\n<!-- clue:index:start -->\n- [P-001](P-001-baseline.md)\n<!-- clue:index:end -->\n",
	"docs/plans/P-001-baseline.md":        "---\nid: P-001\ntype: plan\nstatus: completed\nlinks: [G-001]\ntitle: Baseline\n---\n\n| M-001 | do it | done |\n",
	"docs/architecture/README.md":         "# Architecture\n\nStructure.\n\n<!-- clue:index:start -->\n- [ARCH-001](ARCH-001-shape.md)\n<!-- clue:index:end -->\n",
	"docs/architecture/ARCH-001-shape.md": "---\nid: ARCH-001\ntype: architecture\nstatus: active\nlinks: []\ntitle: Shape\n---\n\n# ARCH-001\n",
	"docs/analysis/README.md":             "# Analysis\n\n<!-- clue:index:start -->\n- [AN-001](AN-001-spike.md)\n<!-- clue:index:end -->\n",
	"docs/analysis/AN-001-spike.md":       "---\nid: AN-001\ntype: analysis\nstatus: active\nlinks: [P-001]\ncarried-by: [ARCH-001]\ntitle: Spike\n---\n\n# AN-001\n",
}

func TestAC155_UnitPositive_AnHonestCarriedByPasses(t *testing.T) {
	if issues := run(t, with(validFiles, analysisFiles), false); len(issues) != 0 {
		t.Fatalf("an analysis naming a durable artifact that resolves should pass, got %v", issues)
	}
}

// An analysis that declares nothing is unaffected. Requiring the field on
// every existing spike would turn a corpus green yesterday red today.
func TestAC155_UnitPositive_AnAnalysisDeclaringNothingPasses(t *testing.T) {
	files := with(validFiles, with(analysisFiles, map[string]string{
		"docs/analysis/AN-001-spike.md": "---\nid: AN-001\ntype: analysis\nstatus: active\nlinks: [P-001]\ntitle: Spike\n---\n\n# AN-001\n",
	}))
	if issues := run(t, files, false); len(issues) != 0 {
		t.Fatalf("an analysis without carried-by should pass, got %v", issues)
	}
}

func TestAC155_UnitNegative_ADishonestCarriedByFails(t *testing.T) {
	cases := map[string]struct{ file, want string }{
		"unresolvable id": {
			"---\nid: AN-001\ntype: analysis\nstatus: active\nlinks: [P-001]\ncarried-by: [ARCH-404]\ntitle: Spike\n---\n\n# AN-001\n",
			"carried-by names ARCH-404, which no artifact declares",
		},
		"empty list": {
			"---\nid: AN-001\ntype: analysis\nstatus: active\nlinks: [P-001]\ncarried-by: []\ntitle: Spike\n---\n\n# AN-001\n",
			"carried-by must name at least one durable artifact",
		},
		"self reference": {
			"---\nid: AN-001\ntype: analysis\nstatus: active\nlinks: [P-001]\ncarried-by: [AN-001]\ntitle: Spike\n---\n\n# AN-001\n",
			"carried-by names the analysis itself",
		},
	}
	for name, tc := range cases {
		files := with(validFiles, with(analysisFiles, map[string]string{"docs/analysis/AN-001-spike.md": tc.file}))
		t.Run(name, func(t *testing.T) { assertIssue(t, run(t, files, false), tc.want) })
	}
}

// The field belongs to analysis. On anything else it would claim a
// lifecycle that artifact does not have.
func TestAC155_UnitNegative_CarriedByOnANonAnalysisFails(t *testing.T) {
	files := with(validFiles, with(analysisFiles, map[string]string{
		"docs/architecture/ARCH-001-shape.md": "---\nid: ARCH-001\ntype: architecture\nstatus: active\nlinks: []\ncarried-by: [P-001]\ntitle: Shape\n---\n\n# ARCH-001\n",
	}))
	assertIssue(t, run(t, files, false), "carried-by is allowed only on analysis artifacts")
}

// One spike carrying another's findings is not durability, it is a
// forwarding address.
func TestAC155_UnitNegative_CarriedByAnotherAnalysisFails(t *testing.T) {
	files := with(validFiles, with(analysisFiles, map[string]string{
		"docs/analysis/README.md":       "# Analysis\n\n<!-- clue:index:start -->\n- [AN-001](AN-001-spike.md)\n- [AN-002](AN-002-other.md)\n<!-- clue:index:end -->\n",
		"docs/analysis/AN-002-other.md": "---\nid: AN-002\ntype: analysis\nstatus: active\nlinks: [P-001]\ntitle: Other\n---\n\n# AN-002\n",
		"docs/analysis/AN-001-spike.md": "---\nid: AN-001\ntype: analysis\nstatus: active\nlinks: [P-001]\ncarried-by: [AN-002]\ntitle: Spike\n---\n\n# AN-001\n",
	}))
	assertIssue(t, run(t, files, false), "which is another analysis")
}

func TestAC156_UnitPositive_ASpentAnalysisIsDerived(t *testing.T) {
	c, issues := Scan(writeCorpus(t, with(validFiles, analysisFiles)))
	if len(issues) != 0 {
		t.Fatalf("scan issues: %v", issues)
	}
	spent := SpentAnalyses(c)
	if len(spent) != 1 {
		t.Fatalf("expected one spent analysis, got %v", spent)
	}
	if spent[0].ID != "AN-001" || len(spent[0].Plans) != 1 || spent[0].Plans[0] != "P-001" {
		t.Fatalf("unexpected derivation: %+v", spent[0])
	}
	if len(spent[0].CarriedBy) != 1 || spent[0].CarriedBy[0] != "ARCH-001" {
		t.Fatalf("unexpected carriers: %+v", spent[0])
	}
}

// Both conditions are required. A completed plan says the campaign ended;
// it says nothing about whether the findings ever reached durable form.
func TestAC156_UnitNegative_OneConditionAloneIsNotSpent(t *testing.T) {
	activePlan := with(validFiles, with(analysisFiles, map[string]string{
		"docs/plans/P-001-baseline.md": "---\nid: P-001\ntype: plan\nstatus: active\nlinks: [G-001]\ntitle: Baseline\n---\n\n| M-001 | do it | todo |\n",
	}))
	noCarrier := with(validFiles, with(analysisFiles, map[string]string{
		"docs/analysis/AN-001-spike.md": "---\nid: AN-001\ntype: analysis\nstatus: active\nlinks: [P-001]\ntitle: Spike\n---\n\n# AN-001\n",
	}))
	noPlan := with(validFiles, with(analysisFiles, map[string]string{
		"docs/analysis/AN-001-spike.md": "---\nid: AN-001\ntype: analysis\nstatus: active\nlinks: [ARCH-001]\ncarried-by: [ARCH-001]\ntitle: Spike\n---\n\n# AN-001\n",
	}))
	for name, files := range map[string]map[string]string{
		"plan still active":   activePlan,
		"no declared carrier": noCarrier,
		"serves no plan":      noPlan,
	} {
		t.Run(name, func(t *testing.T) {
			c, _ := Scan(writeCorpus(t, files))
			if spent := SpentAnalyses(c); len(spent) != 0 {
				t.Fatalf("expected nothing spent, got %+v", spent)
			}
		})
	}
}

// citingCorpus is the spent baseline plus a decisions folder and a
// constraints folder, so a citation can be pointed at AN-001 one artifact
// type at a time.
func citingCorpus(t *testing.T, extra map[string]string) *Corpus {
	t.Helper()
	files := with(validFiles, with(analysisFiles, map[string]string{
		"docs/README.md":             "# Corpus\n\n<!-- clue:index:start -->\n- [goals/](goals/README.md)\n- [plans/](plans/README.md)\n- [analysis/](analysis/README.md)\n- [architecture/](architecture/README.md)\n- [design/](design/README.md)\n- [decisions/](decisions/README.md)\n- [constraints/](constraints/README.md)\n<!-- clue:index:end -->\n",
		"docs/decisions/README.md":   "# Decisions\n\n<!-- clue:index:start -->\n- [ADR-001](ADR-001-a-rule.md)\n<!-- clue:index:end -->\n",
		"docs/constraints/README.md": "# Constraints\n\n<!-- clue:index:start -->\n- [C-001](C-001-a-rule.md)\n<!-- clue:index:end -->\n",
	}))
	c, issues := Scan(writeCorpus(t, with(files, extra)))
	if len(issues) != 0 {
		t.Fatalf("scan issues: %v", issues)
	}
	return c
}

func adr(links string) string {
	return "---\nid: ADR-001\ntype: decision\nstatus: inferred\nlinks: [" + links + "]\ntitle: A rule\nauthor: agent\naccepted-by: []\n---\n\n# ADR-001\n"
}

func constraint(links string) string {
	return "---\nid: C-001\ntype: constraint\nstatus: active\nlinks: [" + links + "]\ntitle: A rule\nsource: G-001\nenforcement: agent\n---\n\n# C-001\n\nThe rule holds.\n"
}

// A decision or constraint citing a spike is neither invalid nor a bar to
// anything else; it only keeps the spike out of the spent report.
func TestAC158_UnitPositive_AnUncitedSpikeIsStillReported(t *testing.T) {
	c := citingCorpus(t, map[string]string{
		"docs/decisions/ADR-001-a-rule.md": adr("P-001"),
		"docs/constraints/C-001-a-rule.md": constraint("P-001"),
		"docs/analysis/README.md":          "# Analysis\n\n<!-- clue:index:start -->\n- [AN-001](AN-001-spike.md)\n- [AN-002](AN-002-other.md)\n<!-- clue:index:end -->\n",
		"docs/analysis/AN-002-other.md":    "---\nid: AN-002\ntype: analysis\nstatus: active\nlinks: [P-001, AN-001]\ntitle: Other\n---\n\n# AN-002\n",
	})
	spent := SpentAnalyses(c)
	if len(spent) != 1 || spent[0].ID != "AN-001" {
		t.Fatalf("a spike named only by a completed plan and another analysis should be reported, got %+v", spent)
	}
	if issues := Validate(c, Options{}); len(issues) != 0 {
		t.Fatalf("a cited spike is not an invalid corpus, got %v", issues)
	}
}

// Both standing rule types gate. Each is read by someone deciding whether
// the rule still holds, which cannot be done once its evidence is gone.
func TestAC158_UnitNegative_ACitedSpikeIsWithheld(t *testing.T) {
	cases := map[string]map[string]string{
		"cited by a decision": {
			"docs/decisions/ADR-001-a-rule.md": adr("P-001, AN-001"),
			"docs/constraints/C-001-a-rule.md": constraint("P-001"),
		},
		"cited by a constraint": {
			"docs/decisions/ADR-001-a-rule.md": adr("P-001"),
			"docs/constraints/C-001-a-rule.md": constraint("P-001, AN-001"),
		},
	}
	for name, extra := range cases {
		t.Run(name, func(t *testing.T) {
			c := citingCorpus(t, extra)
			if spent := SpentAnalyses(c); len(spent) != 0 {
				t.Fatalf("expected nothing spent, got %+v", spent)
			}
		})
	}
}
