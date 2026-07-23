// Package analysis defines the result types produced by table profiling.
// A Result bundles a Table (with row count) and its Columns (with per-column
// statistics), and provides display helpers for tabular output.
package analysis

import (
	"github.com/auho/go-toolkit-mysql/schema"
)

// Table extends schema.Table with the row count of the analysed table.
type Table struct {
	schema.Table
	RowCount int // total number of rows in the table
}

// Column extends schema.Column with per-column statistics gathered during
// profiling.
type Column struct {
	schema.Column
	RowCount int // total rows in the table (same as Table.RowCount)
	Distinct int // number of distinct non-null values
	Empty    int // number of rows whose value is "empty" (0 for numbers, '' for strings)
	Null     int // number of NULL values
}
