package corpus

import "sort"

// ActivationBlocker is a high-cost inferred non-decision joined to an active
// capability by the capability's one-edge activation slice (ADR-035).
type ActivationBlocker struct {
	Artifact   *Artifact
	Capability *Artifact
}

// ProvenanceSummary is the actionable unverified-meaning view. Decisions stay
// non-blocking under merge-binds/approval-signs but are no longer invisible.
type ProvenanceSummary struct {
	Blockers         []ActivationBlocker
	BlockerArtifacts []*Artifact
	Decisions        []*Artifact
}

// ProvenanceBacklog derives activation blockers and unsigned decisions
// directly from the corpus graph; it is never persisted.
func ProvenanceBacklog(c *Corpus) ProvenanceSummary {
	var out ProvenanceSummary
	for _, a := range c.Artifacts {
		if a.Type == "decision" && a.Status == "inferred" {
			out.Decisions = append(out.Decisions, a)
		}
	}
	for _, capability := range c.Artifacts {
		if capability.Type != "capability" || capability.Status != "active" {
			continue
		}
		for _, a := range c.Artifacts {
			if !isHighCostInferred(a) || !inActivationSlice(capability, a) {
				continue
			}
			out.Blockers = append(out.Blockers, ActivationBlocker{Artifact: a, Capability: capability})
		}
	}
	seenBlocker := map[*Artifact]bool{}
	for _, blocker := range out.Blockers {
		if !seenBlocker[blocker.Artifact] {
			seenBlocker[blocker.Artifact] = true
			out.BlockerArtifacts = append(out.BlockerArtifacts, blocker.Artifact)
		}
	}
	sort.Slice(out.Blockers, func(i, j int) bool {
		if out.Blockers[i].Capability.ID != out.Blockers[j].Capability.ID {
			return out.Blockers[i].Capability.ID < out.Blockers[j].Capability.ID
		}
		return out.Blockers[i].Artifact.ID < out.Blockers[j].Artifact.ID
	})
	sort.Slice(out.Decisions, func(i, j int) bool { return out.Decisions[i].ID < out.Decisions[j].ID })
	sort.Slice(out.BlockerArtifacts, func(i, j int) bool { return out.BlockerArtifacts[i].ID < out.BlockerArtifacts[j].ID })
	return out
}

func isHighCostInferred(a *Artifact) bool {
	p, _ := a.Fields["provenance"].(string)
	return a.Type != "decision" && p == "inferred" && a.ReversalCost == "high"
}

func inActivationSlice(capability, a *Artifact) bool {
	if capability == a {
		return true
	}
	for _, l := range capability.Links {
		if l == a.ID {
			return true
		}
	}
	for _, l := range a.Links {
		if l == capability.ID {
			return true
		}
	}
	return false
}
