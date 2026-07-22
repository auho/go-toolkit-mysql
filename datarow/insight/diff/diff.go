package diff

import (
	"fmt"

	"github.com/auho/go-toolkit-mysql/datarow/insight/analysis"
)

func Diff(left, right *analysis.Analysis) *Differ {
	d := &Differ{}
	d.diff(left, right)

	return d
}

type Differ struct {
	results []string
	ok      bool
}

func (d *Differ) IsOK() bool {
	return d.ok
}

func (d *Differ) Differences() []string {
	return d.results
}

func (d *Differ) diff(left, right *analysis.Analysis) {
	d.ok = true

	var results []string

	lColumnsDisplay, lMaxColumnDisplay := left.ColumnsDisplay()
	lMaxDisplayWidth := lMaxColumnDisplay.NameDisplayWidth + 1
	rColumnsDisplay, rMaxColumnDisplay := right.ColumnsDisplay()
	rMaxDisplayWidth := rMaxColumnDisplay.NameDisplayWidth + 1

	maxDisplayWidth := max(lMaxDisplayWidth, rMaxDisplayWidth)

	// table name and amount
	title := fmt.Sprintf("table[%s:%s] amount", left.Table.Name, right.Table.Name)
	if left.Table.Amount == right.Table.Amount {
		results = append(results, d.success(fmt.Sprintf("%s: %d", title, left.Table.Amount)))
	} else {
		results = append(results, d.failure(fmt.Sprintf("%s[%d != %d]", title, left.Table.Amount, right.Table.Amount)))
	}

	// loop left field
	for k, leftFieldName := range left.FieldNames {
		leftColumn := left.Columns[leftFieldName]

		// left field title
		title = fmt.Sprintf(fmt.Sprintf("  %%-%ds", maxDisplayWidth-lColumnsDisplay[k].NameZhLen), leftColumn.Name)

		if rightColumn, ok := right.Columns[leftColumn.Name]; ok {
			// compare amount
			if leftColumn.Amount == rightColumn.Amount {
				results = append(results, d.success(fmt.Sprintf("%s amount: %d", title, leftColumn.Amount)))
			} else {
				results = append(results, d.failure(fmt.Sprintf("%s amount: [%d != %d]", title, leftColumn.Amount, rightColumn.Amount)))
			}

			// compare distinct
			if leftColumn.Distinct != rightColumn.Distinct {
				results = append(results, d.failure(fmt.Sprintf("%s distinct: [%d != %d]", title, leftColumn.Distinct, rightColumn.Distinct)))
			}

			// compare empty
			if leftColumn.Empty != rightColumn.Empty {
				results = append(results, d.failure(fmt.Sprintf("%s empty: [%d != %d]", title, leftColumn.Empty, rightColumn.Empty)))
			}

			// compare null
			if leftColumn.Null != rightColumn.Null {
				results = append(results, d.failure(fmt.Sprintf("%s null: [%d != %d]", title, leftColumn.Null, rightColumn.Null)))
			}

		} else {
			// in left, not in right
			results = append(results, d.onlyInLeft(fmt.Sprintf("%s amount: [%d != 0]", title, leftColumn.Amount)))
		}
	}

	// loop right field
	for k, rightFieldName := range right.FieldNames {
		rightColumn := right.Columns[rightFieldName]

		// right field title
		title = fmt.Sprintf(fmt.Sprintf("  %%-%ds", maxDisplayWidth-rColumnsDisplay[k].NameZhLen), rightColumn.Name)

		// not in left, in right
		if _, ok := left.Columns[rightColumn.Name]; !ok {
			results = append(results, d.onlyInRight(fmt.Sprintf("%s amount: [0 != %d]", title, rightColumn.Amount)))
		}
	}

	results = append(results, "")
	d.results = results
}

func (d *Differ) success(s string) string {
	return "✅  " + s
}

func (d *Differ) warning(s string) string {
	d.ok = false

	return "❎  " + s
}

func (d *Differ) failure(s string) string {
	d.ok = false

	return "❌  " + s
}

func (d *Differ) onlyInLeft(s string) string {
	d.ok = false

	return "❎❌" + s
}

func (d *Differ) onlyInRight(s string) string {
	d.ok = false

	return "❌❎" + s
}
