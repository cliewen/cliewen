package corpus

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

// skillsDir is where a Cliewen repo keeps its agent skills alongside any
// unrelated skills the adopter installs. A Cliewen-managed skill declares
// cliewen-skill: true before its version joins the checked set (ADR-022).
const (
	skillsDir          = ".agents/skills"
	cliewenSkillMarker = "cliewen-skill"
)

// These names predate ADR-022's ownership marker. Reserving the exact legacy
// slots lets a new binary fail toward reinstalling old Cliewen skills instead
// of silently treating a pre-marker installation as an empty managed set.
var legacyCliewenSkillNames = map[string]bool{
	"clue-analysis": true,
	"clue-delta":    true,
	"clue-extract":  true,
	"clue-plan":     true,
	"clue-verify":   true,
}

// checkSkillVersions enforces, in order of what each build can see:
//
//   - only skills marked cliewen-skill: true join the managed set; malformed
//     markers fail, while unmarked third-party skills are ignored (AC-029);
//   - an unmarked canonical legacy slot fails toward reinstall (AC-030);
//   - every managed skill carries a version stamp (AC-031);
//   - managed skills agree on one version — "versioned as a set" via
//     per-skill markers (AC-032);
//   - a released binary (version is neither "" nor "dev") matches the
//     managed skills, else the difference is drift (AC-033). dev/source
//     builds skip this last comparison — they have no release to drift from.
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
		dir := filepath.Join(root, e.Name())
		manifest, ambiguous, derr := findSkillManifest(dir)
		if derr != nil {
			continue // an unreadable directory is not a skill
		}
		if ambiguous {
			// Only a case-sensitive filesystem can hold two entries that
			// fold to skill.md; report it rather than resolve by read order
			// so the verdict does not depend on the filesystem (ADR-028).
			issues = append(issues, Issue{path.Join(skillsDir, e.Name()), "skill directory holds multiple case-variants of skill.md; keep exactly one manifest"})
			continue
		}
		if manifest == "" {
			continue // a folder without a skill.md is not a skill
		}
		rel := path.Join(skillsDir, e.Name(), manifest)
		data, rerr := os.ReadFile(filepath.Join(dir, manifest))
		if rerr != nil {
			continue // a folder without a readable skill.md is not a skill
		}
		text := strings.ReplaceAll(string(data), "\r\n", "\n")
		fields, _, ok, perr := parseFrontmatter(text)
		declaresOwnership := frontmatterDeclares(text, cliewenSkillMarker)
		if perr != nil {
			if legacyCliewenSkillNames[e.Name()] || declaresOwnership {
				issues = append(issues, Issue{rel, "skill frontmatter does not parse: " + perr.Error()})
			}
			continue
		}
		if !ok {
			if declaresOwnership {
				issues = append(issues, Issue{rel, "skill frontmatter does not parse: ownership marker is inside an unterminated YAML block"})
			} else if legacyCliewenSkillNames[e.Name()] {
				issues = append(issues, Issue{rel, "legacy Cliewen skill carries no cliewen-skill: true marker (reinstall the skills from this clue release)"})
			}
			continue
		}
		marker, marked := fields[cliewenSkillMarker]
		if !marked {
			if legacyCliewenSkillNames[e.Name()] {
				issues = append(issues, Issue{rel, "legacy Cliewen skill carries no cliewen-skill: true marker (reinstall the skills from this clue release)"})
			}
			continue
		}
		owns, markerIsBool := marker.(bool)
		if !markerIsBool || !owns {
			issues = append(issues, Issue{rel, fmt.Sprintf("cliewen-skill marker must be boolean true, got %v (%T)", marker, marker)})
			continue
		}
		v := ""
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

	// AC-032: the managed skills must agree on one version. The reference is
	// the version most skills carry (ties go to the earliest-sorted holder),
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

	// AC-033: drift against a released binary. Only meaningful once the
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

// findSkillManifest resolves a skill's manifest inside dir by matching an entry
// whose name case-folds to "skill.md", so a manifest named SKILL.md, Skill.md,
// or skill.md is found identically on a case-sensitive and a case-insensitive
// filesystem alike (ADR-028). Only an entry resolving to a regular file is a
// manifest: a directory, FIFO, or device named skill.md is not one, and
// reading it could block the judge rather than fail it. A symlink to a regular
// file does count — a shared skills tree is a supported shape. It returns the
// actual entry name, whether more than one case-variant is present — only
// reachable on a case-sensitive filesystem, and reported by the caller rather
// than silently resolved — and any directory-read error.
func findSkillManifest(dir string) (name string, ambiguous bool, err error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", false, err
	}
	for _, e := range entries {
		if !strings.EqualFold(e.Name(), "skill.md") {
			continue
		}
		// Stat follows a symlink; an unstattable entry is simply not a
		// manifest, exactly as an unreadable one was before.
		info, serr := os.Stat(filepath.Join(dir, e.Name()))
		if serr != nil || !info.Mode().IsRegular() {
			continue
		}
		if name != "" {
			ambiguous = true
			continue
		}
		name = e.Name()
	}
	return name, ambiguous, nil
}

// frontmatterDeclares is used only when YAML parsing fails: an unparseable
// third-party skill remains outside Cliewen's scope, but a file that visibly
// attempts to declare Cliewen ownership must not escape with a broken block.
func frontmatterDeclares(text, key string) bool {
	if !strings.HasPrefix(text, "---\n") {
		return false
	}
	rest := text[len("---\n"):]
	end := strings.Index(rest, "\n---\n")
	if end < 0 {
		if strings.HasSuffix(rest, "\n---") {
			end = len(rest) - len("\n---")
		} else {
			end = len(rest)
		}
	}
	for _, line := range strings.Split(rest[:end], "\n") {
		if strings.HasPrefix(line, key+":") {
			return true
		}
	}
	return false
}
