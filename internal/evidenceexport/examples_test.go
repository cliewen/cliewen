package evidenceexport

import (
	"encoding/json"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cliewen/cliewen/internal/evidence"
)

func python(t *testing.T) string {
	t.Helper()
	for _, name := range []string{"python3", "python"} {
		if p, err := exec.LookPath(name); err == nil && exec.Command(p, "--version").Run() == nil {
			return p
		}
	}
	t.Fatal("Python is required to verify shipped exporter examples")
	return ""
}

func example(t *testing.T, name string) string {
	t.Helper()
	p, err := filepath.Abs(filepath.Join("../scaffold/templates/evidence", name))
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func producerExample(t *testing.T, script, source, content string) (int, int) {
	t.Helper()
	root := t.TempDir()
	full := filepath.Join(root, source)
	if err := os.WriteFile(full, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	input, _ := json.Marshal(map[string]any{"root": root, "files": []string{source}})
	cmd := exec.Command(python(t), example(t, script))
	cmd.Stdin = strings.NewReader(string(input))
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("producer: %s %v", out, err)
	}
	var result struct {
		References  []evidence.Reference  `json:"references"`
		Diagnostics []evidence.Diagnostic `json:"diagnostics"`
	}
	if err := json.Unmarshal(out, &result); err != nil {
		t.Fatalf("invalid producer result: %s %v", out, err)
	}
	for _, ref := range result.References {
		if ref.Path != source || ref.Type != "Unit" && ref.Type != "Performance" || !evidence.CanonicalID(ref.ID) {
			t.Fatalf("bad metadata: %+v", ref)
		}
	}
	return len(result.References), len(result.Diagnostics)
}

func TestAC207_UnitPositive_ExamplesBindDirectMetadataAndFallback(t *testing.T) {
	for _, tc := range []struct {
		script, source, text string
		want                 int
	}{
		{"python_producer.py", "test_example.py", "class Tests:\n @cliewen(ac='PDO-115', type='Unit', direction='positive')\n def test_positive(self): pass\n @pytest.mark.cliewen(ac='PDO-115', type='Unit', direction='negative')\n def test_negative(self): pass\n def test_PDO_115_UnitPositive_named(self): pass\n def test_helper(self): pass\n", 3},
		{"gatling_producer.py", "ExampleSimulation.java", `class ExampleSimulation extends Simulation { ScenarioBuilder scn = scenario("[PDO-115 Performance Positive] good"); ExampleSimulation(){setUp(scn.injectOpen(atOnceUsers(1))).assertions(global().failedRequests().count().is(0L));}}`, 1},
	} {
		t.Run(tc.script, func(t *testing.T) {
			refs, diags := producerExample(t, tc.script, tc.source, tc.text)
			if refs != tc.want || diags != 0 {
				t.Fatalf("refs=%d diagnostics=%d", refs, diags)
			}
		})
	}
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "example_test.go"), []byte("package example\n// func TestAC999_UnitPositive_comment(t *testing.T) {}\nvar fixture = `func TestAC999_UnitPositive_string(t *testing.T) {}`\nfunc TestAC001_UnitPositive_real(t *testing.T) {}\n"), 0644); err != nil {
		t.Fatal(err)
	}
	p, err := Collect(root, "go")
	if err != nil || len(p.References) != 1 || len(p.Diagnostics) != 0 {
		t.Fatalf("AST attribution: %+v %v", p, err)
	}
}

func TestAC207_UnitNegative_ExamplesRejectContainersAmbiguityAndUnregisteredScenarios(t *testing.T) {
	for _, tc := range []struct {
		script, source, text string
		wantDiagnostics      bool
	}{
		{"python_producer.py", "test_example.py", "# AC: PDO-115\ndef test_comments():\n \"\"\"def test_PDO_115_UnitPositive_fake(): pass\"\"\"\n pass\n", false},
		{"python_producer.py", "test_example.py", "@cliewen(ac='PDO-115', type='Unit', direction='positive')\nclass Tests:\n def test_member(self): pass\n", true},
		{"python_producer.py", "test_example.py", "@cliewen(ac=dynamic, type='Unit', direction='positive')\ndef test_member(): pass\n", true},
		{"python_producer.py", "test_example.py", "@cliewen(ac='PDO-115', type='Unit', direction='positive')\n@pytest.mark.cliewen(ac='PDO-116', type='Unit', direction='positive')\ndef test_member(): pass\n", true},
		{"gatling_producer.py", "ExampleSimulation.java", `class ExampleSimulation extends Simulation { ScenarioBuilder scn = scenario("[PDO-115 Performance Positive] unregistered"); }`, true},
		{"gatling_producer.py", "ExampleSimulation.java", `class ExampleSimulation extends Simulation { /* scenario("[PDO-115 Performance Positive] comment") */ String fake="scenario(\"[PDO-115 Performance Positive] string\")"; }`, false},
	} {
		t.Run(tc.script+tc.text[:12], func(t *testing.T) {
			refs, diags := producerExample(t, tc.script, tc.source, tc.text)
			if refs != 0 || (diags > 0) != tc.wantDiagnostics {
				t.Fatalf("refs=%d diagnostics=%d", refs, diags)
			}
		})
	}
}

