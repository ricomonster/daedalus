package daedalus

import "context"

type (
	Changes struct {
		Diff  string
		Files []string
	}

	StylusApplication interface {
		Commit(ctx context.Context, message string) ([]byte, error)

		GetChanges(ctx context.Context) (*Changes, error)

		GetCommitMessage(ctx context.Context, changes *Changes) (string, error)
	}
)
