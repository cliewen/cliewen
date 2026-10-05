package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// codexAdapter runs Codex headless with the maintainer's own login, kept in a
// directory outside the repository and mounted over the container's .codex.
// Codex's own sandbox cannot start inside the container, so it runs with that
// sandbox off and the container is the boundary.
type codexAdapter struct{}

func init() { register(codexAdapter{}) }

func (codexAdapter) Name() string { return "codex" }

func (codexAdapter) DefaultLogin(home string) string {
	return filepath.Join(home, ".cliewen-trial", "codex")
}

func (codexAdapter) Credentials(o Options) (Credentials, error) {
	dir, err := filepath.Abs(o.Login)
	if err != nil {
		return Credentials{}, err
	}
	if strings.Contains(dir, ",") {
		return Credentials{}, fmt.Errorf("codex login: a comma in %s would alter the mount specification", dir)
	}
	if _, err := os.Stat(filepath.Join(dir, "auth.json")); err != nil {
		return Credentials{}, fmt.Errorf("codex login: no auth.json in %s: %w", dir, err)
	}
	return Credentials{Mounts: []string{"type=bind,source=" + dir + ",target=/home/node/.codex"}}, nil
}

func (codexAdapter) Command(model string) string {
	m := ""
	if model != "" {
		m = " -m " + model
	}
	return "codex exec --json --dangerously-bypass-approvals-and-sandbox --skip-git-repo-check -C /home/node/work" +
		m + " - > /out/events.jsonl 2> /out/stderr.txt"
}

func (codexAdapter) Probe() string {
	// The login directory persists between runs, so its sessions folder holds
	// earlier runs too; the newest session file is this run's.
	return `codex --version > /out/post-agent.txt; ` +
		`f=$(ls -t $(find /home/node/.codex/sessions -name '*.jsonl' 2>/dev/null) 2>/dev/null | head -1); ` +
		`[ -n "$f" ] && grep -ao '"model":"[^"]*"' "$f" | head -1 | cut -d'"' -f4 > /out/post-model.txt`
}

// codexEvent is the part of one Codex exec --json line that is read.
type codexEvent struct {
	Type string `json:"type"`
	Item struct {
		Type    string `json:"type"`
		Text    string `json:"text"`
		Command string `json:"command"`
	} `json:"item"`
}

var shellWrapRe = regexp.MustCompile(`(?s)^\S*(?:ba)?sh\s+-l?c\s+(.*)$`)

// unwrapShell returns the command inside the login-shell wrapper Codex records.
func unwrapShell(c string) string {
	m := shellWrapRe.FindStringSubmatch(strings.TrimSpace(c))
	if m == nil {
		return c
	}
	inner := strings.TrimSpace(m[1])
	if n := len(inner); n >= 2 && (inner[0] == '\'' || inner[0] == '"') && inner[n-1] == inner[0] {
		return inner[1 : n-1]
	}
	return inner
}

func (codexAdapter) Parse(r io.Reader) (Transcript, error) {
	var t Transcript
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 0, 1<<20), 16<<20)
	for sc.Scan() {
		var e codexEvent
		if json.Unmarshal(sc.Bytes(), &e) != nil || e.Type == "" {
			continue
		}
		switch e.Type {
		case "turn.completed":
			t.Finished = true
		case "item.completed":
			switch e.Item.Type {
			case "agent_message":
				t.Steps = append(t.Steps, Step{Kind: "text", Text: e.Item.Text})
			case "command_execution":
				t.Steps = append(t.Steps, Step{Kind: "tool", Tool: "Bash", Command: unwrapShell(e.Item.Command)})
			case "file_change":
				t.Steps = append(t.Steps, Step{Kind: "tool", Tool: "Edit"})
			}
		}
	}
	return t, sc.Err()
}
