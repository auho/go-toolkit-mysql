package analysis

import (
	"cmp"
	"fmt"
	"slices"

	"github.com/auho/go-toolkit-mysql/schema"
)

type Result struct {
	Table      *Table
	FieldNames []string
	Columns    map[string]Column
}

func NewResult() *Result {
	r := &Result{}
	r.Columns = make(map[string]Column)

	return r
}

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

// ColumnsDisplay
// columns display
// max field display width column display
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
