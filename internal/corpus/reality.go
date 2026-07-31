package corpus

import (
	"sort"
	"strings"
)

// RealityGap is one capability contradicted by one or more incident analyses.
type RealityGap struct {
	Capability string
	Analyses   []string
}

// acceptanceCriterionOwners maps each live criterion ID to the capability
// owning its criteria artifact. It deliberately reads declarations only:
// evidence harvesting remains checkACTests' separate concern. "Live" is
// checkACTests' own definition — an active criteria file, minus the
// @retired tombstones — so a criterion in a draft file is not yet a link
// target and cannot carry a reality edge (ADR-035).
func acceptanceCriterionOwners(c *Corpus) map[string]string {
	capabilities := map[string]bool{}
	for _, a := range c.Artifacts {
		if a.Type == "capability" {
			capabilities[a.ID] = true
		}
	}
	out := map[string]string{}
	for _, a := range c.Artifacts {
		if a.Type != "criteria" || a.Status != "active" {
			continue
		}
		owner := ""
		for _, l := range a.Links {
			if capabilities[l] {
				owner = l
				break
			}
		}
		if owner == "" {
			continue
		}
		for _, line := range strings.Split(a.Body, "\n") {
			if strings.Contains(line, "@retired") {
				continue
			}
			for _, id := range canonicalACIDsInLine(line) {
				out[id] = owner
			}
		}
	}
	return out
}

// capabilityIDs is the set of capability IDs a link may name as a failed
// claim directly, rather than through one of its criteria.
func capabilityIDs(c *Corpus) map[string]bool {
	out := map[string]bool{}
	for _, a := range c.Artifacts {
		if a.Type == "capability" {
			out[a.ID] = true
		}
	}
	return out
}

func checkReality(c *Corpus) []Issue {
	owners := acceptanceCriterionOwners(c)
	capabilities := capabilityIDs(c)
	var issues []Issue
	for _, a := range c.Artifacts {
		rawReality, present := a.Fields["reality"]
		if !present {
			continue
		}
		reality, isString := rawReality.(string)
		if a.Type != "analysis" {
			issues = append(issues, Issue{a.Path, "reality is allowed only on analysis artifacts (ADR-035)"})
			continue
		}
		if !isString || reality != "contradicted" {
			issues = append(issues, Issue{a.Path, "reality must be contradicted (ADR-035)"})
			continue
		}
		hasFailedClaim := false
		for _, l := range a.Links {
			if capabilities[l] || owners[l] != "" {
				hasFailedClaim = true
				break
			}
		}
		if !hasFailedClaim {
			issues = append(issues, Issue{a.Path, "reality contradicted requires a links edge to the failed capability or live acceptance criterion (ADR-035)"})
		}
	}
	return issues
}

// RealityGaps derives the capabilities whose once-green claims an incident
// analysis later contradicted. No registry is written to disk.
func RealityGaps(c *Corpus) []RealityGap {
	owners := acceptanceCriterionOwners(c)
	capabilities := capabilityIDs(c)
	byCapability := map[string]map[string]bool{}
	for _, a := range c.Artifacts {
		if a.Type != "analysis" || a.Reality != "contradicted" {
			continue
		}
		for _, l := range a.Links {
			affected := owners[l]
			if capabilities[l] {
				affected = l
			}
			if affected == "" {
				continue
			}
			if byCapability[affected] == nil {
				byCapability[affected] = map[string]bool{}
			}
			byCapability[affected][a.ID] = true
		}
	}
	out := make([]RealityGap, 0, len(byCapability))
	for capability, analyses := range byCapability {
		ids := make([]string, 0, len(analyses))
		for id := range analyses {
			ids = append(ids, id)
		}
		sort.Strings(ids)
		out = append(out, RealityGap{Capability: capability, Analyses: ids})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Capability < out[j].Capability })
	return out
}
