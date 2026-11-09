package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/cliewen/cliewen/internal/evidenceexport"
)

// acceptRepo builds a repository with an accepted base and a digested
// candidate whose history introduces one change proposal, and returns the
// candidate, the base, and helpers bound to the repository.
func acceptRepo(t *testing.T) (candidate, base string, git func(...string) string, write func(string, string)) {
	t.Helper()
	root := t.TempDir()
	t.Chdir(root)
	git = func(args ...string) string {
		t.Helper()
		c := exec.Command("git", args...)
		out, err := c.CombinedOutput()
		if err != nil {
			t.Fatalf("git %v: %v %s", args, err, out)
		}
		return strings.TrimSpace(string(out))
	}
	write = func(path, text string) {
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
	base = git("rev-parse", "HEAD")
	git("switch", "-c", "change")
	write("changes/CH-001-example/proposal.md", "---\nid: CH-001\ntype: change\ntitle: Example\n---\n\nProposal\n")
	git("add", ".")
	git("commit", "-m", "Proposal")
	git("rm", "-r", "changes")
	git("commit", "-m", "Digest")
	candidate = git("rev-parse", "HEAD")
	git("switch", "main")
	return candidate, base, git, write
}

func TestSanity_AcceptCLIExactCandidateHandoff(t *testing.T) {
	candidate, base, git, _ := acceptRepo(t)
	args := []string{candidate, "--base", base}
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
	for _, args := range [][]string{nil, {"--yes"}, {"bad", "--base", "bad"}, {"--base", "abc", "abc", "extra"}} {
		var out, errs bytes.Buffer
		if runAccept(args, strings.NewReader(""), &out, &errs, false) == 0 {
			t.Fatal("invalid arguments passed")
		}
	}
	_ = interactiveInput()
}

func TestAC242_IntegrationNegative_OldBriefOptionIsRefusedPlainly(t *testing.T) {
	candidate, base, git, write := acceptRepo(t)
	write(".git/brief.md", "---\ntype: acceptance-brief\n---\n")
	var out, errs bytes.Buffer
	if code := runAccept([]string{candidate, "--base", base, "--brief", ".git/brief.md"}, strings.NewReader("accept CH-001\n"), &out, &errs, true); code == 0 {
		t.Fatal("old brief option accepted")
	}
	if !strings.Contains(errs.String(), "--brief is no longer used") {
		t.Fatalf("refusal does not say why: %s", errs.String())
	}
	if git("rev-parse", "HEAD") != base {
		t.Fatal("refused acceptance integrated")
	}
}

func TestAC243_IntegrationPositive_AcceptCLIExecutesRecordedHumanDecision(t *testing.T) {
	candidate, base, git, write := acceptRepo(t)
	approval := fmt.Sprintf("type: acceptance-approval\ndecision: accept\ncandidate: %s\nbase: %s\napproved-by: Fixture observer\napproval-source: Synthetic conversation\nstatement: accept this fixture\nrecorded-at: %s\n", candidate, base, time.Now().UTC().Format(time.RFC3339))
	write(".git/approval.yaml", approval)
	approvedArgs := []string{candidate, "--base", base, "--approval", ".git/approval.yaml"}
	var out, errs bytes.Buffer
	if code := runAccept(append(append([]string{}, approvedArgs...), "--check"), strings.NewReader(""), &out, &errs, false); code != 0 || git("rev-parse", "HEAD") != base {
		t.Fatalf("record preflight %d: %s", code, errs.String())
	}
	if code := runAccept(approvedArgs, strings.NewReader(""), &out, &errs, false); code != 0 {
		t.Fatalf("recorded execution %d: %s", code, errs.String())
	}
	if git("show", "-s", "--format=%P", "HEAD") != base+" "+candidate {
		t.Fatal("CLI accepted wrong parents")
	}
}
