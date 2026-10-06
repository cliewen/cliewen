package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func codexStream() string {
	return stream(
		`{"type":"thread.started","thread_id":"t"}`,
		`{"type":"item.completed","item":{"type":"agent_message","text":"Checking the version first."}}`,
		`{"type":"item.completed","item":{"type":"command_execution","command":"/bin/bash -lc 'clue latest --quiet 2>/dev/null'"}}`,
		`{"type":"item.completed","item":{"type":"file_change"}}`,
		`{"type":"item.completed","item":{"type":"agent_message","text":"- Typo: **Recommended route: direct.**\n- CSV export: **Recommended route: tracked.**"}}`,
		`{"type":"turn.completed"}`)
}

func TestUnit_CodexParseGivesTheNeutralTranscript(t *testing.T) {
	tr, err := codexAdapter{}.Parse(strings.NewReader(codexStream()))
	if err != nil {
		t.Fatal(err)
	}
	if !tr.Finished || len(tr.Steps) != 4 {
		t.Fatalf("transcript = %+v", tr)
	}
	if tr.Steps[1].Tool != "Bash" || tr.Steps[1].Command != "clue latest --quiet 2>/dev/null" {
		t.Fatalf("command step = %+v", tr.Steps[1])
	}
	if tr.Steps[2].Tool != "Edit" {
		t.Fatalf("file change step = %+v", tr.Steps[2])
	}
}

func TestUnit_ChecksAgreeAcrossAdaptersOnTheSameBehaviour(t *testing.T) {
	// The same behaviour, version check then an edit then the routes, in each
	// agent's own event stream, must give the same signature.
	claude := parse(t, stream(initEv,
		tool("Bash", `{"command":"clue latest --quiet 2>/dev/null"}`),
		tool("Edit", `{"file_path":"README.md"}`),
		text("- Typo: **Recommended route: direct.**\n- CSV export: **Recommended route: tracked.**"),
		resultEv))
	codex, err := codexAdapter{}.Parse(strings.NewReader(codexStream()))
	if err != nil {
		t.Fatal(err)
	}
	a, b := Check(claude, " M README.md\n"), Check(codex, " M README.md\n")
	if a.Signature() != b.Signature() || !a.VersionCheckFirst || !b.VersionCheckFirst {
		t.Fatalf("claude %q vs codex %q", a.Signature(), b.Signature())
	}
	if a.Signature() != "typo=direct export=tracked edit-before-route=yes readme-edited=yes" {
		t.Fatalf("signature = %q", a.Signature())
	}
}

