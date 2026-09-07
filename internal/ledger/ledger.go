// Package ledger persists Cliewen's corpus-wide identity claims and
// retirement record (ADR-048): a .clue/id-ledger.yaml file replacing
// scan-and-max allocation, so a deleted artifact's ID is never silently
// reissued.
package ledger

import (
	"bytes"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

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

// legacyFile is the version-one on-disk shape retained for migration and
// backwards-compatible reads.
type legacyFile struct {
	Counters map[string]*big.Int `yaml:"counters"`
	Entries  []Entry             `yaml:"entries"`
}

// Coordination selects checkout-local allocation or Git-ref coordination.
type Coordination struct {
	Mode   string `yaml:"mode"`
	Remote string `yaml:"remote,omitempty"`
}

type eventFile struct {
	Version      int          `yaml:"version"`
	Coordination Coordination `yaml:"coordination"`
	Events       []Entry      `yaml:"events"`
}

// Ledger is the in-memory allocator and identity register: byID answers
// "is this ID used" and counters records the greatest permanent claim for
// each prefix, both in O(1) after one Load.
type Ledger struct {
	path     string
	version  int
	coord    Coordination
	events   []Entry
	counters map[string]*big.Int
	byID     map[string]*Entry
}

// DefaultPath is the ledger file's fixed location relative to a repository
// root (ADR-048): an operational registry, not authored corpus prose, so it
// lives outside docs/.
const DefaultPath = ".clue/id-ledger.yaml"

const UnionAttribute = "/.clue/id-ledger.yaml merge=union"

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
		version:  2,
		coord:    Coordination{Mode: "local"},
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
	var header struct {
		Version int `yaml:"version"`
	}
	if err := yaml.Unmarshal(data, &header); err != nil {
		return nil, fmt.Errorf("%s: %w", l.path, err)
	}
	if header.Version == 2 {
		var f eventFile
		if err := yaml.Unmarshal(data, &f); err != nil {
			return nil, fmt.Errorf("%s: %w", l.path, err)
		}
		if f.Coordination.Mode != "local" && f.Coordination.Mode != "git" {
			return nil, fmt.Errorf("%s: coordination mode must be local or git", l.path)
		}
		if f.Coordination.Mode == "git" && f.Coordination.Remote == "" {
			return nil, fmt.Errorf("%s: git coordination requires a remote", l.path)
		}
		l.version = 2
		l.coord = f.Coordination
		for _, e := range f.Events {
			if err := l.foldEvent(e); err != nil {
				return nil, fmt.Errorf("%s: %w", l.path, err)
			}
			l.events = append(l.events, cloneEntry(e))
		}
		return l, nil
	}
	if header.Version != 0 {
		return nil, fmt.Errorf("%s: unsupported ledger version %d", l.path, header.Version)
	}
	var f legacyFile
	if err := yaml.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("%s: %w", l.path, err)
	}
	l.version = 1
	for k, v := range f.Counters {
		if v == nil || v.Sign() < 0 {
			return nil, fmt.Errorf("%s: counter %s is not a non-negative decimal", l.path, k)
		}
		l.counters[k] = v
	}
	for i := range f.Entries {
		e := f.Entries[i]
		if _, exists := l.byID[e.ID]; exists {
			return nil, fmt.Errorf("%s: duplicate entry for id %s", l.path, e.ID)
		}
		l.byID[e.ID] = &e
	}
	for _, e := range l.byID {
		if !ValidNumericEntry(*e) {
			continue
		}
		counter, ok := l.counters[e.Prefix]
		if !ok {
			return nil, fmt.Errorf("%s: missing counter for numeric prefix %s", l.path, e.Prefix)
		}
		if counter.Cmp(e.Component) < 0 {
			return nil, fmt.Errorf("%s: counter %s is below recorded component %s", l.path, e.Prefix, e.Component)
		}
	}
	return l, nil
}

func cloneEntry(e Entry) Entry {
	if e.Component != nil {
		e.Component = new(big.Int).Set(e.Component)
	}
	return e
}

func stateRank(s State) int {
	switch s {
	case StateReserved:
		return 1
	case StateLive:
		return 2
	case StateRetired:
		return 3
	default:
		return 0
	}
}

func sameIdentity(a, b Entry) bool {
	if a.ID != b.ID || a.Kind != b.Kind || a.Prefix != b.Prefix || a.SourceRevision != b.SourceRevision || a.SourceLocation != b.SourceLocation {
		return false
	}
	if a.Component == nil || b.Component == nil {
		return a.Component == nil && b.Component == nil
	}
	return a.Component.Cmp(b.Component) == 0
}

