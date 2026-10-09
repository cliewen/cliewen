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

	"github.com/cliewen/cliewen/internal/corpus"
	"github.com/cliewen/cliewen/internal/parity"
	"github.com/cliewen/cliewen/internal/role"
	"gopkg.in/yaml.v3"
)

const ConfigPath = ".clue/acceptance.yaml"

type config struct {
	Mode   string `yaml:"mode"`
	Branch string `yaml:"branch"`
}

// Brief metadata binds the human-readable claims to an exact candidate.
// Neither this metadata nor its body proves that tests ran or a human observed it.
type Brief struct {
	Type      string `yaml:"type"`
	Title     string `yaml:"title"`
	Change    string `yaml:"change"`
	Candidate string `yaml:"candidate"`
	Base      string `yaml:"base"`
	Reviewed  string `yaml:"reviewed"`
}

type Request struct {
	Root, Candidate, Base, BriefPath, Version string
}

// Candidate is private preflight state. Call Confirm for a fresh state check
// after showing it; callers cannot substitute a different tree or brief.
type Candidate struct {
	request                Request
	root, ref, tree, brief string
	metadata               Brief
}

var oid = regexp.MustCompile(`^(?:[0-9a-f]{40}|[0-9a-f]{64})$`)
var changeID = regexp.MustCompile(`^CH-[0-9]+$`)
var sections = []string{"Intent", "Criteria and evidence", "Binding decisions", "Verification", "Review"}

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

func parseBrief(data []byte, r Request) (Brief, error) {
	text := strings.ReplaceAll(string(data), "\r\n", "\n")
	var b Brief
	if !strings.HasPrefix(text, "---\n") {
		return b, fmt.Errorf("brief needs YAML frontmatter")
	}
	header, body, ok := strings.Cut(text[4:], "\n---\n")
	if !ok {
		return b, fmt.Errorf("brief frontmatter is not closed")
	}
	if err := decode([]byte(header), &b); err != nil {
		return b, fmt.Errorf("brief: %w", err)
	}
	if b.Type != "acceptance-brief" || strings.TrimSpace(b.Title) == "" || !changeID.MatchString(b.Change) {
		return b, fmt.Errorf("brief requires type acceptance-brief, title and canonical CH identity")
	}
	if b.Candidate != r.Candidate || b.Base != r.Base || b.Reviewed != r.Candidate {
		return b, fmt.Errorf("brief candidate, base and reviewed revisions must match the full supplied commit IDs")
	}
	if strings.Contains(text, "<!--") || strings.Contains(text, "<REQUIRED>") {
		return b, fmt.Errorf("brief still contains template placeholders or comments")
	}
	for _, section := range sections {
		marker := "## " + section + "\n"
		if strings.Count(body, marker) != 1 {
			return b, fmt.Errorf("brief requires one %s section", section)
		}
		_, content, _ := strings.Cut(body, marker)
		content, _, _ = strings.Cut(content, "\n## ")
		if strings.TrimSpace(content) == "" {
			return b, fmt.Errorf("brief %s section is empty", section)
		}
	}
	return b, nil
}

func loadConfig(root string) (config, error) {
	var c config
	r, _, err := role.Load(root)
	if err != nil {
		return c, err
	}
	if r == role.Source {
		return c, fmt.Errorf("source repositories require PR acceptance")
	}
	b, err := os.ReadFile(filepath.Join(root, ConfigPath))
	if err != nil {
		return c, fmt.Errorf("local acceptance requires committed %s opt-in: %w", ConfigPath, err)
	}
	if err = decode(b, &c); err != nil {
		return c, err
	}
	if c.Mode != "local" || c.Branch == "" {
		return c, fmt.Errorf("%s requires mode: local and an explicit branch", ConfigPath)
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
	// Load opt-in from the accepted commit, never ignored or assume-unchanged
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
	data, err := os.ReadFile(r.BriefPath)
	if err != nil {
		return nil, err
	}
	p.metadata, err = parseBrief(data, r)
	if err != nil {
		return nil, err
	}
	p.brief = string(data)
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

func (p *Candidate) proposal() error {
	rangeSpec := p.request.Base + ".." + p.request.Candidate
	paths, err := git(p.root, "-c", "core.quotePath=false", "log", "--format=", "--name-only", "--diff-filter=A", rangeSpec, "--", "changes/")
	if err != nil {
		return err
	}
	for _, path := range strings.Split(paths, "\n") {
		if !strings.HasPrefix(path, "changes/"+p.metadata.Change+"-") || !strings.HasSuffix(path, "/proposal.md") {
			continue
		}
		commit, e := git(p.root, "log", "-1", "--format=%H", "--diff-filter=A", rangeSpec, "--", path)
		if e != nil {
			return e
		}
		body, e := git(p.root, "show", commit+":"+path)
		if e != nil {
			return e
		}
		body = strings.ReplaceAll(body, "\r\n", "\n")
		if !strings.HasPrefix(body, "---\n") {
			continue
		}
		header, _, ok := strings.Cut(body[4:], "\n---")
		var meta struct {
			ID   string `yaml:"id"`
			Type string `yaml:"type"`
		}
		if ok && yaml.Unmarshal([]byte(header), &meta) == nil && meta.ID == p.metadata.Change && meta.Type == "change" {
			return nil
		}
	}
	return fmt.Errorf("candidate history lacks the %s proposal introduced after its base", p.metadata.Change)
}

// Confirm repeats all preflight after confirmation, including the brief bytes.
// Interactive is supplied by the CLI's terminal check, never by a CLI flag.
func (p *Candidate) Confirm(input io.Reader, output io.Writer, interactive bool) (string, error) {
	fmt.Fprintf(output, "Candidate: %s\nBase: %s\nBranch: %s\n\n%s\n\nLocal acceptance is procedural: Git identity and this prompt do not prove human presence. Verification and review declarations are recorded claims. No push will be performed.\n", p.request.Candidate, p.request.Base, p.ref, p.brief)
	if !interactive {
		return "", fmt.Errorf("acceptance requires an interactive human terminal; use --check for preflight")
	}
	fmt.Fprintf(output, "Type accept %s to accept this exact candidate: ", p.metadata.Change)
	line, err := bufio.NewReader(input).ReadString('\n')
	if err != nil || strings.TrimSpace(line) != "accept "+p.metadata.Change {
		return "", fmt.Errorf("acceptance cancelled; no integration performed")
	}
	fresh, err := Check(p.request)
	if err != nil {
		return "", err
	}
	if fresh.brief != p.brief || fresh.ref != p.ref || fresh.tree != p.tree {
		return "", fmt.Errorf("acceptance inputs changed during confirmation; rerun and review them")
	}
	message := fmt.Sprintf("Accept %s locally\n\nCandidate: %s\nBase: %s\n\n%s\n", p.metadata.Change, p.request.Candidate, p.request.Base, p.brief)
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
