package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var releaseTagRE = regexp.MustCompile(`^v([0-9]+\.[0-9]+\.[0-9]+)$`)

// readChangelog reads the repository's changelog, the file the release
// workflow extracts every release body from.
func readChangelog(t *testing.T) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "..", "CHANGELOG.md"))
	if err != nil {
		t.Fatalf("CHANGELOG.md not found: %v", err)
	}
	return string(data)
}

// changelogSection returns the body of one version's section, by the same rule
// the release workflow extracts it with: everything after the "## [version]"
// heading, up to the next version heading. Kept deliberately identical to the
// awk in .github/scripts/release-gates.sh — a guard that read the file
// differently from the thing that publishes it would pass on notes that never
// ship, which is the failure it exists to prevent.
func changelogSection(changelog, version string) string {
	header := "## [" + version + "]"
	var body []string
	found := false
	for _, line := range strings.Split(changelog, "\n") {
		if !found {
			found = strings.HasPrefix(line, header)
			continue
		}
		if strings.HasPrefix(line, "## [") {
			break
		}
		body = append(body, line)
	}
	if !found {
		return ""
	}
	return strings.Join(body, "\n")
}

// releasedVersions returns every release tag reachable from this checkout.
// The CI checkout has full history, so this includes releases from before the
// carrier manifest existed. A source archive has no tags and cannot prove the
// repository's release history, so it skips this repository sanity guard.
func releasedVersions(t *testing.T) []string {
	t.Helper()
	root := filepath.Join("..", "..")
	out, err := exec.Command("git", "-C", root, "tag", "--merged", "HEAD", "--list", "v*").Output()
	if err != nil {
		t.Skipf("release tags unavailable in this checkout: %v", err)
	}

	var versions []string
	for _, tag := range strings.Fields(string(out)) {
		match := releaseTagRE.FindStringSubmatch(tag)
		if match != nil {
			versions = append(versions, match[1])
		}
	}
	if len(versions) == 0 {
		t.Skip("no release tags are reachable from this checkout")
	}
	return versions
}

func hasHeading(markdown, heading string) bool {
	for _, line := range strings.Split(markdown, "\n") {
		if strings.TrimSuffix(line, "\r") == heading {
			return true
		}
	}
	return false
}

// Sanity: every release that shipped still has its notes in CHANGELOG.md.
//
// The release gates check the section for the version being cut, and then
// nothing looks again. That leaves a window an ordinary change can walk into:
// immediately after a cut, "## [Unreleased]" is empty and sits directly on top
// of the new version heading, so a digest inserting entries under Unreleased
// can consume the heading line in the same hunk. That is exactly how 0.12.0's
// section was deleted, eleven hours after it was published — leaving its four
// shipped entries inside Unreleased, queued to be published a second time
// under whatever version came next.
//
// The gates cannot catch it: they exit early when the version stamp did not
// change, which is every ordinary change. So the promise belongs here, where
// it is checked on every run rather than only when a release is being made.
//
// A missing section is not cosmetic. It silently rewrites what the next
// release tells users, and it destroys the record of what a version contained
// for anyone reading the file rather than the forge.
func TestSanity_EveryReleasedVersionKeepsItsChangelogSection(t *testing.T) {
	changelog := readChangelog(t)

	for _, v := range releasedVersions(t) {
		if strings.TrimSpace(changelogSection(changelog, v)) == "" {
			t.Errorf("CHANGELOG.md has no '## [%s]' section with content, but %s shipped — "+
				"a released version's notes were removed, and the next release would publish them again as its own", v, v)
		}
	}
}

// Sanity: Unreleased does not still hold a release that already shipped.
//
// The deletion above does not only remove a heading — it merges two sections,
// so the survivor holds both releases' entries. The guard above catches that
// through the missing heading; this one reads the release-only structure in
// the surviving text instead, so the two fail from independent inputs.
func TestSanity_TheUnreleasedSectionIsNotAReleasedOne(t *testing.T) {
	changelog := readChangelog(t)

	if !strings.Contains(changelog, "## [Unreleased]") {
		t.Fatal("CHANGELOG.md has no '## [Unreleased]' section; every change adds its entry there (C-002)")
	}

	// The Install section is written when a release is cut. Its heading inside
	// Unreleased therefore means a released section was merged into Unreleased
	// rather than merely losing its heading. Look for the structure, not an
	// install-command substring that ordinary release prose may legitimately
	// mention.
	unreleased := changelogSection(changelog, "Unreleased")
	if hasHeading(unreleased, "### Install") {
		t.Error("the Unreleased section contains a release-only '### Install' section — " +
			"a released section was merged into Unreleased, and cutting the next one would republish it")
	}
}

func TestSanity_UnreleasedInstallGuardReadsStructure(t *testing.T) {
	legitimateMention := "## [Unreleased]\n\n- Upgrade from `go install example.test/tool@v0.12.0`.\n\n## [0.12.0] - 2026-08-02\n\nnotes"
	if hasHeading(changelogSection(legitimateMention, "Unreleased"), "### Install") {
		t.Fatal("an install command mentioned in ordinary Unreleased prose was mistaken for a released Install section")
	}

	mergedRelease := "## [Unreleased]\n\n### Install\n\n`go install example.test/tool@v0.12.0`\n\n## [0.11.2] - 2026-08-01\n\nnotes"
	if !hasHeading(changelogSection(mergedRelease, "Unreleased"), "### Install") {
		t.Fatal("a released Install section merged into Unreleased was not detected")
	}
}
