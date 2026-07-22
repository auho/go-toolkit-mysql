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
	return i.analyze(ctx, db, table)
}

func (i *Insight) analyze(ctx context.Context, db *simpledb.SimpleDB, table string) (*analysis.Analysis, error) {
	dbColumns, err := db.GetTableColumnsSchema(ctx, table)
	if err != nil {
		return nil, err
	}

	tableAnalysis, err := i.analyzeTable(db, table)
	if err != nil {
		return nil, err
	}

	columnsSchema := schema.NewColumnsFromSimpleDB(dbColumns)
	columnsAnalysis, err := i.analyzeColumns(db, tableAnalysis, columnsSchema)
	if err != nil {
		return nil, err
	}
	a := analysis.NewAnalysis()
	a.Table = tableAnalysis

	for _, columnAnalysis := range columnsAnalysis {
		a.Columns[columnAnalysis.Name] = columnAnalysis
	}

	for _, c := range dbColumns {
		a.FieldNames = append(a.FieldNames, c.Name)
	}

	return a, nil
}

func (i *Insight) analyzeTable(db *simpledb.SimpleDB, table string) (*analysis.Table, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	amount, err := db.RowCount(ctx, table)
	if err != nil {
		return nil, err
	}

	return &analysis.Table{
		Table:  schema.Table{Name: table},
		Amount: amount,
	}, nil
}

func (i *Insight) analyzeColumns(db *simpledb.SimpleDB, tableAnalysis *analysis.Table, columns schema.Columns) ([]analysis.Column, error) {
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

	var rows []map[string]any
	sql := fmt.Sprintf("SELECT %s FROM `%s`", strings.Join(fields, ","), tableAnalysis.Name)
	err := db.GormDB().Raw(sql).Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	row := rows[0]

	var columnsAnalysis []analysis.Column
	for _, column := range columns {
		columnAnalysis := analysis.Column{
			Column: column,
			Amount: tableAnalysis.Amount,
		}
		err = i.toInt(row, column.Name, "distinct", &columnAnalysis.Distinct)
		if err != nil {
			return nil, err
		}

		err = i.toInt(row, column.Name, "empty", &columnAnalysis.Empty)
		if err != nil {
			return nil, err
		}

		err = i.toInt(row, column.Name, "null", &columnAnalysis.Null)
		if err != nil {
			return nil, err
		}

		columnsAnalysis = append(columnsAnalysis, columnAnalysis)
	}

	return columnsAnalysis, nil
}

func (i *Insight) toInt(row map[string]any, name, suffix string, value *int) error {
	if v, ok := row[name+"_"+suffix]; ok {
		var err error
		var n int
		if v == nil {
			v = "0"
		}

		switch v.(type) {
		case string:
			n, err = strconv.Atoi(v.(string))
			if err != nil {
				return err
			}
		case int64:
			n = int(v.(int64))
		case []byte:
			n, err = strconv.Atoi(string(v.([]byte)))
			if err != nil {
				return err
			}
		case float64:
			n = int(v.(float64))
		default:
			return fmt.Errorf("toInt: unknown type %T", v)
		}

		*value = n
	}

	return nil
}
