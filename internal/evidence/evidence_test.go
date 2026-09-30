package evidence

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func put(t *testing.T, root, name, text string) {
	t.Helper()
	full := filepath.Join(root, filepath.FromSlash(name))
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte(text), 0o644); err != nil {
		t.Fatal(err)
	}
}

func producer(t *testing.T, root, id, source string) Producer {
	t.Helper()
	p := Producer{ID: id, Include: []string{source}, References: []Reference{{ID: "PDO-115", Path: source, Subject: "preserves_pages", Type: "Unit", Direction: "positive"}}}
	var err error
	p.Inputs, err = Snapshot(root, p)
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func TestAC203_UnitPositive_CommonContractKeepsCanonicalMetadata(t *testing.T) {
	root := t.TempDir()
	put(t, root, "src/test.any", "opaque test source")
	p := producer(t, root, "unknown-framework", "src/test.any")
	p.Framework = "a framework the judge has never seen"
	p.Language = "a future language"
	m := Manifest{Version: Version, Producers: []Producer{p}}
	if err := Write(root, m); err != nil {
		t.Fatal(err)
	}
	got, err := Load(root)
	if err != nil || got.Producers[0].References[0] != p.References[0] {
		t.Fatalf("got %+v, %v", got, err)
	}
	// Unclassified legacy references omit both fields, not an invented class.
	p.References[0].Type, p.References[0].Direction = "", ""
	if err := CheckShape(Manifest{Version: Version, Producers: []Producer{p}}); err != nil {
		t.Fatal(err)
	}
}

func TestAC203_UnitNegative_InvalidReferencesDoNotPassShape(t *testing.T) {
	root := t.TempDir()
	put(t, root, "test.any", "source")
	for _, tc := range []struct {
		name   string
		change func(*Producer)
	}{
		{"bad id", func(p *Producer) { p.References[0].ID = "pdo-115" }},
		{"no subject", func(p *Producer) { p.References[0].Subject = " " }},
		{"external source", func(p *Producer) { p.References[0].Path = "../other/test.any" }},
		{"unfingerprinted source", func(p *Producer) { p.References[0].Path = "other.any" }},
		{"half classification", func(p *Producer) { p.References[0].Direction = "" }},
		{"bad type", func(p *Producer) { p.References[0].Type = "Human" }},
		{"bad direction", func(p *Producer) { p.References[0].Direction = "Positive" }},
		{"duplicate", func(p *Producer) { p.References = append(p.References, p.References[0]) }},
		{"conflicting id", func(p *Producer) { r := p.References[0]; r.ID = "PDO-116"; p.References = append(p.References, r) }},
		{"bad input", func(p *Producer) { p.Inputs[0].SHA256 = "not a hash" }},
		{"duplicate input", func(p *Producer) { p.Inputs = append(p.Inputs, p.Inputs[0]) }},
		{"bad diagnostic", func(p *Producer) { p.Diagnostics = []Diagnostic{{Path: "../test", Message: "broken"}} }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			p := producer(t, root, "suite", "test.any")
			tc.change(&p)
			if err := CheckShape(Manifest{Version: 1, Producers: []Producer{p}}); err == nil {
				t.Fatal("accepted invalid reference")
			}
		})
	}
}

func TestAC204_UnitPositive_MultipleProducersAreDeterministic(t *testing.T) {
	root := t.TempDir()
	put(t, root, "backend/test.any", "backend")
	put(t, root, "frontend/test.any", "frontend")
	a := producer(t, root, "backend", "backend/test.any")
	b := producer(t, root, "frontend", "frontend/test.any")
	b.References[0].Direction = "negative"
	first, err := Encode(Manifest{Version: 1, Producers: []Producer{b, a}})
	if err != nil {
		t.Fatal(err)
	}
	second, err := Encode(Manifest{Version: 1, Producers: []Producer{a, b}})
	if err != nil || !bytes.Equal(first, second) {
		t.Fatalf("nondeterministic: %v", err)
	}
	if err := Write(root, Manifest{Version: 1, Producers: []Producer{b, a}}); err != nil {
		t.Fatal(err)
	}
	got, err := Load(root)
	if err != nil || len(got.Producers) != 2 {
		t.Fatalf("lost producer: %+v %v", got, err)
	}
	if got.Producers[0].References[0].Subject != got.Producers[1].References[0].Subject {
		t.Fatal("fixture must exercise identical test names")
	}
}

func TestAC204_UnitNegative_FailedAggregatePreservesCompleteExport(t *testing.T) {
	root := t.TempDir()
	put(t, root, "test.any", "source")
	p := producer(t, root, "suite", "test.any")
	m := Manifest{Version: 1, Producers: []Producer{p}}
	if err := Write(root, m); err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(DefaultPath)))
	if err != nil {
		t.Fatal(err)
	}
	for _, change := range []func(*Manifest){
		func(m *Manifest) { m.Producers = append(m.Producers, m.Producers[0]) },
		func(m *Manifest) {
			m.Producers[0].Diagnostics = []Diagnostic{{Path: "test.any", Message: "producer failed"}}
		},
		func(m *Manifest) { m.Producers[0].Inputs[0].SHA256 = strings.Repeat("0", 64) },
	} {
		p := producer(t, root, "suite", "test.any")
		m := Manifest{Version: 1, Producers: []Producer{p}}
		change(&m)
		if err := Write(root, m); err == nil {
			t.Fatal("published failed aggregate")
		}
		after, _ := os.ReadFile(filepath.Join(root, filepath.FromSlash(DefaultPath)))
		if !bytes.Equal(before, after) {
			t.Fatal("failed producer replaced complete evidence")
		}
	}
}