func (l *Ledger) foldEvent(e Entry) error {
	if stateRank(e.State) == 0 {
		return fmt.Errorf("entry %s has invalid state %s", e.ID, e.State)
	}
	if e.Kind == KindNumeric {
		if !ValidNumericEntry(e) {
			return fmt.Errorf("entry %s is numeric-kind but its ID, prefix, and component do not agree", e.ID)
		}
		if current, ok := l.counters[e.Prefix]; !ok || e.Component.Cmp(current) > 0 {
			l.counters[e.Prefix] = new(big.Int).Set(e.Component)
		}
	} else if e.Kind != KindOpaque || e.Component != nil || e.Prefix != "" {
		return fmt.Errorf("entry %s has invalid opaque identity fields", e.ID)
	}
	if current, ok := l.byID[e.ID]; ok {
		if !sameIdentity(*current, e) {
			return fmt.Errorf("conflicting identity metadata for id %s", e.ID)
		}
		if stateRank(e.State) > stateRank(current.State) {
			current.State = e.State
		}
		return nil
	}
	c := cloneEntry(e)
	l.byID[e.ID] = &c
	return nil
}

// Exists reports whether root carries a ledger file at all — the gate that
// leaves a corpus without a ledger yet unaffected by ledger-backed rules.
func Exists(root string) bool {
	_, err := os.Stat(filepath.Join(root, DefaultPath))
	return err == nil
}

// HasUnionMerge reports whether the repository declares the built-in union
// driver that makes one-line ledger events merge without a hot-file conflict.
func HasUnionMerge(root string) bool {
	data, err := os.ReadFile(filepath.Join(root, ".gitattributes"))
	if err != nil {
		return false
	}
	for _, line := range strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n") {
		if strings.TrimSpace(line) == UnionAttribute {
			return true
		}
	}
	return false
}

// Bytes renders the ledger's current in-memory state as the exact YAML
// bytes Save would write, without touching disk — the shape a migration
// plan needs to compare a proposed write against what is already there.
func (l *Ledger) Bytes() ([]byte, error) {
	if l.version == 2 {
		f := eventFile{Version: 2, Coordination: l.coord, Events: l.events}
		var node yaml.Node
		if err := node.Encode(f); err != nil {
			return nil, err
		}
		setEventFlowStyle(&node)
		var out bytes.Buffer
		enc := yaml.NewEncoder(&out)
		enc.SetIndent(4)
		if err := enc.Encode(&node); err != nil {
			return nil, err
		}
		return out.Bytes(), nil
	}
	ids := make([]string, 0, len(l.byID))
	for id := range l.byID {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	f := legacyFile{Counters: l.counters, Entries: make([]Entry, 0, len(ids))}
	for _, id := range ids {
		f.Entries = append(f.Entries, *l.byID[id])
	}
	return yaml.Marshal(f)
}

func setEventFlowStyle(node *yaml.Node) {
	if node.Kind == yaml.DocumentNode && len(node.Content) == 1 {
		setEventFlowStyle(node.Content[0])
		return
	}
	if node.Kind != yaml.MappingNode {
		return
	}
	for i := 0; i+1 < len(node.Content); i += 2 {
		if node.Content[i].Value != "events" {
			continue
		}
		for _, event := range node.Content[i+1].Content {
			event.Style = yaml.FlowStyle
		}
	}
}

// Save persists the ledger atomically, creating .clue/ if needed.
func (l *Ledger) Save() error {
	if err := os.MkdirAll(filepath.Dir(l.path), 0o755); err != nil {
		return err
	}
	data, err := l.Bytes()
	if err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(l.path), ".id-ledger-*.tmp")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)
	if err := tmp.Chmod(0o644); err != nil {
		_ = tmp.Close()
		return err
	}
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmpName, l.path)
}

// Version reports the persisted ledger schema version.
func (l *Ledger) Version() int { return l.version }

// Coordination reports how numeric allocation is coordinated.
func (l *Ledger) Coordination() Coordination { return l.coord }

// ConvertV2 converts the loaded effective state into an append-only event log.
func (l *Ledger) ConvertV2() {
	if l.version == 2 {
		return
	}
	l.version = 2
	l.coord = Coordination{Mode: "local"}
	l.events = nil
	for _, e := range l.Entries() {
		l.events = append(l.events, cloneEntry(e))
	}
}

// SetGitCoordination enables coordinated allocation after its remote journal
// has been initialized successfully.
func (l *Ledger) SetGitCoordination(remote string) error {
	if l.version != 2 {
		return fmt.Errorf("identity ledger must be migrated to version 2")
	}
	if remote == "" {
		return fmt.Errorf("remote name is required")
	}
	l.coord = Coordination{Mode: "git", Remote: remote}
	return nil
}

