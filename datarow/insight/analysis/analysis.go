package analysis

import (
	"cmp"
	"fmt"
	"slices"

	"github.com/auho/go-toolkit-mysql/schema"
)

type Analysis struct {
	Table      *Table
	FieldsName []string
	Columns    map[string]Column
}

func NewAnalysis() *Analysis {
	a := &Analysis{}
	a.Columns = make(map[string]Column)

	return a
}

func (a *Analysis) ToStrings() []string {
	var ss []string

	ss = append(ss, fmt.Sprintf("table[%s]: %d", a.Table.Table.Name, a.Table.Amount))

	columnsShow, maxColumnShow := a.GetColumnsShow()
	maxFieldLen := maxColumnShow.NameShowWidth + 1

	format := fmt.Sprintf("%%-%ds %%-9s %%5s  %%-9s %%-9s", maxFieldLen)
	ss = append(ss, fmt.Sprintf(format, "COLUMN", "TYPE", "FLAG", "IS-EMPTY", "IS-NULL"))

	for k, fieldName := range a.FieldsName {
		column := a.Columns[fieldName]

		warning := ""
		if column.Empty+column.Null >= column.Amount {
			warning = fmt.Sprintf("❗️")
		}

		format = fmt.Sprintf("%%-%ds %%-9s %%5s  %%-9d %%-9d ", maxFieldLen-columnsShow[k].NameZhLen)
		ss = append(ss, fmt.Sprintf(format, column.Column.Name, column.Column.FieldType, warning, column.Empty, column.Null))
	}

	ss = append(ss, "")
	return ss
}

// GetColumnsShow
// columns show
// max filed show width column show
func (a *Analysis) GetColumnsShow() ([]schema.ColumnShow, schema.ColumnShow) {
	var columnsShow []schema.ColumnShow
	for _, fn := range a.FieldsName {
		columnsShow = append(columnsShow, schema.NewColumnShow(a.Columns[fn].Column))
	}

	if len(columnsShow) == 0 {
		return columnsShow, schema.ColumnShow{}
	}

	return columnsShow, slices.MaxFunc(columnsShow, func(i, j schema.ColumnShow) int {
		return cmp.Compare(i.NameShowWidth, j.NameShowWidth)
	})
}
