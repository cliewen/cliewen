// Package carriers implements the operational-carrier reconciliation
// contract (ADR-051, P-011 M-055): a pinned, human/agent-authored carrier
// inventory is compared against the current corpus, so a migration cannot
// silently drop a CI workflow, break a freshness check, lose a
// cross-reference registry, leave a dangling link to a path the migration
// deletes, or drop a diagram asset — none of which is a criterion or a
// piece of acceptance evidence, so neither `clue validate` nor `clue parity`
// would otherwise notice.
package carriers

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/cliewen/cliewen/internal/corpus"
	"gopkg.in/yaml.v3"
)

// Kind is one of the six operational-carrier kinds ADR-051 fixes as the
// inventory's vocabulary. A new migrated kind extends this list only
// through a plan decision, the same route ADR-049's proof-class vocabulary
// uses.
type Kind string

const (
	KindInstruction    Kind = "instruction"
	KindWorkflow       Kind = "workflow"
	KindFreshnessInput Kind = "freshness-input"
	KindRegistry       Kind = "registry"
	KindLink           Kind = "link"
	KindDiagramAsset   Kind = "diagram-asset"
)

var validKinds = map[Kind]bool{
	KindInstruction:    true,
	KindWorkflow:       true,
	KindFreshnessInput: true,
	KindRegistry:       true,
	KindLink:           true,
	KindDiagramAsset:   true,
}

// Entry is one pinned carrier row: an entry either names a mapped target
// path plus the content fingerprint pinned at inventory time, or is an
// explicit `blocked` marker with a reason. It may never combine the two
// (ADR-051).
type Entry struct {
	ID          string `yaml:"id"`
	Kind        Kind   `yaml:"kind"`
	SourcePath  string `yaml:"source-path"`
	TargetPath  string `yaml:"target-path,omitempty"`
	Fingerprint string `yaml:"fingerprint,omitempty"`
	Blocked     bool   `yaml:"blocked,omitempty"`
	Reason      string `yaml:"reason,omitempty"`
}

// Inventory is the pinned, human/agent-authored carrier manifest a
// clue-extract rehearsal writes: one file per source-mapping run, revised
// only when the source is re-read at a new revision (ADR-051).
type Inventory struct {
	SourceRevision string   `yaml:"source-revision"`
	SourceLocation string   `yaml:"source-location"`
	DeletedPaths   []string `yaml:"deleted-paths,omitempty"`
	Entries        []Entry  `yaml:"entries"`
}

// LoadInventory reads and validates a carrier inventory at path.
func LoadInventory(path string) (Inventory, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Inventory{}, err
	}
	var inv Inventory
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	if err := decoder.Decode(&inv); err != nil {
		return Inventory{}, fmt.Errorf("%s: %w", path, err)
	}
	if err := validateInventory(inv); err != nil {
		return Inventory{}, fmt.Errorf("%s: %w", path, err)
	}
	return inv, nil
}

func validateInventory(inv Inventory) error {
	if inv.SourceRevision == "" || inv.SourceLocation == "" {
		return fmt.Errorf("source-revision and source-location are required")
	}
	seen := map[string]bool{}
	for _, e := range inv.Entries {
		if e.ID == "" {
			return fmt.Errorf("entry id is required")
		}
		if seen[e.ID] {
			return fmt.Errorf("entry %q: duplicate id", e.ID)
		}
		seen[e.ID] = true
		if !validKinds[e.Kind] {
			return fmt.Errorf("entry %q: unsupported kind %q", e.ID, e.Kind)
		}
		if e.SourcePath == "" {
			return fmt.Errorf("entry %q: source-path is required", e.ID)
		}
		mapped := e.TargetPath != "" || e.Fingerprint != ""
		switch {
		case e.Blocked:
			if e.Reason == "" {
				return fmt.Errorf("entry %q: blocked entries require a reason", e.ID)
			}
			if mapped {
				return fmt.Errorf("entry %q: blocked entries cannot also carry target-path or fingerprint", e.ID)
			}
		case mapped:
			if e.TargetPath == "" || e.Fingerprint == "" {
				return fmt.Errorf("entry %q: mapped entries require both target-path and fingerprint", e.ID)
			}
			if e.Reason != "" {
				return fmt.Errorf("entry %q: mapped entries cannot carry a reason", e.ID)
			}
		default:
			return fmt.Errorf("entry %q: expected either target-path and fingerprint, or blocked and a reason", e.ID)
		}
	}
	return nil
}

// Fingerprint returns the sha256 hex digest of the file at path, the
// content hash a mapped entry pins and Reconcile re-derives.
func Fingerprint(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}

// Finding classes, matching ADR-051's three required failure classes.
const (
	ClassStaleDeletedPath = "stale-deleted-path"
	ClassLostFingerprint  = "lost-fingerprint"
	ClassMissingAsset     = "missing-asset"
)

// Finding is one reconciliation disagreement between the inventory and the
// current corpus.
type Finding struct {
	Class  string
	ID     string
	Detail string
}

