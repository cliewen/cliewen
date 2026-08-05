// Package ledger persists Cliewen's corpus-wide identity allocator and
// retirement record (ADR-048): a .clue/id-ledger.yaml file replacing
// scan-and-max allocation with O(1) map lookups, so a deleted artifact's
// ID is never silently reissued.
package ledger

import (
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"regexp"
	"sort"

	"gopkg.in/yaml.v3"
)

// Kind distinguishes a numeric, allocator-issued ID from an opaque,
// verbatim-preserved imported one.
type Kind string

const (
	KindNumeric Kind = "numeric"
	KindOpaque  Kind = "opaque"
)

// State is the ledger's identity lifecycle (ADR-048): an ID is claimed
// before an artifact exists (Reserved), backed by a live artifact (Live),
// or its artifact has been deleted (Retired). A retired entry is never
// removed and never reissued.
type State string

const (
	StateReserved State = "reserved"
	StateLive     State = "live"
	StateRetired  State = "retired"
)

// Entry is one ledger row. A numeric entry carries Prefix/Component; an
// opaque entry carries neither. An entry imported from a brownfield source
// additionally carries SourceRevision/SourceLocation.
type Entry struct {
	ID             string   `yaml:"id"`
	Kind           Kind     `yaml:"kind"`
	State          State    `yaml:"state"`
	Prefix         string   `yaml:"prefix,omitempty"`
	Component      *big.Int `yaml:"component,omitempty"`
	SourceRevision string   `yaml:"source-revision,omitempty"`
	SourceLocation string   `yaml:"source-location,omitempty"`
}

// file is the on-disk shape of .clue/id-ledger.yaml.
type file struct {
	Counters map[string]*big.Int `yaml:"counters"`
	Entries  []Entry             `yaml:"entries"`
}

// Ledger is the in-memory allocator and identity register: byID answers
// "is this ID used" and counters answers "what is the next ID for this
// prefix", both in O(1) after one Load.
type Ledger struct {
	path     string
	counters map[string]*big.Int
	byID     map[string]*Entry
}

// DefaultPath is the ledger file's fixed location relative to a repository
// root (ADR-048): an operational registry, not authored corpus prose, so it
// lives outside docs/.
const DefaultPath = ".clue/id-ledger.yaml"

// numericIDRe matches the canonical numeric-ID grammar used by native
// namespaces and criteria: uppercase alphanumeric prefix segments, a decimal
// component, and an optional lowercase suffix (ADR-037).
var (
	numericPrefixRe = regexp.MustCompile(`^[A-Z][A-Z0-9]*(?:-[A-Z][A-Z0-9]*)*$`)
	numericIDRe     = regexp.MustCompile(`^([A-Z][A-Z0-9]*(?:-[A-Z][A-Z0-9]*)*)-(\d+)[a-z]*$`)
)

// Load reads the ledger at root/DefaultPath. A missing file is not an
// error: it returns an empty, usable Ledger, so a corpus without a ledger
// yet is unaffected until the first allocation or backfill.
func Load(root string) (*Ledger, error) {
	l := &Ledger{
		path:     filepath.Join(root, DefaultPath),
		counters: map[string]*big.Int{},
		byID:     map[string]*Entry{},
	}
	data, err := os.ReadFile(l.path)
	if err != nil {
		if os.IsNotExist(err) {
			return l, nil
		}
		return nil, err
	}
	var f file
	if err := yaml.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("%s: %w", l.path, err)
	}
	for k, v := range f.Counters {
		l.counters[k] = v
	}
	for i := range f.Entries {
		e := f.Entries[i]
		if _, exists := l.byID[e.ID]; exists {
			return nil, fmt.Errorf("%s: duplicate entry for id %s", l.path, e.ID)
		}
		l.byID[e.ID] = &e
	}
	return l, nil
}

// Exists reports whether root carries a ledger file at all — the gate that
// leaves a corpus without a ledger yet unaffected by ledger-backed rules.
func Exists(root string) bool {
	_, err := os.Stat(filepath.Join(root, DefaultPath))
	return err == nil
}

