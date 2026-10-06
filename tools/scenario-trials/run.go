package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"time"
)

//go:embed Dockerfile scenarios variants
var assets embed.FS

// effortRe keeps the effort from reaching the container's shell as anything but a word.
var effortRe = regexp.MustCompile(`^[a-z]+$`)

// modelRe keeps the model name from reaching the container's shell as anything but a word.
var modelRe = regexp.MustCompile(`^[A-Za-z0-9._:/-]+$`)

// Options are the parameters of one invocation; none of them is a secret.
type Options struct {
	Scenario, Agent, Commit, Login, OutDir, Model, Variant, Effort string
	Runs                                                           int
}

// Conditions is what every run in an invocation has to be read against.
type Conditions struct {
	Scenario      string    `json:"scenario"`
	Obligation    string    `json:"obligation"`
	Failure       string    `json:"failure"`
	Agent         string    `json:"agent"`
	ModelFlag     string    `json:"modelFlag"`
	Effort        string    `json:"effort"`
	Variant       string    `json:"variant"`
	VariantSHA256 string    `json:"variantSha256,omitempty"`
	Runs          int       `json:"runs"`
	ClueCommit    string    `json:"clueCommit"`
	SkillsVersion string    `json:"skillsVersion"`
	ClueSHA256    string    `json:"clueSha256"`
	ImageTag      string    `json:"imageTag"`
	ImageID       string    `json:"imageId"`
	DockerServer  string    `json:"dockerServer"`
	HostOS        string    `json:"hostOs"`
	Started       time.Time `json:"started"`
}

// RunResult is one run, with what was read back from inside the container.
type RunResult struct {
	Run              int     `json:"run"`
	ExitCode         string  `json:"exitCode"`
	Outcome          Outcome `json:"outcome"`
	Signature        string  `json:"signature"`
	ContainerOS      string  `json:"containerOs"`
	ClueMatchesBuild bool    `json:"clueMatchesBuild"`
	SkillsVersion    string  `json:"skillsVersion"`
	WallMS           int     `json:"wallMs"`
}

// Summary groups the runs by what the checks say they did.
type Summary struct {
	Conditions Conditions     `json:"conditions"`
	Runs       []RunResult    `json:"runs"`
	Spread     map[string]int `json:"spread"`
}

// Run executes the invocation and writes its records under OutDir.
func Run(o Options) error {
	ad, err := adapterFor(o.Agent)
	if err != nil {
		return err
	}
	scn, err := scenarioFor(o.Scenario)
	if err != nil {
		return err
	}
	prompt, err := scn.Prompt()
	if err != nil {
		return fmt.Errorf("scenario %q: %w", o.Scenario, err)
	}
	setup, err := scn.Setup()
	if err != nil {
		return fmt.Errorf("scenario %q: %w", o.Scenario, err)
	}
	variant, variantHash, err := variantScript(o.Variant)
	if err != nil {
		return err
	}
	if o.Model != "" && !modelRe.MatchString(o.Model) {
		return fmt.Errorf("model %q is not a plain model name", o.Model)
	}
	ad, effort, err := withEffort(ad, o.Effort)
	if err != nil {
		return err
	}
	creds, err := ad.Credentials(o)
	if err != nil {
		return err
	}
	cond := Conditions{Scenario: o.Scenario, Obligation: scn.Obligation, Failure: scn.Failure, VariantSHA256: variantHash, Agent: o.Agent, ModelFlag: o.Model, Effort: effort, Variant: o.Variant, Runs: o.Runs,
		HostOS: runtime.GOOS + "/" + runtime.GOARCH, Started: time.Now().UTC()}
	build, err := buildClue(o.Commit)
	if err != nil {
		return err
	}
	defer os.RemoveAll(build.dir)
	cond.ClueCommit, cond.SkillsVersion, cond.ClueSHA256 = build.commit, build.skills, build.sha
	if cond.ImageTag, cond.ImageID, err = ensureImage(build); err != nil {
		return err
	}
	cond.DockerServer = strings.TrimSpace(capture("docker", "version", "--format", "{{.Server.Version}}"))
	dir, err := filepath.Abs(filepath.Join(o.OutDir, cond.Started.Format("20060102-150405")+"-"+o.Scenario+"-"+o.Agent))
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(dir, "conditions.json"), cond); err != nil {
		return err
	}
	sum := Summary{Conditions: cond, Spread: map[string]int{}}
	for i := 1; i <= o.Runs; i++ {
		res, err := runOnce(dir, i, cond, o, ad, creds, prompt, setup, variant)
		if err != nil {
			return fmt.Errorf("run %d: %w", i, err)
		}
		fmt.Printf("run %d: %s turns=%d cost=$%.2f %ds finished=%t config-leak=%t clue-matches-build=%t\n", i, res.Signature, res.Outcome.Turns, res.Outcome.CostUSD, res.WallMS/1000, res.Outcome.Finished, res.Outcome.ConfigLeak, res.ClueMatchesBuild)
		sum.Runs = append(sum.Runs, res)
		sum.Spread[res.Signature]++
	}
	if err := writeJSON(filepath.Join(dir, "summary.json"), sum); err != nil {
		return err
	}
	fmt.Print(FormatSpread(sum))
	fmt.Println("records:", dir)
	return nil
}

