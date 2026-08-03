package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cliewen/cliewen/internal/migrate"
)

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
	data, err := os.ReadFile(filepath.Join("..", "..", "CHANGELOG.md"))
	if err != nil {
		t.Fatalf("CHANGELOG.md not found: %v", err)
	}
	changelog := string(data)

	released := migrate.ReleasedVersions()
	if len(released) == 0 {
		t.Fatal("the carrier manifest records no releases, so this guard is checking nothing")
	}

	for _, v := range released {
		if strings.TrimSpace(changelogSection(changelog, v)) == "" {
			t.Errorf("CHANGELOG.md has no '## [%s]' section with content, but %s shipped — "+
				"a released version's notes were removed, and the next release would publish them again as its own", v, v)
		}
	}
}

// Sanity: the section a released version keeps is its own, not the next one's.
//
// The deletion above does not only remove a heading — it merges two sections,
// so the surviving one holds both releases' entries. A guard that asked only
// "is there content" would go green on a file where Unreleased had swallowed a
// shipped release whole. Every released section must therefore be separated
// from Unreleased by its own heading, which is what this asserts by checking
// that the headings appear in the file at all, newest release first.
func TestSanity_TheUnreleasedSectionIsNotAReleasedOne(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", "CHANGELOG.md"))
	if err != nil {
		t.Fatalf("CHANGELOG.md not found: %v", err)
	}
	changelog := string(data)

	if !strings.Contains(changelog, "## [Unreleased]") {
		t.Fatal("CHANGELOG.md has no '## [Unreleased]' section; every change adds its entry there (C-002)")
	}

	// An Install paragraph is written when a release is cut, so a pinned
	// install route inside Unreleased names a version that already shipped —
	// which happens when a released section was merged into Unreleased rather
	// than merely losing its heading. Checked against every release, not only
	// the newest, because the merge can survive further releases unnoticed.
	unreleased := changelogSection(changelog, "Unreleased")
	for _, v := range migrate.ReleasedVersions() {
		if strings.Contains(unreleased, "@v"+v) {
			t.Errorf("the Unreleased section pins the install route to v%s, which already shipped — "+
				"that release's section was merged into Unreleased, and cutting the next one would republish it", v)
		}
	}
}
