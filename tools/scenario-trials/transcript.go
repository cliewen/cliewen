package main

import (
	"regexp"
	"strings"
)

// Step is one thing an agent said or did, in the order it happened. Every
// adapter turns its own event stream into these, so no check names a vendor.
type Step struct {
	Kind    string // "text" or "tool"
	Text    string // the words, for Kind "text"
	Tool    string // "Bash", "Edit" or "Write", for Kind "tool"
	Command string // the shell command, for Tool "Bash"
}

// Transcript is a run as the checks see it, with what the stream itself reports.
type Transcript struct {
	Steps        []Step
	Model        string
	AgentVersion string
	MCPServers   int
	Turns        int
	CostUSD      float64
	DurationMS   int
	Finished     bool
}

// Outcome is what the deterministic checks say about one run: what every run
// reports, and the fields its scenario observes.
type Outcome struct {
	Observed          map[string]string `json:"observed"`
	VersionCheckFirst bool              `json:"versionCheckFirst"`
	ConfigLeak        bool              `json:"configLeak"`
	MCPServers        int               `json:"mcpServers"`
	Model             string            `json:"model"`
	AgentVersion      string            `json:"agentVersion"`
	Turns             int               `json:"turns"`
	CostUSD           float64           `json:"costUsd"`
	DurationMS        int               `json:"durationMs"`
	Finished          bool              `json:"finished"`
}

var (
	routeRe = regexp.MustCompile(`(?i)recommended route(?:\s+for\s+[^:\n]*)?:?\W*(direct|tracked|simple|full)\b`)
	// looseRe is any statement of a route in plain words, such as "as a direct change".
	looseRe = regexp.MustCompile(`(?i)\b(direct|tracked)\b`)
	// shellWriteRe matches commands that change files whatever the redirections say.
	shellWriteRe = regexp.MustCompile(`(?:^|[\s;&|(])(?:sed\s+(?:-\w*i|--in-place)|perl\s+-\w*i|tee\s|git\s+(?:add|commit|apply|checkout|restore|rm|mv)|rm\s|mv\s|cp\s|touch\s|truncate\s|install\s)`)
	// pyWriteRe matches a script writing a file from an interpreter.
	pyWriteRe = regexp.MustCompile(`open\([^)]*['"][wax]\+?['"]|\.write(?:_text)?\(`)
	// quotedRe, noiseRedirRe and fileRedirRe find a redirection to a real file:
	// quoted text and redirections to /dev/null or between descriptors do not count.
	quotedRe     = regexp.MustCompile(`'[^']*'|"[^"]*"`)
	noiseRedirRe = regexp.MustCompile(`\d*>\s*/dev/null|&>\s*/dev/null|\d*>&\d`)
	fileRedirRe  = regexp.MustCompile(`>>?\s*[^\s>&|;]`)
)

func or(s, d string) string {
	if s == "" {
		return d
	}
	return s
}

// yn renders a boolean the way signatures show it.
func yn(v bool) string {
	if v {
		return "yes"
	}
	return "no"
}

// scan holds what scenarios share about the order of a transcript.
type scan struct {
	EditBeforeRoute   bool // something was changed before any route was stated
	VersionCheckFirst bool // the first tool call was `clue latest`
}

func scanTranscript(t Transcript) scan {
	var sc scan
	toolSeen, routeSeen := false, false
	for _, st := range t.Steps {
		switch st.Kind {
		case "tool":
			if !toolSeen {
				toolSeen = true
				sc.VersionCheckFirst = st.Tool == "Bash" && strings.Contains(st.Command, "clue latest")
			}
			if !routeSeen && mutates(st.Tool, st.Command) {
				sc.EditBeforeRoute = true
			}
		case "text":
			if routeRe.MatchString(st.Text) || looseRe.MatchString(st.Text) {
				routeSeen = true
			}
		}
	}
	return sc
}

func mutates(tool, cmd string) bool {
	switch tool {
	case "Edit", "Write", "NotebookEdit":
		return true
	case "Bash":
		return bashMutates(cmd)
	}
	return false
}

func bashMutates(cmd string) bool {
	if shellWriteRe.MatchString(cmd) || pyWriteRe.MatchString(cmd) {
		return true
	}
	bare := noiseRedirRe.ReplaceAllString(quotedRe.ReplaceAllString(cmd, "''"), " ")
	return fileRedirRe.MatchString(bare)
}

// statusPaths lists the paths in a `git status --porcelain` listing.
func statusPaths(post string) []string {
	var out []string
	for _, l := range strings.Split(post, "\n") {
		if len(l) > 3 {
			p := l[3:]
			if i := strings.Index(p, " -> "); i >= 0 {
				p = p[i+4:]
			}
			out = append(out, strings.Trim(p, `"`))
		}
	}
	return out
}

// touched reports whether the listing names a path with the given prefix.
func touched(post, prefix string) bool {
	for _, p := range statusPaths(post) {
		if strings.HasPrefix(p, prefix) {
			return true
		}
	}
	return false
}

// texts returns what the agent said, in order.
func texts(t Transcript) []string {
	var out []string
	for _, st := range t.Steps {
		if st.Kind == "text" {
			out = append(out, st.Text)
		}
	}
	return out
}

// commands returns the shell commands the agent ran, in order.
func commands(t Transcript) []string {
	var out []string
	for _, st := range t.Steps {
		if st.Kind == "tool" && st.Tool == "Bash" {
			out = append(out, st.Command)
		}
	}
	return out
}

// LeaksConfig reports whether the transcript shows host configuration the
// container should not have: anything that names CodeGraph, which only the
// maintainer's own instructions mention.
func LeaksConfig(raw string) bool {
	return strings.Contains(strings.ToLower(raw), "codegraph")
}
