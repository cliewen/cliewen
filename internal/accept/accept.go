// Package accept implements the procedural, human-controlled local merge
// boundary. It checks Git transitions; corpus validation remains state-only.
package accept

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/cliewen/cliewen/internal/acceptpolicy"
	"github.com/cliewen/cliewen/internal/corpus"
	"github.com/cliewen/cliewen/internal/parity"
	"gopkg.in/yaml.v3"
)

const ConfigPath = acceptpolicy.Path

type config = acceptpolicy.Policy

// Request names the exact candidate and base a human is asked to accept.
type Request struct {
	Root, Candidate, Base, Version string
}

// Candidate is private preflight state. Call Confirm for a fresh state check
// after showing it; callers cannot substitute a different tree or proposal.
type Candidate struct {
	request   Request
	root, ref string
	tree      string
	change    string
	title     string
	evidence  string
}

var oid = regexp.MustCompile(`^(?:[0-9a-f]{40}|[0-9a-f]{64})$`)

func git(root string, args ...string) (string, error) {
	b, err := command(root, args...).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git %s: %w: %s", args[0], err, strings.TrimSpace(string(b)))
	}
	return strings.TrimSpace(string(b)), nil
}

func decode(data []byte, dst any) error {
	d := yaml.NewDecoder(bytes.NewReader(data))
	d.KnownFields(true)
	if err := d.Decode(dst); err != nil {
		return err
	}
	var extra any
	if err := d.Decode(&extra); err != io.EOF {
		return fmt.Errorf("expected one YAML document")
	}
	return nil
}

func loadConfig(root string) (config, error) {
	c, _, err := acceptpolicy.Load(root)
	if err != nil {
		return c, err
	}
	if c.Mode != acceptpolicy.Local {
		return c, fmt.Errorf("PR acceptance is selected; clue accept performs local acceptance only")
	}
	if _, err = git(root, "check-ref-format", "refs/heads/"+c.Branch); err != nil {
		return c, fmt.Errorf("invalid acceptance branch: %w", err)
	}
	return c, nil
}

