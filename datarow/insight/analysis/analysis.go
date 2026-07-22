package analysis

import (
	"cmp"
	"fmt"
	"slices"

	"github.com/auho/go-toolkit-mysql/schema"
)

type Analysis struct {
	Table      *Table
	FieldNames []string
	Columns    map[string]Column
}

func NewAnalysis() *Analysis {
	a := &Analysis{}
	a.Columns = make(map[string]Column)

	return a
}

func (a *Analysis) ToStrings() []string {
	var ss []string

	ss = append(ss, fmt.Sprintf("table[%s]: %d", a.Table.Name, a.Table.Amount))

	columnsDisplay, maxColumnDisplay := a.ColumnsDisplay()
	maxFieldLen := maxColumnDisplay.NameDisplayWidth + 1

	format := fmt.Sprintf("%%-%ds %%-9s %%5s  %%-9s %%-9s", maxFieldLen)
	ss = append(ss, fmt.Sprintf(format, "COLUMN", "TYPE", "FLAG", "IS-EMPTY", "IS-NULL"))

	for k, fieldName := range a.FieldNames {
		column := a.Columns[fieldName]

		warning := ""
		if column.Empty+column.Null >= column.Amount {
			warning = fmt.Sprintf("❗️")
		}

		format = fmt.Sprintf("%%-%ds %%-9s %%5s  %%-9d %%-9d ", maxFieldLen-columnsDisplay[k].NameZhLen)
		ss = append(ss, fmt.Sprintf(format, column.Name, column.FieldType, warning, column.Empty, column.Null))
	}

	ss = append(ss, "")
	return ss
}

// ColumnsDisplay
// columns display
// max field display width column display
func (a *Analysis) ColumnsDisplay() ([]schema.ColumnDisplay, schema.ColumnDisplay) {
	var columnsDisplay []schema.ColumnDisplay
	for _, fn := range a.FieldNames {
		columnsDisplay = append(columnsDisplay, schema.NewColumnDisplay(a.Columns[fn].Column))
	}

	if len(columnsDisplay) == 0 {
		return columnsDisplay, schema.ColumnDisplay{}
	}

	return columnsDisplay, slices.MaxFunc(columnsDisplay, func(i, j schema.ColumnDisplay) int {
		return cmp.Compare(i.NameDisplayWidth, j.NameDisplayWidth)
	})
}
