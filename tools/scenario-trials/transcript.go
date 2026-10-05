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
	routeRe = regexp.MustCompile(`(?i)recommended route(?:\s+for\s+[^:\n]*)?:?\W*(direct|tracked|simple|full)\b`)
	// looseRe is any statement of a route in plain words, such as "as a direct change".
	looseRe = regexp.MustCompile(`(?i)\b(direct|tracked)\b`)
	// sentenceRe splits on line breaks, semicolons, and sentence ends followed by space.
	sentenceRe = regexp.MustCompile(`\n|[.!?]\s+|;\s*`)
	typoRe     = regexp.MustCompile(`(?i)typo|readme|teh`)
	csvRe      = regexp.MustCompile(`(?i)csv|export`)
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

// Check reads one run's transcript and the post-run working-tree listing.
func Check(t Transcript, postStatus string) Outcome {
	o := Outcome{Model: t.Model, AgentVersion: t.AgentVersion, MCPServers: t.MCPServers,
		Turns: t.Turns, CostUSD: t.CostUSD, DurationMS: t.DurationMS, Finished: t.Finished}
	toolSeen := false
	routeSeen := false
	for _, st := range t.Steps {
		switch st.Kind {
		case "tool":
			if !toolSeen {
				toolSeen = true
				o.VersionCheckFirst = st.Tool == "Bash" && strings.Contains(st.Command, "clue latest")
			}
			if !routeSeen && mutates(st.Tool, st.Command) {
				o.EditBeforeRoute = true
			}
		case "text":
			if routeRe.MatchString(st.Text) {
				o.Recommended = true
				routes(st.Text, &o)
			}
			if routeRe.MatchString(st.Text) || looseRe.MatchString(st.Text) {
				routeSeen = true
				looseRoutes(st.Text, &o)
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
	last := ""
	for _, line := range sentenceRe.Split(text, -1) {
		about := topic(line)
		if about == "" {
			about = last
		} else {
			last = about
		}
		m := looseRe.FindStringSubmatch(line)
		if m == nil {
			continue
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
