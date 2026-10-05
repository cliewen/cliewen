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
	"runtime"
	"sort"
	"strings"
	"time"
)

//go:embed Dockerfile scenarios
var assets embed.FS

// Options are the parameters of one invocation; none of them is a secret.
type Options struct {
	Scenario, Agent, Commit, TokenFile, OutDir, Model, Variant string
	Runs                                                       int
}

// Conditions is what every run in an invocation has to be read against.
type Conditions struct {
	Scenario      string    `json:"scenario"`
	Agent         string    `json:"agent"`
	ModelFlag     string    `json:"modelFlag"`
	Variant       string    `json:"variant"`
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
}

// Summary groups the runs by what the checks say they did.
type Summary struct {
	Conditions Conditions     `json:"conditions"`
	Runs       []RunResult    `json:"runs"`
	Spread     map[string]int `json:"spread"`
}

// Run executes the invocation and writes its records under OutDir.
func Run(o Options) error {
	if o.Agent != "claude" {
		return fmt.Errorf("agent %q has no adapter yet", o.Agent)
	}
	prompt, err := assets.ReadFile("scenarios/" + o.Scenario + "/prompt.txt")
	if err != nil {
		return fmt.Errorf("scenario %q: %w", o.Scenario, err)
	}
	setup, _ := assets.ReadFile("scenarios/" + o.Scenario + "/setup.sh")
	token, err := os.ReadFile(o.TokenFile)
	if err != nil {
		return fmt.Errorf("token file: %w", err)
	}
	cond := Conditions{Scenario: o.Scenario, Agent: o.Agent, ModelFlag: o.Model, Variant: o.Variant, Runs: o.Runs,
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
		res, err := runOnce(dir, i, cond, o, string(prompt), string(setup), strings.TrimSpace(string(token)))
		if err != nil {
			return fmt.Errorf("run %d: %w", i, err)
		}
		fmt.Printf("run %d: %s turns=%d cost=$%.2f %ds\n", i, res.Signature, res.Outcome.Turns, res.Outcome.CostUSD, res.Outcome.DurationMS/1000)
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
	sum := Summary{Conditions: cond, Spread: map[string]int{}}
	for i := 1; i <= cond.Runs; i++ {
		runDir := filepath.Join(dir, fmt.Sprintf("run-%d", i))
		var res RunResult
		if b, err := os.ReadFile(filepath.Join(runDir, "result.json")); err == nil {
			_ = json.Unmarshal(b, &res)
		}
		raw, err := os.ReadFile(filepath.Join(runDir, "events.jsonl"))
		if err != nil {
			return err
		}
		events, err := ParseEvents(bytes.NewReader(raw))
		if err != nil {
			return err
		}
		res.Run = i
		res.Outcome = Check(events, read(runDir, "post-status.txt"))
		res.Outcome.ConfigLeak = LeaksConfig(string(raw))
		res.Signature = res.Outcome.Signature()
		if err := writeJSON(filepath.Join(runDir, "result.json"), res); err != nil {
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

func runOnce(dir string, n int, cond Conditions, o Options, prompt, setup, token string) (RunResult, error) {
	runDir := filepath.Join(dir, fmt.Sprintf("run-%d", n))
	if err := os.MkdirAll(runDir, 0o755); err != nil {
		return RunResult{}, err
	}
	model := ""
	if o.Model != "" {
		model = " --model " + o.Model
	}
	script := "( " + setup + " ) && cd /home/node/work && " +
		"claude -p --setting-sources project,local --permission-mode acceptEdits --allowedTools 'Bash Read Write Edit Glob Grep' --max-turns 25 --output-format stream-json --verbose --no-session-persistence" + model +
		" > /out/events.jsonl 2> /out/stderr.txt; echo $? > /out/exit.txt; " +
		"git status --porcelain > /out/post-status.txt; sha256sum /usr/local/bin/clue | cut -d' ' -f1 > /out/post-clue.txt; " +
		"uname -sr > /out/post-os.txt; grep -m1 '^version' .agents/skills/clue-delta/skill.md > /out/post-skills.txt"
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, "docker", "run", "--rm", "-i", "-e", "CLAUDE_CODE_OAUTH_TOKEN",
		"--mount", "type=bind,source="+runDir+",target=/out", cond.ImageTag, "bash", "-c", script)
	cmd.Env = append(os.Environ(), "CLAUDE_CODE_OAUTH_TOKEN="+token)
	cmd.Stdin = strings.NewReader(prompt)
	var docker bytes.Buffer
	cmd.Stdout, cmd.Stderr = &docker, &docker
	runErr := cmd.Run()
	_ = os.WriteFile(filepath.Join(runDir, "docker.txt"), docker.Bytes(), 0o644)
	raw, err := os.ReadFile(filepath.Join(runDir, "events.jsonl"))
	if err != nil {
		return RunResult{}, fmt.Errorf("no event stream (%v); see %s", runErr, filepath.Join(runDir, "docker.txt"))
	}
	events, err := ParseEvents(bytes.NewReader(raw))
	if err != nil {
		return RunResult{}, err
	}
	out := Check(events, read(runDir, "post-status.txt"))
	out.ConfigLeak = LeaksConfig(string(raw))
	res := RunResult{Run: n, Outcome: out, Signature: out.Signature(),
		ExitCode: read(runDir, "exit.txt"), ContainerOS: read(runDir, "post-os.txt"),
		ClueMatchesBuild: read(runDir, "post-clue.txt") == cond.ClueSHA256,
		SkillsVersion:    strings.TrimPrefix(read(runDir, "post-skills.txt"), "version: ")}
	return res, writeJSON(filepath.Join(runDir, "result.json"), res)
}

type built struct{ dir, commit, skills, sha string }

// buildClue compiles a Linux clue from the named commit in a throwaway
// worktree, so the binary and the skills it scaffolds are that commit's.
func buildClue(rev string) (built, error) {
	root := strings.TrimSpace(capture("git", "rev-parse", "--show-toplevel"))
	commit := strings.TrimSpace(capture("git", "-C", root, "rev-parse", rev))
	if len(commit) != 40 {
		return built{}, fmt.Errorf("cannot resolve %q to a commit", rev)
	}
	ctxDir, err := os.MkdirTemp("", "scenario-trials-ctx-")
	if err != nil {
		return built{}, err
	}
	wt := filepath.Join(ctxDir, "src")
	if out, err := exec.Command("git", "-C", root, "worktree", "add", "--detach", wt, commit).CombinedOutput(); err != nil {
		return built{}, fmt.Errorf("worktree: %v: %s", err, out)
	}
	defer func() {
		_ = exec.Command("git", "-C", root, "worktree", "remove", "--force", wt).Run()
	}()
	bin := filepath.Join(ctxDir, "clue")
	cmd := exec.Command("go", "build", "-buildvcs=false", "-o", bin, "./cmd/clue")
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
