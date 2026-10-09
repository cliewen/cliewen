package accept

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Read blobs directly so export-ignore/export-subst attributes and checkout
// filters cannot change what the judge sees. Symlinks and gitlinks require a
// separate materialization contract; v1 refuses them rather than following
// paths outside the isolated snapshot or pretending submodule files exist.
func materialize(root, revision, destination string) error {
	entries, err := command(root, "ls-tree", "-rz", revision).Output()
	if err != nil {
		return err
	}
	c := command(root, "cat-file", "--batch")
	input, err := c.StdinPipe()
	if err != nil {
		return err
	}
	output, err := c.StdoutPipe()
	if err != nil {
		return err
	}
	if err = c.Start(); err != nil {
		return err
	}
	defer func() { input.Close(); _ = c.Wait() }()
	r := bufio.NewReader(output)
	for _, entry := range bytes.Split(entries, []byte{0}) {
		if len(entry) == 0 {
			continue
		}
		meta, path, ok := strings.Cut(string(entry), "\t")
		fields := strings.Fields(meta)
		if !ok || len(fields) != 3 || fields[1] != "blob" || (fields[0] != "100644" && fields[0] != "100755") {
			return fmt.Errorf("local acceptance v1 requires regular tracked files; unsupported entry %q", path)
		}
		if !filepath.IsLocal(filepath.FromSlash(path)) || strings.Contains(path, "\\") {
			return fmt.Errorf("unsafe tracked path %q", path)
		}
		for _, part := range strings.Split(path, "/") {
			if strings.EqualFold(part, ".git") {
				return fmt.Errorf("unsafe tracked path %q", path)
			}
		}
		if _, err = io.WriteString(input, fields[2]+"\n"); err != nil {
			return err
		}
		header, e := r.ReadString('\n')
		if e != nil {
			return e
		}
		info := strings.Fields(header)
		if len(info) != 3 || info[0] != fields[2] || info[1] != "blob" {
			return fmt.Errorf("missing candidate blob %s", fields[2])
		}
		size, e := strconv.ParseInt(info[2], 10, 64)
		if e != nil || size < 0 {
			return fmt.Errorf("invalid candidate blob size")
		}
		name := filepath.Join(destination, filepath.FromSlash(path))
		if err = os.MkdirAll(filepath.Dir(name), 0700); err != nil {
			return err
		}
		f, e := os.OpenFile(name, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
		if e != nil {
			return e
		}
		_, err = io.CopyN(f, r, size)
		closeErr := f.Close()
		if err != nil {
			return err
		}
		if closeErr != nil {
			return closeErr
		}
		if b, e := r.ReadByte(); e != nil || b != '\n' {
			return fmt.Errorf("invalid candidate blob terminator")
		}
	}
	return nil
}
