package corpus

import "sort"

// CapabilityCoverage is one capability's derived proof state (ADR-033):
// computed fresh from the corpus and the evidence tree on every call, never
// written to disk, so it cannot go stale the way a committed registry would.
type CapabilityCoverage struct {
	Capability string // capability ID
	State      string // "covered", "partial", or "gap"
}

// Coverage rolls each declared, non-retired acceptance criterion's evidence
// state up to its owning capability. A criterion counts as satisfied when
// its declared proof class has the evidence ADR-032/ADR-033 requires — a
// Human declaration, a single-direction reference, a classified pair, or a
// legacy unclassified reference — and unsatisfied when it is @draft,
// malformed, or missing evidence. A capability is "covered" when every one
// of its criteria is satisfied, "gap" when none are (including a capability
// with no active criteria at all), and "partial" otherwise.
func Coverage(c *Corpus) []CapabilityCoverage {
	declared, classified, tested, _ := harvestACs(c)

	capabilities := map[string]bool{}
	for _, a := range c.Artifacts {
		if a.Type == "capability" {
			capabilities[a.ID] = true
		}
	}
	capOf := map[string]string{}
	for _, a := range c.Artifacts {
		if a.Type != "criteria" {
			continue
		}
		for _, l := range a.Links {
			if capabilities[l] {
				capOf[a.Path] = l
				break
			}
		}
	}

	total, satisfied := map[string]int{}, map[string]int{}
	for ac, d := range declared {
		if d.status != "active" || d.retired {
			continue
		}
		cap, ok := capOf[d.path]
		if !ok {
			continue
		}
		total[cap]++
		switch {
		case d.draft:
			// Not yet proven by design (ADR-033); does not count as satisfied.
		case d.testType == "Human":
			satisfied[cap]++
		case d.testType == "":
			if tested[ac] {
				satisfied[cap]++
			}
		case d.single:
			if len(classified[ac][d.testType]) > 0 {
				satisfied[cap]++
			}
		default:
			cov := classified[ac][d.testType]
			if cov["positive"] && cov["negative"] {
				satisfied[cap]++
			}
		}
	}

	out := make([]CapabilityCoverage, 0, len(capabilities))
	for cap := range capabilities {
		t, s := total[cap], satisfied[cap]
		state := "gap"
		switch {
		case t > 0 && s == t:
			state = "covered"
		case s > 0:
			state = "partial"
		}
		out = append(out, CapabilityCoverage{Capability: cap, State: state})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Capability < out[j].Capability })
	return out
}
