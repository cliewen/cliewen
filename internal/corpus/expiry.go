package corpus

import "sort"

// Expiry derives which analyses have done their work (PDR-052).
//
// A spike exists to retire a risk. When its findings reach a durable
// artifact it is spent, and keeping it costs every later reader a file that
// records a measurement true only at the revision it names. Retirement is
// deletion (ADR-034), so the corpus has no status for "spent" to rest in —
// what it needs instead is a way to say which spikes are ready to go.
//
// Two conditions are required, and the second is the one that cannot be
// inferred. A completed plan says the campaign ended; it says nothing about
// whether the spike's findings ever reached durable form. In a corpus whose
// plans are all completed, the plan test alone would sweep up every
// analysis at once. The carried-by field is the half a human declares.
//
// Nothing here fails a corpus. An unexpired analysis is not invalid, and
// the judge reads state rather than judgment (ADR-044): this derivation
// feeds a migration notice, and the deletion itself stays a reviewed
// change.

// SpentAnalysis is one analysis whose findings a durable artifact now
// carries, and whose plan is complete.
type SpentAnalysis struct {
	ID        string
	Path      string
	Plans     []string // the completed plans it served
	CarriedBy []string // the durable artifacts now carrying its findings
}

// checkCarriedBy keeps the field honest: it belongs only to analysis, and
// every ID it names must resolve to something a reader can actually open.
// A carrier that does not resolve is worse than an absent field, because it
// claims the findings survived somewhere they did not.
func checkCarriedBy(c *Corpus) []Issue {
	var issues []Issue
	for _, a := range c.Artifacts {
		if _, present := a.Fields["carried-by"]; !present {
			continue
		}
		if a.Type != "analysis" {
			issues = append(issues, Issue{a.Path, "carried-by is allowed only on analysis artifacts (PDR-052)"})
			continue
		}
		if len(a.CarriedBy) == 0 {
			issues = append(issues, Issue{a.Path, "carried-by must name at least one durable artifact, or be omitted (PDR-052)"})
			continue
		}
		for _, id := range a.CarriedBy {
			if id == a.ID {
				issues = append(issues, Issue{a.Path, "carried-by names the analysis itself (PDR-052)"})
				continue
			}
			target, ok := c.ByID[id]
			if !ok || len(target) == 0 {
				issues = append(issues, Issue{a.Path, "carried-by names " + id + ", which no artifact declares (PDR-052)"})
				continue
			}
			if target[0].Type == "analysis" {
				issues = append(issues, Issue{a.Path, "carried-by names " + id + ", which is another analysis; findings are carried by durable artifacts (PDR-052)"})
			}
		}
	}
	return issues
}

// SpentAnalyses derives the analyses ready for a reviewed retirement: every
// plan they serve is completed, and they name a durable carrier that
// resolves. No registry is written to disk.
func SpentAnalyses(c *Corpus) []SpentAnalysis {
	plans := map[string]*Artifact{}
	for _, a := range c.Artifacts {
		if a.Type == "plan" {
			plans[a.ID] = a
		}
	}
	live := map[string]bool{}
	for id, as := range c.ByID {
		if len(as) > 0 {
			live[id] = true
		}
	}
	var out []SpentAnalysis
	for _, a := range c.Artifacts {
		if a.Type != "analysis" || a.Status != "active" || len(a.CarriedBy) == 0 {
			continue
		}
		// A carrier that does not resolve is reported by checkCarriedBy,
		// not silently treated as proof the findings survived.
		carriers := make([]string, 0, len(a.CarriedBy))
		for _, id := range a.CarriedBy {
			if live[id] {
				carriers = append(carriers, id)
			}
		}
		if len(carriers) != len(a.CarriedBy) {
			continue
		}
		var served []string
		complete := true
		for _, l := range a.Links {
			p, ok := plans[l]
			if !ok {
				continue
			}
			served = append(served, p.ID)
			if p.Status != "completed" {
				complete = false
			}
		}
		// An analysis serving no plan is not swept up by a plan test it
		// never entered; a human retires it on its own terms.
		if !complete || len(served) == 0 {
			continue
		}
		sort.Strings(served)
		out = append(out, SpentAnalysis{ID: a.ID, Path: a.Path, Plans: served, CarriedBy: carriers})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Path < out[j].Path })
	return out
}
