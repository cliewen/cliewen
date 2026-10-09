// Package acceptpolicy reads repository-owned acceptance policy without Git
// history, network access, or writes. Init, migration and acceptance share it.
package acceptpolicy

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/cliewen/cliewen/internal/role"
	"gopkg.in/yaml.v3"
)

const Path = ".clue/acceptance.yaml"
const Local = "local"
const PR = "pr"

type Policy struct {
	Mode   string `yaml:"mode"`
	Branch string `yaml:"branch"`
}

// Load defaults to local acceptance on main when no policy is declared.
func Load(root string) (Policy, bool, error) {
	p := Policy{Mode: Local, Branch: "main"}
	r, _, err := role.Load(root)
	if err != nil {
		return p, false, err
	}
	name := filepath.Join(root, Path)
	if _, err = os.Lstat(name); os.IsNotExist(err) {
		return p, false, nil
	} else if err != nil {
		return p, false, err
	}
	data, err := os.ReadFile(name)
	if err != nil {
		return p, true, err
	}
	p = Policy{}
	d := yaml.NewDecoder(bytes.NewReader(data))
	d.KnownFields(true)
	if err = d.Decode(&p); err != nil {
		return p, true, fmt.Errorf("%s: %w", Path, err)
	}
	var extra any
	if err = d.Decode(&extra); err != io.EOF {
		return p, true, fmt.Errorf("%s: expected one YAML document", Path)
	}
	if err = Validate(p); err != nil {
		return p, true, err
	}
	if r == role.Source && p.Mode != PR {
		return p, true, fmt.Errorf("source repositories require explicit PR acceptance")
	}
	return p, true, nil
}

func Validate(p Policy) error {
	if p.Mode != Local && p.Mode != PR {
		return fmt.Errorf("%s requires mode: local or mode: pr", Path)
	}
	// Equivalent to check-ref-format refs/heads/<branch>; no Git executable
	// is needed to materialize a convention into a non-Git directory.
	ref := "refs/heads/" + p.Branch
	if p.Branch == "" || strings.HasSuffix(ref, ".") || strings.Contains(ref, "..") || strings.Contains(ref, "@{") {
		return fmt.Errorf("%s requires a valid explicit branch", Path)
	}
	for _, r := range ref {
		if r <= 32 || r == 127 || strings.ContainsRune("~^:?*[\\", r) {
			return fmt.Errorf("invalid acceptance branch %q", p.Branch)
		}
	}
	for _, part := range strings.Split(ref, "/") {
		if part == "" || strings.HasPrefix(part, ".") || strings.HasSuffix(part, ".lock") {
			return fmt.Errorf("invalid acceptance branch %q", p.Branch)
		}
	}
	return nil
}

func (p Policy) Bytes() []byte {
	return []byte("mode: " + p.Mode + "\nbranch: " + yamlScalar(p.Branch) + "\n")
}

func yamlScalar(s string) string {
	// YAML emits a scalar safely even for otherwise valid Git names such as #.
	b, _ := yaml.Marshal(s)
	return strings.TrimSpace(string(b))
}

// SelectInit resolves policy before any materialization. Explicit options may
// create a missing file, but never replace an existing repository decision.
func SelectInit(root, requested string) (Policy, error) {
	if requested != "" && requested != Local && requested != PR {
		return Policy{}, fmt.Errorf("--acceptance must be local or pr")
	}
	p, present, err := Load(root)
	if err != nil {
		return p, err
	}
	if present {
		if requested != "" && requested != p.Mode {
			return p, fmt.Errorf("existing %s selects %s; init never overwrites policy; change it under the current acceptance workflow", Path, p.Mode)
		}
		return p, nil
	}
	r, _, err := role.Load(root)
	if err != nil {
		return p, err
	}
	if requested != "" {
		p.Mode = requested
	}
	if r == role.Source && p.Mode != PR {
		return p, fmt.Errorf("source repositories require PR acceptance")
	}
	return p, nil
}
