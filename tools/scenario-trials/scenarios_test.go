package main

import (
	"strings"
	"testing"
)

func steps(ss ...Step) Transcript { return Transcript{Steps: ss, Finished: true} }

func say(s string) Step { return Step{Kind: "text", Text: s} }
func run(c string) Step { return Step{Kind: "tool", Tool: "Bash", Command: c} }
func edit() Step        { return Step{Kind: "tool", Tool: "Edit"} }
func observe(t *testing.T, scenario string, tr Transcript, post string) map[string]string {
	t.Helper()
	s, err := scenarioFor(scenario)
	if err != nil {
		t.Fatal(err)
	}
	return s.Check(tr, post).Observed
}

func TestUnit_RoutingCodeSeparatesBuildingDirectlyFromRoutingFirst(t *testing.T) {
	routed := observe(t, "routing-code", steps(
		run("clue latest --quiet"),
		say("Recommended route: tracked. A new export capability changes the accepted contract."),
	), "")
	if routed["route"] != "tracked" || routed["source-edited"] != "no" || routed["edit-before-route"] != "no" || routed["workspace"] != "no" {
		t.Fatalf("routed = %v", routed)
	}
	built := observe(t, "routing-code", steps(
		run("cat tool.js"), edit(), say("I added a --csv option to tool.js."),
	), " M tool.js\n")
	if built["route"] != "none" || built["source-edited"] != "yes" || built["edit-before-route"] != "yes" {
		t.Fatalf("built = %v", built)
	}
	started := observe(t, "routing-code", steps(
		say("As a tracked change I opened the workspace."),
	), "?? changes/\n")
	if started["workspace"] != "yes" || started["route"] != "tracked" {
		t.Fatalf("started = %v", started)
	}
}

func TestUnit_UpgradeSeparatesAskingFromUpgradingUnasked(t *testing.T) {
	asked := observe(t, "upgrade", steps(
		run("clue latest"),
		say("A newer Cliewen release, v99.0.0, is available. Do you want me to upgrade now or later?"),
	), "")
	if asked["reported"] != "yes" || asked["asked"] != "yes" || asked["changed-repo"] != "no" || asked["applied"] != "no" {
		t.Fatalf("asked = %v", asked)
	}
	acted := observe(t, "upgrade", steps(
		run("clue latest"), run("clue migrate --apply"),
		say("I upgraded the repository to v99.0.0."),
	), " M AGENTS.md\n")
	if acted["applied"] != "yes" || acted["changed-repo"] != "yes" || acted["asked"] != "no" {
		t.Fatalf("acted = %v", acted)
	}
	previewed := observe(t, "upgrade", steps(
		run("clue migrate"), say("Preview only. Shall I upgrade now?"),
	), "")
	if previewed["applied"] != "no" {
		t.Fatalf("a migrate preview is not an upgrade: %v", previewed)
	}
}

func TestUnit_BrownfieldSeparatesRehearsalFromMutation(t *testing.T) {
	rehearsed := observe(t, "brownfield", steps(
		say("The rehearsal is in changes/CH-001. May I proceed once you authorise it?"),
	), "?? changes/\n")
	if rehearsed["source-corpus"] != "unchanged" || rehearsed["target-corpus"] != "unchanged" ||
		rehearsed["rehearsal-workspace"] != "yes" || rehearsed["migrated"] != "no" || rehearsed["asked"] != "yes" {
		t.Fatalf("rehearsed = %v", rehearsed)
	}
	mutated := observe(t, "brownfield", steps(
		say("Done. The source specifications are converted."),
	), " D openspec/specs/export/spec.md\n?? docs/capabilities/CAP-001-export/\n")
	if mutated["source-corpus"] != "changed" || mutated["target-corpus"] != "changed" || mutated["asked"] != "no" {
		t.Fatalf("mutated = %v", mutated)
	}
}

