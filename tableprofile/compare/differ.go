// Package compare diffs two analysis Results and reports the differences
// as human-readable lines.
package compare

import (
	"fmt"

	"github.com/auho/go-toolkit-mysql/tableprofile/analysis"
)

// Diff compares two analysis Results and returns a Differ containing the
// differences. The Differ reports mismatches in row count, per-column
// distinct/empty/null counts, and columns that exist only on one side.
func Diff(left, right *analysis.Result) *Differ {
	d := &Differ{}
	d.run(left, right)

	return d
}

// Differ holds the outcome of comparing two analysis Results. Call IsOK to
// check whether the two results are identical, and Differences to retrieve
// the human-readable diff lines.
type Differ struct {
	results []string
	ok      bool
}

// IsOK returns true when the two compared results have no differences.
func (d *Differ) IsOK() bool {
	return d.ok
}

// Differences returns the human-readable diff lines. Each line is prefixed
// with an emoji indicating the outcome (success, failure, or column only on
// one side).
func (d *Differ) Differences() []string {
	return d.results
}

// run performs the comparison of left and right analysis Results, populating
// d.results and d.ok.
func (d *Differ) run(left, right *analysis.Result) {
	d.ok = true

	var results []string

	lColumnsDisplay, lMaxColumnDisplay := left.ColumnsDisplay()
	lMaxDisplayWidth := lMaxColumnDisplay.NameDisplayWidth + 1
	rColumnsDisplay, rMaxColumnDisplay := right.ColumnsDisplay()
	rMaxDisplayWidth := rMaxColumnDisplay.NameDisplayWidth + 1

	maxDisplayWidth := max(lMaxDisplayWidth, rMaxDisplayWidth)

	// table name and row count
	title := fmt.Sprintf("table[%s:%s] amount", left.Table.Name, right.Table.Name)
	if left.Table.RowCount == right.Table.RowCount {
		results = append(results, d.success(fmt.Sprintf("%s: %d", title, left.Table.RowCount)))
	} else {
		results = append(results, d.failure(fmt.Sprintf("%s[%d != %d]", title, left.Table.RowCount, right.Table.RowCount)))
	}

	// compare columns present in the left result
	for k, leftFieldName := range left.FieldNames {
		leftColumn := left.Columns[leftFieldName]

		title = fmt.Sprintf(fmt.Sprintf("  %%-%ds", maxDisplayWidth-lColumnsDisplay[k].NameZhLen), leftColumn.Name)

		if rightColumn, ok := right.Columns[leftColumn.Name]; ok {
			// compare row count
			if leftColumn.RowCount == rightColumn.RowCount {
				results = append(results, d.success(fmt.Sprintf("%s amount: %d", title, leftColumn.RowCount)))
			} else {
				results = append(results, d.failure(fmt.Sprintf("%s amount: [%d != %d]", title, leftColumn.RowCount, rightColumn.RowCount)))
			}

			// compare distinct count
			if leftColumn.Distinct != rightColumn.Distinct {
				results = append(results, d.failure(fmt.Sprintf("%s distinct: [%d != %d]", title, leftColumn.Distinct, rightColumn.Distinct)))
			}

			// compare empty count
			if leftColumn.Empty != rightColumn.Empty {
				results = append(results, d.failure(fmt.Sprintf("%s empty: [%d != %d]", title, leftColumn.Empty, rightColumn.Empty)))
			}

			// compare null count
			if leftColumn.Null != rightColumn.Null {
				results = append(results, d.failure(fmt.Sprintf("%s null: [%d != %d]", title, leftColumn.Null, rightColumn.Null)))
			}

		} else {
			// column exists only in the left result
			results = append(results, d.onlyInLeft(fmt.Sprintf("%s amount: [%d != 0]", title, leftColumn.RowCount)))
		}
	}

	// detect columns present only in the right result
	for k, rightFieldName := range right.FieldNames {
		rightColumn := right.Columns[rightFieldName]

		title = fmt.Sprintf(fmt.Sprintf("  %%-%ds", maxDisplayWidth-rColumnsDisplay[k].NameZhLen), rightColumn.Name)

		if _, ok := left.Columns[rightColumn.Name]; !ok {
			results = append(results, d.onlyInRight(fmt.Sprintf("%s amount: [0 != %d]", title, rightColumn.RowCount)))
		}
	}

	results = append(results, "")
	d.results = results
}

// success formats a line for a matching value.
func (d *Differ) success(s string) string {
	return "✅  " + s
}

// warning formats a line for a warning-level mismatch.
func (d *Differ) warning(s string) string {
	d.ok = false

	return "❎  " + s
}

// failure formats a line for a mismatch.
func (d *Differ) failure(s string) string {
	d.ok = false

	return "❌  " + s
}

// onlyInLeft formats a line for a column that exists only in the left result.
func (d *Differ) onlyInLeft(s string) string {
	d.ok = false

	return "❎❌" + s
}

// onlyInRight formats a line for a column that exists only in the right result.
func (d *Differ) onlyInRight(s string) string {
	d.ok = false

	return "❌❎" + s
}
