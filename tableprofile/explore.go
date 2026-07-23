package tableprofile

import (
	"context"

	"github.com/auho/go-toolkit-mysql/tableprofile/analysis"
	"github.com/auho/go-toolkit-mysql/tableprofile/explore"
)

func Explore(ctx context.Context, source Source) (*analysis.Result, error) {
	return explore.New(source.DB).Analyze(ctx, source.Name)
}
