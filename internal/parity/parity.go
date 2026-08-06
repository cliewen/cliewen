// Package parity implements the migration-proof parity contract (ADR-049,
// P-011 M-053): a pinned, human/agent-authored source manifest is compared
// against a target manifest clue derives deterministically from the current
// corpus and identity ledger, so a migration PR's claimed coverage can be
// checked against what the source actually declared.
package parity

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"sort"

	"github.com/cliewen/cliewen/internal/corpus"
	"github.com/cliewen/cliewen/internal/ledger"
	"gopkg.in/yaml.v3"
)

// Disposition names why a criterion carries no proof-class evidence in the
// source manifest: a deliberate @draft, Human, or retired outcome, each
// requiring a justification (ADR-049, PDR-024).
type Disposition string

const (
	DispositionDraft   Disposition = "draft"
	DispositionHuman   Disposition = "human"
	DispositionRetired Disposition = "retired"
)

// SourceEntry is one pinned source-side outcome. A proof outcome names one
// classified source reference; several proof outcomes may share an ID when a
// criterion has several directions or locations. An exclusion or disposition
// is instead the sole outcome for its ID (ADR-049).
type SourceEntry struct {
	ID               string      `yaml:"id"`
	ProofClass       string      `yaml:"proof-class,omitempty"`
	Direction        string      `yaml:"direction,omitempty"`
	EvidenceLocation string      `yaml:"evidence-location,omitempty"`
	Excluded         bool        `yaml:"excluded,omitempty"`
	Reason           string      `yaml:"reason,omitempty"`
	Disposition      Disposition `yaml:"disposition,omitempty"`
	Justification    string      `yaml:"justification,omitempty"`
	SourceLocation   string      `yaml:"disposition-source-location,omitempty"`
	PlanDoor         string      `yaml:"plan-door,omitempty"`
}

// SourceManifest is the pinned manifest a source mapping emits during the
// clue-extract rehearsal (ADR-049): one file per mapping run, naming the
// exact source revision and location it was read from.
type SourceManifest struct {
	SourceRevision string        `yaml:"source-revision"`
	SourceLocation string        `yaml:"source-location"`
	Entries        []SourceEntry `yaml:"entries"`
}

// LoadSourceManifest reads and parses a source manifest at path.
func LoadSourceManifest(path string) (SourceManifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return SourceManifest{}, err
	}
	var m SourceManifest
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	if err := decoder.Decode(&m); err != nil {
		return SourceManifest{}, fmt.Errorf("%s: %w", path, err)
	}
	if err := validateSourceManifest(m); err != nil {
		return SourceManifest{}, fmt.Errorf("%s: %w", path, err)
	}
	return m, nil
}

func validateSourceManifest(m SourceManifest) error {
	if m.SourceRevision == "" || m.SourceLocation == "" {
		return fmt.Errorf("source-revision and source-location are required")
	}
	byID := make(map[string][]SourceEntry, len(m.Entries))
	for _, e := range m.Entries {
		if e.ID == "" {
			return fmt.Errorf("entry id is required")
		}
		proof := e.ProofClass != "" || e.Direction != "" || e.EvidenceLocation != ""
		disposition := e.Disposition != "" || e.Justification != "" || e.SourceLocation != "" || e.PlanDoor != ""
		switch {
		case e.Excluded:
			if e.Reason == "" || proof || disposition {
				return fmt.Errorf("entry %q: excluded entries require only a reason", e.ID)
			}
		case disposition:
			if proof || e.Reason != "" || e.Disposition == "" || e.Justification == "" || e.SourceLocation == "" || e.PlanDoor == "" {
				return fmt.Errorf("entry %q: dispositions require a disposition, justification, disposition-source-location, and plan-door only", e.ID)
			}
			if e.Disposition != DispositionDraft && e.Disposition != DispositionHuman && e.Disposition != DispositionRetired {
				return fmt.Errorf("entry %q: unsupported disposition %q", e.ID, e.Disposition)
			}
		case proof:
			if e.ProofClass == "" || e.Direction == "" || e.EvidenceLocation == "" || e.Reason != "" {
				return fmt.Errorf("entry %q: proof entries require proof-class, direction, and evidence-location only", e.ID)
			}
		default:
			return fmt.Errorf("entry %q: expected a proof, exclusion, or disposition", e.ID)
		}
		byID[e.ID] = append(byID[e.ID], e)
	}
	for id, entries := range byID {
		if len(entries) < 2 {
			continue
		}
		proofClass := entries[0].ProofClass
		for _, e := range entries {
			if e.Excluded || e.Disposition != "" {
				return fmt.Errorf("entry %q: exclusions and dispositions cannot share an ID with another entry", id)
			}
			if e.ProofClass != proofClass {
				return fmt.Errorf("entry %q: proof rows must declare one proof-class", id)
			}
		}
	}
	return nil
}

