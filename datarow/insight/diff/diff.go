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
	ss []string
	ok bool
}

func (d *Differ) IsOk() bool {
	return d.ok
}

func (d *Differ) DifferenceToStrings() []string {
	return d.ss
}

func (d *Differ) diff(left, right *analysis.Analysis) {
	d.ok = true

	var ss []string

	lColumnsShow, lMaxColumnShow := left.GetColumnsShow()
	lMaxShow := lMaxColumnShow.NameShowWidth + 1
	rColumnsShow, rMaxColumnShow := right.GetColumnsShow()
	rMaxShow := rMaxColumnShow.NameShowWidth + 1

	maxShow := max(lMaxShow, rMaxShow)

	// table name and amount
	title := fmt.Sprintf("table[%s:%s] amount", left.Table.Table.Name, right.Table.Table.Name)
	if left.Table.Amount == right.Table.Amount {
		ss = append(ss, d.success(fmt.Sprintf("%s: %d", title, left.Table.Amount)))
	} else {
		ss = append(ss, d.failure(fmt.Sprintf("%s[%d != %d]", title, left.Table.Amount, right.Table.Amount)))
	}

	// loop left field
	for k, lfn := range left.FieldsName {
		lc := left.Columns[lfn]

		// left field title
		title = fmt.Sprintf(fmt.Sprintf("  %%-%ds", maxShow-lColumnsShow[k].NameZhLen), lc.Column.Name)

		if rc, ok := right.Columns[lc.Column.Name]; ok {
			// compare amount
			if lc.Amount == rc.Amount {
				ss = append(ss, d.success(fmt.Sprintf("%s amount: %d", title, lc.Amount)))
			} else {
				ss = append(ss, d.failure(fmt.Sprintf("%s amount: [%d != %d]", title, lc.Amount, rc.Amount)))
			}

			// compare distinct
			if lc.Distinct != rc.Distinct {
				ss = append(ss, d.failure(fmt.Sprintf("%s distinct: [%d != %d]", title, lc.Distinct, rc.Distinct)))
			}

			// compare empty
			if lc.Empty != rc.Empty {
				ss = append(ss, d.failure(fmt.Sprintf("%s empty: [%d != %d]", title, lc.Empty, rc.Empty)))
			}

			// compare null
			if lc.Null != rc.Null {
				ss = append(ss, d.failure(fmt.Sprintf("%s null: [%d != %d]", title, lc.Null, rc.Null)))
			}

		} else {
			// in left, not in right
			ss = append(ss, d.warningAndFailure(fmt.Sprintf("%s amount: [%d != 0]", title, lc.Amount)))
		}
	}

	// loop right field
	for k, rfn := range right.FieldsName {
		rc := right.Columns[rfn]

		// right field title
		title = fmt.Sprintf(fmt.Sprintf("  %%-%ds", maxShow-rColumnsShow[k].NameZhLen), rc.Column.Name)

		// not in left, in right
		if _, ok := left.Columns[rc.Column.Name]; !ok {
			ss = append(ss, d.failureAndWarning(fmt.Sprintf("%s amount: [0 != %d]", title, rc.Amount)))
		}
	}

	ss = append(ss, "")
	d.ss = ss
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

func (d *Differ) warningAndFailure(s string) string {
	d.ok = false

	return "❎❌" + s
}

func (d *Differ) failureAndWarning(s string) string {
	d.ok = false

	return "❌❎" + s
}
