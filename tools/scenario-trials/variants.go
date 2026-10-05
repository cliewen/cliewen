package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/fs"
	"sort"
	"strings"
)

// A method variant is a named removal applied to the container's copy of the
// instructions after the fixture is built, so a scenario can be run with one
// obligation absent. It never touches the repository's own skills or hub: it
// runs inside the throwaway fixture repository only. "baseline" is no removal.
const baselineVariant = "baseline"

// variantScript returns the shell script of a variant and its hash, which the
// conditions record. The script runs in the fixture repository's directory.
func variantScript(name string) (script, hash string, err error) {
	if name == baselineVariant {
		return "", "", nil
	}
	b, err := assets.ReadFile("variants/" + name + ".sh")
	if err != nil {
		return "", "", fmt.Errorf("method variant %q does not exist; available: %s", name, strings.Join(variantNames(), ", "))
	}
	sum := sha256.Sum256(b)
	return string(b), hex.EncodeToString(sum[:]), nil
}

func variantNames() []string {
	names := []string{baselineVariant}
	entries, _ := fs.ReadDir(assets, "variants")
	for _, e := range entries {
		if n := strings.TrimSuffix(e.Name(), ".sh"); n != e.Name() {
			names = append(names, n)
		}
	}
	sort.Strings(names)
	return names
}