// Recheck evaluates the stored transcripts of an earlier invocation again with
// the current checks and rewrites its run results and summary. It runs no agent.
func Recheck(dir string) error {
	var cond Conditions
	b, err := os.ReadFile(filepath.Join(dir, "conditions.json"))
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, &cond); err != nil {
		return err
	}
	ad, err := adapterFor(cond.Agent)
	if err != nil {
		return err
	}
	sum := Summary{Conditions: cond, Spread: map[string]int{}}
	for i := 1; i <= cond.Runs; i++ {
		runDir := filepath.Join(dir, fmt.Sprintf("run-%d", i))
		var prior RunResult
		if b, err := os.ReadFile(filepath.Join(runDir, "result.json")); err == nil {
			_ = json.Unmarshal(b, &prior)
		}
		res, err := evaluate(runDir, i, cond, ad, prior)
		if err != nil {
			return err
		}
		fmt.Printf("run %d: %s\n", i, res.Signature)
		sum.Runs = append(sum.Runs, res)
		sum.Spread[res.Signature]++
	}
	if err := writeJSON(filepath.Join(dir, "summary.json"), sum); err != nil {
		return err
	}
	fmt.Print(FormatSpread(sum))
	return nil
}

// FormatSpread lists the distinct behaviours and how many runs showed each.
func FormatSpread(s Summary) string {
	keys := make([]string, 0, len(s.Spread))
	for k := range s.Spread {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return s.Spread[keys[i]] > s.Spread[keys[j]] || (s.Spread[keys[i]] == s.Spread[keys[j]] && keys[i] < keys[j])
	})
	var b strings.Builder
	fmt.Fprintf(&b, "%d distinct behaviours over %d runs:\n", len(keys), len(s.Runs))
	for _, k := range keys {
		fmt.Fprintf(&b, "  %d x %s\n", s.Spread[k], k)
	}
	return b.String()
}