// TargetEntry is one criterion's state as derived from the current corpus
// and ledger — never authored, always recomputed (ADR-049).
type TargetEntry struct {
	ID                string
	LedgerState       string // "live" | "reserved" | "retired", empty when no ledger entry
	SourceRevision    string // from the ledger entry, when the ID was imported
	ProofClass        string
	Directions        []string
	EvidenceLocations []string
	Draft             bool
	Human             bool
	Retired           bool
}

// TargetManifest is the corpus-derived comparison side of a parity run,
// keyed by criterion ID.
type TargetManifest struct {
	Entries   map[string]TargetEntry
	PlanDoors map[string]bool
}

var milestoneRowRe = regexp.MustCompile(`(?m)^\s*\|\s*(M-\d+)\s*\|`)

// DeriveTargetManifest scans root's corpus and ledger and derives one target
// entry per declared criterion, reusing the same declaration and evidence
// harvest clue validate already runs (corpus.AcceptanceEvidence) rather than
// re-parsing the tree a second way.
func DeriveTargetManifest(root string) (TargetManifest, error) {
	c, issues := corpus.Scan(root)
	if len(issues) > 0 {
		return TargetManifest{}, fmt.Errorf("corpus scan issues: %v", issues)
	}
	declared, locations := corpus.AcceptanceEvidence(c)

	var led *ledger.Ledger
	if ledger.Exists(root) {
		l, err := ledger.Load(root)
		if err != nil {
			return TargetManifest{}, err
		}
		led = l
	}

	entries := make(map[string]TargetEntry, len(declared))
	planDoors := map[string]bool{}
	for _, a := range c.Artifacts {
		if a.Type != "plan" {
			continue
		}
		for _, row := range milestoneRowRe.FindAllStringSubmatch(a.Body, -1) {
			planDoors[row[1]] = true
		}
	}
	for id, d := range declared {
		te := TargetEntry{
			ID:         id,
			ProofClass: d.TestType,
			Draft:      d.Draft,
			Human:      d.TestType == "Human",
			Retired:    d.Retired,
		}
		if led != nil {
			if e, ok := led.Lookup(id); ok {
				te.LedgerState = string(e.State)
				te.SourceRevision = e.SourceRevision
			}
		}
		dirSet, locationSet := map[string]bool{}, map[string]bool{}
		for _, ref := range locations[id] {
			// Only a classified reference (both type and direction known)
			// counts as the evidence ADR-049's orphaned-tag and
			// changed-evidence classes compare; an unclassified legacy
			// reference proves the criterion to checkACTests but says
			// nothing about direction or location parity.
			if ref.Type == "" || ref.Direction == "" {
				continue
			}
			locationSet[ref.Path] = true
			dirSet[ref.Direction] = true
		}
		for dir := range dirSet {
			te.Directions = append(te.Directions, dir)
		}
		for location := range locationSet {
			te.EvidenceLocations = append(te.EvidenceLocations, location)
		}
		sort.Strings(te.Directions)
		sort.Strings(te.EvidenceLocations)
		entries[id] = te
	}
	return TargetManifest{Entries: entries, PlanDoors: planDoors}, nil
}

// Finding classes, matching M-053's five required failure classes (ADR-049).
const (
	ClassMissingCriterion         = "missing-criterion"
	ClassOrphanedTag              = "orphaned-tag"
	ClassChangedEvidence          = "changed-evidence"
	ClassStaleFingerprint         = "stale-fingerprint"
	ClassUnjustifiedDisposition   = "unjustified-disposition"
	ClassUnaccountableDisposition = "unaccountable-disposition"
)

// Finding is one parity disagreement between the source and target manifest.
type Finding struct {
	Class  string
	ID     string
	Detail string
}

// Report is Compare's deterministic result: a clean report has no findings,
// and repeated runs against the same inputs produce the same findings in the
// same order, so a CI artifact built from it is reproducible (ADR-049).
type Report struct {
	Findings []Finding
	Deferred int
}

// Failed reports whether the parity run must exit non-zero.
func (r Report) Failed() bool { return len(r.Findings) > 0 }

