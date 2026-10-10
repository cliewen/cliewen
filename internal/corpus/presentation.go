package corpus

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

// ReferenceNames is a local, read-only presentation index. It never resolves
// Cliewen source rules or foreign references using an adopter's identities.
type ReferenceNames struct {
	names       map[string][]string
	owners      map[string][]*Artifact
	localIssues map[Issue][]string
}

// HumanLine removes terminal controls and normalizes whitespace for a label.
// Stored metadata and machine-facing identifiers are not rewritten.
func HumanLine(value string) string {
	value = strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return ' '
		}
		return r
	}, value)
	return strings.Join(strings.Fields(value), " ")
}

// HumanReference renders a supplied name with its secondary identity.
func HumanReference(name, id string) string {
	name, id = HumanLine(name), HumanLine(id)
	if name == "" || name == id {
		name = "unnamed reference"
	}
	return fmt.Sprintf("%s (%s)", name, id)
}

func NewReferenceNames(c *Corpus) *ReferenceNames {
	n := &ReferenceNames{names: map[string][]string{}, owners: map[string][]*Artifact{}}
	if c == nil {
		return n
	}
	n.localIssues = map[Issue][]string{}
	for issue, ids := range c.localIssueReferences {
		n.localIssues[issue] = append([]string(nil), ids...)
	}
	for _, a := range c.Artifacts {
		if a.ID != "" {
			n.names[a.ID] = append(n.names[a.ID], a.Title)
		}
		n.owners[a.Path] = append(n.owners[a.Path], a)
		if a.Type == "plan" {
			for _, m := range PlanMilestones(a) {
				n.names[m.ID] = append(n.names[m.ID], m.Name)
			}
		}
	}
	// Parse declarations only: naming does not load or validate evidence again.
	declarations, _ := criterionDeclarations(c)
	for id, d := range declarations {
		n.names[id] = append(n.names[id], d.name)
		if d.ambiguous {
			n.names[id] = append(n.names[id], d.name)
		}
	}
	// Context also recognizes legacy prose milestone references. A known
	// identity without a dedicated name is unnamed, rather than invented.
	for id, owners := range contextOwners(c) {
		if _, present := n.names[id]; present {
			continue
		}
		for range owners {
			n.names[id] = append(n.names[id], "")
		}
	}

	return n
}

func (n *ReferenceNames) Label(id string) string {
	if strings.HasPrefix(id, "clue:") || strings.Contains(id, "://") {
		return HumanReference("external reference, not resolved locally", id)
	}
	names := n.names[id]
	switch len(names) {
	case 0:
		return HumanReference("unknown reference", id)
	case 1:
		return HumanReference(names[0], id)
	default:
		return HumanReference("ambiguous reference", id)
	}
}

// Issue names only references explicitly registered by the producing check.
// It does not infer reference origin by finding ID-looking text in a message.
// Raw Issue values, their order and validation outcomes remain unchanged.
func (n *ReferenceNames) Issue(issue Issue) string {
	message := HumanLine(issue.Msg)
	ids := append([]string(nil), n.localIssues[issue]...)
	sort.Slice(ids, func(i, j int) bool { return len(ids[i]) > len(ids[j]) || len(ids[i]) == len(ids[j]) && ids[i] < ids[j] })
	// Locate tokens in the original message; never rescan inserted names.
	type replacement struct {
		start, end int
		value      string
	}
	var changes []replacement
	for _, id := range ids {
		pattern := regexp.MustCompile(`(^|[^A-Za-z0-9_-])(` + regexp.QuoteMeta(id) + `)([^A-Za-z0-9_-]|$)`)
		match := pattern.FindStringSubmatchIndex(message)
		if match == nil {
			continue
		}
		changes = append(changes, replacement{match[4], match[5], n.Label(id)})
	}
	sort.Slice(changes, func(i, j int) bool { return changes[i].start > changes[j].start })
	for _, change := range changes {
		message = message[:change.start] + change.value + message[change.end:]
	}
	result := HumanLine(issue.Path) + ": " + message
	if owners := n.owners[issue.Path]; len(owners) == 1 && len(ids) == 0 && owners[0].ID != "" {
		result += "\n  Document: " + n.Label(owners[0].ID)
	}
	return result
}

// localIssue records provenance separately from the stable raw diagnostic.
func localIssue(c *Corpus, path, message string, ids ...string) Issue {
	issue := Issue{path, message}
	if c.localIssueReferences == nil {
		c.localIssueReferences = map[Issue][]string{}
	}
	for _, id := range ids {
		if id == "" {
			continue
		}
		present := false
		for _, existing := range c.localIssueReferences[issue] {
			if existing == id {
				present = true
				break
			}
		}
		if !present {
			c.localIssueReferences[issue] = append(c.localIssueReferences[issue], id)
		}
	}
	return issue
}

func registerLocalIssues(c *Corpus, issues []Issue, ids ...string) []Issue {
	for _, issue := range issues {
		localIssue(c, issue.Path, issue.Msg, ids...)
	}
	return issues
}
