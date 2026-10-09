package accept

import (
	"io"
	"os"
	"os/exec"
	"strings"
)

// Each command owns an empty hooks directory for its complete lifetime,
// including prepared reference transactions. No repository config is changed.
// Filesystem-monitor callbacks are disabled independently of ordinary hooks.
type gitCommand struct {
	*exec.Cmd
	hooks   string
	initErr error
}

func command(root string, args ...string) *gitCommand {
	hooks, err := os.MkdirTemp("", "clue-accept-hooks-")
	options := []string{"-C", root, "-c", "core.hooksPath=" + hooks, "-c", "core.fsmonitor=false"}
	c := &gitCommand{Cmd: exec.Command("git", append(options, args...)...), hooks: hooks, initErr: err}
	// Fail closed if a private directory cannot be created; never execute with
	// an empty hooks path, which could resolve hooks relative to the checkout.
	for _, e := range os.Environ() {
		if !strings.HasPrefix(e, "GIT_") {
			c.Env = append(c.Env, e)
		}
	}
	c.Env = append(c.Env, "GIT_NO_REPLACE_OBJECTS=1", "GIT_OPTIONAL_LOCKS=0")
	return c
}

func (c *gitCommand) cleanup() {
	if c.hooks != "" {
		_ = os.RemoveAll(c.hooks)
	}
}

func (c *gitCommand) Output() ([]byte, error) {
	defer c.cleanup()
	if c.initErr != nil {
		return nil, c.initErr
	}
	return c.Cmd.Output()
}

func (c *gitCommand) CombinedOutput() ([]byte, error) {
	defer c.cleanup()
	if c.initErr != nil {
		return nil, c.initErr
	}
	return c.Cmd.CombinedOutput()
}

func (c *gitCommand) StdinPipe() (io.WriteCloser, error) {
	if c.initErr != nil {
		return nil, c.initErr
	}
	return c.Cmd.StdinPipe()
}

func (c *gitCommand) StdoutPipe() (io.ReadCloser, error) {
	if c.initErr != nil {
		return nil, c.initErr
	}
	return c.Cmd.StdoutPipe()
}

func (c *gitCommand) Start() error {
	if c.initErr != nil {
		return c.initErr
	}
	if err := c.Cmd.Start(); err != nil {
		c.cleanup()
		return err
	}
	return nil
}

func (c *gitCommand) Wait() error {
	defer c.cleanup()
	return c.Cmd.Wait()
}