// Compare diffs source against target and reports every unmatched or
// altered entry, classified by the five failure classes ADR-049 requires.
// The report stays derived: Compare never mutates source, target, or any
// file on disk.
func Compare(source SourceManifest, target TargetManifest) Report {
	bySource := make(map[string][]SourceEntry, len(source.Entries))
	for _, e := range source.Entries {
		bySource[e.ID] = append(bySource[e.ID], e)
	}

	sourceIDs := make([]string, 0, len(bySource))
	for id := range bySource {
		sourceIDs = append(sourceIDs, id)
	}
	sort.Strings(sourceIDs)

	var findings []Finding
	deferred := map[string]bool{}
	for _, id := range sourceIDs {
		entries := bySource[id]
		se := entries[0]
		if se.Excluded {
			continue
		}
		te, ok := target.Entries[id]
		if !ok {
			findings = append(findings, Finding{ClassMissingCriterion, id, "source entry has no matching target criterion"})
			continue
		}
		if te.SourceRevision != "" && te.SourceRevision != source.SourceRevision {
			findings = append(findings, Finding{ClassStaleFingerprint, id, fmt.Sprintf("manifest source-revision %q disagrees with ledger revision %q", source.SourceRevision, te.SourceRevision)})
		}
		if se.Disposition != "" {
			deferred[id] = true
			if se.Justification == "" || !matchesDisposition(se.Disposition, te) {
				findings = append(findings, Finding{ClassUnjustifiedDisposition, id, fmt.Sprintf("source disposition %q does not match the target's draft, Human, or retired state", se.Disposition)})
			}
			if target.PlanDoors != nil && !target.PlanDoors[se.PlanDoor] {
				findings = append(findings, Finding{ClassUnaccountableDisposition, id, fmt.Sprintf("plan door %q is not a milestone in the target corpus", se.PlanDoor)})
			}
			continue
		}
		if te.Draft || te.Human || te.Retired {
			findings = append(findings, Finding{ClassUnjustifiedDisposition, id, "target is @draft, Human, or retired but the source manifest declares a proof class instead of a disposition"})
			continue
		}
		sourceDirections, sourceLocations := sourceEvidence(entries)
		mismatch := se.ProofClass != te.ProofClass || !sameStrings(sourceDirections, te.Directions) || !sameStrings(sourceLocations, te.EvidenceLocations)
		if mismatch {
			findings = append(findings, Finding{ClassChangedEvidence, id, fmt.Sprintf("source proof-class=%q directions=%v locations=%v; target proof-class=%q directions=%v locations=%v", se.ProofClass, sourceDirections, sourceLocations, te.ProofClass, te.Directions, te.EvidenceLocations)})
		}
	}

	targetIDs := make([]string, 0, len(target.Entries))
	for id := range target.Entries {
		targetIDs = append(targetIDs, id)
	}
	sort.Strings(targetIDs)
	for _, id := range targetIDs {
		if _, present := bySource[id]; present {
			continue // every present entry, excluded or not, was already handled above
		}
		te := target.Entries[id]
		switch {
		case te.Draft || te.Human || te.Retired:
			findings = append(findings, Finding{ClassUnjustifiedDisposition, id, "target is @draft, Human, or retired with no source manifest entry"})
		case len(te.EvidenceLocations) > 0:
			findings = append(findings, Finding{ClassOrphanedTag, id, "target carries classified evidence with no source manifest entry"})
		}
	}

	sort.Slice(findings, func(i, j int) bool {
		if findings[i].ID != findings[j].ID {
			return findings[i].ID < findings[j].ID
		}
		return findings[i].Class < findings[j].Class
	})
	return Report{Findings: findings, Deferred: len(deferred)}
}

func matchesDisposition(d Disposition, target TargetEntry) bool {
	return (d == DispositionDraft && target.Draft) || (d == DispositionHuman && target.Human) || (d == DispositionRetired && target.Retired)
}

func sourceEvidence(entries []SourceEntry) (directions, locations []string) {
	directionsSeen, locationsSeen := map[string]bool{}, map[string]bool{}
	for _, entry := range entries {
		directionsSeen[entry.Direction] = true
		locationsSeen[entry.EvidenceLocation] = true
	}
	for direction := range directionsSeen {
		directions = append(directions, direction)
	}
	for location := range locationsSeen {
		locations = append(locations, location)
	}
	sort.Strings(directions)
	sort.Strings(locations)
	return directions, locations
}

func sameStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
