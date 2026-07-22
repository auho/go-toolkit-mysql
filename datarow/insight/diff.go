package insight

import (
	"github.com/auho/go-toolkit/v2/mysql/datarow/insight/diff"
)

func Diff(tables ...diff.Source) (*diff.Differ, error) {
	_left := tables[0]
	_right := tables[1]

	_leftAly, err := Explore(_left.DB, _left.Name)
	if err != nil {
		return nil, err
	}

	_rightAly, err := Explore(_right.DB, _right.Name)
	if err != nil {
		return nil, err
	}

	return diff.Diff(_leftAly, _rightAly), nil
}