// Claims returns every numeric identity normalized to a permanent reservation.
func (l *Ledger) Claims() []Entry {
	var out []Entry
	for _, e := range l.Entries() {
		if e.Kind != KindNumeric {
			continue
		}
		e.State = StateReserved
		out = append(out, e)
	}
	return out
}

// MergeClaims adds previously unseen remote numeric claims as reservations.
// Existing lifecycle state is never downgraded.
func (l *Ledger) MergeClaims(claims []Entry) error {
	for _, e := range claims {
		e.State = StateReserved
		if current, ok := l.byID[e.ID]; ok {
			if !sameIdentity(*current, e) {
				return fmt.Errorf("conflicting identity metadata for id %s", e.ID)
			}
			continue
		}
		if err := l.foldEvent(e); err != nil {
			return err
		}
		if l.version == 2 {
			l.events = append(l.events, cloneEntry(e))
		}
	}
	return nil
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
	e := Entry{ID: id, Kind: KindNumeric, State: StateReserved, Prefix: prefix, Component: new(big.Int).Set(n)}
	if l.version == 2 {
		if err := l.foldEvent(e); err != nil {
			return "", err
		}
		l.events = append(l.events, cloneEntry(e))
	} else {
		l.counters[prefix] = new(big.Int).Set(n)
		l.byID[id] = &e
	}
	return id, nil
}

// NextNumericN reserves count consecutive numeric identities atomically in
// memory. Callers persist the resulting ledger only after the whole batch has
// been allocated.
func (l *Ledger) NextNumericN(prefix string, count int) ([]string, error) {
	if count < 1 {
		return nil, fmt.Errorf("count must be positive")
	}
	ids := make([]string, 0, count)
	for range count {
		id, err := l.NextNumeric(prefix)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// ReserveOpaque records id as reserved for an opaque namespace, rejecting
// reuse of an ID already present under any state: an opaque ID is
// preserved verbatim and never reissued (ADR-048).
func (l *Ledger) ReserveOpaque(id, sourceRevision, sourceLocation string) error {
	if l.byID[id] != nil {
		return fmt.Errorf("id %s already used in the ledger", id)
	}
	e := Entry{ID: id, Kind: KindOpaque, State: StateReserved, SourceRevision: sourceRevision, SourceLocation: sourceLocation}
	if l.version == 2 {
		if err := l.foldEvent(e); err != nil {
			return err
		}
		l.events = append(l.events, cloneEntry(e))
	} else {
		l.byID[id] = &e
	}
	return nil
}

// MarkLive records id as backed by a live artifact. An ID with no existing
// entry gets one classified by shape — used by backfill, which seeds one
// live entry per currently-live corpus ID.
func (l *Ledger) MarkLive(id string) {
	if e, ok := l.byID[id]; ok {
		if l.version == 2 {
			if stateRank(e.State) >= stateRank(StateLive) {
				return
			}
			next := cloneEntry(*e)
			next.State = StateLive
			_ = l.foldEvent(next)
			l.events = append(l.events, next)
		} else {
			e.State = StateLive
		}
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
	if l.version == 2 {
		_ = l.foldEvent(*e)
		l.events = append(l.events, cloneEntry(*e))
	} else {
		l.byID[id] = e
	}
}

// MarkRetired records id as a tombstoned or deleted identity. Like MarkLive,
// it classifies a previously unseen ID by shape so migration backfill retains
// a criterion tombstone without ever making that ID allocatable again.
func (l *Ledger) MarkRetired(id string) {
	l.MarkLive(id)
	if l.version == 2 {
		e := cloneEntry(*l.byID[id])
		e.State = StateRetired
		_ = l.foldEvent(e)
		l.events = append(l.events, e)
	} else {
		l.byID[id].State = StateRetired
	}
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
	if l.version == 2 {
		next := cloneEntry(*e)
		next.State = StateLive
		if err := l.foldEvent(next); err != nil {
			return err
		}
		l.events = append(l.events, next)
	} else {
		e.State = StateLive
	}
	return nil
}

// Retire transitions an entry to Retired. A retired entry is never removed
// and never reissued (ADR-048).
func (l *Ledger) Retire(id string) {
	if e, ok := l.byID[id]; ok {
		if l.version == 2 {
			next := cloneEntry(*e)
			next.State = StateRetired
			_ = l.foldEvent(next)
			l.events = append(l.events, next)
		} else {
			e.State = StateRetired
		}
	}
}

func parseComponent(s string) (*big.Int, error) {
	n, ok := new(big.Int).SetString(s, 10)
	if !ok {
		return nil, fmt.Errorf("invalid decimal component %q", s)
	}
	return n, nil
}
