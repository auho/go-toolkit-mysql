package insight

import (
	"context"

	"github.com/auho/go-toolkit-mysql/datarow/insight/diff"
)

func Diff(ctx context.Context, tables ...diff.Source) (*diff.Differ, error) {
	_left := tables[0]
	_right := tables[1]

	_leftAly, err := Explore(ctx, _left.DB, _left.Name)
	if err != nil {
		return nil, err
	}

	_rightAly, err := Explore(ctx, _right.DB, _right.Name)
	if err != nil {
		return nil, err
	}

	return diff.Diff(_leftAly, _rightAly), nil
}
