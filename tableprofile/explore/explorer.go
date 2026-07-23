// Package explore profiles a MySQL table by querying its schema and computing
// per-column statistics (row count, distinct values, empty values, and null
// values). The Explorer struct encapsulates the database connection and
// exposes an Analyze method that returns an analysis.Result.
package explore

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	simpledb "github.com/auho/go-simple-db/v3"
	"github.com/auho/go-toolkit-mysql/schema"
	"github.com/auho/go-toolkit-mysql/tableprofile/analysis"
)

// Explorer profiles MySQL tables using a SimpleDB connection.
type Explorer struct {
	db *simpledb.SimpleDB
}

// New creates an Explorer bound to the given database connection.
func New(db *simpledb.SimpleDB) *Explorer {
	return &Explorer{db: db}
}

// Analyze profiles a single table and returns the analysis Result containing
// row count and per-column statistics.
func (e *Explorer) Analyze(ctx context.Context, table string) (*analysis.Result, error) {
	return e.analyze(ctx, table)
}

// analyze is the internal implementation of Analyze. It fetches the table's
// column schema, computes table-level and column-level statistics, and
// assembles them into an analysis.Result.
func (e *Explorer) analyze(ctx context.Context, table string) (*analysis.Result, error) {
	dbColumns, err := e.db.GetTableColumnsSchema(ctx, table)
	if err != nil {
		return nil, err
	}

	tableAnalysis, err := e.analyzeTable(ctx, table)
	if err != nil {
		return nil, err
	}

	columnsSchema := schema.NewColumnsFromSimpleDB(dbColumns)
	columnsAnalysis, err := e.analyzeColumns(tableAnalysis, columnsSchema)
	if err != nil {
		return nil, err
	}
	r := analysis.NewResult()
	r.Table = tableAnalysis

	for _, columnAnalysis := range columnsAnalysis {
		r.Columns[columnAnalysis.Name] = columnAnalysis
	}

	for _, c := range dbColumns {
		r.FieldNames = append(r.FieldNames, c.Name)
	}

	return r, nil
}

// analyzeTable queries the row count of the given table with a 10-second
// timeout derived from the parent context.
func (e *Explorer) analyzeTable(ctx context.Context, table string) (*analysis.Table, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	rowCount, err := e.db.RowCount(ctx, table)
	if err != nil {
		return nil, err
	}

	return &analysis.Table{
		Table:    schema.Table{Name: table},
		RowCount: rowCount,
	}, nil
}

// analyzeColumns builds and executes a single aggregate SQL query that
// computes, for every column, the distinct count, empty count, and null count.
// The empty-count criterion depends on the column's DataType: 0 for numeric
// types, '' for string types.
func (e *Explorer) analyzeColumns(tableAnalysis *analysis.Table, columns schema.Columns) ([]analysis.Column, error) {
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
	err := e.db.GormDB().Raw(sql).Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	row := rows[0]

	var columnsAnalysis []analysis.Column
	for _, column := range columns {
		columnAnalysis := analysis.Column{
			Column:   column,
			RowCount: tableAnalysis.RowCount,
		}
		err = e.toInt(row, column.Name, "distinct", &columnAnalysis.Distinct)
		if err != nil {
			return nil, err
		}

		err = e.toInt(row, column.Name, "empty", &columnAnalysis.Empty)
		if err != nil {
			return nil, err
		}

		err = e.toInt(row, column.Name, "null", &columnAnalysis.Null)
		if err != nil {
			return nil, err
		}

		columnsAnalysis = append(columnsAnalysis, columnAnalysis)
	}

	return columnsAnalysis, nil
}

// toInt extracts a value keyed by "name_suffix" from a query result row and
// converts it to int. The value may be nil (treated as 0), string, int64,
// []byte, or float64. It is a no-op when the key is absent.
func (e *Explorer) toInt(row map[string]any, name, suffix string, value *int) error {
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
