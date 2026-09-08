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
	"slices"
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

// eventFile is the read shape. It still understands an inline coordination
// block so a ledger written before the split loads unchanged; ledgerFile is
// what Save writes, and never emits one.
type eventFile struct {
	Version      int          `yaml:"version"`
	Coordination Coordination `yaml:"coordination"`
	Events       []Entry      `yaml:"events"`
	HighWater    []Entry      `yaml:"high-water,omitempty"`
}

type ledgerFile struct {
	Version   int     `yaml:"version"`
	Events    []Entry `yaml:"events"`
	HighWater []Entry `yaml:"high-water,omitempty"`
}

// Ledger is the in-memory allocator and identity register: byID answers
// "is this ID used" and counters records the greatest permanent claim for
// each prefix, both in O(1) after one Load.
type Ledger struct {
	root      string
	path      string
	version   int
	coord     Coordination
	events    []Entry
	highWater []Entry
	counters  map[string]*big.Int
	byID      map[string]*Entry
	damage    string
}

// DefaultPath is the ledger file's fixed location relative to a repository
// root (ADR-048): an operational registry, not authored corpus prose, so it
// lives outside docs/.
const DefaultPath = ".clue/id-ledger.yaml"

const UnionAttribute = "/.clue/id-ledger.yaml merge=union"

