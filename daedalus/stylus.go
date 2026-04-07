package daedalus

import "context"

type (
	Changes struct {
		Diff  string
		Files []string
	}

	StylusApplication interface {
		Commit(ctx context.Context, message string) error
		GetChanges(ctx context.Context) (*Changes, error)

		GetCommitMessage(ctx context.Context, changes *Changes) (string, error)
	}
)

func TruncateDiff(diff string, max int) string {
	if len(diff) > max {
		return diff[:max]
	}

	return diff
}
