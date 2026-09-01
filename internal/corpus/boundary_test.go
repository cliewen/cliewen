package corpus

import (
	"strings"
	"testing"
)

// boundaryFiles adds a decisions folder holding one record that declares it
// binds adopter behaviour but names no shipped carrier.
var boundaryFiles = map[string]string{
	"docs/README.md":                   "# Corpus\n\n<!-- clue:index:start -->\n- [goals/](goals/README.md)\n- [plans/](plans/README.md)\n- [decisions/](decisions/README.md)\n- [architecture/](architecture/README.md)\n- [design/](design/README.md)\n<!-- clue:index:end -->\n",
	"docs/architecture/README.md":      "# Architecture\n\nStructure.\n\n<!-- clue:index:start -->\n<!-- clue:index:end -->\n",
	"docs/design/README.md":            "# Design\n\nCross-cutting behaviour.\n\n<!-- clue:index:start -->\n<!-- clue:index:end -->\n",
	"docs/decisions/README.md":         "# Decisions\n\n<!-- clue:index:start -->\n- [ADR-001](ADR-001-a-rule.md)\n<!-- clue:index:end -->\n",
	"docs/decisions/ADR-001-a-rule.md": "---\nid: ADR-001\ntype: decision\nstatus: inferred\nlinks: []\nbinds: adopter\ntitle: A rule\nauthor: agent\naccepted-by: []\n---\n\n# ADR-001\n\nAdopters must do the thing.\n",
}

const shippedCarrierBody = "---\nid: ADR-001\ntype: decision\nstatus: inferred\nlinks: []\nbinds: adopter\ntitle: A rule\nauthor: agent\naccepted-by: []\n---\n\n# ADR-001\n\nAdopters must do the thing, stated in `internal/skills/source/shared/local-conventions.md.tmpl`.\n"

func sourceRepo(extra map[string]string) map[string]string {
	return with(validFiles, with(boundaryFiles, with(map[string]string{".clue/role.yaml": "role: source\n"}, extra)))
}

// In the source repository, a rule that binds adopters must name a carrier
// on the surface adopters actually receive.
func TestAC154_UnitNegative_AdopterRuleWithoutAShippedCarrierFails(t *testing.T) {
	assertIssue(t, run(t, sourceRepo(nil), false), "binds adopter but names no carrier on the shipped surface")
}

func TestAC154_UnitPositive_AdopterRuleNamingAShippedCarrierPasses(t *testing.T) {
	files := sourceRepo(map[string]string{"docs/decisions/ADR-001-a-rule.md": shippedCarrierBody})
	if issues := run(t, files, false); len(issues) != 0 {
		t.Fatalf("a rule naming a canonical skill source should pass, got %v", issues)
	}
}

// A template path is equally a shipped surface: clue init materializes it
// into the adopter's repository.
func TestAC154_UnitPositive_AScaffoldedTemplateIsAShippedCarrier(t *testing.T) {
	body := "---\nid: ADR-001\ntype: decision\nstatus: inferred\nlinks: []\nbinds: adopter\ntitle: A rule\nauthor: agent\naccepted-by: []\n---\n\n# ADR-001\n\nCarried by `internal/scaffold/templates/AGENTS.md`.\n"
	files := sourceRepo(map[string]string{"docs/decisions/ADR-001-a-rule.md": body})
	if issues := run(t, files, false); len(issues) != 0 {
		t.Fatalf("a rule naming a scaffolded template should pass, got %v", issues)
	}
}

// The rule never judges a repository that ships nothing. An adopter's own
// decisions have no shipped surface to name, and an undeclared repository
// is an adopter.
func TestAC154_UnitPositive_AnAdopterAndAnUndeclaredRepoAreNotJudged(t *testing.T) {
	adopter := with(validFiles, with(boundaryFiles, map[string]string{".clue/role.yaml": "role: adopter\n"}))
	undeclared := with(validFiles, boundaryFiles)
	for name, files := range map[string]map[string]string{
		"adopter":    adopter,
		"undeclared": undeclared,
	} {
		t.Run(name, func(t *testing.T) {
			for _, issue := range run(t, files, false) {
				if strings.Contains(issue.String(), "shipped surface") {
					t.Fatalf("%s corpus was judged by the source-only rule: %v", name, issue)
				}
			}
		})
	}
}

// The role is the only switch on the carrier rule, so a marker that cannot
// be read must say so. Defaulting a malformed marker to "not source" would
// let a one-character typo disable the rule with no output anywhere.
func TestAC153_UnitNegative_AnUnreadableRoleMarkerIsReportedNotDefaulted(t *testing.T) {
	for name, marker := range map[string]string{
		"unknown role": "role: sorce\n",
		"empty role":   "role:\n",
		"not yaml":     "role: [source\n",
	} {
		files := with(validFiles, with(boundaryFiles, map[string]string{".clue/role.yaml": marker}))
		t.Run(name, func(t *testing.T) {
			assertIssue(t, run(t, files, false), "role marker cannot be read")
		})
	}
}

func TestAC154_UnitNegative_AnUnknownOrMisplacedBindsFails(t *testing.T) {
	unknown := sourceRepo(map[string]string{
		"docs/decisions/ADR-001-a-rule.md": "---\nid: ADR-001\ntype: decision\nstatus: inferred\nlinks: []\nbinds: everyone\ntitle: A rule\nauthor: agent\naccepted-by: []\n---\n\n# ADR-001\n",
	})
	assertIssue(t, run(t, unknown, false), "binds must be adopter or repo")

	misplaced := sourceRepo(map[string]string{
		"docs/plans/P-001-baseline.md": "---\nid: P-001\ntype: plan\nstatus: active\nlinks: [G-001]\nbinds: adopter\ntitle: Baseline\n---\n\n| M-001 | do it | todo |\n",
	})
	assertIssue(t, run(t, misplaced, false), "binds is allowed only on decision and constraint artifacts")
}