// Report is Reconcile's deterministic result: a clean report has no
// findings, and repeated runs against the same inputs produce the same
// findings in the same order, so a CI artifact built from it is
// reproducible (ADR-051).
type Report struct {
	Findings []Finding
}

// Failed reports whether the reconciliation run must exit non-zero.
func (r Report) Failed() bool { return len(r.Findings) > 0 }

// mdLinkTargetRe matches the target half of a markdown link.
var mdLinkTargetRe = regexp.MustCompile(`\]\(([^)]+)\)`)

// mdReferenceTargetRe matches a reference-style Markdown link definition.
var mdReferenceTargetRe = regexp.MustCompile(`^\s*\[[^\]]+\]:\s*<?([^\s>]+)`)

// Reconcile recomputes every mapped entry's target fingerprint from what is
// actually in root right now and compares it against the pinned value, and
// scans the corpus under root for a Markdown link still resolving to a path
// inv names as deleted. It never writes back to inv or to root (ADR-051).
func Reconcile(inv Inventory, root string) (Report, error) {
	var findings []Finding

	entries := make([]Entry, len(inv.Entries))
	copy(entries, inv.Entries)
	sort.Slice(entries, func(i, j int) bool { return entries[i].ID < entries[j].ID })

	for _, e := range entries {
		if e.Blocked {
			// A blocked entry has no target to reconcile; its presence is
			// itself the record that the carrier is a known gap.
			continue
		}
		full := filepath.Join(root, filepath.FromSlash(e.TargetPath))
		info, err := os.Stat(full)
		if err != nil || info.IsDir() {
			findings = append(findings, Finding{ClassMissingAsset, e.ID, fmt.Sprintf("target path %q does not exist", e.TargetPath)})
			continue
		}
		fp, err := Fingerprint(full)
		if err != nil {
			return Report{}, err
		}
		if fp != e.Fingerprint {
			findings = append(findings, Finding{ClassLostFingerprint, e.ID, fmt.Sprintf("target %q fingerprint %s disagrees with pinned %s", e.TargetPath, fp, e.Fingerprint)})
		}
	}

	if len(inv.DeletedPaths) > 0 {
		staleFindings, err := staleDeletedPathFindings(root, inv.DeletedPaths)
		if err != nil {
			return Report{}, err
		}
		findings = append(findings, staleFindings...)
	}

	sort.Slice(findings, func(i, j int) bool {
		if findings[i].ID != findings[j].ID {
			return findings[i].ID < findings[j].ID
		}
		return findings[i].Class < findings[j].Class
	})
	return Report{Findings: findings}, nil
}

// staleDeletedPathFindings scans every Markdown file under root's corpus
// for a link target that still resolves to a path deletedPaths names.
func staleDeletedPathFindings(root string, deletedPaths []string) ([]Finding, error) {
	c, issues := corpus.Scan(root)
	if len(issues) > 0 {
		return nil, fmt.Errorf("corpus scan issues: %v", issues)
	}

	deleted := map[string]bool{}
	for _, p := range deletedPaths {
		deleted[strings.TrimSuffix(filepath.ToSlash(p), "/")] = true
	}

	paths := make([]string, 0, len(c.Contents))
	for p := range c.Contents {
		paths = append(paths, p)
	}
	sort.Strings(paths)

	var findings []Finding
	for _, path := range paths {
		text := c.Contents[path]
		for lineNo, line := range strings.Split(text, "\n") {
			matches := mdLinkTargetRe.FindAllStringSubmatch(line, -1)
			matches = append(matches, mdReferenceTargetRe.FindAllStringSubmatch(line, -1)...)
			for _, m := range matches {
				target := normalizeLinkTarget(m[1])
				if target == "" {
					continue
				}
				resolved := filepath.ToSlash(filepath.Clean(filepath.Join(filepath.Dir(path), target)))
				if deleted[target] || (!isExternalLinkTarget(target) && deleted[resolved]) {
					findings = append(findings, Finding{
						ClassStaleDeletedPath,
						path,
						fmt.Sprintf("line %d: still references deleted path %q", lineNo+1, target),
					})
				}
			}
		}
	}
	return findings, nil
}

// normalizeLinkTarget strips a title, a fragment, and surrounding
// whitespace from a raw Markdown link target, leaving the bare path a
// deleted-paths entry can be compared against.
func normalizeLinkTarget(raw string) string {
	target := strings.TrimSpace(raw)
	if idx := strings.IndexAny(target, " \t"); idx >= 0 {
		target = target[:idx]
	}
	target = strings.SplitN(target, "#", 2)[0]
	return filepath.ToSlash(target)
}

// isExternalLinkTarget reports whether target has a URL scheme or is a
// scheme-relative URL, neither of which can resolve to a deleted repository
// path by walking from the referencing Markdown file.
func isExternalLinkTarget(target string) bool {
	return strings.Contains(target, "://") || strings.HasPrefix(target, "//")
}
