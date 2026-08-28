package git

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Client struct {
	path string
}

func New(path ...string) *Client {
	c := &Client{}
	if len(path) > 0 {
		c.path = path[0]
	}

	return c
}

func (c *Client) Commit(message string) ([]byte, error) {
	cmd := c.exec("commit", "-m", message)

	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return output.Bytes(), fmt.Errorf("commit failed: %w", err)
	}

	return output.Bytes(), nil
}

func (c *Client) GetChangedFiles() ([]string, error) {
	cmd := c.exec("diff", "--staged", "--name-only")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get changed files: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")

	// filter empty lines
	var files []string
	for _, line := range lines {
		if line != "" {
			files = append(files, line)
		}
	}
	return files, nil
}

func (c *Client) GetStagedDiff() (string, error) {
	cmd := c.exec("diff", "--staged")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get staged diff: %w", err)
	}

	// []byte → string (utf-8 by default in Go)
	diff := strings.TrimSpace(string(out))
	if diff == "" {
		return "", fmt.Errorf("no staged changes found. Run `git add` first.")
	}

	return diff, nil
}

func (c *Client) Tags() ([]string, error) {
	cmd := c.exec("tag", "--list", "v*", "--sort=-v:refname")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list git tags: %w", err)
	}

	if strings.TrimSpace(string(out)) == "" {
		return []string{}, nil
	}

	return strings.Fields(string(out)), nil
}

func (c *Client) Push() error {
	cmd := c.exec("push")
	cmd.Stdout, cmd.Stderr, cmd.Stdin = os.Stdout, os.Stderr, os.Stdin

	return cmd.Run()
}

func (c *Client) PushTag(tag string) error {
	// Create tag locally
	cmd := c.exec("tag", tag)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(
			"failed to create tag %s: %w: %s",
			tag,
			err,
			strings.TrimSpace(string(out)),
		)
	}

	// Push tag to the origin
	cmd = c.exec("push", "origin", "refs/tags/"+tag)
	out, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(
			"failed to push tag %s: %w: %s",
			tag,
			err,
			strings.TrimSpace(string(out)),
		)
	}

	return cmd.Run()
}

func (c *Client) Validate() error {
	// Check if git is installed
	_, err := exec.LookPath("git")
	if err != nil {
		return fmt.Errorf("git is not installed or not in PATH.")
	}

	// Check if cwd is inside a valid git repository
	cmd := c.exec("rev-parse", "--is-inside-work-tree")
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		cwd := c.path
		if cwd == "" {
			cwd, _ = os.Getwd()
		}
		return fmt.Errorf("not a git repository: %s\n run: git init (to initialize one)", cwd)
	}

	// Check if there is at least one commit
	cmd = c.exec("rev-parse", "HEAD")
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("no commits yet in this repository")
	}

	return nil
}

func (c *Client) exec(args ...string) *exec.Cmd {
	cmd := exec.Command("git", args...)
	if c.path != "" {
		cmd.Dir = c.path
	}

	return cmd
}
