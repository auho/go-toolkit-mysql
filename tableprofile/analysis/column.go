package analysis

import "github.com/auho/go-toolkit-mysql/schema"

type Column struct {
	schema.Column
	RowCount   int
	Distinct int
	Empty    int
	Null     int
}