func TestAC205_UnitPositive_LineEndingsAndExcludedOutputsKeepEvidenceCurrent(t *testing.T) {
	root := t.TempDir()
	put(t, root, "src/test.any", "line\nnext\n")
	p := producer(t, root, "suite", "src/test.any")
	p.Include = []string{"src/**"}
	p.Exclude = []string{"src/build/**"}
	if err := Write(root, Manifest{Version: 1, Producers: []Producer{p}}); err != nil {
		t.Fatal(err)
	}
	put(t, root, "src/test.any", "line\r\nnext\r\n")
	put(t, root, "src/build/test.any", "generated")
	if _, err := Load(root); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		pattern, name string
		want          bool
	}{{"src/**/*.any", "src/test.any", true}, {"src/**/*.any", "src/nested/test.any", true}, {"src/*.any", "src/nested/test.any", false}, {"src/?est.[ab]ny", "src/test.any", true}} {
		if got := Match(tc.pattern, tc.name); got != tc.want {
			t.Errorf("%s %s: %v", tc.pattern, tc.name, got)
		}
	}
}

func TestAC205_UnitNegative_ChangedAddedDeletedOrUnsafeInputsInvalidateEvidence(t *testing.T) {
	for _, kind := range []string{"changed", "added", "deleted", "exporter"} {
		t.Run(kind, func(t *testing.T) {
			root := t.TempDir()
			put(t, root, "src/test.any", "source")
			put(t, root, "tools/export.any", "exporter")
			p := producer(t, root, "suite", "src/test.any")
			p.Include = []string{"src/**", "tools/export.any"}
			p.Inputs, _ = Snapshot(root, p)
			if err := Write(root, Manifest{Version: 1, Producers: []Producer{p}}); err != nil {
				t.Fatal(err)
			}
			switch kind {
			case "changed":
				put(t, root, "src/test.any", "changed")
			case "added":
				put(t, root, "src/other.any", "added")
			case "deleted":
				if err := os.Remove(filepath.Join(root, "src/test.any")); err != nil {
					t.Fatal(err)
				}
			case "exporter":
				put(t, root, "tools/export.any", "changed")
			}
			if _, err := Load(root); err == nil {
				t.Fatal("credited stale evidence")
			}
		})
	}
	root := t.TempDir()
	for _, pattern := range []string{"../outside/**", "/absolute/**", "C:/outside", "src/**bad", "src/[bad", DefaultPath} {
		if _, err := Snapshot(root, Producer{ID: "suite", Include: []string{pattern}}); err == nil {
			t.Errorf("accepted unsafe scope %s", pattern)
		}
	}
	if _, err := ReadFile(root, "../outside"); err == nil {
		t.Fatal("accepted unsafe source")
	}
	if _, err := Snapshot(root, Producer{ID: "suite", Include: []string{"missing-exact.file"}}); err == nil {
		t.Fatal("ignored missing exact input")
	}
}

func TestAC206_UnitPositive_UnknownFrameworkNeedsNoInstalledRunner(t *testing.T) {
	root := t.TempDir()
	put(t, root, "tests/test.future", "does not execute")
	p := producer(t, root, "future", "tests/test.future")
	if err := Write(root, Manifest{Version: 1, Producers: []Producer{p}}); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(root); err != nil {
		t.Fatal(err)
	}
}

func TestAC206_UnitNegative_InvalidExchangesGiveNoReferences(t *testing.T) {
	root := t.TempDir()
	put(t, root, "test.any", "source")
	if _, err := Load(root); err == nil {
		t.Fatal("accepted missing export")
	}
	for _, data := range []string{"version: 99\nproducers: []\n", "version: 1\nversion: 1\n", "version: 1\nunknown: true\n", "version: 1\nproducers: []\n---\nversion: 1\n", "version: [broken\n"} {
		put(t, root, DefaultPath, data)
		if m, err := Load(root); err == nil || len(m.Producers) > 0 {
			t.Fatalf("credited invalid exchange %q: %+v %v", data, m, err)
		}
	}
	p := producer(t, root, "suite", "test.any")
	p.Diagnostics = []Diagnostic{{Path: "test.any", Subject: "test", Message: "ambiguous AC identity"}}
	data, err := Encode(Manifest{Version: 1, Producers: []Producer{p}})
	if err != nil {
		t.Fatal(err)
	}
	put(t, root, DefaultPath, string(data))
	if m, err := Load(root); err == nil || !strings.Contains(err.Error(), "ambiguous AC identity") || len(m.Producers) > 0 {
		t.Fatalf("credited diagnostic: %+v %v", m, err)
	}
}

func TestUnit_EvidenceSymlinksNeverEscapeRepository(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	put(t, outside, "test.any", "outside")
	if err := os.Symlink(outside, filepath.Join(root, "linked")); err != nil {
		t.Skipf("symlink creation unavailable: %v", err)
	}
	if _, err := ReadFile(root, "linked/test.any"); err == nil {
		t.Fatal("read through symlink")
	}
	for _, pattern := range []string{"linked/**", "**/*.any"} {
		if _, err := Snapshot(root, Producer{ID: "suite", Include: []string{pattern}}); err == nil {
			t.Fatalf("walked symlink scope %s", pattern)
		}
	}
}