func TestAC204_UnitPositive_ShippedAggregateProducesJudgeReadableDeterministicBytes(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "tools"), 0755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"aggregate.py", "python_producer.py"} {
		data, err := os.ReadFile(example(t, name))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, "tools", name), data, 0644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(root, "test_example.py"), []byte("@cliewen(ac='PDO-115',type='Unit',direction='positive')\ndef test_example(): pass\n"), 0644); err != nil {
		t.Fatal(err)
	}
	config := map[string]any{"producers": []any{map[string]any{"id": "python", "include": []string{"test_*.py", "tools/*.py"}, "command": []string{python(t), "tools/python_producer.py"}}}}
	data, _ := json.Marshal(config)
	if err := os.WriteFile(filepath.Join(root, "producers.json"), data, 0644); err != nil {
		t.Fatal(err)
	}
	run := func() {
		cmd := exec.Command(python(t), "tools/aggregate.py", "producers.json")
		cmd.Dir = root
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("aggregate: %s %v", out, err)
		}
	}
	run()
	before, _ := os.ReadFile(filepath.Join(root, ".clue/evidence.yaml"))
	run()
	after, _ := os.ReadFile(filepath.Join(root, ".clue/evidence.yaml"))
	if string(before) != string(after) {
		t.Fatal("aggregate includes unstable state")
	}
	if m, err := evidence.Load(root); err != nil || len(m.Producers[0].References) != 1 {
		t.Fatalf("export incompatible with judge: %+v %v", m, err)
	}
	// A new source is visible even if no active criterion refers to it yet.
	if err := os.WriteFile(filepath.Join(root, "test_added.py"), []byte("def test_helper():pass\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := evidence.Load(root); err == nil {
		t.Fatal("missed added source")
	}
}

func TestAC204_UnitNegative_ShippedAggregateDoesNotPublishFailedProducer(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".clue"), 0755); err != nil {
		t.Fatal(err)
	}
	previous := []byte("previous complete export")
	if err := os.WriteFile(filepath.Join(root, ".clue/evidence.yaml"), previous, 0644); err != nil {
		t.Fatal(err)
	}
	config := `{"producers":[{"id":"failed","include":["producers.json"],"command":["` + strings.ReplaceAll(python(t), "\\", "\\\\") + `","-c","raise RuntimeError('producer failed')"]}]}`
	if err := os.WriteFile(filepath.Join(root, "producers.json"), []byte(config), 0644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(python(t), example(t, "aggregate.py"), "producers.json")
	cmd.Dir = root
	if out, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("accepted failed producer: %s", out)
	}
	after, _ := os.ReadFile(filepath.Join(root, ".clue/evidence.yaml"))
	if string(previous) != string(after) {
		t.Fatal("published partial evidence")
	}
}

func TestAC207_UnitPositive_VitestNativeTagsAndParameterizedDeclarations(t *testing.T) {
	vitestExample(t, false)
}
func TestAC207_UnitNegative_VitestSuiteAndAmbiguousMetadataGiveDiagnostics(t *testing.T) {
	vitestExample(t, true)
}

func vitestExample(t *testing.T, negative bool) {
	t.Helper()
	node, err := exec.LookPath("node")
	if err != nil {
		t.Fatal("Node is required for Vitest example checks")
	}
	modulePath := filepath.ToSlash(example(t, "vitest_producer.mjs"))
	if !strings.HasPrefix(modulePath, "/") {
		modulePath = "/" + modulePath
	}
	moduleURL := (&url.URL{Scheme: "file", Path: modulePath}).String()

	script := `import {collectModules} from ` + fmtJSON(moduleURL) + `; const root=process.cwd(); const file=root+'/test.spec.ts'; const test={type:'test',name:'example',options:{tags:['PDO-115','Unit','positive']}};`
	if negative {
		script += `const children=[{type:'suite',name:'container',options:{tags:['PDO-115']},children:[test]},{...test,name:'ambiguous',options:{tags:['PDO-115','PDO-116','Unit','positive']}}];`
	} else {
		script += `const each={...test,options:{...test.options,each:true},location:{line:4,column:1}};const children=[test,each,{...each,name:'second instance'}];`
	}
	script += `const result=collectModules([{moduleId:file,children}],root,['test.spec.ts']); console.log(JSON.stringify(result));`
	cmd := exec.Command(node, "--input-type=module", "-e", script)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Vitest example: %s %v", out, err)
	}
	var result struct {
		References  []evidence.Reference
		Diagnostics []evidence.Diagnostic
	}
	if err := json.Unmarshal(out, &result); err != nil {
		t.Fatal(err)
	}
	if negative {
		if len(result.Diagnostics) != 2 {
			t.Fatalf("missing attribution failures: %s", out)
		}
	} else if len(result.References) != 2 || len(result.Diagnostics) != 0 {
		t.Fatalf("bad native/each evidence: %s", out)
	}
}

func fmtJSON(value string) string { data, _ := json.Marshal(value); return string(data) }
