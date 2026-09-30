// Package evidence defines the framework-neutral, repository-owned executable
// evidence exchange. It reads current files and never starts a producer.
package evidence

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

const DefaultPath = ".clue/evidence.yaml"
const Version = 1

type Manifest struct {
	Version   int        `yaml:"version"`
	Producers []Producer `yaml:"producers"`
}

// Producer owns one input scope and its source-qualified executable references.
// Framework and language describe the producer, not validator dispatch keys.
type Producer struct {
	ID          string       `yaml:"id"`
	Framework   string       `yaml:"framework,omitempty"`
	Language    string       `yaml:"language,omitempty"`
	Include     []string     `yaml:"include"`
	Exclude     []string     `yaml:"exclude,omitempty"`
	Inputs      []Input      `yaml:"inputs"`
	References  []Reference  `yaml:"references"`
	Diagnostics []Diagnostic `yaml:"diagnostics,omitempty"`
}

type Input struct {
	Path   string `yaml:"path"`
	SHA256 string `yaml:"sha256"`
}

type Reference struct {
	ID        string `yaml:"id"`
	Path      string `yaml:"path"`
	Subject   string `yaml:"subject"`
	Type      string `yaml:"type,omitempty"`
	Direction string `yaml:"direction,omitempty"`
}

type Diagnostic struct {
	Path    string `yaml:"path"`
	Subject string `yaml:"subject,omitempty"`
	Message string `yaml:"message"`
}