func TestUnit_UnwrapShell(t *testing.T) {
	for in, want := range map[string]string{
		`/bin/bash -lc 'git status'`: "git status",
		`/bin/bash -lc "echo hi"`:    "echo hi",
		`bash -c 'ls -la'`:           "ls -la",
		`sh -lc cat README.md`:       "cat README.md",
		`python3 script.py`:          "python3 script.py",
		`/bin/bash -lc 'echo "a b"'`: `echo "a b"`,
	} {
		if got := unwrapShell(in); got != want {
			t.Errorf("unwrapShell(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestUnit_CredentialsComeFromOutsideTheRepositoryAndStayOutOfArguments(t *testing.T) {
	dir := t.TempDir()
	tok := filepath.Join(dir, "token")
	if err := os.WriteFile(tok, []byte(" secret-value\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	c, err := claudeAdapter{}.Credentials(Options{Login: tok})
	if err != nil {
		t.Fatal(err)
	}
	if len(c.Pass) != 1 || c.Pass[0] != "CLAUDE_CODE_OAUTH_TOKEN" || c.Env[0] != "CLAUDE_CODE_OAUTH_TOKEN=secret-value" {
		t.Fatalf("credentials = %+v", c)
	}
	for _, m := range c.Mounts {
		if strings.Contains(m, "secret-value") {
			t.Fatal("secret in a mount")
		}
	}
	if _, err := (claudeAdapter{}).Credentials(Options{Login: filepath.Join(dir, "missing")}); err == nil {
		t.Fatal("a missing token file must be an error")
	}
	if _, err := (codexAdapter{}).Credentials(Options{Login: dir}); err == nil {
		t.Fatal("a Codex home without auth.json must be an error")
	}
	if err := os.WriteFile(filepath.Join(dir, "auth.json"), []byte("{}"), 0o600); err != nil {
		t.Fatal(err)
	}
	cx, err := codexAdapter{}.Credentials(Options{Login: dir})
	if err != nil || len(cx.Env) != 0 || len(cx.Mounts) != 1 || !strings.HasSuffix(cx.Mounts[0], "target=/home/node/.codex") {
		t.Fatalf("codex credentials = %+v, %v", cx, err)
	}
}

func TestUnit_AdapterRegistry(t *testing.T) {
	for _, n := range []string{"claude", "codex", "opencode"} {
		a, err := adapterFor(n)
		if err != nil || a.Name() != n {
			t.Fatalf("adapter %q: %v", n, err)
		}
		if !strings.Contains(a.Command("m1"), "m1") || a.Probe() == "" {
			t.Errorf("%s command or probe is empty", n)
		}
	}
	if _, err := adapterFor("nope"); err == nil || !strings.Contains(err.Error(), "claude, codex, opencode") {
		t.Fatalf("error = %v", err)
	}
}

func TestUnit_UnwrapShellKeepsMultiLineBodies(t *testing.T) {
	wrapped := "/bin/bash -lc 'sed -i s/Teh/The/ README.md\ngit diff'"
	if got := unwrapShell(wrapped); got != "sed -i s/Teh/The/ README.md\ngit diff" {
		t.Fatalf("unwrapShell = %q", got)
	}
	tr := Transcript{Steps: []Step{
		{Kind: "tool", Tool: "Bash", Command: unwrapShell(wrapped)},
		{Kind: "text", Text: "Recommended route: direct for the typo."},
	}}
	if !Check(tr, "").EditBeforeRoute {
		t.Fatal("a multi-line command that starts with sed -i is an edit")
	}
}

func TestUnit_CodexProbeReadsTheNewestSessionAndLoginPathsWithCommasAreRefused(t *testing.T) {
	if p := (codexAdapter{}).Probe(); !strings.Contains(p, "ls -t") || strings.Contains(p, "grep -rhao") {
		t.Fatalf("probe must pick the newest session file, got %q", p)
	}
	dir := filepath.Join(t.TempDir(), "a,b")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "auth.json"), []byte("{}"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := (codexAdapter{}).Credentials(Options{Login: dir}); err == nil {
		t.Fatal("a comma in the login path must be refused")
	}
}

func TestUnit_OpencodeParseGivesTheNeutralTranscript(t *testing.T) {
	stream := `{"type":"step_start","timestamp":1000,"part":{"type":"step-start"}}
{"type":"tool_use","timestamp":2000,"part":{"type":"tool","tool":"bash","state":{"input":{"command":"clue latest --quiet"}}}}
{"type":"tool_use","timestamp":2500,"part":{"type":"tool","tool":"edit","state":{"input":{}}}}
{"type":"step_finish","timestamp":2600,"part":{"type":"step-finish","reason":"tool-calls","cost":0.5}}
{"type":"text","timestamp":3000,"part":{"type":"text","text":"Recommended route: direct"}}
{"type":"step_finish","timestamp":3100,"part":{"type":"step-finish","reason":"stop","cost":0.25}}
`
	tr, err := opencodeAdapter{}.Parse(strings.NewReader(stream))
	if err != nil {
		t.Fatal(err)
	}
	if len(tr.Steps) != 3 || tr.Steps[0].Command != "clue latest --quiet" || tr.Steps[1].Tool != "Edit" || tr.Steps[2].Text != "Recommended route: direct" {
		t.Fatalf("steps = %+v", tr.Steps)
	}
	if !tr.Finished || tr.Turns != 2 || tr.CostUSD != 0.75 || tr.DurationMS != 2100 {
		t.Fatalf("transcript = %+v", tr)
	}
}

func TestUnit_OpencodeCredentialsMountTheDataDirectory(t *testing.T) {
	dir := t.TempDir()
	if _, err := (opencodeAdapter{}).Credentials(Options{Login: dir}); err == nil {
		t.Fatal("a data directory without auth.json must be an error")
	}
	if err := os.WriteFile(filepath.Join(dir, "auth.json"), []byte("{}"), 0o600); err != nil {
		t.Fatal(err)
	}
	c, err := opencodeAdapter{}.Credentials(Options{Login: dir})
	if err != nil || len(c.Env) != 0 || len(c.Mounts) != 1 || !strings.HasSuffix(c.Mounts[0], "target=/home/node/.local/share/opencode") {
		t.Fatalf("credentials = %+v, %v", c, err)
	}
}
