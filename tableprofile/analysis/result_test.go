package analysis

import (
	"strings"
	"testing"

	"github.com/auho/go-toolkit-mysql/schema"
)

// helper to build a Result for testing
func buildTestResult() *Result {
	r := NewResult()
	r.Table = &Table{
		Table:    schema.Table{Name: "users"},
		RowCount: 100,
	}
	r.FieldNames = []string{"id", "name", "中文字段"}
	r.Columns = map[string]Column{
		"id": {
			Column:   schema.Column{Name: "id", FieldType: schema.FieldTypeInt, DataType: schema.DataTypeInt},
			RowCount: 100,
			Distinct: 100,
			Empty:    0,
			Null:     0,
		},
		"name": {
			Column:   schema.Column{Name: "name", FieldType: schema.FieldTypeVarchar, DataType: schema.DataTypeString},
			RowCount: 100,
			Distinct: 80,
			Empty:    10,
			Null:     5,
		},
		"中文字段": {
			Column:   schema.Column{Name: "中文字段", FieldType: schema.FieldTypeVarchar, DataType: schema.DataTypeString},
			RowCount: 100,
			Distinct: 0,
			Empty:    100,
			Null:     0,
		},
	}

	return r
}

func TestNewResult(t *testing.T) {
	r := NewResult()

	if r == nil {
		t.Fatal("NewResult returned nil")
	}
	if r.Columns == nil {
		t.Error("Columns map is nil")
	}
	if len(r.Columns) != 0 {
		t.Errorf("Columns map len = %d, want 0", len(r.Columns))
	}
	if r.Table != nil {
		t.Error("Table should be nil")
	}
	if len(r.FieldNames) != 0 {
		t.Errorf("FieldNames len = %d, want 0", len(r.FieldNames))
	}
}

func TestResult_ToStrings(t *testing.T) {
	r := buildTestResult()
	lines := r.ToStrings()

	// last line should be empty
	if lines[len(lines)-1] != "" {
		t.Errorf("last line = %q, want empty", lines[len(lines)-1])
	}

	// first line should contain table name and row count
	if !strings.Contains(lines[0], "users") {
		t.Errorf("first line %q should contain table name \"users\"", lines[0])
	}
	if !strings.Contains(lines[0], "100") {
		t.Errorf("first line %q should contain row count 100", lines[0])
	}

	// second line should be the header
	if !strings.Contains(lines[1], "COLUMN") {
		t.Errorf("header line %q should contain \"COLUMN\"", lines[1])
	}
	if !strings.Contains(lines[1], "TYPE") {
		t.Errorf("header line %q should contain \"TYPE\"", lines[1])
	}

	// there should be a line per column (3 columns)
	// total lines = 1 (table) + 1 (header) + 3 (columns) + 1 (empty) = 6
	if len(lines) != 6 {
		t.Errorf("lines count = %d, want 6", len(lines))
	}

	// the "中文字段" column has empty+null >= rowCount, so should show warning
	foundWarning := false
	for _, line := range lines {
		if strings.Contains(line, "中文字段") && strings.Contains(line, "❗️") {
			foundWarning = true
		}
	}
	if !foundWarning {
		t.Error("expected warning emoji for column where empty+null >= rowCount")
	}
}

func TestResult_ToStrings_EmptyColumns(t *testing.T) {
	r := NewResult()
	r.Table = &Table{
		Table:    schema.Table{Name: "empty"},
		RowCount: 0,
	}

	lines := r.ToStrings()

	// should have table line, header line, and trailing empty line
	if len(lines) != 3 {
		t.Errorf("lines count = %d, want 3", len(lines))
	}
}

func TestResult_ColumnsDisplay(t *testing.T) {
	r := buildTestResult()
	displays, maxDisplay := r.ColumnsDisplay()

	if len(displays) != 3 {
		t.Fatalf("displays len = %d, want 3", len(displays))
	}

	// "中文字段" has 4 CJK chars, display width = 8; others are shorter
	if maxDisplay.NameDisplayWidth != 8 {
		t.Errorf("maxDisplay.NameDisplayWidth = %d, want 8", maxDisplay.NameDisplayWidth)
	}

	// display order should match FieldNames
	if displays[0].Name != "id" {
		t.Errorf("displays[0].Name = %q, want \"id\"", displays[0].Name)
	}
	if displays[1].Name != "name" {
		t.Errorf("displays[1].Name = %q, want \"name\"", displays[1].Name)
	}
	if displays[2].Name != "中文字段" {
		t.Errorf("displays[2].Name = %q, want \"中文字段\"", displays[2].Name)
	}
}

func TestResult_ColumnsDisplay_Empty(t *testing.T) {
	r := NewResult()
	r.Table = &Table{
		Table:    schema.Table{Name: "empty"},
		RowCount: 0,
	}

	displays, maxDisplay := r.ColumnsDisplay()

	if len(displays) != 0 {
		t.Errorf("displays len = %d, want 0", len(displays))
	}
	// maxDisplay should be zero value
	if maxDisplay.NameDisplayWidth != 0 {
		t.Errorf("maxDisplay.NameDisplayWidth = %d, want 0", maxDisplay.NameDisplayWidth)
	}
}
