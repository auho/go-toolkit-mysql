package analysis

import (
	"github.com/auho/go-toolkit-mysql/schema"
)

type Table struct {
	schema.Table
	RowCount int
}

type Column struct {
	schema.Column
	RowCount int
	Distinct int
	Empty    int
	Null     int
}
