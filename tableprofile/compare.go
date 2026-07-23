package tableprofile

import (
	"context"

	simpledb "github.com/auho/go-simple-db/v3"
	"github.com/auho/go-toolkit-mysql/tableprofile/compare"
)

type Source struct {
	Name string // table name
	DB   *simpledb.SimpleDB
}

func CompareTables(ctx context.Context, left, right Source) (*compare.Differ, error) {
	leftAnalysis, err := Explore(ctx, left.DB, left.Name)
	if err != nil {
		return nil, err
	}

	rightAnalysis, err := Explore(ctx, right.DB, right.Name)
	if err != nil {
		return nil, err
	}

	return compare.Diff(leftAnalysis, rightAnalysis), nil
}
