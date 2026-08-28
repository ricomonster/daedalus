package daedalus

import (
	"context"
	"fmt"
	"strings"

	"github.com/ricomonster/daedalus/gemini"
	"github.com/ricomonster/daedalus/git"
)

type (
	Changes struct {
		Diff  string
		Files []string
	}
)

func NewStylusService(gi *git.Client, ge *gemini.Client) *StylusService {
	return &StylusService{gi, ge}
}

type StylusService struct {
	git    *git.Client
	gemini *gemini.Client
}

func (c *StylusService) Commit(ctx context.Context, message string) ([]byte, error) {
	return c.git.Commit(message)
}

func (c *StylusService) GetChanges(ctx context.Context) (*Changes, error) {
	// Validate first
	err := c.git.Validate()
	if err != nil {
		return nil, err
	}

	// Get staged diff
	diff, err := c.git.GetStagedDiff()
	if err != nil {
		return nil, err
	}

	files, err := c.git.GetChangedFiles()
	if err != nil {
		return nil, err
	}

	return &Changes{Diff: diff, Files: files}, nil
}

func (c *StylusService) GetCommitMessage(ctx context.Context, changes *Changes) (string, error) {
	prompt := fmt.Sprintf(`Generate a conventional commit message for these changes.
Files: %s 
Diff: %s 
Respond with ONLY the commit message.`,
		strings.Join(changes.Files, " "),
		changes.Diff,
	)

	// Run gemini
	return c.gemini.Prompt(ctx, prompt)
}