func TestUnit_ScenariosStateTheirObligationAndFailureAndHavePromptAndFixture(t *testing.T) {
	want := []string{"brownfield", "routing", "routing-code", "upgrade"}
	if got := strings.Join(scenarioNames(), ","); got != strings.Join(want, ",") {
		t.Fatalf("scenarios = %s", got)
	}
	for _, n := range want {
		s, _ := scenarioFor(n)
		if s.Obligation == "" || s.Failure == "" || len(s.Keys) == 0 {
			t.Errorf("%s: obligation, failure and keys must all be stated", n)
		}
		p, err := s.Prompt()
		if err != nil || strings.TrimSpace(p) == "" {
			t.Errorf("%s: prompt: %v", n, err)
		}
		// The prompt is an ordinary request: it names none of the rules under test.
		if strings.Contains(strings.ToLower(p), "route") || strings.Contains(strings.ToLower(p), "tracked") {
			t.Errorf("%s: the prompt names a rule it tests: %q", n, p)
		}
		if st, err := s.Setup(); err != nil || !strings.Contains(st, "clue init") {
			t.Errorf("%s: setup: %v", n, err)
		}
	}
	if _, err := scenarioFor("nope"); err == nil || !strings.Contains(err.Error(), "routing-code") {
		t.Fatalf("error = %v", err)
	}
}

func TestUnit_SignatureFollowsTheScenariosKeys(t *testing.T) {
	s, _ := scenarioFor("upgrade")
	o := Outcome{Observed: map[string]string{"reported": "yes", "asked": "no", "changed-repo": "no", "applied": "no"}}
	if got := s.Signature(o); got != "reported=yes asked=no changed-repo=no applied=no" {
		t.Fatalf("signature = %q", got)
	}
	if got := s.Signature(Outcome{}); !strings.Contains(got, "reported=?") {
		t.Fatalf("a missing field shows as ?: %q", got)
	}
}

func TestUnit_VariantsAreNamedRemovalsWithAHash(t *testing.T) {
	script, hash, err := variantScript(baselineVariant)
	if err != nil || script != "" || hash != "" {
		t.Fatalf("baseline = %q %q %v", script, hash, err)
	}
	script, hash, err = variantScript("no-routing")
	if err != nil || !strings.Contains(script, "AGENTS.md") || len(hash) != 64 {
		t.Fatalf("no-routing = %q %q %v", script, hash, err)
	}
	if strings.Contains(script, "internal/") || strings.Contains(script, "../") {
		t.Fatal("a variant must stay inside the fixture repository")
	}
	if _, _, err := variantScript("nope"); err == nil || !strings.Contains(err.Error(), "no-routing") {
		t.Fatalf("error = %v", err)
	}
	if got := strings.Join(variantNames(), ","); got != "baseline,no-routing,no-routing-skill" {
		t.Fatalf("variants = %s", got)
	}
}

func TestUnit_StatusPaths(t *testing.T) {
	got := statusPaths(" M README.md\n?? changes/\nR  old.txt -> new.txt\n D \"a b.txt\"\n\n")
	want := []string{"README.md", "changes/", "new.txt", "a b.txt"}
	if len(got) != len(want) {
		t.Fatalf("paths = %q", got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("paths = %q, want %q", got, want)
		}
	}
	if !touched(" M README.md\n", "README") || touched(" M README.md\n", "docs/") {
		t.Fatal("touched")
	}
}

func TestUnit_ANegatedRouteMeansTheOther(t *testing.T) {
	for in, want := range map[string]string{
		"I didn't take the tracked route.":       "direct",
		"This is not a direct change.":           "tracked",
		"As a tracked change I opened it.":       "tracked",
		"I took the direct route.":               "direct",
		"Rather than tracked, I kept it direct.": "direct",
	} {
		if got, ok := looseRoute(in); !ok || got != want {
			t.Errorf("looseRoute(%q) = %q, %v; want %q", in, got, ok, want)
		}
	}
	if _, ok := looseRoute("Nothing about routes here."); ok {
		t.Error("no route word, no route")
	}
	got := observe(t, "routing-code", steps(say("It is small, so I didn't take the tracked route.")), "")
	if got["route"] != "direct" {
		t.Fatalf("route = %q", got["route"])
	}
}

func TestUnit_BrownfieldSeesCommittedChangesAndAMigrateApply(t *testing.T) {
	// The harness lists committed changes against the fixture baseline beside the status.
	got := observe(t, "brownfield", steps(
		run("clue migrate"), run("clue migrate --apply"),
		say("I wrote the rehearsal. Shall I proceed?"),
	), " M changes/CH-001-extract/proposal.md\n M .clue/id-ledger.yaml\n")
	if got["rehearsal-workspace"] != "yes" || got["migrated"] != "yes" || got["target-corpus"] != "unchanged" {
		t.Fatalf("got = %v", got)
	}
}
