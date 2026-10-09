package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cliewen/cliewen/internal/evidenceexport"
)

func TestSanity_AcceptCLIExactCandidateHandoff(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	git := func(args ...string) string {
		t.Helper()
		c := exec.Command("git", args...)
		out, err := c.CombinedOutput()
		if err != nil {
			t.Fatalf("git %v: %v %s", args, err, out)
		}
		return strings.TrimSpace(string(out))
	}
	write := func(path, text string) {
		t.Helper()
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(text), 0644); err != nil {
			t.Fatal(err)
		}
	}
	git("init", "-b", "main")
	git("config", "user.name", "Test human")
	git("config", "user.email", "test@example.invalid")
	git("config", "commit.gpgsign", "false")
	write(".clue/acceptance.yaml", "mode: local\nbranch: main\n")
	write("docs/README.md", "---\ntype: index\ntitle: Corpus\n---\n\n<!-- clue:index:start -->\n<!-- clue:index:end -->\n")
	if err := evidenceexport.Fixture(root); err != nil {
		t.Fatal(err)
	}
	git("add", ".")
	git("commit", "-m", "Base")
	base := git("rev-parse", "HEAD")
	git("switch", "-c", "change")
	write("changes/CH-001-example/proposal.md", "---\nid: CH-001\ntype: change\n---\n\nProposal\n")
	git("add", ".")
	git("commit", "-m", "Proposal")
	git("rm", "-r", "changes")
	git("commit", "-m", "Digest")
	candidate := git("rev-parse", "HEAD")
	git("switch", "main")
	brief := ".git/brief.md"
	text := fmt.Sprintf("---\ntype: acceptance-brief\ntitle: Example\nchange: CH-001\ncandidate: %s\nbase: %s\nreviewed: %s\n---\n", candidate, base, candidate)
	for _, section := range []string{"Intent", "Criteria and evidence", "Binding decisions", "Verification", "Review"} {
		text += "\n## " + section + "\n\nCompleted.\n"
	}
	write(brief, text)
	args := []string{candidate, "--base", base, "--brief", brief}
	var out, errs bytes.Buffer
	if code := runAccept(append(append([]string{}, args...), "--check"), strings.NewReader(""), &out, &errs, false); code != 0 {
		t.Fatalf("preflight %d: %s", code, errs.String())
	}
	if !strings.Contains(out.String(), "no acceptance performed") || git("rev-parse", "HEAD") != base {
		t.Fatal("preflight integrated")
	}
	if code := runAccept(args, strings.NewReader("accept CH-001\n"), &out, &errs, false); code != 1 {
		t.Fatal("piped acceptance allowed")
	}
	if code := runAccept(args, strings.NewReader("accept CH-001\n"), &out, &errs, true); code != 0 {
		t.Fatalf("acceptance %d: %s", code, errs.String())
	}
	if git("show", "-s", "--format=%P", "HEAD") != base+" "+candidate {
		t.Fatal("CLI accepted wrong parents")
	}
}

func TestUnit_AcceptCLIRejectsInvalidArguments(t *testing.T) {
	for _, args := range [][]string{nil, {"--yes"}, {"bad", "--base", "bad", "--brief", "missing"}, {"--base", "abc", "--brief", "missing", "abc", "extra"}} {
		var out, errs bytes.Buffer
		if runAccept(args, strings.NewReader(""), &out, &errs, false) == 0 {
			t.Fatal("invalid arguments passed")
		}
	}
	_ = interactiveInput()
}
