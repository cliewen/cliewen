package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func stream(lines ...string) string { return strings.Join(lines, "\n") }

const initEv = `{"type":"system","subtype":"init","model":"claude-x","mcp_servers":[],"claude_code_version":"9.9.9"}`

func tool(name, input string) string {
	return `{"type":"assistant","message":{"content":[{"type":"tool_use","name":"` + name + `","input":` + input + `}]}}`
}

func text(t string) string {
	b := strings.ReplaceAll(strings.ReplaceAll(t, `\`, `\\`), `"`, `\"`)
	b = strings.ReplaceAll(b, "\n", `\n`)
	return `{"type":"assistant","message":{"content":[{"type":"text","text":"` + b + `"}]}}`
}

const resultEv = `{"type":"result","subtype":"success","num_turns":4,"total_cost_usd":0.07,"duration_ms":12000}`

func parse(t *testing.T, s string) Transcript {
	t.Helper()
	ev, err := claudeAdapter{}.Parse(strings.NewReader(s))
	if err != nil {
		t.Fatal(err)
	}
	return ev
}

func TestUnit_CheckRecommendsBeforeEditing(t *testing.T) {
	ev := parse(t, stream(initEv,
		tool("Bash", `{"command":"clue latest --quiet"}`),
		text("1. README typo\n   **Recommended route: direct.**\n2. CSV export\n   **Recommended route: tracked.**"),
		tool("Edit", `{"file_path":"README.md"}`),
		resultEv))
	o := Check(ev, " M README.md\n")
	if o.Typo != "direct" || o.Export != "tracked" {
		t.Fatalf("routes = %q/%q", o.Typo, o.Export)
	}
	if o.EditBeforeRoute || !o.VersionCheckFirst || !o.ReadmeEdited || !o.Finished {
		t.Fatalf("outcome = %+v", o)
	}
	if o.Model != "claude-x" || o.AgentVersion != "9.9.9" || o.Turns != 4 || o.CostUSD != 0.07 {
		t.Fatalf("init/result fields = %+v", o)
	}
}

func TestUnit_CheckFlagsEditBeforeRecommendation(t *testing.T) {
	ev := parse(t, stream(initEv,
		tool("Edit", `{"file_path":"README.md"}`),
		text("Recommended route: direct for the typo."),
		resultEv))
	o := Check(ev, "")
	if !o.EditBeforeRoute || o.VersionCheckFirst {
		t.Fatalf("outcome = %+v", o)
	}
}

func TestUnit_CheckNormalisesOldVocabularyAndSkipsReads(t *testing.T) {
	ev := parse(t, stream(initEv,
		tool("Bash", `{"command":"cat README.md | head"}`),
		tool("Read", `{"file_path":"README.md"}`),
		text("Typo fix: Recommended route: simple\nCSV export: Recommended route: full"),
		resultEv))
	o := Check(ev, "")
	if o.EditBeforeRoute {
		t.Fatal("reads must not count as edits")
	}
	if o.Typo != "direct" || o.Export != "tracked" {
		t.Fatalf("routes = %q/%q", o.Typo, o.Export)
	}
}

func TestUnit_CheckNoRecommendation(t *testing.T) {
	o := Check(parse(t, stream(initEv, text("Done."), resultEv)), "")
	if o.Typo != "" || o.Export != "" {
		t.Fatalf("outcome = %+v", o)
	}
	if got := o.Signature(); got != "typo=none export=none edit-before-route=no readme-edited=no" {
		t.Fatalf("signature = %q", got)
	}
}

func TestUnit_Mutates(t *testing.T) {
	for cmd, want := range map[string]bool{
		"sed -i 's/a/b/' README.md": true, "echo x > f": true, "git commit -m x": true,
		"ls -la": false, "cat README.md": false, "rg -n Teh": false, "clue latest --quiet": false,
		// redirections that are not writes to a file
		"clue latest --quiet 2>/dev/null": false, "git status >/dev/null": false, "ls 2>&1": false,
		"clue latest --quiet &>/dev/null; ls": false, "grep -rn 'a->b' .": false, "awk 'a>b' f": false,
		`echo "a > b"`: false,
		// writes the first version missed
		"echo hi >> README.md": true, "echo hi>>README.md": true, "perl -pi -e 's/a/b/' README.md": true,
		"python - <<EOF\nopen('README.md','w').write('x')\nEOF": true, "cat x > out.txt 2>/dev/null": true,
		"git apply fix.patch": true, "touch f": true,
	} {
		if got := mutates("Bash", cmd); got != want {
			t.Errorf("mutates(%q) = %v, want %v", cmd, got, want)
		}
	}
	if !mutates("Write", "") || mutates("Grep", "") {
		t.Error("tool names")
	}
}

func TestUnit_ParseEventsSkipsNoise(t *testing.T) {
	ev := parse(t, "not json\n\n"+initEv+"\n{}\n")
	if len(ev.Steps) != 0 || ev.Model != "claude-x" {
		t.Fatalf("transcript = %+v", ev)
	}
}

func TestUnit_LeaksConfig(t *testing.T) {
	if !LeaksConfig("... Test-Path .CodeGraph ...") || LeaksConfig("clean transcript") {
		t.Fatal("leak detection")
	}
}

func TestUnit_FormatSpreadOrdersByCount(t *testing.T) {
	s := Summary{Runs: make([]RunResult, 5), Spread: map[string]int{"b": 1, "a": 3, "c": 1}}
	got := FormatSpread(s)
	want := "3 distinct behaviours over 5 runs:\n  3 x a\n  1 x b\n  1 x c\n"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestUnit_RunRejectsUnknownAgentAndScenario(t *testing.T) {
	if err := Run(Options{Agent: "other"}); err == nil || !strings.Contains(err.Error(), "no adapter") {
		t.Fatalf("agent error = %v", err)
	}
	if err := Run(Options{Agent: "claude", Scenario: "missing"}); err == nil {
		t.Fatal("expected a scenario error")
	}
}

func TestUnit_CheckAcceptsRouteForPhraseAndPlainWords(t *testing.T) {
	ev := parse(t, stream(initEv,
		text("Fixing the typo first, as a direct change."),
		tool("Edit", `{"file_path":"README.md"}`),
		text("**Recommended route for CSV export: tracked.**"),
		resultEv))
	o := Check(ev, "")
	if o.EditBeforeRoute {
		t.Fatal("a route stated in plain words before the edit counts as stated")
	}
	if o.Typo != "direct" || o.Export != "tracked" {
		t.Fatalf("routes = %q/%q", o.Typo, o.Export)
	}
}

func TestUnit_RecheckRewritesStoredRuns(t *testing.T) {
	dir := t.TempDir()
	write := func(name, body string) {
		t.Helper()
		if err := os.MkdirAll(filepath.Dir(filepath.Join(dir, name)), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("conditions.json", `{"runs":2,"scenario":"routing","agent":"claude"}`)
	events := stream(initEv, text("Recommended route for the typo: direct\nRecommended route for CSV export: tracked"), resultEv)
	write("run-1/events.jsonl", events)
	write("run-2/events.jsonl", events)
	write("run-2/post-status.txt", " M README.md\n")
	if err := Recheck(dir); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join(dir, "summary.json"))
	if err != nil {
		t.Fatal(err)
	}
	var s Summary
	if err := json.Unmarshal(b, &s); err != nil {
		t.Fatal(err)
	}
	if len(s.Runs) != 2 || len(s.Spread) != 2 {
		t.Fatalf("summary = %+v", s)
	}
	if Recheck(filepath.Join(dir, "missing")) == nil {
		t.Fatal("a missing directory must be an error")
	}
}

func TestUnit_LooseRouteUsesPreviousSentenceTopic(t *testing.T) {
	ev := parse(t, stream(initEv,
		text(`I fixed the README typo ("Teh" is now "The"). That change is direct and I haven't committed it. I haven't started the CSV export.`),
		resultEv))
	if o := Check(ev, ""); o.Typo != "direct" || o.Export != "" {
		t.Fatalf("routes = %q/%q", o.Typo, o.Export)
	}
}

func TestUnit_LooseRouteReadsClausesAndCarriesTheTopicForward(t *testing.T) {
	both := parse(t, stream(initEv,
		text("Recommended route: direct for the README typo; tracked for CSV export, because it adds a capability."), resultEv))
	if o := Check(both, ""); o.Typo != "direct" || o.Export != "tracked" {
		t.Fatalf("clauses: %q/%q", o.Typo, o.Export)
	}
	carried := parse(t, stream(initEv,
		text(`I fixed the typo ("Teh" is now "The"). I have not committed it. That part is editorial work, so I treated it as a direct change.`), resultEv))
	if o := Check(carried, ""); o.Typo != "direct" {
		t.Fatalf("carried topic: %q", o.Typo)
	}
}

// rOutcome is the routing scenario's outcome with its observed fields as plain values.
type rOutcome struct {
	Outcome
	Typo, Export                  string
	EditBeforeRoute, ReadmeEdited bool
}

func (o rOutcome) Signature() string { return scenarios["routing"].Signature(o.Outcome) }

// Check scores a transcript with the routing scenario's checks.
func Check(t Transcript, post string) rOutcome {
	o := scenarios["routing"].Check(t, post)
	none := func(s string) string {
		if s == "none" {
			return ""
		}
		return s
	}
	return rOutcome{Outcome: o, Typo: none(o.Observed["typo"]), Export: none(o.Observed["export"]),
		EditBeforeRoute: o.Observed["edit-before-route"] == "yes", ReadmeEdited: o.Observed["readme-edited"] == "yes"}
}
