package guide

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

const (
	configPath = ".vitepress/config.mts"
	homePage   = "index.md"
)

var (
	internalConfigLink = regexp.MustCompile(`link:\s*"(/[^"]*)"`)
	nextTargetLink     = regexp.MustCompile(`(?s)## Next\n\n\[[^\]]+\]\(([^\)]+)\)\n\z`)
)

// walkNextChain follows one page's single Next action to the next page, starting
// at start. It returns the pages visited in order and the target that leaves the
// chain. A revisited page is a cycle, which is the defect this reports.
func walkNextChain(start string, chain map[string]string) ([]string, string, error) {
	visited := map[string]bool{}
	var order []string
	page := start
	for {
		target, ok := chain[page]
		if !ok {
			return order, "", fmt.Errorf("%s has no Next action", page)
		}
		if visited[page] {
			return order, "", fmt.Errorf("the Next chain returns to %s, so the reading path contains a cycle", page)
		}
		visited[page] = true
		order = append(order, page)
		if _, isPage := chain[target]; !isPage {
			return order, target, nil
		}
		page = target
	}
}

func guidePages(t *testing.T) []string {
	t.Helper()
	pages, err := filepath.Glob("*.md")
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(pages)
	return pages
}

// pageFor resolves a guide link — "/corpus", "./corpus", or "./corpus#anchor" —
// to the Markdown file that serves it. An external link resolves to itself.
func pageFor(link string) string {
	if strings.HasPrefix(link, "http") {
		return link
	}
	path := strings.TrimSuffix(strings.SplitN(link, "#", 2)[0], "/")
	path = strings.TrimPrefix(strings.TrimPrefix(path, "."), "/")
	if path == "" {
		return homePage
	}
	return path + ".md"
}

func TestAC036_UnitPositive_GuideNavigationIsAcyclicAndComplete(t *testing.T) {
	pages := guidePages(t)
	isPage := map[string]bool{}
	for _, page := range pages {
		isPage[page] = true
	}

	config, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatal(err)
	}
	inSidebar := map[string]bool{}
	for _, match := range internalConfigLink.FindAllStringSubmatch(string(config), -1) {
		page := pageFor(match[1])
		if strings.HasSuffix(page, ".md") && !isPage[page] {
			t.Errorf("%s links to %q, but %s does not exist", configPath, match[1], page)
			continue
		}
		inSidebar[page] = true
	}
	for _, page := range pages {
		if page == homePage {
			continue // the home page is reached through the site logo and title.
		}
		if !inSidebar[page] {
			t.Errorf("%s exists but no navigation entry in %s links to it", page, configPath)
		}
	}

	chain := map[string]string{}
	for _, page := range pages {
		content, err := os.ReadFile(page)
		if err != nil {
			t.Fatal(err)
		}
		match := nextTargetLink.FindSubmatch(content)
		if match == nil {
			t.Errorf("%s must end with exactly one Next action", page)
			continue
		}
		chain[page] = pageFor(string(match[1]))
	}

	var terminals []string
	for page, target := range chain {
		if _, isChained := chain[target]; !isChained {
			terminals = append(terminals, fmt.Sprintf("%s → %s", page, target))
		}
	}
	sort.Strings(terminals)
	if len(terminals) != 1 {
		t.Errorf("the guide needs exactly one page whose Next leaves the reading path, found %d: %s", len(terminals), strings.Join(terminals, ", "))
	}

	order, exit, err := walkNextChain(homePage, chain)
	if err != nil {
		t.Fatalf("walking Next from %s failed after %v: %v", homePage, order, err)
	}
	if len(order) != len(pages) {
		visited := map[string]bool{}
		for _, page := range order {
			visited[page] = true
		}
		var unreachable []string
		for _, page := range pages {
			if !visited[page] {
				unreachable = append(unreachable, page)
			}
		}
		t.Errorf("walking Next from %s reaches %d of %d pages; never reached: %s", homePage, len(order), len(pages), strings.Join(unreachable, ", "))
	}
	if !strings.HasPrefix(exit, "http") {
		t.Errorf("the reading path should end by leaving the guide, it ends at %q", exit)
	}
}

func TestAC036_UnitNegative_GuideNavigationRejectsACycle(t *testing.T) {
	cyclic := map[string]string{
		"index.md":           "getting-started.md",
		"getting-started.md": "skills.md",
		"skills.md":          "getting-started.md",
	}
	order, exit, err := walkNextChain("index.md", cyclic)
	if err == nil {
		t.Fatalf("a Next chain that returns to an earlier page must be rejected, walked %v and exited at %q", order, exit)
	}
	if !strings.Contains(err.Error(), "cycle") {
		t.Errorf("the reported defect must name the cycle, got %q", err)
	}

	orphaned := map[string]string{
		"index.md":  "a.md",
		"a.md":      "https://example.invalid",
		"orphan.md": "https://example.invalid",
	}
	order, _, err = walkNextChain("index.md", orphaned)
	if err != nil {
		t.Fatal(err)
	}
	if len(order) == len(orphaned) {
		t.Errorf("a page nobody links to must stay unreached, walked %v", order)
	}
}
