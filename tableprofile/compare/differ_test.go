package compare

import (
	"strings"
	"testing"

	"github.com/auho/go-toolkit-mysql/schema"
	"github.com/auho/go-toolkit-mysql/tableprofile/analysis"
)

// helper to build a Result for testing
func buildResult(tableName string, rowCount int, cols map[string]analysis.Column, fieldNames []string) *analysis.Result {
	r := analysis.NewResult()
	r.Table = &analysis.Table{
		Table:    schema.Table{Name: tableName},
		RowCount: rowCount,
	}
	r.FieldNames = fieldNames
	r.Columns = cols
	return r
}

func col(name string, ft schema.FieldType, dt schema.DataType, rowCount, distinct, empty, null int) analysis.Column {
	return analysis.Column{
		Column:   schema.Column{Name: name, FieldType: ft, DataType: dt},
		RowCount: rowCount,
		Distinct: distinct,
		Empty:    empty,
		Null:     null,
	}
}

func TestDiff_Identical(t *testing.T) {
	cols := map[string]analysis.Column{
		"id":   col("id", schema.FieldTypeInt, schema.DataTypeInt, 100, 100, 0, 0),
		"name": col("name", schema.FieldTypeVarchar, schema.DataTypeString, 100, 80, 10, 5),
	}
	left := buildResult("t1", 100, cols, []string{"id", "name"})
	right := buildResult("t2", 100, cols, []string{"id", "name"})

	d := Diff(left, right)

	if !d.IsOK() {
		t.Error("IsOK = false, want true for identical results")
	}
}

func TestDiff_RowCountMismatch(t *testing.T) {
	left := buildResult("t1", 100, nil, nil)
	right := buildResult("t2", 200, nil, nil)

	d := Diff(left, right)

	if d.IsOK() {
		t.Error("IsOK = true, want false for different row counts")
	}

	found := false
	for _, line := range d.Differences() {
		if strings.Contains(line, "100 != 200") {
			found = true
		}
	}
	if !found {
		t.Error("expected a line showing 100 != 200")
	}
}

func TestDiff_ColumnStatMismatch(t *testing.T) {
	leftCols := map[string]analysis.Column{
		"name": col("name", schema.FieldTypeVarchar, schema.DataTypeString, 100, 80, 10, 5),
	}
	rightCols := map[string]analysis.Column{
		"name": col("name", schema.FieldTypeVarchar, schema.DataTypeString, 100, 70, 20, 10),
	}
	left := buildResult("t1", 100, leftCols, []string{"name"})
	right := buildResult("t2", 100, rightCols, []string{"name"})

	d := Diff(left, right)

	if d.IsOK() {
		t.Error("IsOK = true, want false for different column stats")
	}

	diffText := strings.Join(d.Differences(), "\n")

	// distinct mismatch
	if !strings.Contains(diffText, "distinct") || !strings.Contains(diffText, "80 != 70") {
		t.Errorf("expected distinct mismatch 80 != 70 in:\n%s", diffText)
	}
	// empty mismatch
	if !strings.Contains(diffText, "empty") || !strings.Contains(diffText, "10 != 20") {
		t.Errorf("expected empty mismatch 10 != 20 in:\n%s", diffText)
	}
	// null mismatch
	if !strings.Contains(diffText, "null") || !strings.Contains(diffText, "5 != 10") {
		t.Errorf("expected null mismatch 5 != 10 in:\n%s", diffText)
	}
}

func TestDiff_ColumnOnlyInLeft(t *testing.T) {
	leftCols := map[string]analysis.Column{
		"id":   col("id", schema.FieldTypeInt, schema.DataTypeInt, 100, 100, 0, 0),
		"extra": col("extra", schema.FieldTypeVarchar, schema.DataTypeString, 100, 50, 10, 0),
	}
	rightCols := map[string]analysis.Column{
		"id": col("id", schema.FieldTypeInt, schema.DataTypeInt, 100, 100, 0, 0),
	}
	left := buildResult("t1", 100, leftCols, []string{"id", "extra"})
	right := buildResult("t2", 100, rightCols, []string{"id"})

	d := Diff(left, right)

	if d.IsOK() {
		t.Error("IsOK = true, want false when column exists only in left")
	}

	found := false
	for _, line := range d.Differences() {
		if strings.Contains(line, "extra") && strings.Contains(line, "❎❌") {
			found = true
		}
	}
	if !found {
		t.Error("expected a line for column only in left with ❎❌ prefix")
	}
}

func TestDiff_ColumnOnlyInRight(t *testing.T) {
	leftCols := map[string]analysis.Column{
		"id": col("id", schema.FieldTypeInt, schema.DataTypeInt, 100, 100, 0, 0),
	}
	rightCols := map[string]analysis.Column{
		"id":    col("id", schema.FieldTypeInt, schema.DataTypeInt, 100, 100, 0, 0),
		"extra": col("extra", schema.FieldTypeVarchar, schema.DataTypeString, 100, 50, 10, 0),
	}
	left := buildResult("t1", 100, leftCols, []string{"id"})
	right := buildResult("t2", 100, rightCols, []string{"id", "extra"})

	d := Diff(left, right)

	if d.IsOK() {
		t.Error("IsOK = true, want false when column exists only in right")
	}

	found := false
	for _, line := range d.Differences() {
		if strings.Contains(line, "extra") && strings.Contains(line, "❌❎") {
			found = true
		}
	}
	if !found {
		t.Error("expected a line for column only in right with ❌❎ prefix")
	}
}

func TestDiff_Differences_EndsWithEmptyLine(t *testing.T) {
	left := buildResult("t1", 1, nil, nil)
	right := buildResult("t2", 1, nil, nil)

	d := Diff(left, right)

	diffs := d.Differences()
	if diffs[len(diffs)-1] != "" {
		t.Errorf("last diff line = %q, want empty", diffs[len(diffs)-1])
	}
}
