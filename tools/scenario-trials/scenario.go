package main

import (
	"fmt"
	"sort"
	"strings"
)

// Scenario is one situation the method has to handle. Besides its prompt and
// fixture, it states the obligation it exercises and the failure that
// obligation prevents, so a run can be read as evidence about that obligation.
// Its prompt is an ordinary request and names none of the rules it tests; the
// checks live here, outside the container, where the agent under test cannot
// read them.
type Scenario struct {
	Name       string
	Obligation string
	Failure    string
	// Keys fixes the order of the observed fields in a run's signature.
	Keys    []string
	Observe func(t Transcript, post string, sc scan) map[string]string
}

var scenarios = map[string]Scenario{}

func registerScenario(s Scenario) { scenarios[s.Name] = s }

func scenarioFor(name string) (Scenario, error) {
	if s, ok := scenarios[name]; ok {
		return s, nil
	}
	return Scenario{}, fmt.Errorf("scenario %q does not exist; available: %s", name, strings.Join(scenarioNames(), ", "))
}

func scenarioNames() []string {
	names := make([]string, 0, len(scenarios))
	for n := range scenarios {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

// Prompt is the request given to the agent, on standard input.
func (s Scenario) Prompt() (string, error) {
	b, err := assets.ReadFile("scenarios/" + s.Name + "/prompt.txt")
	return string(b), err
}

// Setup is the shell script that builds the fixture repository in the container.
func (s Scenario) Setup() (string, error) {
	b, err := assets.ReadFile("scenarios/" + s.Name + "/setup.sh")
	return string(b), err
}

// Check scores one run: what every run reports, and what this scenario observes.
func (s Scenario) Check(t Transcript, postStatus string) Outcome {
	sc := scanTranscript(t)
	return Outcome{
		Observed: s.Observe(t, postStatus, sc), VersionCheckFirst: sc.VersionCheckFirst,
		Model: t.Model, AgentVersion: t.AgentVersion, MCPServers: t.MCPServers,
		Turns: t.Turns, CostUSD: t.CostUSD, DurationMS: t.DurationMS, Finished: t.Finished,
	}
}

// Signature groups runs that behaved the same way on what the scenario asks.
func (s Scenario) Signature(o Outcome) string {
	parts := make([]string, 0, len(s.Keys))
	for _, k := range s.Keys {
		parts = append(parts, k+"="+or(o.Observed[k], "?"))
	}
	return strings.Join(parts, " ")
}