// CoordinationPath holds how this repository allocates — checkout-local, or
// serialized through a named Git remote. It is deliberately not the ledger.
// The ledger is merged with Git's union driver so parallel branches can each
// append events without a conflict, and that driver keeps both sides of every
// difference: two branches that each enabled coordination against their own
// remote would silently produce a ledger claiming both, which no reading can
// settle and which reaches main already broken. Kept in its own file, the same
// disagreement is an ordinary merge conflict that Git raises at merge time,
// for the person doing the merge to resolve before anything is committed
// (ADR-069).
const CoordinationPath = ".clue/id-coordination.yaml"

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
		root:     root,
		path:     filepath.Join(root, DefaultPath),
		version:  2,
		coord:    Coordination{Mode: "local"},
		counters: map[string]*big.Int{},
		byID:     map[string]*Entry{},
	}
	data, err := os.ReadFile(l.path)
	if err != nil {
		if os.IsNotExist(err) {
			// Settings are deliberately not read here. Nothing allocates
			// before a ledger exists, so there is no reinterpretation to
			// prevent, and failing on this path would stop the migration that
			// seeds the ledger in the first place.
			return l, nil
		}
		return nil, err
	}
	f, damage, err := readEventFile(data)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", l.path, err)
	}
	if f.Version == 2 {
		coord, coordErr := resolveCoordination(root, f.Coordination)
		if coordErr != nil {
			return nil, coordErr
		}
		l.version = 2
		l.coord = coord
		l.damage = damage
		for _, e := range f.HighWater {
			if e.Kind != KindNumeric || e.State != StateReserved || !ValidNumericEntry(e) {
				return nil, fmt.Errorf("%s: invalid high-water claim %s", l.path, e.ID)
			}
			if current, ok := l.counters[e.Prefix]; !ok || e.Component.Cmp(current) > 0 {
				l.counters[e.Prefix] = new(big.Int).Set(e.Component)
			}
			l.highWater = append(l.highWater, cloneEntry(e))
		}
		for _, e := range f.Events {
			if err := l.foldEvent(e); err != nil {
				return nil, fmt.Errorf("%s: %w", l.path, err)
			}
			l.events = append(l.events, cloneEntry(e))
		}
		return l, nil
	}
	if f.Version != 0 {
		return nil, fmt.Errorf("%s: unsupported ledger version %d", l.path, f.Version)
	}
	var legacy legacyFile
	if err := yaml.Unmarshal(data, &legacy); err != nil {
		return nil, fmt.Errorf("%s: %w", l.path, err)
	}
	l.version = 1
	for k, v := range legacy.Counters {
		if v == nil || v.Sign() < 0 {
			return nil, fmt.Errorf("%s: counter %s is not a non-negative decimal", l.path, k)
		}
		l.counters[k] = v
	}
	for i := range legacy.Entries {
		e := legacy.Entries[i]
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

// readEventFile decodes the on-disk ledger through the YAML node API rather
// than straight into eventFile, because Git's union merge can leave the file
// with its header repeated. The driver combines both branches' copies of a
// conflicting hunk, and while that is exactly right for the entry list, a
// hunk covering the top of the file yields two `version` and `coordination`
// keys, which a struct decode rejects outright — wedging every command,
// including the sync that documentation offers as the recovery.
//
// Node decoding keeps every key and value in order, so a combined file can be
// folded back into the one ledger both branches meant. Repeated sequences
// concatenate, which is sound because events fold idempotently and advance
// state monotonically: their union is the same answer whatever the order.
// Repeated scalars must agree, because a genuine disagreement — one branch
// coordinating through Git while the other stays local — is a decision only a
// human can make, and guessing either way would silently move a repository
// off the allocation mode its team chose.
//
// It returns the file, a description of any repetition it had to reconcile
// (empty when the file was clean), and an error only for content no reading
// can rescue.
func readEventFile(data []byte) (eventFile, string, error) {
	var doc yaml.Node
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return eventFile{}, "", err
	}
	if doc.Kind != yaml.DocumentNode || len(doc.Content) == 0 {
		return eventFile{}, "", nil // an empty file is an empty legacy ledger
	}
	root := doc.Content[0]
	if root.Kind != yaml.MappingNode {
		return eventFile{}, "", fmt.Errorf("ledger must be a mapping")
	}
	var f eventFile
	var versionSet, coordSet bool
	repeated := map[string]bool{}
	seen := map[string]bool{}
	for i := 0; i+1 < len(root.Content); i += 2 {
		key := root.Content[i].Value
		value := root.Content[i+1]
		if seen[key] {
			repeated[key] = true
		}
		seen[key] = true
		switch key {
		case "version":
			var v int
			if err := value.Decode(&v); err != nil {
				return eventFile{}, "", err
			}
			if versionSet && v != f.Version {
				return eventFile{}, "", fmt.Errorf("the ledger declares two different versions (%d and %d); it was combined from incompatible copies and must be repaired by hand", f.Version, v)
			}
			f.Version, versionSet = v, true
		case "coordination":
			c, err := decodeCoordination(value)
			if err != nil {
				return eventFile{}, "", err
			}
			if coordSet && c != f.Coordination {
				return eventFile{}, "", fmt.Errorf("the ledger carries two different coordination settings (%s and %s); Git's union merge combined two branches that each chose one, and only a person can decide which the team meant — set the surviving one and remove the other", describeCoordination(f.Coordination), describeCoordination(c))
			}
			f.Coordination, coordSet = c, true
		case "events":
			var events []Entry
			if err := value.Decode(&events); err != nil {
				return eventFile{}, "", err
			}
			f.Events = append(f.Events, events...)
		case "high-water":
			var highWater []Entry
			if err := value.Decode(&highWater); err != nil {
				return eventFile{}, "", err
			}
			f.HighWater = append(f.HighWater, highWater...)
		}
	}
	return f, describeRepetition(repeated), nil
}

