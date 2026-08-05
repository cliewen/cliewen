// Package parity implements the migration-proof parity contract (ADR-049,
// P-011 M-053): a pinned, human/agent-authored source manifest is compared
// against a target manifest clue derives deterministically from the current
// corpus and identity ledger, so a migration PR's claimed coverage can be
// checked against what the source actually declared.
package parity

import (
	"fmt"
	"os"
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

// SourceEntry is one criterion's pinned source-side state: a declared proof
// class and evidence location, a declared exclusion, or a disposition
// justification — exactly one outcome per entry (ADR-049).
type SourceEntry struct {
	ID               string      `yaml:"id"`
	ProofClass       string      `yaml:"proof-class,omitempty"`
	Direction        string      `yaml:"direction,omitempty"`
	EvidenceLocation string      `yaml:"evidence-location,omitempty"`
	Excluded         bool        `yaml:"excluded,omitempty"`
	Reason           string      `yaml:"reason,omitempty"`
	Disposition      Disposition `yaml:"disposition,omitempty"`
	Justification    string      `yaml:"justification,omitempty"`
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
	if err := yaml.Unmarshal(data, &m); err != nil {
		return SourceManifest{}, fmt.Errorf("%s: %w", path, err)
	}
	return m, nil
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
	Entries map[string]TargetEntry
}

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
		dirSet := map[string]bool{}
		for _, ref := range locations[id] {
			// Only a classified reference (both type and direction known)
			// counts as the evidence ADR-049's orphaned-tag and
			// changed-evidence classes compare; an unclassified legacy
			// reference proves the criterion to checkACTests but says
			// nothing about direction or location parity.
			if ref.Type == "" || ref.Direction == "" {
				continue
			}
			te.EvidenceLocations = append(te.EvidenceLocations, ref.Path)
			dirSet[ref.Direction] = true
		}
		for dir := range dirSet {
			te.Directions = append(te.Directions, dir)
		}
		sort.Strings(te.Directions)
		sort.Strings(te.EvidenceLocations)
		entries[id] = te
	}
	return TargetManifest{Entries: entries}, nil
}

// Finding classes, matching M-053's five required failure classes (ADR-049).
const (
	ClassMissingCriterion       = "missing-criterion"
	ClassOrphanedTag            = "orphaned-tag"
	ClassChangedEvidence        = "changed-evidence"
	ClassStaleFingerprint       = "stale-fingerprint"
	ClassUnjustifiedDisposition = "unjustified-disposition"
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
}

// Failed reports whether the parity run must exit non-zero.
func (r Report) Failed() bool { return len(r.Findings) > 0 }

// Compare diffs source against target and reports every unmatched or
// altered entry, classified by the five failure classes ADR-049 requires.
// The report stays derived: Compare never mutates source, target, or any
// file on disk.
func Compare(source SourceManifest, target TargetManifest) Report {
	bySource := make(map[string]SourceEntry, len(source.Entries))
	for _, e := range source.Entries {
		bySource[e.ID] = e
	}

	sourceIDs := make([]string, 0, len(bySource))
	for id := range bySource {
		sourceIDs = append(sourceIDs, id)
	}
	sort.Strings(sourceIDs)

	var findings []Finding
	for _, id := range sourceIDs {
		se := bySource[id]
		if se.Excluded {
			continue
		}
		te, ok := target.Entries[id]
		if !ok {
			findings = append(findings, Finding{ClassMissingCriterion, id, "source entry has no matching target criterion"})
			continue
		}
		if se.Disposition != "" {
			if se.Justification == "" {
				findings = append(findings, Finding{ClassUnjustifiedDisposition, id, "source entry declares disposition " + string(se.Disposition) + " with no justification"})
			}
			continue
		}
		if te.SourceRevision != "" && te.SourceRevision != source.SourceRevision {
			findings = append(findings, Finding{ClassStaleFingerprint, id, fmt.Sprintf("manifest source-revision %q disagrees with ledger revision %q", source.SourceRevision, te.SourceRevision)})
		}
		if te.Draft || te.Human || te.Retired {
			findings = append(findings, Finding{ClassUnjustifiedDisposition, id, "target is @draft, Human, or retired but the source manifest declares a proof class instead of a disposition"})
			continue
		}
		mismatch := se.ProofClass != te.ProofClass
		if se.Direction != "" && !containsString(te.Directions, se.Direction) {
			mismatch = true
		}
		if se.EvidenceLocation != "" && !containsString(te.EvidenceLocations, se.EvidenceLocation) {
			mismatch = true
		}
		if mismatch {
			findings = append(findings, Finding{ClassChangedEvidence, id, fmt.Sprintf("source proof-class=%q direction=%q location=%q; target proof-class=%q directions=%v locations=%v", se.ProofClass, se.Direction, se.EvidenceLocation, te.ProofClass, te.Directions, te.EvidenceLocations)})
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
	return Report{Findings: findings}
}

func containsString(list []string, want string) bool {
	for _, v := range list {
		if v == want {
			return true
		}
	}
	return false
}
