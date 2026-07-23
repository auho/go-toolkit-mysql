package analysis

import (
	"cmp"
	"fmt"
	"slices"

	"github.com/auho/go-toolkit-mysql/schema"
)

// Result is the outcome of profiling a single table. It holds the table-level
// metadata, the ordered list of field names, and per-column statistics.
type Result struct {
	Table      *Table
	FieldNames []string
	Columns    map[string]Column
}

// NewResult creates an empty Result with an initialised Columns map.
func NewResult() *Result {
	r := &Result{}
	r.Columns = make(map[string]Column)

	return r
}

// ToStrings renders the profiling result as a slice of human-readable lines
// suitable for terminal output. Each column line shows its name, field type,
// a warning flag (when all values are empty or null), and empty/null counts.
func (r *Result) ToStrings() []string {
	var lines []string

	lines = append(lines, fmt.Sprintf("table[%s]: %d", r.Table.Name, r.Table.RowCount))

	columnsDisplay, maxColumnDisplay := r.ColumnsDisplay()
	maxFieldLen := maxColumnDisplay.NameDisplayWidth + 1

	format := fmt.Sprintf("%%-%ds %%-9s %%5s  %%-9s %%-9s", maxFieldLen)
	lines = append(lines, fmt.Sprintf(format, "COLUMN", "TYPE", "FLAG", "IS-EMPTY", "IS-NULL"))

	for k, fieldName := range r.FieldNames {
		column := r.Columns[fieldName]

		warning := ""
		if column.Empty+column.Null >= column.RowCount {
			warning = fmt.Sprintf("❗️")
		}

		format = fmt.Sprintf("%%-%ds %%-9s %%5s  %%-9d %%-9d ", maxFieldLen-columnsDisplay[k].NameZhLen)
		lines = append(lines, fmt.Sprintf(format, column.Name, column.FieldType, warning, column.Empty, column.Null))
	}

	lines = append(lines, "")
	return lines
}

// ColumnsDisplay returns display metadata for every column in FieldNames order,
// plus the ColumnDisplay with the largest name display width (used for column
// alignment). If there are no columns, the second return value is the zero
// ColumnDisplay.
func (r *Result) ColumnsDisplay() ([]schema.ColumnDisplay, schema.ColumnDisplay) {
	var columnsDisplay []schema.ColumnDisplay
	for _, fn := range r.FieldNames {
		columnsDisplay = append(columnsDisplay, schema.NewColumnDisplay(r.Columns[fn].Column))
	}

	if len(columnsDisplay) == 0 {
		return columnsDisplay, schema.ColumnDisplay{}
	}

	return columnsDisplay, slices.MaxFunc(columnsDisplay, func(i, j schema.ColumnDisplay) int {
		return cmp.Compare(i.NameDisplayWidth, j.NameDisplayWidth)
	})
}
