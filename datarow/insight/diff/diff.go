package diff

import (
	"fmt"

	"github.com/auho/go-toolkit/v2/mysql/datarow/insight/analysis"
)

func Diff(as ...*analysis.Analysis) *Differ {
	d := &Differ{}
	d.diff(as...)

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

func (d *Differ) diff(as ...*analysis.Analysis) {
	d.ok = true

	var ss []string

	_las := as[0]
	_ras := as[1]

	_lColumnsShow, _lMaxColumnShow := _las.GetColumnsShow()
	_lMaxShow := _lMaxColumnShow.NameShowWidth + 1
	_rColumnsShow, _rMaxColumnShow := _ras.GetColumnsShow()
	_rMaxShow := _rMaxColumnShow.NameShowWidth + 1

	_maxShow := _lMaxShow
	if _maxShow < _rMaxShow {
		_maxShow = _rMaxShow
	}

	// table name and amount
	_title := fmt.Sprintf("table[%s:%s] amount", _las.Table.Table.Name, _ras.Table.Table.Name)
	if _las.Table.Amount == _ras.Table.Amount {
		ss = append(ss, d.success(fmt.Sprintf("%s: %d", _title, _las.Table.Amount)))
	} else {
		ss = append(ss, d.failure(fmt.Sprintf("%s[%d != %d]", _title, _las.Table.Amount, _ras.Table.Amount)))
	}

	// loop left field
	for _k, _lfn := range _las.FieldsName {
		_lc := _las.Columns[_lfn]

		// left field title
		_title = fmt.Sprintf(fmt.Sprintf("  %%-%ds", _maxShow-_lColumnsShow[_k].NameZhLen), _lc.Column.Name)

		if _rc, ok := _ras.Columns[_lc.Column.Name]; ok {
			// compare amount
			if _lc.Amount == _rc.Amount {
				ss = append(ss, d.success(fmt.Sprintf("%s amount: %d", _title, _lc.Amount)))
			} else {
				ss = append(ss, d.failure(fmt.Sprintf("%s amount: [%d != %d]", _title, _lc.Amount, _rc.Amount)))
			}

			// compare distinct
			if _lc.Distinct != _rc.Distinct {
				ss = append(ss, d.failure(fmt.Sprintf("%s distinct: [%d != %d]", _title, _lc.Distinct, _rc.Distinct)))
			}

			// compare empty
			if _lc.Empty != _rc.Empty {
				ss = append(ss, d.failure(fmt.Sprintf("%s empty: [%d != %d]", _title, _lc.Empty, _rc.Empty)))
			}

			// compare null
			if _lc.Null != _rc.Null {
				ss = append(ss, d.failure(fmt.Sprintf("%s null: [%d != %d]", _title, _lc.Null, _rc.Null)))
			}

		} else {
			// in left, not in right
			ss = append(ss, d.warningAndFailure(fmt.Sprintf("%s amount: [%d != 0]", _title, _lc.Amount)))
		}
	}

	// loop right field
	for _k, _rfn := range _ras.FieldsName {
		_rc := _ras.Columns[_rfn]

		// right field title
		_title = fmt.Sprintf(fmt.Sprintf("  %%-%ds", _maxShow-_rColumnsShow[_k].NameZhLen), _rc.Column.Name)

		// not in left, in right
		if _, ok := _las.Columns[_rc.Column.Name]; !ok {
			ss = append(ss, d.failureAndWarning(fmt.Sprintf("%s amount: [0 != %d]", _title, _rc.Amount)))
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
