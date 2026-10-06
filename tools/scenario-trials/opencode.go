package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// opencodeAdapter runs OpenCode headless against a model reached through a
// provider key kept in a data directory outside the repository and mounted over
// the container's OpenCode data directory. The container has no OpenCode
// configuration of its own, so what the run reads is what the fixture holds.
// The model is the -model flag in OpenCode's provider/model form; the reasoning
// effort is fixed here so every run states the same condition.
type opencodeAdapter struct{ effort string }

func init() { register(opencodeAdapter{}) }

const opencodeDefaultEffort = "medium"

func (opencodeAdapter) Name() string { return "opencode" }

// Effort is the reasoning effort the run uses, which the conditions record.
func (a opencodeAdapter) Effort() string {
	if a.effort != "" {
		return a.effort
	}
	return opencodeDefaultEffort
}

// WithEffort returns the adapter set to run at the given reasoning effort.
func (a opencodeAdapter) WithEffort(effort string) Adapter { a.effort = effort; return a }

func (opencodeAdapter) DefaultLogin(home string) string {
	return filepath.Join(home, ".cliewen-trial", "opencode")
}

func (opencodeAdapter) Credentials(o Options) (Credentials, error) {
	dir, err := filepath.Abs(o.Login)
	if err != nil {
		return Credentials{}, err
	}
	if strings.Contains(dir, ",") {
		return Credentials{}, fmt.Errorf("opencode login: a comma in %s would alter the mount specification", dir)
	}
	if _, err := os.Stat(filepath.Join(dir, "auth.json")); err != nil {
		return Credentials{}, fmt.Errorf("opencode login: no auth.json in %s: %w", dir, err)
	}
	return Credentials{Mounts: []string{"type=bind,source=" + dir + ",target=/home/node/.local/share/opencode"}}, nil
}

func (a opencodeAdapter) Command(model string) string {
	m := ""
	if model != "" {
		m = " -m " + model
	}
	// Docker creates the parents of the mount as root, so OpenCode's state, cache
	// and config directories go to /tmp, where the user can write.
	return "XDG_STATE_HOME=/tmp/oc-state XDG_CACHE_HOME=/tmp/oc-cache XDG_CONFIG_HOME=/tmp/oc-config " +
		"opencode run --format json --variant " + a.Effort() + m + " > /out/events.jsonl 2> /out/stderr.txt"
}

func (opencodeAdapter) Probe() string { return "opencode --version > /out/post-agent.txt" }

// opencodeEvent is the part of one OpenCode run --format json line that is read.
type opencodeEvent struct {
	Type      string `json:"type"`
	Timestamp int64  `json:"timestamp"`
	Part      struct {
		Type   string  `json:"type"`
		Text   string  `json:"text"`
		Tool   string  `json:"tool"`
		Reason string  `json:"reason"`
		Cost   float64 `json:"cost"`
		State  struct {
			Input struct {
				Command string `json:"command"`
			} `json:"input"`
		} `json:"state"`
	} `json:"part"`
}

func (opencodeAdapter) Parse(r io.Reader) (Transcript, error) {
	var t Transcript
	var first, last int64
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 0, 1<<20), 16<<20)
	for sc.Scan() {
		var e opencodeEvent
		if json.Unmarshal(sc.Bytes(), &e) != nil || e.Type == "" {
			continue
		}
		if first == 0 {
			first = e.Timestamp
		}
		last = e.Timestamp
		switch e.Type {
		case "text":
			t.Steps = append(t.Steps, Step{Kind: "text", Text: e.Part.Text})
		case "tool_use":
			switch e.Part.Tool {
			case "bash":
				t.Steps = append(t.Steps, Step{Kind: "tool", Tool: "Bash", Command: e.Part.State.Input.Command})
			case "write":
				t.Steps = append(t.Steps, Step{Kind: "tool", Tool: "Write"})
			case "edit", "multiedit", "patch", "apply_patch":
				t.Steps = append(t.Steps, Step{Kind: "tool", Tool: "Edit"})
			}
		case "step_finish":
			t.Turns++
			t.CostUSD += e.Part.Cost
			t.Finished = e.Part.Reason == "stop"
		}
	}
	t.DurationMS = int(last - first)
	return t, sc.Err()
}
