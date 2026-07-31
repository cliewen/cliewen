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

// Context returns the deterministic outgoing-link closure rooted at id.
// Artifact IDs resolve directly; milestone and acceptance-criterion IDs resolve
// to the plan or criteria artifact that declares them. The root is first and
// each breadth-first layer is ordered by repository-relative path.
//
// A link that resolves to no artifact, or to more than one, is reported as an
// unfollowed-edge issue instead of ending the slice: focused reading stays
// available while `clue validate` remains the judge of graph health. Only an
// unresolvable requested id is an error.
func Context(c *Corpus, id string) ([]*Artifact, []Issue, error) {
	owners := contextOwners(c)
	root, err := owners.resolve(id)
	if err != nil {
		return nil, nil, err
	}

	var unfollowed []Issue
	result := []*Artifact{root}
	seen := map[string]bool{root.Path: true}
	layer := []*Artifact{root}
	for len(layer) > 0 {
		nextByPath := map[string]*Artifact{}
		for _, artifact := range layer {
			for _, link := range artifact.Links {
				target, err := owners.resolve(link)
				if err != nil {
					unfollowed = append(unfollowed, Issue{artifact.Path, "link " + link + " not followed: " + err.Error()})
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
			result = append(result, artifact)
			layer = append(layer, artifact)
		}
	}
	return result, unfollowed, nil
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
