package main

import (
	"regexp"
	"strings"
)

var (
	// sentenceRe splits on line breaks, semicolons, and sentence ends followed by space.
	sentenceRe = regexp.MustCompile(`\n|[.!?]\s+|;\s*`)
	typoRe     = regexp.MustCompile(`(?i)typo|readme|teh`)
	csvRe      = regexp.MustCompile(`(?i)csv|export`)
)

// routeState is the route an agent stated for each of the routing scenario's
// two requests.
type routeState struct{ Typo, Export string }

// observeRouting reads the routing scenario: which route was stated for the
// typo fix and for the new capability, and whether anything was edited first.
func observeRouting(t Transcript, post string, sc scan) map[string]string {
	var rs routeState
	for _, text := range texts(t) {
		if routeRe.MatchString(text) {
			routes(text, &rs)
		}
		if routeRe.MatchString(text) || looseRe.MatchString(text) {
			looseRoutes(text, &rs)
		}
	}
	return map[string]string{
		"typo": or(rs.Typo, "none"), "export": or(rs.Export, "none"),
		"edit-before-route": yn(sc.EditBeforeRoute), "readme-edited": yn(touched(post, "README.md")),
	}
}

// routes assigns each "Recommended route" mention to the request it is about,
// by looking at its own line and the line before it.
func routes(text string, rs *routeState) {
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
			rs.Export = route
		case "typo":
			rs.Typo = route
		}
	}
}

// looseRoutes fills a route the fixed phrase did not give, from any clause that
// names a route in plain words; a clause that names neither request takes the
// last request named before it.
func looseRoutes(text string, rs *routeState) {
	last := ""
	for _, line := range sentenceRe.Split(text, -1) {
		about := topic(line)
		if about == "" {
			about = last
		} else {
			last = about
		}
		route, ok := looseRoute(line)
		if !ok {
			continue
		}
		switch about {
		case "export":
			if rs.Export == "" {
				rs.Export = route
			}
		case "typo":
			if rs.Typo == "" {
				rs.Typo = route
			}
		}
	}
}

var negationRe = regexp.MustCompile(`(?i)\b(?:not|never|no|without|rather than|instead of)\b|n't`)

// looseRoute reads a route stated in plain words in one clause. A route that is
// negated before the word ("I didn't take the tracked route") means the other
// one, since there are only two.
func looseRoute(clause string) (string, bool) {
	loc := looseRe.FindStringSubmatchIndex(clause)
	if loc == nil {
		return "", false
	}
	route := strings.ToLower(clause[loc[2]:loc[3]])
	if negationRe.MatchString(clause[:loc[0]]) {
		if route == "tracked" {
			return "direct", true
		}
		return "tracked", true
	}
	return route, true
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
