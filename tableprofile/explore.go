package tableprofile

import (
	"context"

	"github.com/auho/go-toolkit-mysql/tableprofile/analysis"
	"github.com/auho/go-toolkit-mysql/tableprofile/explore"
)

// Explore analyses a single table and returns its profiling result, including
// row count and per-column statistics (distinct, empty, null counts).
func Explore(ctx context.Context, source Source) (*analysis.Result, error) {
	return explore.New(source.DB).Analyze(ctx, source.Name)
}
