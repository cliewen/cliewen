package corpus

import (
	"sort"
	"strings"
)

// Milestone is one declared row in a plan's milestone table. The parser is
// intentionally limited to the table shape the validator already recognizes:
// an ID column and a Status column. Other columns are carried when their
// headers identify them, so a caller can show the exit criterion without
// imposing a new schema on existing plans.
type Milestone struct {
	Plan          *Artifact
	ID            string
	Name          string
	ExitCriterion string
	Status        string
	Evidence      string
	Row           int
}

// PlanMilestones returns declared milestone rows in source order. Fenced and
// indented examples are ignored by the same scanner used by validation.
func PlanMilestones(plan *Artifact) []Milestone {
	if plan == nil || plan.Type != "plan" {
		return nil
	}
	var out []Milestone
	idCol, statusCol := -1, -1
	nameCol, criterionCol, evidenceCol := -1, -1, -1
	blocks := blockScanner{verbatimOnly: true}
	var header []string
	for row, line := range strings.Split(plan.Body, "\n") {
		if blocks.next(line) {
			idCol, statusCol, nameCol, criterionCol, evidenceCol, header = -1, -1, -1, -1, -1, nil
			continue
		}
		t := strings.TrimSpace(line)
		if tableDelimRe.MatchString(t) {
			idCol, statusCol, nameCol, criterionCol, evidenceCol = -1, -1, -1, -1, -1
			for i, cell := range header {
				normal := strings.ToLower(strings.TrimSpace(cell))
				switch {
				case normal == "id":
					idCol = i
				case normal == "status":
					statusCol = i
				case strings.HasPrefix(normal, "milestone"):
					nameCol = i
				case strings.HasPrefix(normal, "exit criterion"):
					criterionCol = i
				case normal == "evidence":
					evidenceCol = i
				}
			}
			continue
		}
		if !strings.Contains(t, "|") {
			idCol, statusCol, nameCol, criterionCol, evidenceCol, header = -1, -1, -1, -1, -1, nil
			continue
		}
		cells := tableCells(t)
		if idCol < 0 || statusCol < 0 {
			header = cells
			continue
		}
		if idCol >= len(cells) || statusCol >= len(cells) {
			continue
		}
		id := strings.TrimSpace(cells[idCol])
		if !strings.HasPrefix(id, "M-") || !planItemRe.MatchString(id) {
			continue
		}
		cell := func(col int) string {
			if col < 0 || col >= len(cells) {
				return ""
			}
			return strings.TrimSpace(cells[col])
		}
		name := cell(nameCol)
		criterion := cell(criterionCol)
		if criterion == "" {
			criterion = name
		}
		out = append(out, Milestone{Plan: plan, ID: id, Name: name, ExitCriterion: criterion, Status: strings.Trim(strings.TrimSpace(cells[statusCol]), "`"), Evidence: cell(evidenceCol), Row: row + 1})
	}
	return out
}

// UnfinishedMilestones returns actionable candidates from active plans. A
// doing milestone is preferred to todo, then stable path and source row order
// make the result repeatable across clean contexts.
func UnfinishedMilestones(c *Corpus) []Milestone {
	if c == nil {
		return nil
	}
	var out []Milestone
	for _, artifact := range c.Artifacts {
		if artifact.Type != "plan" || artifact.Status != "active" {
			continue
		}
		for _, milestone := range PlanMilestones(artifact) {
			if milestone.Status == "todo" || milestone.Status == "doing" {
				out = append(out, milestone)
			}
		}
	}
	sort.SliceStable(out, func(i, j int) bool {
		priority := func(status string) int {
			if status == "doing" {
				return 0
			}
			return 1
		}
		if priority(out[i].Status) != priority(out[j].Status) {
			return priority(out[i].Status) < priority(out[j].Status)
		}
		if out[i].Plan.Path != out[j].Plan.Path {
			return out[i].Plan.Path < out[j].Plan.Path
		}
		return out[i].Row < out[j].Row
	})
	return out
}

// DraftUnfinishedMilestones reports unfinished rows in draft plans so a
// caller can explain why visible work is not currently actionable.
func DraftUnfinishedMilestones(c *Corpus) []Milestone {
	if c == nil {
		return nil
	}
	var out []Milestone
	for _, artifact := range c.Artifacts {
		if artifact.Type != "plan" || artifact.Status != "draft" {
			continue
		}
		for _, milestone := range PlanMilestones(artifact) {
			if milestone.Status == "todo" || milestone.Status == "doing" {
				out = append(out, milestone)
			}
		}
	}
	return out
}
