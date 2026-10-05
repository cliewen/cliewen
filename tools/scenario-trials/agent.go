package main

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// Adapter is everything specific to one agent. The runner, the checks, the
// conditions record and the summary know only this interface and Transcript,
// so a vendor's name appears nowhere outside its own adapter file.
type Adapter interface {
	Name() string
	// DefaultLogin is where the login lives when none is given, under the maintainer's home.
	DefaultLogin(home string) string
	// Credentials returns what logs the agent in inside the container. Secret
	// values travel only in Env, never in an argument list.
	Credentials(o Options) (Credentials, error)
	// Command is run in the fixture repository inside the container. It reads
	// the prompt on stdin and writes the agent's event stream to /out/events.jsonl.
	Command(model string) string
	// Probe is a shell fragment that writes the agent's version to
	// /out/post-agent.txt and its model, if it can tell, to /out/post-model.txt.
	Probe() string
	// Parse turns the agent's event stream into a Transcript.
	Parse(r io.Reader) (Transcript, error)
}

// Credentials is how an adapter supplies a login from outside the repository.
// The login Options.Login names is whatever the adapter needs: a file, a directory.
type Credentials struct {
	Env    []string // NAME=value for the docker child process only
	Pass   []string // names from Env that docker forwards into the container
	Mounts []string // docker --mount specifications
}

var adapters = map[string]Adapter{}

func register(a Adapter) { adapters[a.Name()] = a }

func adapterFor(name string) (Adapter, error) {
	if a, ok := adapters[name]; ok {
		return a, nil
	}
	return nil, fmt.Errorf("agent %q has no adapter; available: %s", name, strings.Join(adapterNames(), ", "))
}

func adapterNames() []string {
	names := make([]string, 0, len(adapters))
	for n := range adapters {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}
