package corpus

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

// skillsDir is where a Cliewen repo keeps its agent skills. Each skill.md
// carries a version stamp in its frontmatter (G-002/M-004); clue compares
// those stamps to the running binary's release version to make drift
// lintable (ADR-011).
const skillsDir = ".agents/skills"

// checkSkillVersions enforces, in order of what each build can see:
//
//   - every skill carries a version stamp (AC-020);
//   - the skills agree on one version — "versioned as a set" via per-skill
//     markers (AC-021);
//   - a released binary (version is neither "" nor "dev") matches the
//     skills, else the difference is drift (AC-022). dev/source builds skip
//     this last comparison — they have no release to drift from.
//
// A repo with no skills folder has nothing to check.
func checkSkillVersions(c *Corpus, binaryVersion string) []Issue {
	root := filepath.Join(c.Root, filepath.FromSlash(skillsDir))
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil
	}

	type skill struct{ path, version string }
	var skills []skill
	var issues []Issue
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		rel := path.Join(skillsDir, e.Name(), "skill.md")
		data, rerr := os.ReadFile(filepath.Join(root, e.Name(), "skill.md"))
		if rerr != nil {
			continue // a folder without a skill.md is not a skill
		}
		text := strings.ReplaceAll(string(data), "\r\n", "\n")
		fields, _, ok, perr := parseFrontmatter(text)
		if perr != nil {
			issues = append(issues, Issue{rel, "skill frontmatter does not parse: " + perr.Error()})
			continue
		}
		v := ""
		if ok {
			switch t := fields["version"].(type) {
			case nil:
			case string:
				v = t
			default:
				// YAML reads e.g. `version: 1.0` as a number; the stamp
				// must be a semver string so it compares against the
				// binary's.
				issues = append(issues, Issue{rel, fmt.Sprintf("skill version %v is a YAML %T, not a string — stamp bare semver (e.g. 0.1.0)", t, t)})
				continue
			}
		}
		if v == "" {
			issues = append(issues, Issue{rel, "skill carries no version stamp"})
			continue
		}
		skills = append(skills, skill{rel, v})
	}
	if len(skills) == 0 {
		return issues
	}
	sort.Slice(skills, func(i, j int) bool { return skills[i].path < skills[j].path })

	// AC-021: the skills must agree on one version. The reference is the
	// version most skills carry (ties go to the earliest-sorted holder),
	// so the report names the actual outlier rather than whichever skill
	// happens to sort first.
	count := make(map[string]int, len(skills))
	for _, s := range skills {
		count[s.version]++
	}
	ref := skills[0].version
	for _, s := range skills[1:] {
		if count[s.version] > count[ref] {
			ref = s.version
		}
	}
	consistent := true
	for _, s := range skills {
		if s.version != ref {
			consistent = false
			issues = append(issues, Issue{s.path, "skill version " + s.version + " disagrees with the set's version " + ref})
		}
	}

	// AC-022: drift against a released binary. Only meaningful once the
	// skills agree; otherwise the disagreement above is the finding.
	if consistent && binaryVersion != "" && binaryVersion != "dev" {
		for _, s := range skills {
			if s.version != binaryVersion {
				issues = append(issues, Issue{s.path, "skill version " + s.version + " != clue " + binaryVersion + " (drift — reinstall the skills or clue)"})
			}
		}
	}
	return issues
}