// Bytes renders the ledger's current in-memory state as the exact YAML
// bytes Save would write, without touching disk — the shape a migration
// plan needs to compare a proposed write against what is already there.
func (l *Ledger) Bytes() ([]byte, error) {
	ids := make([]string, 0, len(l.byID))
	for id := range l.byID {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	f := file{Counters: l.counters, Entries: make([]Entry, 0, len(ids))}
	for _, id := range ids {
		f.Entries = append(f.Entries, *l.byID[id])
	}
	return yaml.Marshal(f)
}

// Save persists the ledger, creating .clue/ if needed. Entries are sorted
// by ID so the file diffs cleanly across runs.
func (l *Ledger) Save() error {
	if err := os.MkdirAll(filepath.Dir(l.path), 0o755); err != nil {
		return err
	}
	data, err := l.Bytes()
	if err != nil {
		return err
	}
	return os.WriteFile(l.path, data, 0o644)
}

// IsUsed reports whether id already has a ledger entry under any state —
// the check every allocator, native or opaque, must pass before minting.
func (l *Ledger) IsUsed(id string) bool {
	_, ok := l.byID[id]
	return ok
}

// Lookup returns the entry for id, if any.
func (l *Ledger) Lookup(id string) (Entry, bool) {
	e, ok := l.byID[id]
	if !ok {
		return Entry{}, false
	}
	return *e, true
}

// Entries returns every ledger entry, sorted by ID.
func (l *Ledger) Entries() []Entry {
	ids := make([]string, 0, len(l.byID))
	for id := range l.byID {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	out := make([]Entry, 0, len(ids))
	for _, id := range ids {
		out = append(out, *l.byID[id])
	}
	return out
}

// ValidNumericEntry reports whether a numeric ledger entry preserves its
// canonical ID's prefix and decimal component exactly as structured fields.
func ValidNumericEntry(e Entry) bool {
	if e.Kind != KindNumeric || e.Prefix == "" || e.Component == nil || e.Component.Sign() < 0 {
		return false
	}
	m := numericIDRe.FindStringSubmatch(e.ID)
	if m == nil || m[1] != e.Prefix {
		return false
	}
	n, err := parseComponent(m[2])
	return err == nil && n.Cmp(e.Component) == 0
}

// ValidNumericPrefix reports whether prefix can name a native numeric
// namespace under the canonical segmented-prefix grammar (ADR-037).
func ValidNumericPrefix(prefix string) bool {
	return numericPrefixRe.MatchString(prefix)
}

// NextNumeric allocates and reserves the next numeric ID for prefix: an
// increment of the stored counter, never a corpus scan. It skips past any
// component the ledger already holds for that prefix — live, reserved, or
// retired — so a retired ID's number is never reissued even if it once
// exceeded every live artifact's number.
func (l *Ledger) NextNumeric(prefix string) (string, error) {
	if !ValidNumericPrefix(prefix) {
		return "", fmt.Errorf("prefix %q is not a canonical numeric prefix", prefix)
	}
	n := new(big.Int)
	if current, ok := l.counters[prefix]; ok {
		n.Set(current)
	}
	var id string
	for {
		n.Add(n, big.NewInt(1))
		id = fmt.Sprintf("%s-%03d", prefix, n)
		if l.byID[id] == nil {
			break
		}
	}
	l.counters[prefix] = new(big.Int).Set(n)
	l.byID[id] = &Entry{ID: id, Kind: KindNumeric, State: StateReserved, Prefix: prefix, Component: new(big.Int).Set(n)}
	return id, nil
}

// ReserveOpaque records id as reserved for an opaque namespace, rejecting
// reuse of an ID already present under any state: an opaque ID is
// preserved verbatim and never reissued (ADR-048).
func (l *Ledger) ReserveOpaque(id, sourceRevision, sourceLocation string) error {
	if l.byID[id] != nil {
		return fmt.Errorf("id %s already used in the ledger", id)
	}
	l.byID[id] = &Entry{ID: id, Kind: KindOpaque, State: StateReserved, SourceRevision: sourceRevision, SourceLocation: sourceLocation}
	return nil
}

// MarkLive records id as backed by a live artifact. An ID with no existing
// entry gets one classified by shape — used by backfill, which seeds one
// live entry per currently-live corpus ID.
func (l *Ledger) MarkLive(id string) {
	if e, ok := l.byID[id]; ok {
		e.State = StateLive
		return
	}
	e := &Entry{ID: id, State: StateLive, Kind: KindOpaque}
	if m := numericIDRe.FindStringSubmatch(id); m != nil {
		if n, err := parseComponent(m[2]); err == nil {
			e.Kind = KindNumeric
			e.Prefix = m[1]
			e.Component = n
			if current, ok := l.counters[m[1]]; !ok || n.Cmp(current) > 0 {
				l.counters[m[1]] = new(big.Int).Set(n)
			}
		}
	}
	l.byID[id] = e
}

// PromoteReserved marks a previously allocated ID as live once its artifact
// has been created. It deliberately refuses unreserved IDs so this transition
// cannot bypass the allocator or revive a retired identity.
func (l *Ledger) PromoteReserved(id string) error {
	e, ok := l.byID[id]
	if !ok {
		return fmt.Errorf("id %s is not reserved in the ledger", id)
	}
	if e.State != StateReserved {
		return fmt.Errorf("id %s is %s, not reserved", id, e.State)
	}
	e.State = StateLive
	return nil
}

// Retire transitions an entry to Retired. A retired entry is never removed
// and never reissued (ADR-048).
func (l *Ledger) Retire(id string) {
	if e, ok := l.byID[id]; ok {
		e.State = StateRetired
	}
}

func parseComponent(s string) (*big.Int, error) {
	n, ok := new(big.Int).SetString(s, 10)
	if !ok {
		return nil, fmt.Errorf("invalid decimal component %q", s)
	}
	return n, nil
}