// decodeCoordination reads the coordination block through the node API for the
// same reason the file itself is read that way, one level down. Union merge
// duplicates only the lines that genuinely differ, so two branches that each
// enabled coordination against their own remote do not leave a repeated header
// at all — `version` and `mode` are identical on both sides and merge as
// context, while `remote` appears twice inside this block. That is the damage
// a real merge produces, and a struct decode rejected it with a parser message
// about a line nobody wrote.
//
// A nested duplicate is always a disagreement, because identical lines never
// duplicate, so it is always the case only a person can settle.
func decodeCoordination(value *yaml.Node) (Coordination, error) {
	var c Coordination
	if value.Kind != yaml.MappingNode {
		return c, value.Decode(&c)
	}
	seen := map[string]string{}
	for i := 0; i+1 < len(value.Content); i += 2 {
		key := value.Content[i].Value
		text := value.Content[i+1].Value
		if previous, ok := seen[key]; ok {
			if previous == text {
				continue
			}
			return Coordination{}, fmt.Errorf("the ledger names two different %ss (%s and %s); two branches each coordinated to their own, and only a person can decide which the team meant — keep one and remove the other", coordinationNoun(key), previous, text)
		}
		seen[key] = text
		switch key {
		case "mode":
			c.Mode = text
		case "remote":
			c.Remote = text
		}
	}
	return c, nil
}

func coordinationNoun(key string) string {
	switch key {
	case "remote":
		return "allocator remote"
	case "mode":
		return "coordination mode"
	default:
		return key + " setting"
	}
}

func describeCoordination(c Coordination) string {
	if c.Remote == "" {
		return c.Mode
	}
	return c.Mode + " through " + c.Remote
}

// describeRepetition names what the file repeated, in the file's own order, so
// the message a user sees points at the lines they can go and look at.
func describeRepetition(repeated map[string]bool) string {
	if len(repeated) == 0 {
		return ""
	}
	var keys []string
	for _, key := range []string{"version", "coordination", "events", "high-water"} {
		if repeated[key] {
			keys = append(keys, key)
		}
	}
	var what string
	switch len(keys) {
	case 1:
		what = keys[0]
	case 2:
		what = keys[0] + " and " + keys[1]
	default:
		what = strings.Join(keys[:len(keys)-1], ", ") + ", and " + keys[len(keys)-1]
	}
	return "the identity ledger repeats " + what + ", which is what Git's union merge leaves behind when two branches both change the top of the file; its entries are intact"
}

// resolveCoordination prefers the standalone file and falls back to a
// coordination block still written inline in the ledger, which is how every
// version-two ledger recorded it before the split. A ledger that carries
// neither allocates locally, so a file written by Save — which never emits the
// inline block — reads back correctly.
func resolveCoordination(root string, inline Coordination) (Coordination, error) {
	coord := inline
	path := filepath.Join(root, filepath.FromSlash(CoordinationPath))
	data, err := os.ReadFile(path)
	switch {
	case err == nil:
		if bytes.Contains(data, []byte("<<<<<<<")) {
			return Coordination{}, fmt.Errorf("%s: unresolved merge conflict — two branches chose different allocation settings, and Git is asking which the team meant; keep one and remove the markers", path)
		}
		var fromFile Coordination
		if unmarshalErr := yaml.Unmarshal(data, &fromFile); unmarshalErr != nil {
			return Coordination{}, fmt.Errorf("%s: %w", path, unmarshalErr)
		}
		if fromFile.Mode == "" {
			return Coordination{}, fmt.Errorf("%s: exists but names no allocation mode; a settings file that cannot be read as local or git is not evidence of either, and guessing would move the repository off the mode its team chose", path)
		}
		coord = fromFile
	case !os.IsNotExist(err):
		return Coordination{}, err
	}
	if coord.Mode == "" {
		coord.Mode = "local"
	}
	if coord.Mode != "local" && coord.Mode != "git" {
		return Coordination{}, fmt.Errorf("%s: coordination mode must be local or git", path)
	}
	if coord.Mode == "git" && coord.Remote == "" {
		return Coordination{}, fmt.Errorf("%s: git coordination requires a remote", path)
	}
	return coord, nil
}

// CoordinationBytes renders allocation settings as the file's exact contents,
// so a caller that writes the ledger itself — migration plans a change per
// path rather than saving a Ledger — can carry the settings across with it.
func CoordinationBytes(c Coordination) ([]byte, error) { return yaml.Marshal(c) }

