package main

import (
	"bufio"
	"encoding/json"
	"io"
	"regexp"
	"strings"
)

// Event is the part of one Claude Code stream-json line the checks read.
type Event struct {
	Type    string `json:"type"`
	Subtype string `json:"subtype"`

	// system/init
	Model        string            `json:"model"`
	MCPServers   []json.RawMessage `json:"mcp_servers"`
	Plugins      []json.RawMessage `json:"plugins"`
	Skills       []string          `json:"skills"`
	MemoryPaths  map[string]string `json:"memory_paths"`
	CodeVersion  string            `json:"claude_code_version"`
	APIKeySource string            `json:"apiKeySource"`

	// assistant
	Message struct {
		Content []Block `json:"content"`
	} `json:"message"`

	// result
	NumTurns int     `json:"num_turns"`
	CostUSD  float64 `json:"total_cost_usd"`
	Duration int     `json:"duration_ms"`
}

// Block is one content block of an assistant message.
type Block struct {
	Type  string          `json:"type"`
	Text  string          `json:"text"`
	Name  string          `json:"name"`
	Input json.RawMessage `json:"input"`
}

// ParseEvents reads a stream-json transcript, skipping lines that are not JSON.
func ParseEvents(r io.Reader) ([]Event, error) {
	var out []Event
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 0, 1<<20), 16<<20)
	for sc.Scan() {
		var e Event
		if json.Unmarshal(sc.Bytes(), &e) == nil && e.Type != "" {
			out = append(out, e)
		}
	}
	return out, sc.Err()
}

// Outcome is what the deterministic checks say about one run.
type Outcome struct {
	TypoRoute         string  `json:"typoRoute"`
	ExportRoute       string  `json:"exportRoute"`
	VersionCheckFirst bool    `json:"versionCheckFirst"`
	EditBeforeRoute   bool    `json:"editBeforeRoute"`
	ReadmeEdited      bool    `json:"readmeEdited"`
	Recommended       bool    `json:"recommended"`
	ConfigLeak        bool    `json:"configLeak"`
	MCPServers        int     `json:"mcpServers"`
	Model             string  `json:"model"`
	AgentVersion      string  `json:"agentVersion"`
	Turns             int     `json:"turns"`
	CostUSD           float64 `json:"costUsd"`
	DurationMS        int     `json:"durationMs"`
	Finished          bool    `json:"finished"`
}

var (
	routeRe = regexp.MustCompile(`(?i)recommended route(?:\s+for\s+[^:\n]*)?:?\W*(direct|tracked|simple|full)`)
	// looseRe is any statement of a route in plain words, such as "as a direct change".
	looseRe = regexp.MustCompile(`(?i)\b(direct|tracked)\b`)
	// sentenceRe splits on line breaks and on sentence ends followed by space.
	sentenceRe = regexp.MustCompile(`\n|[.!?]\s+`)
	typoRe     = regexp.MustCompile(`(?i)typo|readme|teh`)
	csvRe      = regexp.MustCompile(`(?i)csv|export`)
	editRe     = regexp.MustCompile(`(?:^|[\s;&|])(?:sed\s+-i|tee\s|git\s+(?:add|commit)|rm\s|mv\s|cp\s)|[^>]>[^>&]`)
)

// Signature groups runs that behaved the same way on what the scenario asks.
func (o Outcome) Signature() string {
	b := func(v bool) string {
		if v {
			return "yes"
		}
		return "no"
	}
	return "typo=" + or(o.TypoRoute, "none") + " export=" + or(o.ExportRoute, "none") +
		" edit-before-route=" + b(o.EditBeforeRoute) + " readme-edited=" + b(o.ReadmeEdited)
}

func or(s, d string) string {
	if s == "" {
		return d
	}
	return s
}

// Check reads one run's events and the post-run working-tree listing.
func Check(events []Event, postStatus string) Outcome {
	var o Outcome
	toolSeen := false
	routeSeen := false
	for _, e := range events {
		switch e.Type {
		case "system":
			if e.Subtype == "init" {
				o.Model, o.AgentVersion, o.MCPServers = e.Model, e.CodeVersion, len(e.MCPServers)
			}
		case "result":
			o.Turns, o.CostUSD, o.DurationMS, o.Finished = e.NumTurns, e.CostUSD, e.Duration, e.Subtype == "success"
		case "assistant":
			for _, b := range e.Message.Content {
				switch b.Type {
				case "tool_use":
					cmd := toolCommand(b)
					if !toolSeen {
						toolSeen = true
						o.VersionCheckFirst = b.Name == "Bash" && strings.Contains(cmd, "clue latest")
					}
					if !routeSeen && mutates(b.Name, cmd) {
						o.EditBeforeRoute = true
					}
				case "text":
					if routeRe.MatchString(b.Text) {
						o.Recommended = true
						routes(b.Text, &o)
					}
					if routeRe.MatchString(b.Text) || looseRe.MatchString(b.Text) {
						routeSeen = true
						looseRoutes(b.Text, &o)
					}
				}
			}
		}
	}
	for _, line := range strings.Split(postStatus, "\n") {
		if strings.Contains(line, "README.md") {
			o.ReadmeEdited = true
		}
	}
	return o
}

func toolCommand(b Block) string {
	var in struct {
		Command string `json:"command"`
	}
	_ = json.Unmarshal(b.Input, &in)
	return in.Command
}

func mutates(tool, cmd string) bool {
	switch tool {
	case "Edit", "Write", "NotebookEdit":
		return true
	case "Bash":
		return editRe.MatchString(cmd)
	}
	return false
}

// routes assigns each "Recommended route" mention to the request it is about,
// by looking at its own line and the line before it.
func routes(text string, o *Outcome) {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		m := routeRe.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		route := normalize(m[1])
		ctx := line
		if topic(line) == "" && i > 0 {
			ctx = lines[i-1]
		}
		switch topic(ctx) {
		case "export":
			o.ExportRoute = route
		case "typo":
			o.TypoRoute = route
		}
	}
}

// looseRoutes fills a route the fixed phrase did not give, from any line that
// names a route in plain words and is about exactly one of the two requests.
func looseRoutes(text string, o *Outcome) {
	parts := sentenceRe.Split(text, -1)
	for i, line := range parts {
		m := looseRe.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		about := topic(line)
		if about == "" && i > 0 {
			about = topic(parts[i-1])
		}
		switch about {
		case "export":
			if o.ExportRoute == "" {
				o.ExportRoute = strings.ToLower(m[1])
			}
		case "typo":
			if o.TypoRoute == "" {
				o.TypoRoute = strings.ToLower(m[1])
			}
		}
	}
}

// topic says which request a line is about, or "" when it names neither or both.
func topic(s string) string {
	csv, typo := csvRe.MatchString(s), typoRe.MatchString(s)
	switch {
	case csv && !typo:
		return "export"
	case typo && !csv:
		return "typo"
	}
	return ""
}

func normalize(r string) string {
	switch strings.ToLower(r) {
	case "simple":
		return "direct"
	case "full":
		return "tracked"
	}
	return strings.ToLower(r)
}

// LeaksConfig reports whether the transcript shows host configuration the
// container should not have: anything that names CodeGraph, which only the
// maintainer's own instructions mention.
func LeaksConfig(raw string) bool {
	return strings.Contains(strings.ToLower(raw), "codegraph")
}
