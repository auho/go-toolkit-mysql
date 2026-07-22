package insight

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	simpledb "github.com/auho/go-simple-db/v3"
	"github.com/auho/go-toolkit-mysql/datarow/insight/analysis"
	"github.com/auho/go-toolkit-mysql/schema"
)

func Explore(ctx context.Context, db *simpledb.SimpleDB, table string) (*analysis.Analysis, error) {
	return (&Insight{}).Explore(ctx, db, table)
}

type Insight struct{}

func (i *Insight) Explore(ctx context.Context, db *simpledb.SimpleDB, table string) (*analysis.Analysis, error) {
	return i.analyse(ctx, db, table)
}

func (i *Insight) analyse(ctx context.Context, db *simpledb.SimpleDB, table string) (*analysis.Analysis, error) {
	_cs, err := db.GetTableColumnsSchema(ctx, table)
	if err != nil {
		return nil, err
	}

	tableAly, err := i.analyseTable(db, table)
	if err != nil {
		return nil, err
	}

	columnsSchema := schema.NewColumnsFromsimpledb(_cs)
	columnsAly, err := i.analyseColumns(db, tableAly, columnsSchema)
	if err != nil {
		return nil, err
	}
	a := analysis.NewAnalysis()
	a.Table = tableAly

	for _, _ca := range columnsAly {
		a.Columns[_ca.Column.Name] = _ca
	}

	for _, _c := range _cs {
		a.FieldsName = append(a.FieldsName, _c.Name)
	}

	return a, nil
}

func (i *Insight) analyseTable(db *simpledb.SimpleDB, table string) (*analysis.Table, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	amount, err := db.RowCount(ctx, table)
	if err != nil {
		return nil, err
	}

	return &analysis.Table{
		Table:  schema.Table{Name: table},
		Amount: amount,
	}, err
}

func (i *Insight) analyseColumns(db *simpledb.SimpleDB, tableAly *analysis.Table, columns schema.Columns) ([]analysis.Column, error) {
	var fields []string

	fields = append(fields, "COUNT(*) AS 'amount'")
	for _, column := range columns {
		switch column.DataType {
		case schema.DataTypeInt, schema.DataTypeFloat:
			fields = append(fields,
				fmt.Sprintf("SUM(IF(`%s` = 0, 1, 0)) AS `%s_empty`", column.Name, column.Name),
				fmt.Sprintf("SUM(IF(`%s` IS NULL, 1, 0)) AS `%s_null`", column.Name, column.Name),
			)
		case schema.DataTypeString:
			fields = append(fields,
				fmt.Sprintf("SUM(IF(`%s` = '', 1, 0)) AS `%s_empty`", column.Name, column.Name),
				fmt.Sprintf("SUM(IF(`%s` IS NULL, 1, 0)) AS `%s_null`", column.Name, column.Name),
			)
		default:

		}

		fields = append(fields,
			fmt.Sprintf("COUNT(DISTINCT `%s`) as `%s_distinct`", column.Name, column.Name),
		)
	}

	var retRows []map[string]any
	sql := fmt.Sprintf("SELECT %s FROM `%s`", strings.Join(fields, ","), tableAly.Table.Name)
	err := db.GormDB().Raw(sql).Scan(&retRows).Error
	if err != nil {
		return nil, err
	}

	ret := retRows[0]

	var columnsAly []analysis.Column
	for _, column := range columns {
		_ca := analysis.Column{
			Column: column,
			Amount: tableAly.Amount,
		}
		err = i.toNum(ret, column.Name, "distinct", &_ca.Distinct)
		if err != nil {
			return nil, err
		}

		err = i.toNum(ret, column.Name, "empty", &_ca.Empty)
		if err != nil {
			return nil, err
		}

		err = i.toNum(ret, column.Name, "null", &_ca.Null)
		if err != nil {
			return nil, err
		}

		columnsAly = append(columnsAly, _ca)
	}

	return columnsAly, nil
}

func (i *Insight) toNum(ret map[string]any, name, suffix string, value *int) error {
	if v, ok := ret[name+"_"+suffix]; ok {
		var err error
		var _v int
		if v == nil {
			v = "0"
		}

		switch v.(type) {
		case string:
			_v, err = strconv.Atoi(v.(string))
			if err != nil {
				return err
			}
		case int64:
			_v = int(v.(int64))
		default:
			panic("unknown type")
		}

		*value = _v
	}

	return nil
}