func runOnce(dir string, n int, cond Conditions, o Options, ad Adapter, creds Credentials, prompt, setup, variant string) (RunResult, error) {
	runDir := filepath.Join(dir, fmt.Sprintf("run-%d", n))
	if err := os.MkdirAll(runDir, 0o755); err != nil {
		return RunResult{}, err
	}
	apply := ""
	if variant != "" {
		apply = "( cd /home/node/work && " + variant + "\n ) && "
	}
	script := "( " + setup + " ) && " + apply + "( cd /home/node/work && git tag -f trial-base >/dev/null ) && cd /home/node/work && " + ad.Command(o.Model) + "; echo $? > /out/exit.txt; " +
		"( git status --porcelain; git diff --name-only trial-base HEAD 2>/dev/null | sed 's/^/ M /' ) > /out/post-status.txt; sha256sum /usr/local/bin/clue | cut -d' ' -f1 > /out/post-clue.txt; " +
		"uname -sr > /out/post-os.txt; grep -m1 '^version' .agents/skills/clue-delta/skill.md > /out/post-skills.txt; " + ad.Probe()
	args := []string{"run", "--rm", "-i"}
	for _, name := range creds.Pass {
		args = append(args, "-e", name)
	}
	args = append(args, "--mount", "type=bind,source="+runDir+",target=/out")
	for _, m := range creds.Mounts {
		args = append(args, "--mount", m)
	}
	args = append(args, cond.ImageTag, "bash", "-c", script)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, "docker", args...)
	cmd.Env = append(os.Environ(), creds.Env...)
	cmd.Stdin = strings.NewReader(prompt)
	var docker bytes.Buffer
	cmd.Stdout, cmd.Stderr = &docker, &docker
	start := time.Now()
	runErr := cmd.Run()
	wall := int(time.Since(start).Milliseconds())
	_ = os.WriteFile(filepath.Join(runDir, "docker.txt"), docker.Bytes(), 0o644)
	if _, err := os.Stat(filepath.Join(runDir, "events.jsonl")); err != nil {
		return RunResult{}, fmt.Errorf("no event stream (%v); see %s", runErr, filepath.Join(runDir, "docker.txt"))
	}
	return evaluate(runDir, n, cond, ad, RunResult{WallMS: wall})
}

// evaluate reads a run directory's stored files and scores them with the
// current checks. It is the one place a result is made, for a live run and for
// a re-check of stored ones, so the two cannot disagree.
func evaluate(runDir string, n int, cond Conditions, ad Adapter, prior RunResult) (RunResult, error) {
	raw, err := os.ReadFile(filepath.Join(runDir, "events.jsonl"))
	if err != nil {
		return RunResult{}, err
	}
	t, err := ad.Parse(bytes.NewReader(raw))
	if err != nil {
		return RunResult{}, err
	}
	scn, err := scenarioFor(cond.Scenario)
	if err != nil {
		return RunResult{}, err
	}
	out := scn.Check(t, readRaw(runDir, "post-status.txt"))
	out.ConfigLeak = LeaksConfig(string(raw))
	if out.AgentVersion == "" {
		out.AgentVersion = read(runDir, "post-agent.txt")
	}
	if out.Model == "" {
		out.Model = read(runDir, "post-model.txt")
	}
	res := prior
	res.Run, res.Outcome, res.Signature = n, out, scn.Signature(out)
	res.ExitCode, res.ContainerOS = read(runDir, "exit.txt"), read(runDir, "post-os.txt")
	res.ClueMatchesBuild = read(runDir, "post-clue.txt") == cond.ClueSHA256
	res.SkillsVersion = strings.TrimPrefix(read(runDir, "post-skills.txt"), "version: ")
	return res, writeJSON(filepath.Join(runDir, "result.json"), res)
}

type built struct{ dir, commit, skills, sha string }

