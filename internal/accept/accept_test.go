package accept

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cliewen/cliewen/internal/evidenceexport"
)

func write(t *testing.T, root, path, body string) {
	t.Helper()
	name := filepath.Join(root, path)
	if err := os.MkdirAll(filepath.Dir(name), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(name, []byte(body), 0644); err != nil {
		t.Fatal(err)
	}
}

func mustGit(t *testing.T, root string, args ...string) string {
	t.Helper()
	out, err := git(root, args...)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func fixture(t *testing.T, modify func(string)) Request {
	t.Helper()
	root := t.TempDir()
	mustGit(t, root, "init", "-b", "main")
	mustGit(t, root, "config", "user.name", "Acceptance tester")
	mustGit(t, root, "config", "user.email", "accept@example.invalid")
	mustGit(t, root, "config", "commit.gpgsign", "false")
	mustGit(t, root, "config", "core.autocrlf", "false")
	write(t, root, ConfigPath, "mode: local\nbranch: main\n")
	write(t, root, "docs/README.md", "---\ntype: index\ntitle: Corpus\n---\n\n# Corpus\n\n<!-- clue:index:start -->\n<!-- clue:index:end -->\n")
	write(t, root, "app.txt", "before\n")
	if err := evidenceexport.Fixture(root); err != nil {
		t.Fatal(err)
	}
	mustGit(t, root, "add", ".")
	mustGit(t, root, "commit", "-m", "Accepted base")
	base := mustGit(t, root, "rev-parse", "HEAD")
	mustGit(t, root, "switch", "-c", "change")
	write(t, root, "changes/CH-001-example/proposal.md", "---\nid: CH-001\ntype: change\nstatus: open\nlinks: []\ntitle: Example\n---\n\nA plan-less change.\n")
	mustGit(t, root, "add", ".")
	mustGit(t, root, "commit", "-m", "Propose CH-001")
	mustGit(t, root, "rm", "-r", "changes")
	write(t, root, "app.txt", "after\n")
	if modify != nil {
		modify(root)
	}
	mustGit(t, root, "add", ".")
	mustGit(t, root, "commit", "-m", "Implement and digest")
	candidate := mustGit(t, root, "rev-parse", "HEAD")
	mustGit(t, root, "switch", "main")
	path := filepath.Join(root, ".git", "clue", "acceptance", "CH-001.md")
	r := Request{Root: root, Candidate: candidate, Base: base, BriefPath: path, Version: "dev"}
	writeBrief(t, r)
	return r
}

func writeBrief(t *testing.T, r Request) {
	t.Helper()
	body := fmt.Sprintf("---\ntype: acceptance-brief\ntitle: Accept example\nchange: CH-001\ncandidate: %s\nbase: %s\nreviewed: %s\n---\n", r.Candidate, r.Base, r.Candidate)
	for _, section := range sections {
		body += "\n## " + section + "\n\nRecorded result for " + section + ".\n"
	}
	write(t, filepath.Dir(r.BriefPath), filepath.Base(r.BriefPath), body)
}

func TestAC219_IntegrationPositive_ExactTreeAndRecoverableHistory(t *testing.T) {
	r := fixture(t, nil)
	p, err := Check(r)
	if err != nil {
		t.Fatal(err)
	}
	brief, _ := os.ReadFile(r.BriefPath)
	commit, err := p.Confirm(strings.NewReader("accept CH-001\n"), io.Discard, true)
	if err != nil {
		t.Fatal(err)
	}
	if got := mustGit(t, r.Root, "show", "-s", "--format=%P", commit); got != r.Base+" "+r.Candidate {
		t.Fatalf("parents %s", got)
	}
	if got := mustGit(t, r.Root, "rev-parse", "HEAD^{tree}"); got != p.tree {
		t.Fatalf("wrong tree %s", got)
	}
	if got := mustGit(t, r.Root, "status", "--porcelain"); got != "" {
		t.Fatalf("dirty after acceptance: %s", got)
	}
	mustGit(t, r.Root, "branch", "-D", "change")
	if got := mustGit(t, r.Root, "show", "-s", "--format=%B", "HEAD"); !strings.Contains(got, strings.TrimSpace(string(brief))) {
		t.Fatal("brief not retained")
	}
	proposal := mustGit(t, r.Root, "show", r.Candidate+"^:changes/CH-001-example/proposal.md")
	if !strings.Contains(proposal, "id: CH-001") {
		t.Fatal("proposal lost")
	}
	if got := mustGit(t, r.Root, "log", "--format=%H", "HEAD"); !strings.Contains(got, r.Candidate) {
		t.Fatal("candidate not reachable")
	}
}

func TestAC219_IntegrationNegative_IncompleteOrMismatchedBrief(t *testing.T) {
	r := fixture(t, nil)
	original, _ := os.ReadFile(r.BriefPath)
	for _, mutation := range []func(string) string{
		func(s string) string { return strings.Replace(s, "reviewed: "+r.Candidate, "reviewed: "+r.Base, 1) },
		func(s string) string { return strings.Replace(s, "candidate: "+r.Candidate, "candidate: "+r.Base, 1) },
		func(s string) string { return strings.Replace(s, "base: "+r.Base, "base: "+r.Candidate, 1) },
		func(s string) string { return strings.Replace(s, "Recorded result for Verification.", "", 1) },
		func(s string) string { return s + "\n<!-- REQUIRED -->\n" },
		func(s string) string { return strings.Replace(s, "change: CH-001", "change: CH-002", 1) },
		func(s string) string { return "no header" },
	} {
		write(t, filepath.Dir(r.BriefPath), filepath.Base(r.BriefPath), mutation(string(original)))
		if _, err := Check(r); err == nil {
			t.Fatal("invalid brief accepted")
		}
		if got := mustGit(t, r.Root, "rev-parse", "HEAD"); got != r.Base {
			t.Fatal("ref changed")
		}
	}
}

func TestSanity_LinkedWorktreeAcceptance(t *testing.T) {
	r := fixture(t, nil)
	mustGit(t, r.Root, "switch", "change")
	other := filepath.Join(t.TempDir(), "integration")
	mustGit(t, r.Root, "worktree", "add", other, "main")
	r.Root = other
	path := mustGit(t, other, "rev-parse", "--git-path", "clue/acceptance/CH-001.md")
	if !filepath.IsAbs(path) {
		path = filepath.Join(other, path)
	}
	r.BriefPath = path
	writeBrief(t, r)
	p, err := Check(r)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = p.Confirm(strings.NewReader("accept CH-001\n"), io.Discard, true); err != nil {
		t.Fatal(err)
	}
	if got := mustGit(t, other, "status", "--porcelain"); got != "" {
		t.Fatalf("dirty linked checkout: %s", got)
	}
}

func TestAC220_IntegrationNegative_BaseAdvanceDuringConfirmation(t *testing.T) {
	r := fixture(t, nil)
	p, err := Check(r)
	if err != nil {
		t.Fatal(err)
	}
	var advanced string
	reader := &duringConfirmation{reader: strings.NewReader("accept CH-001\n"), action: func() {
		mustGit(t, r.Root, "commit", "--allow-empty", "-m", "Another accepted change")
		advanced = mustGit(t, r.Root, "rev-parse", "HEAD")
	}}
	if _, err = p.Confirm(reader, io.Discard, true); err == nil || !strings.Contains(err.Error(), "tip changed") {
		t.Fatalf("stale confirmation: %v", err)
	}
	if got := mustGit(t, r.Root, "rev-parse", "HEAD"); got != advanced {
		t.Fatal("overwrote newer accepted work")
	}
}

func TestAC220_IntegrationPositive_GitCallbacksNeverExecute(t *testing.T) {
	r := fixture(t, nil)
	callbacks := filepath.Join(r.Root, ".git", "configured-hooks")
	marker := filepath.Join(r.Root, ".git", "callback-ran")
	for _, name := range []string{"reference-transaction", "post-index-change", "fsmonitor"} {
		write(t, callbacks, name, "#!/bin/sh\nprintf '%s\\n' '"+name+"' >> '"+filepath.ToSlash(marker)+"'\nexit 0\n")
		if err := os.Chmod(filepath.Join(callbacks, name), 0755); err != nil {
			t.Fatal(err)
		}
	}
	mustGit(t, r.Root, "config", "core.hooksPath", callbacks)
	mustGit(t, r.Root, "config", "core.fsmonitor", filepath.ToSlash(filepath.Join(callbacks, "fsmonitor")))
	// Prove this fixture would invoke its callbacks without our overrides.
	for _, args := range [][]string{{"status", "--porcelain"}, {"update-ref", "refs/heads/hook-control", r.Base}} {
		out, err := exec.Command("git", append([]string{"-C", r.Root}, args...)...).CombinedOutput()
		if err != nil {
			t.Fatalf("hook control: %v: %s", err, out)
		}
	}
	control, err := os.ReadFile(marker)
	if err != nil || !strings.Contains(string(control), "fsmonitor") || !strings.Contains(string(control), "reference-transaction") {
		t.Fatalf("callbacks were not active: %s (%v)", control, err)
	}
	if err = os.Remove(marker); err != nil {
		t.Fatal(err)
	}
	p, err := Check(r)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = p.Confirm(strings.NewReader("accept CH-001\n"), io.Discard, true); err != nil {
		t.Fatal(err)
	}
	if _, err = os.Stat(marker); !os.IsNotExist(err) {
		t.Fatalf("acceptance invoked a configured callback: %v", err)
	}
}

func TestAC220_IntegrationPositive_CheckAndCancellationAreReadOnly(t *testing.T) {
	r := fixture(t, nil)
	refs := mustGit(t, r.Root, "show-ref")
	index := mustGit(t, r.Root, "write-tree")
	p, err := Check(r)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		line        string
		interactive bool
	}{{"no\n", true}, {"", true}, {"accept CH-001\n", false}} {
		if _, err = p.Confirm(strings.NewReader(tc.line), io.Discard, tc.interactive); err == nil {
			t.Fatal("unexpected acceptance")
		}
	}
	if mustGit(t, r.Root, "show-ref") != refs || mustGit(t, r.Root, "write-tree") != index || mustGit(t, r.Root, "status", "--porcelain") != "" {
		t.Fatal("preflight or cancellation mutated repository")
	}
}

type duringConfirmation struct {
	action func()
	reader io.Reader
}

func (d *duringConfirmation) Read(b []byte) (int, error) {
	if d.action != nil {
		d.action()
		d.action = nil
	}
	return d.reader.Read(b)
}

func TestAC220_IntegrationNegative_UnsafeOrChangedState(t *testing.T) {
	for _, name := range []string{"dirty", "untracked", "stale", "workspace", "invalid", "missing-history", "changed-brief", "changed-during-confirmation", "ref-lock", "unsupported-link", "export-ignore"} {
		t.Run(name, func(t *testing.T) {
			r := fixture(t, func(root string) {
				switch name {
				case "workspace":
					write(t, root, "changes/left.txt", "unfinished")
				case "invalid", "export-ignore":
					write(t, root, "docs/broken.md", "no frontmatter")
					if name == "export-ignore" {
						write(t, root, ".gitattributes", "docs/broken.md export-ignore\n")
					}
				}
			})
			switch name {
			case "dirty":
				write(t, r.Root, "app.txt", "edited")
			case "untracked":
				write(t, r.Root, "untracked.txt", "edited")
			case "stale":
				mustGit(t, r.Root, "commit", "--allow-empty", "-m", "Advance base")
			case "missing-history":
				r.Candidate = mustGit(t, r.Root, "commit-tree", r.Candidate+"^{tree}", "-p", r.Base, "-m", "No proposal")
				writeBrief(t, r)
			case "unsupported-link":
				mustGit(t, r.Root, "switch", "change")
				blob := mustGit(t, r.Root, "rev-parse", r.Candidate+":app.txt")
				mustGit(t, r.Root, "update-index", "--add", "--cacheinfo", "120000,"+blob+",link")
				mustGit(t, r.Root, "commit", "-m", "Add symlink")
				r.Candidate = mustGit(t, r.Root, "rev-parse", "HEAD")
				mustGit(t, r.Root, "switch", "-f", "main")
				writeBrief(t, r)
			}
			before := mustGit(t, r.Root, "rev-parse", "HEAD")
			p, err := Check(r)
			if name == "changed-brief" || name == "changed-during-confirmation" || name == "ref-lock" {
				if err != nil {
					t.Fatal(err)
				}
				reader := &duringConfirmation{reader: strings.NewReader("accept CH-001\n"), action: func() {
					switch name {
					case "changed-brief":
						b, _ := os.ReadFile(r.BriefPath)
						write(t, filepath.Dir(r.BriefPath), filepath.Base(r.BriefPath), string(b)+"\nChanged meaning.\n")
					case "changed-during-confirmation":
						write(t, r.Root, "app.txt", "concurrent edit")
					case "ref-lock":
						write(t, r.Root, ".git/refs/heads/main.lock", "locked")
					}
				}}
				_, err = p.Confirm(reader, io.Discard, true)
			}
			if err == nil {
				t.Fatal("unsafe acceptance allowed")
			}
			if got := mustGit(t, r.Root, "rev-parse", "HEAD"); got != before {
				t.Fatal("acceptance moved ref on refusal")
			}
		})
	}
}

func TestAC230_IntegrationPositive_AllocationAndProceduralBoundary(t *testing.T) {
	r := fixture(t, func(root string) {
		write(t, root, ".clue/id-ledger.yaml", "version: 2\nevents: []\n")
		write(t, root, ".clue/id-coordination.yaml", "mode: git\nremote: origin\n")
		write(t, root, ".gitattributes", "/.clue/id-ledger.yaml merge=union\n")
	})
	p, err := Check(r)
	if err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if _, err = p.Confirm(strings.NewReader("accept CH-001\n"), &out, true); err != nil {
		t.Fatal(err)
	}
	for _, phrase := range []string{"procedural", "do not prove human presence", "recorded claims"} {
		if !strings.Contains(out.String(), phrase) {
			t.Fatalf("missing %s", phrase)
		}
	}
	for _, path := range []string{".clue/id-ledger.yaml", ".clue/id-coordination.yaml"} {
		if got := mustGit(t, r.Root, "diff", r.Candidate, "HEAD", "--", path); got != "" {
			t.Fatal("identity state changed")
		}
	}
}

func TestAC230_IntegrationNegative_ChangedPolicyCannotAcceptLocally(t *testing.T) {
	for _, name := range []string{"candidate-mode", "candidate-branch", "accepted-pr"} {
		t.Run(name, func(t *testing.T) {
			r := fixture(t, func(root string) {
				switch name {
				case "candidate-mode":
					write(t, root, ConfigPath, "mode: pr\nbranch: main\n")
				case "candidate-branch":
					write(t, root, ConfigPath, "mode: local\nbranch: other\n")
				}
			})
			if name == "accepted-pr" {
				write(t, r.Root, ConfigPath, "mode: pr\nbranch: main\n")
				mustGit(t, r.Root, "add", ConfigPath)
				mustGit(t, r.Root, "commit", "-m", "Require PR acceptance")
				r.Base = mustGit(t, r.Root, "rev-parse", "HEAD")
				mustGit(t, r.Root, "switch", "change")
				mustGit(t, r.Root, "merge", "--no-edit", "main")
				write(t, r.Root, ConfigPath, "mode: local\nbranch: main\n")
				mustGit(t, r.Root, "add", ConfigPath)
				mustGit(t, r.Root, "commit", "-m", "Unaccepted local policy")
				r.Candidate = mustGit(t, r.Root, "rev-parse", "HEAD")
				mustGit(t, r.Root, "switch", "main")
				writeBrief(t, r)
			}
			if _, err := Check(r); err == nil {
				t.Fatal("unsupported local acceptance passed")
			}
		})
	}
}

func TestAC230_IntegrationPositive_MissingPolicyAcceptsLocally(t *testing.T) {
	for _, explicitCandidate := range []bool{false, true} {
		r := fixture(t, nil)
		mustGit(t, r.Root, "rm", ConfigPath)
		mustGit(t, r.Root, "commit", "-m", "Use default local acceptance")
		r.Base = mustGit(t, r.Root, "rev-parse", "HEAD")
		mustGit(t, r.Root, "switch", "change")
		mustGit(t, r.Root, "merge", "--no-edit", "main")
		if explicitCandidate {
			write(t, r.Root, ConfigPath, "mode: local\nbranch: main\n")
			mustGit(t, r.Root, "add", ConfigPath)
			mustGit(t, r.Root, "commit", "-m", "Record unchanged local policy")
		}
		r.Candidate = mustGit(t, r.Root, "rev-parse", "HEAD")
		mustGit(t, r.Root, "switch", "main")
		writeBrief(t, r)
		p, err := Check(r)
		if err != nil {
			t.Fatal(err)
		}
		if _, err = p.Confirm(strings.NewReader("accept CH-001\n"), io.Discard, true); err != nil {
			t.Fatal(err)
		}
	}
}

func TestAC230_IntegrationPositive_SourceUsesLocalAcceptance(t *testing.T) {
	r := fixture(t, nil)
	write(t, r.Root, ".clue/role.yaml", "role: source\n")
	mustGit(t, r.Root, "add", ".clue/role.yaml")
	mustGit(t, r.Root, "commit", "-m", "Declare source role with local policy")
	r.Base = mustGit(t, r.Root, "rev-parse", "HEAD")
	mustGit(t, r.Root, "switch", "change")
	mustGit(t, r.Root, "merge", "--no-edit", "main")
	r.Candidate = mustGit(t, r.Root, "rev-parse", "HEAD")
	mustGit(t, r.Root, "switch", "main")
	writeBrief(t, r)
	p, err := Check(r)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = p.Confirm(strings.NewReader("accept CH-001\n"), io.Discard, true); err != nil {
		t.Fatal(err)
	}
}
