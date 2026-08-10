package corpus

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// milestoneRowRe matches a milestone declared in a plan's milestone table:
// the ID is the first cell of the row. A mention in prose is a reference, not
// a declaration, so it must not make the ID look ambiguous.
var milestoneRowRe = regexp.MustCompile(`(?m)^\s*\|\s*(M-\d+)\s*\|`)

// DepthAll follows every outgoing edge to exhaustion: the whole reachable
// closure, which is what a slice meant before it carried a bound.
const DepthAll = -1

// ContextOptions bounds a slice.
type ContextOptions struct {
	// Depth is how many link hops beyond the root to include: 0 is the root
	// alone, 1 adds what it links to directly, DepthAll follows every edge to
	// exhaustion.
	//
	// Unlike this package's other options, the zero value is not "use the
	// default". A depth has a natural zero, and root-only is the safest thing
	// an unset bound can mean; the command supplies the default a human gets.
	Depth int
}

// Frontier is one artifact a slice reached but did not include, with the hop
// count at which it was reached. Reporting it is what keeps a bounded slice
// from being a silent omission: a reader sees the edge exists and can widen to
// it by name. PDR-034 bounds where reading starts, never what it may reach.
type Frontier struct {
	Artifact *Artifact
	Hops     int
}

// Context returns the deterministic outgoing-link slice rooted at id, bounded
// by opts.Depth, together with the artifacts beyond that bound.
// Artifact IDs resolve directly; milestone and acceptance-criterion IDs resolve
// to the plan or criteria artifact that declares them. The root is first and
// each breadth-first layer is ordered by repository-relative path; the frontier
// is ordered the same way, by hop count and then path.
//
// A link that resolves to no artifact, or to more than one, is reported as an
// unfollowed-edge issue instead of ending the slice: focused reading stays
// available while `clue validate` remains the judge of graph health. Only an
// unresolvable requested id is an error. An edge leaving an artifact the slice
// did not include is not reported — it describes a part of the graph the reader
// never sees, and naming it turns a bounded read into someone else's backlog.
func Context(c *Corpus, id string, opts ContextOptions) ([]*Artifact, []Frontier, []Issue, error) {
	owners := contextOwners(c)
	root, err := owners.resolve(id)
	if err != nil {
		return nil, nil, nil, err
	}
	included := func(hops int) bool { return opts.Depth == DepthAll || hops <= opts.Depth }

	var unfollowed []Issue
	var frontier []Frontier
	result := []*Artifact{root}
	seen := map[string]bool{root.Path: true}
	layer := []*Artifact{root}
	for hops := 1; len(layer) > 0; hops++ {
		// layer holds the artifacts at hops-1, so their edges reach hops.
		reportEdges := included(hops - 1)
		nextByPath := map[string]*Artifact{}
		for _, artifact := range layer {
			for _, link := range artifact.Links {
				target, err := owners.resolve(link)
				if err != nil {
					if reportEdges {
						unfollowed = append(unfollowed, Issue{artifact.Path, "link " + link + " not followed: " + err.Error()})
					}
					continue
				}
				if !seen[target.Path] {
					nextByPath[target.Path] = target
				}
			}
		}

		paths := make([]string, 0, len(nextByPath))
		for path := range nextByPath {
			paths = append(paths, path)
		}
		sort.Strings(paths)
		layer = layer[:0]
		for _, path := range paths {
			artifact := nextByPath[path]
			seen[path] = true
			if included(hops) {
				result = append(result, artifact)
			} else {
				frontier = append(frontier, Frontier{artifact, hops})
			}
			layer = append(layer, artifact)
		}
	}
	return result, frontier, unfollowed, nil
}

// contextIndex maps every resolvable identity — artifact IDs plus the
// milestone and acceptance-criterion IDs declared inside artifact bodies — to
// the artifacts that declare it, ordered by repository-relative path.
type contextIndex map[string][]*Artifact

func contextOwners(c *Corpus) contextIndex {
	byID := map[string]map[string]*Artifact{}
	declare := func(id string, a *Artifact) {
		if byID[id] == nil {
			byID[id] = map[string]*Artifact{}
		}
		byID[id][a.Path] = a
	}
	for _, artifact := range c.Artifacts {
		declare(artifact.ID, artifact)
		switch artifact.Type {
		case "plan":
			for _, row := range milestoneRowRe.FindAllStringSubmatch(artifact.Body, -1) {
				declare(row[1], artifact)
			}
		case "criteria":
			for _, line := range strings.Split(artifact.Body, "\n") {
				for _, id := range canonicalACIDsInLine(line) {
					declare(id, artifact)
				}
			}
		}
	}

	index := contextIndex{}
	for id, byPath := range byID {
		paths := make([]string, 0, len(byPath))
		for path := range byPath {
			paths = append(paths, path)
		}
		sort.Strings(paths)
		owners := make([]*Artifact, 0, len(paths))
		for _, path := range paths {
			owners = append(owners, byPath[path])
		}
		index[id] = owners
	}
	return index
}

func (index contextIndex) resolve(id string) (*Artifact, error) {
	owners := index[id]
	switch len(owners) {
	case 0:
		return nil, fmt.Errorf("ID %s not found", id)
	case 1:
		return owners[0], nil
	default:
		paths := make([]string, 0, len(owners))
		for _, owner := range owners {
			paths = append(paths, owner.Path)
		}
		return nil, fmt.Errorf("ID %s is ambiguous (declared by %v)", id, paths)
	}
}
