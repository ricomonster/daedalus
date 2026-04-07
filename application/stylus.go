package application

import (
	"context"
	"fmt"
	"strings"

	"github.com/ricomonster/daedalus/daedalus"
	"github.com/ricomonster/daedalus/git"
)

type stylusApplication struct {
	git *git.Client
	// oracleApp daedalus.OracleApplication
	gemini daedalus.LLMApplication
}

func NewStylusApplication(gi *git.Client, ge daedalus.LLMApplication) daedalus.StylusApplication {
	return &stylusApplication{gi, ge}
}

func (sa *stylusApplication) Commit(ctx context.Context, message string) error {
	return sa.git.Commit(message)
}

func (sa *stylusApplication) GetChanges(ctx context.Context) (*daedalus.Changes, error) {
	// Validate first
	err := sa.git.Validate()
	if err != nil {
		return nil, err
	}

	// Get staged diff
	diff, err := sa.git.GetStagedDiff()
	if err != nil {
		return nil, err
	}

	files, err := sa.git.GetChangedFiles()
	if err != nil {
		return nil, err
	}

	return &daedalus.Changes{Diff: diff, Files: files}, nil
}

func (sa *stylusApplication) GetCommitMessage(
	ctx context.Context,
	changes *daedalus.Changes,
) (string, error) {
	// This is only for Gemini. For other LLM implementation, this might be different
	prompt := fmt.Sprintf(`Generate a conventional commit message for these changes.
Files: %s 
Diff: %s 
Respond with ONLY the commit message.`,
		strings.Join(changes.Files, " "),
		daedalus.TruncateDiff(changes.Diff, 8000))

	// Run gemini
	return sa.gemini.Prompt(ctx, prompt)
}
