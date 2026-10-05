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

// claudeAdapter runs Claude Code headless with a subscription token.
type claudeAdapter struct{}

func init() { register(claudeAdapter{}) }

func (claudeAdapter) Name() string { return "claude" }

func (claudeAdapter) DefaultLogin(home string) string {
	return filepath.Join(home, ".cliewen-trial", "token")
}

func (claudeAdapter) Credentials(o Options) (Credentials, error) {
	token, err := os.ReadFile(o.Login)
	if err != nil {
		return Credentials{}, fmt.Errorf("token file: %w", err)
	}
	return Credentials{
		Env:  []string{"CLAUDE_CODE_OAUTH_TOKEN=" + strings.TrimSpace(string(token))},
		Pass: []string{"CLAUDE_CODE_OAUTH_TOKEN"},
	}, nil
}

func (claudeAdapter) Command(model string) string {
	m := ""
	if model != "" {
		m = " --model " + model
	}
	return "claude -p --setting-sources project,local --permission-mode acceptEdits " +
		"--allowedTools 'Bash Read Write Edit Glob Grep' --max-turns 25 --output-format stream-json " +
		"--verbose --no-session-persistence" + m + " > /out/events.jsonl 2> /out/stderr.txt"
}

func (claudeAdapter) Probe() string { return "claude --version > /out/post-agent.txt" }

// claudeEvent is the part of one Claude Code stream-json line that is read.
type claudeEvent struct {
	Type       string            `json:"type"`
	Subtype    string            `json:"subtype"`
	Model      string            `json:"model"`
	MCPServers []json.RawMessage `json:"mcp_servers"`
	CodeVers   string            `json:"claude_code_version"`
	Message    struct {
		Content []claudeBlock `json:"content"`
	} `json:"message"`
	NumTurns int     `json:"num_turns"`
	CostUSD  float64 `json:"total_cost_usd"`
	Duration int     `json:"duration_ms"`
}

type claudeBlock struct {
	Type  string          `json:"type"`
	Text  string          `json:"text"`
	Name  string          `json:"name"`
	Input json.RawMessage `json:"input"`
}

func (claudeAdapter) Parse(r io.Reader) (Transcript, error) {
	var t Transcript
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 0, 1<<20), 16<<20)
	for sc.Scan() {
		var e claudeEvent
		if json.Unmarshal(sc.Bytes(), &e) != nil || e.Type == "" {
			continue
		}
		switch e.Type {
		case "system":
			if e.Subtype == "init" {
				t.Model, t.AgentVersion, t.MCPServers = e.Model, e.CodeVers, len(e.MCPServers)
			}
		case "result":
			t.Turns, t.CostUSD, t.DurationMS, t.Finished = e.NumTurns, e.CostUSD, e.Duration, e.Subtype == "success"
		case "assistant":
			for _, b := range e.Message.Content {
				switch b.Type {
				case "text":
					t.Steps = append(t.Steps, Step{Kind: "text", Text: b.Text})
				case "tool_use":
					var in struct {
						Command string `json:"command"`
					}
					_ = json.Unmarshal(b.Input, &in)
					t.Steps = append(t.Steps, Step{Kind: "tool", Tool: b.Name, Command: in.Command})
				}
			}
		}
	}
	return t, sc.Err()
}
