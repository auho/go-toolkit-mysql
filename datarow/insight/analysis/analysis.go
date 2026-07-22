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
	FieldWidth []struct {
		width int
		zhLen int
	}
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
	_maxFieldLen := maxColumnShow.NameShowWidth + 1

	_format := fmt.Sprintf("%%-%ds %%-9s %%5s  %%-9s %%-9s", _maxFieldLen)
	ss = append(ss, fmt.Sprintf(_format, "COLUMN", "TYPE", "FLAG", "IS-EMPTY", "IS-NULL"))

	for _k, fieldName := range a.FieldsName {
		column := a.Columns[fieldName]

		warning := ""
		if column.Empty+column.Null >= column.Amount {
			warning = fmt.Sprintf("❗️")
		}

		_format = fmt.Sprintf("%%-%ds %%-9s %%5s  %%-9d %%-9d ", _maxFieldLen-columnsShow[_k].NameZhLen)
		ss = append(ss, fmt.Sprintf(_format, column.Column.Name, column.Column.FieldType, warning, column.Empty, column.Null))
	}

	ss = append(ss, "")
	return ss
}

// GetColumnsShow
// columns show
// max filed show width column show
func (a *Analysis) GetColumnsShow() ([]schema.ColumnShow, schema.ColumnShow) {
	var columnsShow []schema.ColumnShow
	for _, _fn := range a.FieldsName {
		columnsShow = append(columnsShow, schema.NewColumnShow(a.Columns[_fn].Column))
	}

	return columnsShow, slices.MaxFunc(columnsShow, func(i, j schema.ColumnShow) int {
		return cmp.Compare(i.NameShowWidth, j.NameShowWidth)
	})
}
