package corpus

import (
	"path"
	"regexp"
	"sort"
	"strings"
)

// VisionPath is the one address a corpus's vision lives at (ADR-065). The
// identity is still what links name; the fixed path is an additional
// guarantee, so an agent orienting can read the direction without first
// scanning frontmatter to find out whether there is one.
const VisionPath = "docs/vision.md"

// VisionBootstrapMarker marks a scaffolded vision an agent must replace with
// the repository's own direction before validation can pass.
const VisionBootstrapMarker = "<!-- clue:vision:bootstrap -->"

// UseCaseDir is the folder optional use cases live in. A corpus with no use
// cases has no obligation here, including no obligation to have the folder.
const UseCaseDir = "docs/use-cases"

// useCaseSections are the structural headings a use case carries. They are
// the parts that make a journey readable as a journey; preconditions,
// alternative flows, and open questions are written when they carry meaning
// and are not required, because demanding them would produce headings with
// "none" under them.
var useCaseSections = []string{"Actors", "Trigger", "Main flow", "Outcome"}

var useCaseFileRe = regexp.MustCompile(`^UC-\d+-.+\.md$`)

// checkIntent validates the form of the intent artifacts and nothing about
// their meaning (ADR-065, ADR-066, ADR-067). Neither artifact is ever
// required to exist: a corpus with no vision and no use cases is valid, and
// no goal, capability, or criterion is judged for lacking either.
//
// Whether a vision states the right direction, whether a use case represents
// real users, and whether every actor has been found are human judgements.
// A check that appeared to make them would be worse than no check at all,
// because a reader would stop making them.
func checkIntent(c *Corpus) []Issue {
	var issues []Issue
	issues = append(issues, checkVision(c)...)
	issues = append(issues, checkUseCases(c)...)
	return issues
}

func checkVision(c *Corpus) []Issue {
	var issues []Issue
	var visions []*Artifact
	for _, a := range c.Artifacts {
		if a.Type == "vision" {
			visions = append(visions, a)
		}
	}
	if len(visions) > 1 {
		var paths []string
		for _, a := range visions {
			paths = append(paths, a.Path)
		}
		sort.Strings(paths)
		for _, a := range visions {
			issues = append(issues, Issue{a.Path, "a corpus has one vision, and this corpus has " + strings.Join(paths, ", ") + " — merge them into " + VisionPath + " (ADR-065)"})
		}
	}
	for _, a := range visions {
		if a.Path != VisionPath {
			issues = append(issues, Issue{a.Path, "a vision lives at " + VisionPath + " (ADR-065)"})
		}
	}
	// An artifact sitting at the vision's address that is not one is a
	// different defect from a vision in the wrong place, and it is reported
	// as itself: the repair is to move that artifact, not to move a vision.
	for _, a := range c.Artifacts {
		if a.Path == VisionPath && a.Type != "vision" && a.Type != "" {
			issues = append(issues, Issue{a.Path, "this address belongs to the corpus vision, and this artifact is a " + a.Type + " (ADR-065)"})
		}
	}
	// The bootstrap rule fires only on a file that exists. That is what lets
	// a corpus with no vision stay green while a repository that ran
	// clue init and stopped is told to finish (ADR-067).
	if text, ok := c.Contents[VisionPath]; ok && strings.Contains(text, VisionBootstrapMarker) {
		issues = append(issues, Issue{VisionPath, "vision is still the scaffold bootstrap — replace it with this repository's own direction, or delete it if the repository is not stating one"})
	}
	return issues
}

func checkUseCases(c *Corpus) []Issue {
	var issues []Issue
	for _, a := range c.Artifacts {
		if a.Type != "use-case" {
			continue
		}
		base := path.Base(a.Path)
		if path.Dir(a.Path) != UseCaseDir || !useCaseFileRe.MatchString(base) || !strings.HasPrefix(base, a.ID+"-") {
			issues = append(issues, Issue{a.Path, "use cases live in " + UseCaseDir + " with a UC-<number>-<slug>.md filename matching their id"})
		}
		// Link resolution is checkLinks's job for every artifact; what is
		// checked here is only that the journey declares the two ends it is
		// a journey between (ADR-066).
		var goal, capability bool
		for _, l := range a.Links {
			switch {
			case strings.HasPrefix(l, "G-"):
				goal = true
			case strings.HasPrefix(l, "CAP-"):
				capability = true
			}
		}
		if !goal {
			issues = append(issues, Issue{a.Path, "use case names no goal it serves — add the G-xxx identity to links (ADR-066)"})
		}
		if !capability {
			issues = append(issues, Issue{a.Path, "use case names no capability it crosses — add the CAP-xxx identities to links (ADR-066)"})
		}
		for _, section := range useCaseSections {
			if !strings.Contains(a.Body, "\n## "+section+"\n") && !strings.HasPrefix(a.Body, "## "+section+"\n") {
				issues = append(issues, Issue{a.Path, "use case has no \"## " + section + "\" section"})
			}
		}
	}
	return issues
}

// VisionState is what a corpus has to say about its own direction: the
// artifact if there is one, and whether a human has confirmed its meaning.
type VisionState struct {
	Present  bool
	ID       string
	Title    string
	Status   string
	Path     string
	Inferred bool // provenance: inferred — an agent drafted it and no human has confirmed it
}

// UseCaseState is one use case and the capabilities it crosses.
type UseCaseState struct {
	ID           string
	Title        string
	Status       string
	Path         string
	Capabilities []string
}

// IntentState is the derived answer to "what direction does this corpus
// state, and what journeys has it written down".
//
// It carries no ratio and no coverage figure, and that is a decision rather
// than an omission (PDR-054): a percentage over an optional artifact reads
// as a target, and the only way to move it is to write use cases nobody
// needs.
type IntentState struct {
	Vision   VisionState
	UseCases []UseCaseState
}

// Intent derives the intent state from a scanned corpus.
func Intent(c *Corpus) IntentState {
	var state IntentState
	for _, a := range c.Artifacts {
		switch a.Type {
		case "vision":
			if state.Vision.Present {
				continue // a second vision is checkVision's finding, not this report's guess
			}
			provenance, _ := a.Fields["provenance"].(string)
			state.Vision = VisionState{Present: true, ID: a.ID, Title: a.Title, Status: a.Status, Path: a.Path, Inferred: provenance == "inferred"}
		case "use-case":
			use := UseCaseState{ID: a.ID, Title: a.Title, Status: a.Status, Path: a.Path}
			for _, l := range a.Links {
				if strings.HasPrefix(l, "CAP-") {
					use.Capabilities = append(use.Capabilities, l)
				}
			}
			state.UseCases = append(state.UseCases, use)
		}
	}
	sort.Slice(state.UseCases, func(i, j int) bool { return state.UseCases[i].Path < state.UseCases[j].Path })
	return state
}

// UseCasesNaming returns the use cases whose links name id, ordered by
// repository-relative path.
//
// This is the one direction intent links do not run in (ADR-066). A use case
// names its goal and its capabilities; a capability names neither back, so
// the edge is written once and cannot drift. Reading it the other way is a
// scan restricted to one artifact type, and the caller emits names alone —
// no edge is followed and no content is expanded, so a bounded slice stays
// exactly as bounded as it was.
func UseCasesNaming(c *Corpus, id string) []*Artifact {
	var out []*Artifact
	for _, a := range c.Artifacts {
		if a.Type != "use-case" || a.ID == id {
			continue
		}
		for _, l := range a.Links {
			if l == id {
				out = append(out, a)
				break
			}
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Path < out[j].Path })
	return out
}
