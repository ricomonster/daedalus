package git

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Client struct{}

func New() *Client {
	return &Client{}
}

func (c *Client) Commit(message string) ([]byte, error) {
	cmd := exec.Command("git", "commit", "-m", message)

	var output bytes.Buffer
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return output.Bytes(), fmt.Errorf("commit failed: %w", err)
	}

	return output.Bytes(), nil
}

func (c *Client) GetChangedFiles() ([]string, error) {
	cmd := exec.Command("git", "diff", "--staged", "--name-only")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get changed files.")
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
	cmd := exec.Command("git", "diff", "--staged")
	out, err := cmd.Output() // captures stdout as []byte
	if err != nil {
		return "", fmt.Errorf("failed to get staged diff.")
	}

	// []byte → string (utf-8 by default in Go)
	diff := strings.TrimSpace(string(out))
	if diff == "" {
		return "", fmt.Errorf("no staged changes found. Run `git add` first.")
	}

	return diff, nil
}

func (c *Client) Push() error {
	cmd := exec.Command("git", "push")
	cmd.Stdout, cmd.Stderr, cmd.Stdin = os.Stdout, os.Stderr, os.Stdin

	return cmd.Run()
}

func (c *Client) Validate() error {
	// Check if git is installed
	_, err := exec.LookPath("git")
	if err != nil {
		return fmt.Errorf("git is not installed or not in PATH.")
	}

	// Check if cwd is inside a valid git repository
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		cwd, _ := os.Getwd()
		return fmt.Errorf("not a git repository: %s\n run: git init (to initialize one)", cwd)
	}

	// Check if there is at least one commit
	cmd = exec.Command("git", "rev-parse", "HEAD")
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("no commits yet in this repository")
	}

	return nil
}
