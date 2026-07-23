package tableprofile

import (
	"context"

	"github.com/auho/go-toolkit-mysql/tableprofile/compare"
)

// CompareTables profiles two tables and returns a Differ describing the
// differences between their analysis results.
func CompareTables(ctx context.Context, left, right Source) (*compare.Differ, error) {
	leftAnalysis, err := Explore(ctx, left)
	if err != nil {
		return nil, err
	}

	rightAnalysis, err := Explore(ctx, right)
	if err != nil {
		return nil, err
	}

	return compare.Diff(leftAnalysis, rightAnalysis), nil
}