// saveCoordination writes the file only for a coordinated repository. Local
// allocation is the absence of the file rather than a file saying so, which
// keeps a repository that never coordinates free of a setting it does not
// have, and removes the file when a repository stops coordinating.
func (l *Ledger) saveCoordination() error {
	path := filepath.Join(l.root, filepath.FromSlash(CoordinationPath))
	if l.coord.Mode != "git" {
		if l.version != 2 {
			return nil // settings this ledger never read are not ours to delete
		}
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			return err
		}
		return nil
	}
	data, err := CoordinationBytes(l.coord)
	if err != nil {
		return err
	}
	return writeFileAtomically(path, data)
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
	return declaresUnionMerge(data)
}

// WithUnionAttribute returns .gitattributes content carrying the union rule,
// and reports whether it had to append it. The append is strictly additive —
// nothing already in the file is read for meaning, reordered, or rewritten —
// because both callers are installing one missing line into a file the
// repository owns: `clue migrate` for a repository that already has a ledger,
// and `clue init` for one scaffolded over an existing `* text=auto`. They
// share this so the two paths cannot drift into different files.
func WithUnionAttribute(data []byte) ([]byte, bool) {
	if declaresUnionMerge(data) {
		return data, false
	}
	out := append([]byte(nil), data...)
	if len(out) > 0 && out[len(out)-1] != '\n' {
		out = append(out, '\n')
	}
	out = append(out, UnionAttribute...)
	out = append(out, '\n')
	return out, true
}

// ledgerAttributePatterns are the .gitattributes patterns that name the ledger
// file. Git resolves far more than these, but validate reads only repository
// bytes, so this recognizes the spellings a repository actually writes rather
// than reimplementing pattern matching or asking Git. Matching one exact line
// used to report a repository that wrote the rule without the leading slash as
// non-compliant and refuse to enable coordination for it, when the rule it had
// was working.
var ledgerAttributePatterns = []string{"/" + DefaultPath, DefaultPath}

func declaresUnionMerge(data []byte) bool {
	for _, line := range strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 || strings.HasPrefix(fields[0], "#") {
			continue
		}
		if slices.Contains(ledgerAttributePatterns, fields[0]) && slices.Contains(fields[1:], "merge=union") {
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
		f := ledgerFile{Version: 2, Events: l.events, HighWater: l.highWater}
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
		if node.Content[i].Value != "events" && node.Content[i].Value != "high-water" {
			continue
		}
		for _, event := range node.Content[i+1].Content {
			event.Style = yaml.FlowStyle
		}
	}
}

// Save persists the ledger atomically, creating .clue/ if needed.
func (l *Ledger) Save() error {
	data, err := l.Bytes()
	if err != nil {
		return err
	}
	if err := writeFileAtomically(l.path, data); err != nil {
		return err
	}
	return l.saveCoordination()
}

// writeFileAtomically replaces path through a temporary file in the same
// directory, so a reader never observes a half-written registry.
func writeFileAtomically(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), "."+filepath.Base(path)+"-*.tmp")
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
	return os.Rename(tmpName, path)
}

// Damage describes any repetition Load had to reconcile in the file it read,
// and is empty for a clean ledger. Saving the ledger rewrites it whole, so any
// command that writes repairs the file as a side effect of succeeding; a
// command that only reads reports this instead.
func (l *Ledger) Damage() string { return l.damage }

// Version reports the persisted ledger schema version.
func (l *Ledger) Version() int { return l.version }

// Coordination reports how numeric allocation is coordinated.
func (l *Ledger) Coordination() Coordination { return l.coord }

