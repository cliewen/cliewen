package corpus

import (
	"strings"

	"github.com/cliewen/cliewen/internal/role"
)

// The ships-generic-vs-repo-local boundary (ADR-013), made checkable in the
// one repository that can be checked (ADR-062).
//
// Cliewen's own repository and every adopter carry the same corpus shape, so
// a rule written here can bind adopter behaviour and never reach the surface
// an adopter receives — or reach it and never be recorded — with nothing
// failing either way. A record states whose behaviour it binds, and in the
// source repository a record binding adopters must name a carrier on the
// shipped surface.
//
// The declaration is opt-in. Deriving "this text binds adopters" from prose
// would be a judgment, and the judge reads state (ADR-044); requiring the
// field on every existing record would turn a corpus green yesterday red
// today for a field its authors could not have written.

// shippedSurfaces are the paths whose contents reach an adopter: the
// canonical skill sources that generate the managed skills, and the
// templates clue init and clue scaffold materialize.
var shippedSurfaces = []string{
	"internal/skills/source/",
	"internal/scaffold/templates/",
}

// bindsValues are the audiences a record may declare.
var bindsValues = map[string]bool{"adopter": true, "repo": true}

func checkBoundary(c *Corpus) []Issue {
	var issues []Issue
	for _, a := range c.Artifacts {
		if _, present := a.Fields["binds"]; !present {
			continue
		}
		if a.Type != "decision" && a.Type != "constraint" {
			issues = append(issues, Issue{a.Path, "binds is allowed only on decision and constraint artifacts (ADR-062)"})
			continue
		}
		if !bindsValues[a.Binds] {
			issues = append(issues, Issue{a.Path, "binds must be adopter or repo (ADR-062)"})
		}
	}
	// A marker that cannot be read is reported, never defaulted. The role is
	// the only switch on the rule below, so swallowing the error would let a
	// one-character typo turn off a validator rule with no output anywhere —
	// exactly the silent drift the marker exists to end.
	declared, _, err := role.Load(c.Root)
	if err != nil {
		return append(issues, Issue{role.DefaultPath, "role marker cannot be read, so the repository's role is unknown: " + err.Error() + " (ADR-062)"})
	}
	// Only the source repository ships anything, so only it can be asked
	// whether an adopter-binding rule reached the shipped surface. An
	// adopter's corpus, and an undeclared one, are never judged by this.
	if declared != role.Source {
		return issues
	}
	for _, a := range c.Artifacts {
		if a.Binds != "adopter" {
			continue
		}
		if a.Type != "decision" && a.Type != "constraint" {
			continue
		}
		if namesShippedSurface(a.Body) {
			continue
		}
		issues = append(issues, Issue{a.Path, "binds adopter but names no carrier on the shipped surface; name a path under " + strings.Join(shippedSurfaces, " or ") + " (ADR-062)"})
	}
	return issues
}

// namesShippedSurface reports whether body cites a path whose contents
// reach an adopter.
func namesShippedSurface(body string) bool {
	for _, prefix := range shippedSurfaces {
		if strings.Contains(body, prefix) {
			return true
		}
	}
	return false
}
