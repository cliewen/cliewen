package corpus

import "sort"

// OpenChanges returns the change proposals in the current workspace. An open
// change is the strongest orientation signal because it records work already
// started on this branch rather than a choice of new work.
func OpenChanges(c *Corpus) []*Artifact {
	return artifactsByTypeAndStatus(c, "change", "open")
}

// ProposedGoals returns the repository inbox in stable path order. Proposed
// goals are choices for human and agent review, never actionable work.
func ProposedGoals(c *Corpus) []*Artifact {
	return artifactsByTypeAndStatus(c, "goal", "proposed")
}

func artifactsByTypeAndStatus(c *Corpus, artifactType, status string) []*Artifact {
	if c == nil {
		return nil
	}
	var out []*Artifact
	for _, artifact := range c.Artifacts {
		if artifact.Type == artifactType && artifact.Status == status {
			out = append(out, artifact)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Path < out[j].Path })
	return out
}