var canonicalID = regexp.MustCompile(`^[A-Z][A-Z0-9]*(?:-[A-Z][A-Z0-9]*)*-[0-9]+[a-z]*$`)
var producerID = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_.-]*$`)
var digest = regexp.MustCompile(`^[0-9a-f]{64}$`)

func CanonicalID(id string) bool { return canonicalID.MatchString(id) }

// Parse refuses unknown fields, duplicate mapping keys and additional documents.
func Parse(data []byte) (Manifest, error) {
	var m Manifest
	d := yaml.NewDecoder(bytes.NewReader(data))
	d.KnownFields(true)
	if err := d.Decode(&m); err != nil {
		return m, err
	}
	var extra any
	if err := d.Decode(&extra); err != io.EOF {
		return m, fmt.Errorf("evidence must contain exactly one YAML document")
	}
	return m, CheckShape(m)
}

func CheckShape(m Manifest) error {
	if m.Version != Version {
		return fmt.Errorf("unsupported evidence version %d; expected %d", m.Version, Version)
	}
	seen := map[string]bool{}
	for _, p := range m.Producers {
		if !producerID.MatchString(p.ID) || seen[p.ID] {
			return fmt.Errorf("invalid or duplicate evidence producer %q", p.ID)
		}
		seen[p.ID] = true
		if len(p.Include) == 0 {
			return fmt.Errorf("producer %s has no include scope", p.ID)
		}
		for _, pattern := range append(append([]string{}, p.Include...), p.Exclude...) {
			if err := checkPattern(pattern); err != nil {
				return fmt.Errorf("producer %s: %w", p.ID, err)
			}
		}
		inputs := map[string]bool{}
		for _, in := range p.Inputs {
			if !safePath(in.Path) || !digest.MatchString(in.SHA256) || inputs[in.Path] {
				return fmt.Errorf("producer %s has invalid or duplicate input %q", p.ID, in.Path)
			}
			inputs[in.Path] = true
		}
		executables := map[string]Reference{}
		for _, ref := range p.References {
			if !CanonicalID(ref.ID) || !safePath(ref.Path) || strings.TrimSpace(ref.Subject) == "" || !inputs[ref.Path] {
				return fmt.Errorf("producer %s has invalid reference %q at %s (source must be an input)", p.ID, ref.ID, ref.Path)
			}
			if (ref.Type == "") != (ref.Direction == "") || (ref.Type != "" && !validType(ref.Type)) || (ref.Direction != "" && ref.Direction != "positive" && ref.Direction != "negative") {
				return fmt.Errorf("producer %s reference %s has malformed classification", p.ID, ref.ID)
			}
			key := ref.Path + "\x00" + ref.Subject
			if prev, exists := executables[key]; exists {
				return fmt.Errorf("producer %s has duplicate or conflicting executable %s %s (%s and %s)", p.ID, ref.Path, ref.Subject, prev.ID, ref.ID)
			}
			executables[key] = ref
		}
		for _, d := range p.Diagnostics {
			if !safePath(d.Path) || strings.TrimSpace(d.Message) == "" {
				return fmt.Errorf("producer %s has an invalid diagnostic", p.ID)
			}
		}
	}
	return nil
}

func validType(t string) bool {
	return t == "Unit" || t == "Integration" || t == "E2E" || t == "Performance"
}

func safePath(p string) bool {
	return p != "" && p != "." && !strings.ContainsAny(p, "\\:\x00") && !strings.HasPrefix(p, "/") && path.Clean(p) == p && p != ".." && !strings.HasPrefix(p, "../")
}

func checkPattern(p string) error {
	if !safePath(p) || p == DefaultPath {
		return fmt.Errorf("invalid repository-relative evidence pattern %q", p)
	}
	for _, segment := range strings.Split(p, "/") {
		if strings.Contains(segment, "**") && segment != "**" {
			return fmt.Errorf("** must be a complete path segment in %q", p)
		}
		if segment != "**" {
			if _, err := path.Match(segment, ""); err != nil {
				return fmt.Errorf("invalid evidence pattern %q: %w", p, err)
			}
		}
	}
	return nil
}

// Match supports slash-separated path.Match segments and recursive ** segments.
func Match(pattern, name string) bool {
	var match func([]string, []string) bool
	match = func(p, n []string) bool {
		if len(p) == 0 {
			return len(n) == 0
		}
		if p[0] == "**" {
			if match(p[1:], n) {
				return true
			}
			return len(n) > 0 && match(p, n[1:])
		}
		if len(n) == 0 {
			return false
		}
		ok, _ := path.Match(p[0], n[0])
		return ok && match(p[1:], n[1:])
	}
	return match(strings.Split(pattern, "/"), strings.Split(name, "/"))
}

func selected(patterns []string, name string) bool {
	for _, pattern := range patterns {
		if Match(pattern, name) {
			return true
		}
	}
	return false
}

func excluded(patterns []string, name string) bool {
	if selected(patterns, name) {
		return true
	}
	for parent := path.Dir(name); parent != "."; parent = path.Dir(parent) {
		if selected(patterns, parent) {
			return true
		}
	}
	return false
}

// Fingerprint normalizes CRLF, making a text checkout's digest platform-neutral.
func Fingerprint(data []byte) string {
	h := sha256.Sum256(bytes.ReplaceAll(data, []byte("\r\n"), []byte("\n")))
	return hex.EncodeToString(h[:])
}

// ReadFile confines a source to the repository, including every symlink ancestor.
func ReadFile(root, relative string) ([]byte, error) {
	if !safePath(relative) {
		return nil, fmt.Errorf("unsafe evidence path %q", relative)
	}
	current := root
	for _, segment := range strings.Split(relative, "/") {
		current = filepath.Join(current, segment)
		info, err := os.Lstat(current)
		if err != nil {
			return nil, err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return nil, fmt.Errorf("evidence input %s is a symlink", relative)
		}
	}
	return os.ReadFile(current)
}

// Snapshot enumerates complete declared scopes. Excludes also prune directories.
// The manifest itself cannot be an input: its digest would be self-referential.
func Snapshot(root string, p Producer) ([]Input, error) {
	for _, pattern := range append(append([]string{}, p.Include...), p.Exclude...) {
		if err := checkPattern(pattern); err != nil {
			return nil, err
		}
	}
	files := map[string]Input{}
	for _, pattern := range p.Include {
		base := []string{}
		for _, segment := range strings.Split(pattern, "/") {
			if strings.ContainsAny(segment, "*?[") {
				break
			}
			base = append(base, segment)
		}
		start := root
		if len(base) > 0 {
			start = filepath.Join(root, filepath.FromSlash(strings.Join(base, "/")))
		}
		// Missing exact inputs are errors; missing wildcard directories are an
		// empty scope, whose later creation still changes the file inventory.
		if _, err := os.Lstat(start); err != nil {
			if errors.Is(err, os.ErrNotExist) && len(base) < len(strings.Split(pattern, "/")) {
				continue
			}
			return nil, fmt.Errorf("producer %s scope %s: %w", p.ID, pattern, err)
		}
		// WalkDir must not follow a symlink in the static prefix either.
		if len(base) > 0 {
			current := root
			for _, segment := range base {
				current = filepath.Join(current, segment)
				info, err := os.Lstat(current)
				if err != nil {
					return nil, err
				}
				if info.Mode()&os.ModeSymlink != 0 {
					return nil, fmt.Errorf("symlink in evidence scope %s", pattern)
				}
			}
		}
		err := filepath.WalkDir(start, func(full string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			rel, err := filepath.Rel(root, full)
			if err != nil {
				return err
			}
			name := filepath.ToSlash(rel)
			if name == "." {
				return nil
			}
			if excluded(p.Exclude, name) {
				if d.IsDir() {
					return fs.SkipDir
				}
				return nil
			}
			if d.Type()&os.ModeSymlink != 0 {
				return fmt.Errorf("symlink in evidence scope: %s", name)
			}
			if d.IsDir() || !selected(p.Include, name) {
				return nil
			}
			if name == DefaultPath {
				return fmt.Errorf("evidence manifest cannot fingerprint itself")
			}
			data, err := ReadFile(root, name)
			if err != nil {
				return err
			}
			files[name] = Input{name, Fingerprint(data)}
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("producer %s: %w", p.ID, err)
		}
	}
	result := []Input{}
	for _, file := range files {
		result = append(result, file)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Path < result[j].Path })
	return result, nil
}

// Load verifies all producers before exposing any references. Errors and
// diagnostics fail the whole exchange; a partial snapshot never earns credit.
func Load(root string) (Manifest, error) {
	data, err := ReadFile(root, DefaultPath)
	if err != nil {
		return Manifest{}, fmt.Errorf("%s: %w; establish repository-owned evidence export", DefaultPath, err)
	}
	m, err := Parse(data)
	if err != nil {
		return Manifest{}, fmt.Errorf("%s: %w", DefaultPath, err)
	}
	for _, p := range m.Producers {
		actual, err := Snapshot(root, p)
		if err != nil {
			return Manifest{}, err
		}
		expected := append([]Input{}, p.Inputs...)
		sort.Slice(expected, func(i, j int) bool { return expected[i].Path < expected[j].Path })
		if len(actual) != len(expected) {
			return Manifest{}, fmt.Errorf("producer %s has stale evidence: input files added or deleted; regenerate export", p.ID)
		}
		for i := range actual {
			if actual[i] != expected[i] {
				return Manifest{}, fmt.Errorf("producer %s has stale evidence at %s; regenerate export", p.ID, actual[i].Path)
			}
		}
		if len(p.Diagnostics) > 0 {
			var messages []string
			for _, d := range p.Diagnostics {
				messages = append(messages, d.Path+": "+d.Subject+" "+d.Message)
			}
			sort.Strings(messages)
			return Manifest{}, fmt.Errorf("producer %s diagnostics: %s", p.ID, strings.Join(messages, "; "))
		}
	}
	return m, nil
}

// Encode imposes stable order without modifying its caller's slices.
func Encode(m Manifest) ([]byte, error) {
	if err := CheckShape(m); err != nil {
		return nil, err
	}
	m.Producers = append([]Producer{}, m.Producers...)
	for i, original := range m.Producers {
		p := original
		p.Include = append([]string{}, p.Include...)
		p.Exclude = append([]string{}, p.Exclude...)
		p.Inputs = append([]Input{}, p.Inputs...)
		p.References = append([]Reference{}, p.References...)
		p.Diagnostics = append([]Diagnostic{}, p.Diagnostics...)
		sort.Strings(p.Include)
		sort.Strings(p.Exclude)
		sort.Slice(p.Inputs, func(i, j int) bool { return p.Inputs[i].Path < p.Inputs[j].Path })
		sort.Slice(p.References, func(i, j int) bool {
			a, b := p.References[i], p.References[j]
			return a.Path+"\x00"+a.Subject+"\x00"+a.ID < b.Path+"\x00"+b.Subject+"\x00"+b.ID
		})
		sort.Slice(p.Diagnostics, func(i, j int) bool {
			return p.Diagnostics[i].Path+p.Diagnostics[i].Subject+p.Diagnostics[i].Message < p.Diagnostics[j].Path+p.Diagnostics[j].Subject+p.Diagnostics[j].Message
		})
		m.Producers[i] = p
	}
	sort.Slice(m.Producers, func(i, j int) bool { return m.Producers[i].ID < m.Producers[j].ID })
	var b bytes.Buffer
	e := yaml.NewEncoder(&b)
	e.SetIndent(2)
	if err := e.Encode(m); err != nil {
		return nil, err
	}
	if err := e.Close(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// Write atomically replaces the complete exchange only after all producers
// pass shape, freshness and diagnostic checks. Exporters supply all producers.
func Write(root string, m Manifest) error {
	data, err := Encode(m)
	if err != nil {
		return err
	}
	for _, p := range m.Producers {
		if len(p.Diagnostics) != 0 {
			return fmt.Errorf("producer %s failed; previous evidence preserved", p.ID)
		}
		actual, err := Snapshot(root, p)
		if err != nil {
			return err
		}
		expected := append([]Input{}, p.Inputs...)
		sort.Slice(expected, func(i, j int) bool { return expected[i].Path < expected[j].Path })
		if !slices.Equal(actual, expected) {
			return fmt.Errorf("producer %s inputs changed during export", p.ID)
		}
	}
	dir := filepath.Join(root, ".clue")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	info, err := os.Lstat(dir)
	if err != nil || info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("unsafe evidence output directory")
	}
	f, err := os.CreateTemp(dir, ".evidence-*")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	if _, err := f.Write(data); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(f.Name(), filepath.Join(root, filepath.FromSlash(DefaultPath)))
}
