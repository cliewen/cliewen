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
// evidence harvesting remains checkACTests' separate concern.
func acceptanceCriterionOwners(c *Corpus) map[string]string {
	capabilities := map[string]bool{}
	for _, a := range c.Artifacts {
		if a.Type == "capability" {
			capabilities[a.ID] = true
		}
	}
	out := map[string]string{}
	for _, a := range c.Artifacts {
		if a.Type != "criteria" {
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
			for _, match := range acTagRe.FindAllStringSubmatch(line, -1) {
				out[match[1]+"-"+match[2]] = owner
			}
		}
	}
	return out
}

func checkReality(c *Corpus) []Issue {
	owners := acceptanceCriterionOwners(c)
	capabilities := map[string]bool{}
	for _, a := range c.Artifacts {
		if a.Type == "capability" {
			capabilities[a.ID] = true
		}
	}
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
	capabilities := map[string]bool{}
	for _, a := range c.Artifacts {
		if a.Type == "capability" {
			capabilities[a.ID] = true
		}
	}
	byCap := map[string]map[string]bool{}
	for _, a := range c.Artifacts {
		if a.Type != "analysis" || a.Reality != "contradicted" {
			continue
		}
		for _, l := range a.Links {
			cap := ""
			if capabilities[l] {
				cap = l
			} else {
				cap = owners[l]
			}
			if cap == "" {
				continue
			}
			if byCap[cap] == nil {
				byCap[cap] = map[string]bool{}
			}
			byCap[cap][a.ID] = true
		}
	}
	out := make([]RealityGap, 0, len(byCap))
	for cap, analyses := range byCap {
		ids := make([]string, 0, len(analyses))
		for id := range analyses {
			ids = append(ids, id)
		}
		sort.Strings(ids)
		out = append(out, RealityGap{Capability: cap, Analyses: ids})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Capability < out[j].Capability })
	return out
}
