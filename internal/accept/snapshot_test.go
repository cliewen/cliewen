package accept

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func snapshotFixture(t *testing.T, links map[string]string) (string, string) {
	t.Helper()
	root := t.TempDir()
	mustGit(t, root, "init", "-b", "main")
	mustGit(t, root, "config", "user.name", "Snapshot fixture")
	mustGit(t, root, "config", "user.email", "snapshot@example.invalid")
	mustGit(t, root, "config", "commit.gpgsign", "false")
	write(t, root, "regular/data.txt", "committed bytes\n")
	write(t, root, "regular/nested/note.md", "committed note\n")
	mustGit(t, root, "add", ".")
	for name, target := range links {
		c := command(root, "hash-object", "-w", "--stdin")
		c.Stdin = strings.NewReader(target)
		out, e := c.Output()
		if e != nil {
			t.Fatal(e)
		}
		mustGit(t, root, "update-index", "--add", "--cacheinfo", "120000,"+strings.TrimSpace(string(out))+","+name)
	}
	mustGit(t, root, "commit", "-m", "Snapshot fixture")
	return root, mustGit(t, root, "rev-parse", "HEAD")
}

func TestAC236_IntegrationPositive_InternalFileDirectoryAndChainLinksUseCommittedContent(t *testing.T) {
	root, rev := snapshotFixture(t, map[string]string{"file-link": "regular/data.txt", "mirror": "regular", "chain": "mirror/nested/note.md", "relative/entry": "../regular/data.txt"})
	write(t, root, "regular/data.txt", "uncommitted replacement")
	destination := t.TempDir()
	if e := materialize(root, rev, destination); e != nil {
		t.Fatal(e)
	}
	for name, want := range map[string]string{"file-link": "committed bytes\n", "mirror/data.txt": "committed bytes\n", "mirror/nested/note.md": "committed note\n", "chain": "committed note\n", "relative/entry": "committed bytes\n"} {
		p := filepath.Join(destination, filepath.FromSlash(name))
		got, e := os.ReadFile(p)
		if e != nil || string(got) != want {
			t.Fatalf("%s: %q %v", name, got, e)
		}
		info, e := os.Lstat(p)
		if e != nil || info.Mode()&os.ModeSymlink != 0 {
			t.Fatalf("snapshot follows an OS link: %s", name)
		}
	}
	if got := mustGit(t, root, "ls-tree", rev, "file-link"); !strings.HasPrefix(got, "120000") {
		t.Fatal("candidate link mode changed")
	}
}

func TestAC236_IntegrationNegative_ExternalDanglingMetadataAndCyclesFailClosed(t *testing.T) {
	cases := map[string]map[string]string{
		"escape": {"link": "../outside"}, "absolute": {"link": "/outside"}, "drive": {"link": "C:/outside"}, "backslash": {"link": "..\\outside"}, "metadata": {"link": ".git/config"}, "empty": {"link": ""}, "dangling": {"link": "missing"}, "file-parent": {"link": "regular/data.txt/child"}, "cycle": {"a": "b", "b": "a"}, "ancestor": {"regular/self": "."}, "metadata-parent": {"link": "regular/../.GiT/config"},
	}
	for name, links := range cases {
		t.Run(name, func(t *testing.T) {
			root, rev := snapshotFixture(t, links)
			if e := materialize(root, rev, t.TempDir()); e == nil {
				t.Fatal("unsafe snapshot accepted")
			}
		})
	}
	root, rev := snapshotFixture(t, nil)
	mustGit(t, root, "update-index", "--add", "--cacheinfo", "160000,"+rev+",submodule")
	mustGit(t, root, "commit", "-m", "Add unmaterialized submodule")
	if e := materialize(root, mustGit(t, root, "rev-parse", "HEAD"), t.TempDir()); e == nil {
		t.Fatal("submodule accepted")
	}
}

func TestAC236_IntegrationNegative_ExcessiveAcyclicLinkExpansionIsBounded(t *testing.T) {
	links := map[string]string{}
	// Each level doubles the prior committed subtree without a cycle.
	for level := 1; level <= 11; level++ {
		previous := "regular"
		if level > 1 {
			previous = fmt.Sprintf("level%d", level-1)
		}
		links[fmt.Sprintf("level%d/a", level)] = "../" + previous
		links[fmt.Sprintf("level%d/b", level)] = "../" + previous
	}
	root, rev := snapshotFixture(t, links)
	if e := materialize(root, rev, t.TempDir()); e == nil || !strings.Contains(e.Error(), "budget") {
		t.Fatalf("unbounded expansion: %v", e)
	}
}

func TestAC236_IntegrationNegative_OversizedLinkBlobFailsWithoutBlockedGitReader(t *testing.T) {
	root, rev := snapshotFixture(t, map[string]string{"link": strings.Repeat("x", 100000)})
	if e := materialize(root, rev, t.TempDir()); e == nil {
		t.Fatal("oversized link accepted")
	}
}

func TestAC236_IntegrationPositive_ParentTraversalFollowsCommittedDirectoryLinks(t *testing.T) {
	root, _ := snapshotFixture(t, map[string]string{"alias": "regular/nested", "link": "alias/../data.txt"})
	write(t, root, "data.txt", "wrong root bytes\n")
	mustGit(t, root, "add", "data.txt")
	mustGit(t, root, "commit", "-m", "Distinct root target")
	destination := t.TempDir()
	if e := materialize(root, mustGit(t, root, "rev-parse", "HEAD"), destination); e != nil {
		t.Fatal(e)
	}
	got, e := os.ReadFile(filepath.Join(destination, "link"))
	if e != nil || string(got) != "committed bytes\n" {
		t.Fatalf("wrong resolved bytes: %q %v", got, e)
	}
}

func TestAC236_IntegrationNegative_ParentTraversalCannotEraseInvalidComponents(t *testing.T) {
	for _, target := range []string{"missing/../regular/data.txt", "regular/data.txt/../data.txt", "alias/../../../outside", "alias/../.git/../data.txt"} {
		t.Run(target, func(t *testing.T) {
			root, rev := snapshotFixture(t, map[string]string{"alias": "regular/nested", "link": target})
			if e := materialize(root, rev, t.TempDir()); e == nil {
				t.Fatal("invalid traversal accepted")
			}
		})
	}
}

func TestAC236_IntegrationNegative_ExcessiveAcyclicResolutionIsBounded(t *testing.T) {
	links := map[string]string{"b0": "regular"}
	for level := 1; level <= 30; level++ {
		previous := fmt.Sprintf("b%d", level-1)
		links[fmt.Sprintf("b%d", level)] = previous + "/../" + previous
	}
	root, rev := snapshotFixture(t, links)
	if e := materialize(root, rev, t.TempDir()); e == nil || !strings.Contains(e.Error(), "resolution exceeds snapshot budget") {
		t.Fatalf("unbounded resolution: %v", e)
	}
}
