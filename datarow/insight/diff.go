package insight

import (
	"context"

	"github.com/auho/go-toolkit-mysql/datarow/insight/diff"
)

func Diff(ctx context.Context, left, right diff.Source) (*diff.Differ, error) {
	leftAly, err := Explore(ctx, left.DB, left.Name)
	if err != nil {
		return nil, err
	}

	rightAly, err := Explore(ctx, right.DB, right.Name)
	if err != nil {
		return nil, err
	}

	return diff.Diff(leftAly, rightAly), nil
}