// buildClue compiles a Linux clue from the named commit in a throwaway
// worktree, so the binary and the skills it scaffolds are that commit's.
func buildClue(rev string) (_ built, retErr error) {
	root := strings.TrimSpace(capture("git", "rev-parse", "--show-toplevel"))
	commit := strings.TrimSpace(capture("git", "-C", root, "rev-parse", rev))
	if len(commit) != 40 {
		return built{}, fmt.Errorf("cannot resolve %q to a commit", rev)
	}
	ctxDir, err := os.MkdirTemp("", "scenario-trials-ctx-")
	if err != nil {
		return built{}, err
	}
	defer func() {
		if retErr != nil {
			_ = os.RemoveAll(ctxDir)
		}
	}()
	wt := filepath.Join(ctxDir, "src")
	if out, err := exec.Command("git", "-C", root, "worktree", "add", "--detach", wt, commit).CombinedOutput(); err != nil {
		return built{}, fmt.Errorf("worktree: %v: %s", err, out)
	}
	defer func() {
		_ = exec.Command("git", "-C", root, "worktree", "remove", "--force", wt).Run()
	}()
	bin := filepath.Join(ctxDir, "clue")
	cmd := exec.Command("go", "build", "-buildvcs=false", "-trimpath", "-o", bin, "./cmd/clue")
	cmd.Dir = wt
	cmd.Env = append(os.Environ(), "GOOS=linux", "GOARCH=amd64", "CGO_ENABLED=0")
	if out, err := cmd.CombinedOutput(); err != nil {
		return built{}, fmt.Errorf("go build: %v: %s", err, out)
	}
	skills := ""
	if b, err := os.ReadFile(filepath.Join(wt, ".agents", "skills", "clue-delta", "skill.md")); err == nil {
		for _, l := range strings.Split(string(b), "\n") {
			if strings.HasPrefix(l, "version:") {
				skills = strings.TrimSpace(strings.TrimPrefix(l, "version:"))
				break
			}
		}
	}
	data, _ := os.ReadFile(bin)
	sum := sha256.Sum256(data)
	return built{dir: ctxDir, commit: commit, skills: skills, sha: hex.EncodeToString(sum[:])}, nil
}

// ensureImage builds the container image if this exact clue and Dockerfile
// have none yet, and returns its tag and identity.
func ensureImage(b built) (tag, id string, err error) {
	df, _ := assets.ReadFile("Dockerfile")
	if err = os.WriteFile(filepath.Join(b.dir, "Dockerfile"), df, 0o644); err != nil {
		return
	}
	h := sha256.Sum256(append([]byte(b.sha), df...))
	tag = "cliewen-trial:" + hex.EncodeToString(h[:])[:12]
	if exec.Command("docker", "image", "inspect", tag).Run() != nil {
		build := exec.Command("docker", "build", "-q", "-t", tag, "-f", filepath.Join(b.dir, "Dockerfile"), b.dir)
		if out, e := build.CombinedOutput(); e != nil {
			return "", "", fmt.Errorf("docker build: %v: %s", e, out)
		}
	}
	id = strings.TrimSpace(capture("docker", "image", "inspect", "--format", "{{.Id}}", tag))
	return
}

func capture(name string, args ...string) string {
	out, _ := exec.Command(name, args...).Output()
	return string(out)
}

// readRaw returns a file as it is, without trimming: a porcelain status line starts with a space that means something.
func readRaw(dir, name string) string {
	b, _ := os.ReadFile(filepath.Join(dir, name))
	return string(b)
}

func read(dir, name string) string {
	b, _ := os.ReadFile(filepath.Join(dir, name))
	return strings.TrimSpace(string(b))
}

func writeJSON(path string, v any) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(b, '\n'), 0o644)
}

// withEffort applies a requested reasoning effort to an adapter that takes one
// and returns the effort the conditions record: the adapter's own, or "default"
// for an agent that is run at its default.
func withEffort(ad Adapter, requested string) (Adapter, string, error) {
	if requested != "" {
		if !effortRe.MatchString(requested) {
			return nil, "", fmt.Errorf("effort %q is not a plain word", requested)
		}
		ea, ok := ad.(interface{ WithEffort(string) Adapter })
		if !ok {
			return nil, "", fmt.Errorf("agent %q does not take a reasoning effort", ad.Name())
		}
		ad = ea.WithEffort(requested)
	}
	if e, ok := ad.(interface{ Effort() string }); ok {
		return ad, e.Effort(), nil
	}
	return ad, "default", nil
}
