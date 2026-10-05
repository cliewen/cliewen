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

func parse(t *testing.T, s string) []Event {
	t.Helper()
	ev, err := ParseEvents(strings.NewReader(s))
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
	if o.TypoRoute != "direct" || o.ExportRoute != "tracked" {
		t.Fatalf("routes = %q/%q", o.TypoRoute, o.ExportRoute)
	}
	if o.EditBeforeRoute || !o.VersionCheckFirst || !o.ReadmeEdited || !o.Recommended || !o.Finished {
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
	if o.TypoRoute != "direct" || o.ExportRoute != "tracked" {
		t.Fatalf("routes = %q/%q", o.TypoRoute, o.ExportRoute)
	}
}

func TestUnit_CheckNoRecommendation(t *testing.T) {
	o := Check(parse(t, stream(initEv, text("Done."), resultEv)), "")
	if o.Recommended || o.TypoRoute != "" || o.ExportRoute != "" {
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
	if len(ev) != 1 || ev[0].Type != "system" {
		t.Fatalf("events = %+v", ev)
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
	if o.TypoRoute != "direct" || o.ExportRoute != "tracked" {
		t.Fatalf("routes = %q/%q", o.TypoRoute, o.ExportRoute)
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
	write("conditions.json", `{"runs":2,"scenario":"routing"}`)
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
	if o := Check(ev, ""); o.TypoRoute != "direct" || o.ExportRoute != "" {
		t.Fatalf("routes = %q/%q", o.TypoRoute, o.ExportRoute)
	}
}
