// Package role records whether a repository is Cliewen's own source
// repository or an adopter of it (PDR-052): a .clue/role.yaml file naming
// a fact the tooling would otherwise infer from the incidental presence of
// cmd/clue, so a rule that binds only one kind of repository can be applied
// without guessing.
//
// The marker records the role and nothing else. A repository's Cliewen
// version already has two canonical carriers — the version stamp each
// generated skill.md carries and the CI caller's clue-version input — and a
// third copy would be exactly the drift ADR-013 warns about.
package role

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Role is the kind of repository a corpus belongs to.
type Role string

const (
	// Source is Cliewen's own repository: the one that generates the
	// shipped skills and templates an adopter receives.
	Source Role = "source"
	// Adopter is a repository that has adopted Cliewen. It is the default
	// for an undeclared repository, because every repository that predates
	// the marker is one.
	Adopter Role = "adopter"
)

// DefaultPath is the marker's fixed location relative to a repository root.
// It sits beside the identity ledger under .clue/ rather than in docs/ for
// the same reason: it is derived operational state, not authored corpus
// prose, and it configures nothing.
const DefaultPath = ".clue/role.yaml"

// file is the on-disk shape of .clue/role.yaml.
type file struct {
	Role string `yaml:"role"`
}

// Valid reports whether r is a role the marker may carry.
func Valid(r Role) bool {
	return r == Source || r == Adopter
}

// Load reads the marker at root/DefaultPath. A missing file is not an
// error: it returns Adopter with declared false, so a repository adopted
// before the marker existed keeps working and is never judged by a rule
// that binds only the source repository.
func Load(root string) (r Role, declared bool, err error) {
	path := filepath.Join(root, DefaultPath)
	data, readErr := os.ReadFile(path)
	if readErr != nil {
		if os.IsNotExist(readErr) {
			return Adopter, false, nil
		}
		return Adopter, false, readErr
	}
	var f file
	if err := yaml.Unmarshal(data, &f); err != nil {
		return Adopter, false, fmt.Errorf("%s: %w", DefaultPath, err)
	}
	value := Role(strings.TrimSpace(f.Role))
	if !Valid(value) {
		return Adopter, false, fmt.Errorf("%s: role must be %s or %s, not %q", DefaultPath, Source, Adopter, f.Role)
	}
	return value, true, nil
}

// Exists reports whether root carries a role marker at all.
func Exists(root string) bool {
	_, err := os.Stat(filepath.Join(root, DefaultPath))
	return err == nil
}

// Bytes renders a marker's exact file content without touching disk — the
// shape a migration plan needs to compare a proposed write against what is
// already there.
func Bytes(r Role) ([]byte, error) {
	if !Valid(r) {
		return nil, fmt.Errorf("role must be %s or %s, not %q", Source, Adopter, r)
	}
	return []byte("role: " + string(r) + "\n"), nil
}

// Write creates or replaces the marker at root/DefaultPath, creating .clue/
// if it is absent.
func Write(root string, r Role) error {
	data, err := Bytes(r)
	if err != nil {
		return err
	}
	path := filepath.Join(root, DefaultPath)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
