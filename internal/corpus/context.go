package corpus

import (
	"fmt"
	"sort"
)

// Context returns the deterministic outgoing-link closure rooted at id.
// Artifact IDs resolve directly; milestone and acceptance-criterion IDs resolve
// to the plan or criteria artifact that declares them. The root is first and
// each breadth-first layer is ordered by repository-relative path.
func Context(c *Corpus, id string) ([]*Artifact, error) {
	root, err := contextOwner(c, id)
	if err != nil {
		return nil, err
	}

	result := []*Artifact{root}
	seen := map[string]bool{root.Path: true}
	layer := []*Artifact{root}
	for len(layer) > 0 {
		nextByPath := map[string]*Artifact{}
		for _, artifact := range layer {
			for _, link := range artifact.Links {
				target, err := contextOwner(c, link)
				if err != nil {
					return nil, fmt.Errorf("%s links %s: %w", artifact.ID, link, err)
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
	return result, nil
}

func contextOwner(c *Corpus, id string) (*Artifact, error) {
	ownersByPath := map[string]*Artifact{}
	for _, artifact := range c.ByID[id] {
		ownersByPath[artifact.Path] = artifact
	}
	for _, artifact := range c.Artifacts {
		switch artifact.Type {
		case "plan":
			for _, declared := range milestoneRe.FindAllString(artifact.Body, -1) {
				if declared == id {
					ownersByPath[artifact.Path] = artifact
				}
			}
		case "criteria":
			for _, match := range acTagRe.FindAllStringSubmatch(artifact.Body, -1) {
				if match[1]+"-"+match[2] == id {
					ownersByPath[artifact.Path] = artifact
				}
			}
		}
	}

	paths := make([]string, 0, len(ownersByPath))
	for path := range ownersByPath {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	switch len(paths) {
	case 0:
		return nil, fmt.Errorf("ID %s not found", id)
	case 1:
		return ownersByPath[paths[0]], nil
	default:
		return nil, fmt.Errorf("ID %s is ambiguous (declared by %v)", id, paths)
	}
}