// Check reads committed objects into a private snapshot, never a checkout of
// the candidate. Git attributes, filters, and hooks cannot alter that snapshot.
func Check(r Request) (*Candidate, error) {
	if !oid.MatchString(r.Candidate) || !oid.MatchString(r.Base) || r.Candidate == r.Base {
		return nil, fmt.Errorf("supply distinct full lowercase candidate and base commit IDs")
	}
	root, err := git(r.Root, "rev-parse", "--show-toplevel")
	if err != nil {
		return nil, err
	}
	r.Root = root
	p := &Candidate{request: r, root: root}
	if err = p.cleanBase(); err != nil {
		return nil, err
	}
	// Resolve policy from the accepted commit, never ignored or assume-unchanged
	// working files. An unaccepted candidate cannot authorize itself.
	baseSnapshot, err := os.MkdirTemp("", "clue-accept-base-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(baseSnapshot)
	if err = materialize(root, r.Base, baseSnapshot); err != nil {
		return nil, err
	}
	baseConfig, err := loadConfig(baseSnapshot)
	if err != nil {
		return nil, err
	}
	p.ref = "refs/heads/" + baseConfig.Branch
	if err = p.cleanBase(); err != nil {
		return nil, err
	}
	for _, rev := range []string{r.Base, r.Candidate} {
		kind, e := git(root, "cat-file", "-t", rev)
		if e != nil || kind != "commit" {
			return nil, fmt.Errorf("%s is not an available commit", rev)
		}
	}
	if _, err = git(root, "merge-base", "--is-ancestor", r.Base, r.Candidate); err != nil {
		return nil, fmt.Errorf("candidate must incorporate the current accepted base: %w", err)
	}
	if err = p.proposal(); err != nil {
		return nil, err
	}
	snapshot, err := os.MkdirTemp("", "clue-accept-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(snapshot)
	if err = materialize(root, r.Candidate, snapshot); err != nil {
		return nil, err
	}
	candidateConfig, err := loadConfig(snapshot)
	if err != nil {
		return nil, err
	}
	if candidateConfig != baseConfig {
		return nil, fmt.Errorf("base and candidate must select the same acceptance mode and branch")
	}
	c, issues := corpus.Scan(snapshot)
	issues = append(issues, corpus.Validate(c, corpus.Options{ForbidChanges: true, Version: r.Version})...)
	issues = append(issues, parity.CheckReports(snapshot)...)
	if len(issues) > 0 {
		var messages []string
		for _, issue := range issues {
			messages = append(messages, issue.String())
		}
		return nil, fmt.Errorf("candidate validation failed:\n%s", strings.Join(messages, "\n"))
	}
	p.tree, err = git(root, "rev-parse", r.Candidate+"^{tree}")
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (p *Candidate) cleanBase() error {
	head, err := git(p.root, "rev-parse", "HEAD")
	if err != nil {
		return err
	}
	if head != p.request.Base {
		return fmt.Errorf("integration tip changed; expected base %s", p.request.Base)
	}
	ref, err := git(p.root, "symbolic-ref", "HEAD")
	if err != nil {
		return fmt.Errorf("acceptance requires a checked-out integration branch")
	}
	if p.ref != "" && ref != p.ref {
		return fmt.Errorf("checkout %s before acceptance", p.ref)
	}
	status, err := git(p.root, "status", "--porcelain=v1", "--untracked-files=all")
	if err != nil {
		return err
	}
	if status != "" {
		return fmt.Errorf("acceptance requires a clean index and working tree, including untracked files")
	}
	for _, name := range []string{"MERGE_HEAD", "CHERRY_PICK_HEAD", "REVERT_HEAD", "rebase-merge", "rebase-apply", "sequencer", "BISECT_START"} {
		path, e := git(p.root, "rev-parse", "--git-path", name)
		if e != nil {
			return e
		}
		if !filepath.IsAbs(path) {
			path = filepath.Join(p.root, path)
		}
		if _, e = os.Stat(path); e == nil {
			return fmt.Errorf("finish or abort the existing Git operation (%s)", name)
		} else if !os.IsNotExist(e) {
			return e
		}
	}
	return nil
}

// proposal finds the one change whose proposal the candidate's history
// introduces after the base, and reads that proposal's last committed text:
// the digest deletes the workspace, so the candidate tree holds none.
func (p *Candidate) proposal() error {
	rangeSpec := p.request.Base + ".." + p.request.Candidate
	paths, err := git(p.root, "-c", "core.quotePath=false", "log", "--format=", "--name-only", "--diff-filter=A", rangeSpec, "--", "changes/")
	if err != nil {
		return err
	}
	var ids []string
	for _, path := range strings.Split(paths, "\n") {
		if !strings.HasPrefix(path, "changes/") || !strings.HasSuffix(path, "/proposal.md") {
			continue
		}
		commit, e := git(p.root, "log", "-1", "--format=%H", "--diff-filter=AM", rangeSpec, "--", path)
		if e != nil {
			return e
		}
		body, e := git(p.root, "show", commit+":"+path)
		if e != nil {
			return e
		}
		id, title, evidence, ok := parseProposal(body)
		if !ok || !strings.HasPrefix(path, "changes/"+id+"-") {
			continue
		}
		ids = append(ids, id)
		p.change, p.title, p.evidence = id, title, evidence
	}
	switch len(ids) {
	case 0:
		return fmt.Errorf("candidate history introduces no tracked change proposal after its base")
	case 1:
		return nil
	default:
		return fmt.Errorf("candidate history introduces more than one tracked change proposal after its base (%s); accept one change at a time", strings.Join(ids, ", "))
	}
}

var changeID = regexp.MustCompile(`^CH-[0-9]+$`)

// parseProposal reads a change proposal's identity, title and optional
// human evidence section, which carries Human-class proof (PDR-073).
func parseProposal(body string) (id, title, evidence string, ok bool) {
	body = strings.ReplaceAll(body, "\r\n", "\n")
	if !strings.HasPrefix(body, "---\n") {
		return "", "", "", false
	}
	header, text, found := strings.Cut(body[4:], "\n---")
	var meta struct {
		ID    string `yaml:"id"`
		Type  string `yaml:"type"`
		Title string `yaml:"title"`
	}
	if !found || yaml.Unmarshal([]byte(header), &meta) != nil || meta.Type != "change" || !changeID.MatchString(meta.ID) {
		return "", "", "", false
	}
	if _, rest, has := strings.Cut(text, "\n## Human evidence\n"); has {
		rest, _, _ = strings.Cut(rest, "\n## ")
		evidence = strings.TrimSpace(rest)
	}
	return meta.ID, strings.TrimSpace(meta.Title), evidence, true
}

// summary is what the human is shown and what the merge message records.
func (p *Candidate) summary() string {
	text := fmt.Sprintf("Change: %s — %s\nCandidate: %s\nBase: %s\n", p.change, p.title, p.request.Candidate, p.request.Base)
	if p.evidence != "" {
		text += "\n## Human evidence\n\n" + p.evidence + "\n"
	}
	return text
}

// Confirm repeats all preflight after confirmation.
// Interactive is supplied by the CLI's terminal check, never by a CLI flag.
func (p *Candidate) Confirm(input io.Reader, output io.Writer, interactive bool) (string, error) {
	fmt.Fprintf(output, "%sBranch: %s\n\nLocal acceptance is procedural: Git identity and this prompt do not prove human presence. Verification and review are recorded claims. No push will be performed.\n", p.summary(), p.ref)
	if !interactive {
		return "", fmt.Errorf("acceptance requires an interactive human terminal; use --check for preflight")
	}
	fmt.Fprintf(output, "Type accept %s to accept this exact candidate: ", p.change)
	line, err := bufio.NewReader(input).ReadString('\n')
	if err != nil || strings.TrimSpace(line) != "accept "+p.change {
		return "", fmt.Errorf("acceptance cancelled; no integration performed")
	}
	return p.complete("")
}

func (p *Candidate) complete(approvalRecord string) (string, error) {
	fresh, err := Check(p.request)
	if err != nil {
		return "", err
	}
	if fresh.summary() != p.summary() || fresh.ref != p.ref || fresh.tree != p.tree {
		return "", fmt.Errorf("acceptance inputs changed during confirmation; rerun and review them")
	}
	message := fmt.Sprintf("Accept %s locally\n\n%s", p.change, p.summary())
	if approvalRecord != "" {
		message += "\n## Recorded human approval\n\n" + approvalRecord + "\n"
	}
	c := command(p.root, "commit-tree", p.tree, "-p", p.request.Base, "-p", p.request.Candidate)
	c.Stdin = strings.NewReader(message)
	var errors bytes.Buffer
	c.Stderr = &errors
	out, err := c.Output()
	if err != nil {
		return "", fmt.Errorf("create acceptance commit: %w: %s", err, errors.String())
	}
	commit := strings.TrimSpace(string(out))
	if err = p.integrate(commit); err != nil {
		return "", err
	}
	return commit, nil
}

// integrate locks the target ref with Git's prepared transaction before
// updating the checkout. read-tree refuses overwriting local edits. CAS binds
// acceptance to the reviewed base even if another process advances the branch.
func (p *Candidate) integrate(commit string) error {
	if err := p.cleanBase(); err != nil {
		return err
	}
	c := command(p.root, "update-ref", "-m", "clue: local human acceptance", "--stdin")
	stdin, err := c.StdinPipe()
	if err != nil {
		return err
	}
	stdout, err := c.StdoutPipe()
	if err != nil {
		return err
	}
	var errors bytes.Buffer
	c.Stderr = &errors
	if err = c.Start(); err != nil {
		return err
	}
	waited := false
	defer func() {
		stdin.Close()
		if !waited {
			_ = c.Wait()
		}
	}()
	reader := bufio.NewReader(stdout)
	send := func(cmd, want string) error {
		if _, e := io.WriteString(stdin, cmd+"\n"); e != nil {
			return e
		}
		line, e := reader.ReadString('\n')
		if e != nil {
			return e
		}
		if strings.TrimSpace(line) != want {
			return fmt.Errorf("Git transaction: %s", line)
		}
		return nil
	}
	if err = send("start", "start: ok"); err == nil {
		err = send(fmt.Sprintf("update %s %s %s\nprepare", p.ref, commit, p.request.Base), "prepare: ok")
	}
	if err != nil {
		return fmt.Errorf("cannot lock acceptance base; no integration performed: %w", err)
	}
	if err = p.cleanBase(); err != nil {
		return err
	}
	if _, err = git(p.root, "read-tree", "-m", "-u", p.request.Base, p.request.Candidate); err != nil {
		return err
	}
	if err = send("commit", "commit: ok"); err != nil {
		// Preserve edits if an external writer appeared: read-tree never forces.
		_, rollback := git(p.root, "read-tree", "-m", "-u", p.request.Candidate, p.request.Base)
		return fmt.Errorf("acceptance ref transaction failed: %w; checkout rollback: %v; inspect git status before retrying", err, rollback)
	}
	stdin.Close()
	err = c.Wait()
	waited = true
	if err != nil {
		return fmt.Errorf("acceptance transaction: %w: %s", err, errors.String())
	}
	return nil
}
