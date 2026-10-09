// Package acceptpolicy reads repository-owned acceptance policy without Git
// history, network access, or writes. Init, migration and acceptance share it.
package acceptpolicy

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
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

// Load preserves the historical PR convention when the policy is absent.
// Fresh init materializes local explicitly, so absence never upgrades a repo
// into a different acceptance workflow.
func Load(root string) (Policy, bool, error) {
	p := Policy{Mode: PR, Branch: "main"}
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
	} else {
		prior, e := PreviouslyAdopted(root)
		if e != nil {
			return p, e
		}
		if !prior {
			p.Mode = Local
		}
	}
	if r == role.Source && p.Mode != PR {
		return p, fmt.Errorf("source repositories require PR acceptance")
	}
	return p, nil
}

// PreviouslyAdopted is conservative: Cliewen machine state, canonical skills,
// corpus artifacts or index markers retain PR. Ordinary Git/project files do
// not count as adoption. Symlinked convention trees never become fresh roots.
func PreviouslyAdopted(root string) (bool, error) {
	for _, name := range []string{".clue/role.yaml", ".clue/id-ledger.yaml", ".clue/evidence.yaml"} {
		if _, err := os.Lstat(filepath.Join(root, name)); err == nil {
			return true, nil
		} else if !os.IsNotExist(err) {
			return false, err
		}
	}
	for _, dir := range []string{".clue", ".agents/skills", ".claude/skills", "docs", "changes"} {
		if info, err := os.Lstat(filepath.Join(root, dir)); err == nil && info.Mode()&os.ModeSymlink != 0 {
			return true, nil
		} else if err != nil && !os.IsNotExist(err) {
			return false, err
		}
	}
	for _, name := range []string{"clue-analysis", "clue-delta", "clue-extract", "clue-plan", "clue-upgrade", "clue-verify"} {
		for _, file := range []string{filepath.Join(".agents/skills", name, "skill.md"), filepath.Join(".claude/skills", name, "SKILL.md")} {
			if _, err := os.Lstat(filepath.Join(root, file)); err == nil {
				return true, nil
			} else if !os.IsNotExist(err) {
				return false, err
			}
		}
	}
	for _, folder := range []string{"docs", "changes"} {
		found := false
		err := filepath.WalkDir(filepath.Join(root, folder), func(name string, d fs.DirEntry, err error) error {
			if os.IsNotExist(err) {
				return nil
			}
			if err != nil {
				return err
			}
			if found {
				return fs.SkipAll
			}
			if d.IsDir() {
				return nil
			}
			if d.Type()&os.ModeSymlink != 0 {
				found = true
				return fs.SkipAll
			}
			if !strings.HasSuffix(name, ".md") {
				return nil
			}
			data, e := os.ReadFile(name)
			if e != nil {
				return e
			}
			text := strings.ReplaceAll(strings.TrimPrefix(string(data), "\ufeff"), "\r\n", "\n")
			if strings.Contains(text, "<!-- clue:index:start -->") {
				found = true
				return fs.SkipAll
			}
			if !strings.HasPrefix(text, "---\n") {
				return nil
			}
			header, _, ok := strings.Cut(text[4:], "\n---")
			if !ok {
				return nil
			}
			var meta map[string]any
			if yaml.Unmarshal([]byte(header), &meta) != nil {
				// Malformed Cliewen-shaped data is ambiguous, never fresh.
				found = true
				for _, key := range []string{"id:", "type:", "status:", "links:", "title:"} {
					found = found && strings.Contains(header, key)
				}
			} else {
				found = true
				for _, key := range []string{"id", "type", "status", "links", "title"} {
					_, ok := meta[key]
					found = found && ok
				}
			}
			if found {
				return fs.SkipAll
			}
			return nil
		})
		if err != nil {
			return false, err
		}
		if found {
			return true, nil
		}
	}
	return false, nil
}