// ConvertV2 converts the loaded effective state into an append-only event log.
func (l *Ledger) ConvertV2() {
	if l.version == 2 {
		return
	}
	entries := l.Entries()
	highWater := make(map[string]*big.Int, len(l.counters))
	for prefix, component := range l.counters {
		highWater[prefix] = new(big.Int).Set(component)
	}
	l.version = 2
	l.coord = Coordination{Mode: "local"}
	l.events = nil
	l.highWater = nil
	l.counters = map[string]*big.Int{}
	l.byID = map[string]*Entry{}
	for _, e := range entries {
		_ = l.foldEvent(e)
		l.events = append(l.events, cloneEntry(e))
	}
	// The legacy counter is the numeric allocator's high-water mark. Opaque
	// identities are permanent claims, but their spelling must not change the
	// next numeric result merely because it resembles a numeric ID.
	l.counters = map[string]*big.Int{}
	for prefix, component := range highWater {
		l.counters[prefix] = new(big.Int).Set(component)
	}
	prefixes := make([]string, 0, len(highWater))
	for prefix := range highWater {
		prefixes = append(prefixes, prefix)
	}
	sort.Strings(prefixes)
	for _, prefix := range prefixes {
		component := highWater[prefix]
		if component.Sign() == 0 || !ValidNumericPrefix(prefix) {
			continue
		}
		e := Entry{ID: fmt.Sprintf("%s-%03d", prefix, component), Kind: KindNumeric, State: StateReserved, Prefix: prefix, Component: new(big.Int).Set(component)}
		l.highWater = append(l.highWater, cloneEntry(e))
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

// HighWater returns migrated numeric allocator high-water claims that have no
// corresponding lifecycle identity event, such as a legacy counter whose
// earlier artifact was removed. They are allocator metadata, not artifacts.
func (l *Ledger) HighWater() []Entry {
	out := make([]Entry, len(l.highWater))
	for i, e := range l.highWater {
		out[i] = cloneEntry(e)
	}
	return out
}

// MergeHighWater imports allocator metadata without creating lifecycle
// identities. It advances the local numeric boundary monotonically and keeps
// an event only when that boundary was previously absent locally.
func (l *Ledger) MergeHighWater(highWater []Entry) error {
	for _, e := range highWater {
		if e.Kind != KindNumeric || e.State != StateReserved || !ValidNumericEntry(e) {
			return fmt.Errorf("high-water claim %s is not a valid numeric identity", e.ID)
		}
		if current, ok := l.counters[e.Prefix]; ok && current.Cmp(e.Component) >= 0 {
			continue
		}
		l.counters[e.Prefix] = new(big.Int).Set(e.Component)
		if l.version == 2 {
			l.highWater = append(l.highWater, cloneEntry(e))
		}
	}
	return nil
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
func (l *Ledger) MarkLive(id string) { l.mark(id, StateLive) }

// mark records id at state, classifying an unseen ID by shape. It is the
// shared body of MarkLive and MarkRetired: exactly one event, at the state
// asked for, and none at all when the entry already stands there or beyond.
// Recording an identity that was retired before this ledger existed must not
// invent a live event first — the log is read by people, and a transition
// that never happened is noise the union merge then has to reconcile.
func (l *Ledger) mark(id string, state State) {
	if e, ok := l.byID[id]; ok {
		if stateRank(e.State) >= stateRank(state) {
			return
		}
		if l.version == 2 {
			next := cloneEntry(*e)
			next.State = state
			_ = l.foldEvent(next)
			l.events = append(l.events, next)
		} else {
			e.State = state
		}
		return
	}
	e := &Entry{ID: id, State: state, Kind: KindOpaque}
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
func (l *Ledger) MarkRetired(id string) { l.mark(id, StateRetired) }

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
	if _, ok := l.byID[id]; ok {
		l.mark(id, StateRetired)
	}
}

func parseComponent(s string) (*big.Int, error) {
	n, ok := new(big.Int).SetString(s, 10)
	if !ok {
		return nil, fmt.Errorf("invalid decimal component %q", s)
	}
	return n, nil
}
