package schema

import (
	"testing"

 simplesdb "github.com/auho/go-simple-db/v3/schema"
)

func TestFieldTypeToDataType(t *testing.T) {
	tests := []struct {
		name string
		ft   FieldType
		want DataType
	}{
		// integer types
		{"bit", FieldTypeBit, DataTypeInt},
		{"tinyint", FieldTypeTinyint, DataTypeInt},
		{"smallint", FieldTypeSmallint, DataTypeInt},
		{"mediumint", FieldTypeMediumint, DataTypeInt},
		{"int", FieldTypeInt, DataTypeInt},
		{"integer", FieldTypeInteger, DataTypeInt},
		{"bigint", FieldTypeBigint, DataTypeInt},

		// float types
		{"decimal", FieldTypeDecimal, DataTypeFloat},
		{"float", FieldTypeFloat, DataTypeFloat},
		{"double", FieldTypeDouble, DataTypeFloat},

		// bool types
		{"bool", FieldTypeBool, DataTypeBool},
		{"boolean", FieldTypeBoolean, DataTypeBool},

		// time types
		{"date", FieldTypeDate, DataTypeTime},
		{"time", FieldTypeTime, DataTypeTime},
		{"datetime", FieldTypeDatetime, DataTypeTime},
		{"timestamp", FieldTypeTimestamp, DataTypeTime},
		{"year", FieldTypeYear, DataTypeTime},

		// string types
		{"char", FieldTypeChar, DataTypeString},
		{"varchar", FieldTypeVarchar, DataTypeString},
		{"text", FieldTypeText, DataTypeString},
		{"enum", FieldTypeEnum, DataTypeString},
		{"set", FieldTypeSet, DataTypeString},

		// bytes types
		{"binary", FieldTypeBinary, DataTypeBytes},
		{"varbinary", FieldTypeVarbinary, DataTypeBytes},
		{"blob", FieldTypeBlob, DataTypeBytes},

		// unknown
		{"unknown", FieldType("json"), DataTypeUnknown},
		{"empty", FieldType(""), DataTypeUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FieldTypeToDataType(tt.ft)
			if got != tt.want {
				t.Errorf("FieldTypeToDataType(%q) = %v, want %v", tt.ft, got, tt.want)
			}
		})
	}
}

func TestNewColumnDisplay_ASCII(t *testing.T) {
	c := Column{Name: "id", FieldType: FieldTypeInt, DataType: DataTypeInt}
	d := NewColumnDisplay(c)

	if d.NameZhLen != 0 {
		t.Errorf("NameZhLen = %d, want 0 for ASCII name", d.NameZhLen)
	}
	if d.NameDisplayWidth != 2 {
		t.Errorf("NameDisplayWidth = %d, want 2 for \"id\"", d.NameDisplayWidth)
	}
	if d.FieldTypeLen != len(FieldTypeInt) {
		t.Errorf("FieldTypeLen = %d, want %d", d.FieldTypeLen, len(FieldTypeInt))
	}
}

func TestNewColumnDisplay_CJK(t *testing.T) {
	// "中文字段" = 4 CJK characters, each 3 bytes in UTF-8
	// byteLen = 12, runeLen = 4, zhLen = (12-4)/2 = 4
	// displayWidth = 4 + 4 = 8 (each CJK char takes 2 display cells)
	c := Column{Name: "中文字段", FieldType: FieldTypeVarchar, DataType: DataTypeString}
	d := NewColumnDisplay(c)

	if d.NameZhLen != 4 {
		t.Errorf("NameZhLen = %d, want 4 for 4 CJK characters", d.NameZhLen)
	}
	if d.NameDisplayWidth != 8 {
		t.Errorf("NameDisplayWidth = %d, want 8", d.NameDisplayWidth)
	}
}

func TestNewColumnDisplay_Mixed(t *testing.T) {
	// "d1" = 2 ASCII chars, "中文" = 2 CJK chars
	// byteLen = 2 + 6 = 8, runeLen = 4, zhLen = (8-4)/2 = 2
	// displayWidth = 4 + 2 = 6
	c := Column{Name: "d1中文", FieldType: FieldTypeVarchar}
	d := NewColumnDisplay(c)

	if d.NameZhLen != 2 {
		t.Errorf("NameZhLen = %d, want 2", d.NameZhLen)
	}
	if d.NameDisplayWidth != 6 {
		t.Errorf("NameDisplayWidth = %d, want 6", d.NameDisplayWidth)
	}
}

func TestNewColumnsFromSimpleDB(t *testing.T) {
	input := []simplesdb.Column{
		{Name: "id", FieldType: "INT"},
		{Name: "name", FieldType: "VARCHAR"},
		{Name: "中文字段", FieldType: "Text"},
	}

	cols := NewColumnsFromSimpleDB(input)

	if len(cols) != 3 {
		t.Fatalf("len = %d, want 3", len(cols))
	}

	// field types should be lower-cased
	if cols[0].FieldType != FieldTypeInt {
		t.Errorf("cols[0].FieldType = %q, want %q", cols[0].FieldType, FieldTypeInt)
	}
	if cols[1].FieldType != FieldTypeVarchar {
		t.Errorf("cols[1].FieldType = %q, want %q", cols[1].FieldType, FieldTypeVarchar)
	}
	if cols[2].FieldType != FieldTypeText {
		t.Errorf("cols[2].FieldType = %q, want %q", cols[2].FieldType, FieldTypeText)
	}

	// data types should be derived
	if cols[0].DataType != DataTypeInt {
		t.Errorf("cols[0].DataType = %v, want %v", cols[0].DataType, DataTypeInt)
	}
	if cols[1].DataType != DataTypeString {
		t.Errorf("cols[1].DataType = %v, want %v", cols[1].DataType, DataTypeString)
	}
	if cols[2].DataType != DataTypeString {
		t.Errorf("cols[2].DataType = %v, want %v", cols[2].DataType, DataTypeString)
	}
}

func TestNewColumnsFromSimpleDB_Empty(t *testing.T) {
	cols := NewColumnsFromSimpleDB(nil)
	if len(cols) != 0 {
		t.Errorf("len = %d, want 0", len(cols))
	}
}

func TestColumns_ToColumnDisplay(t *testing.T) {
	cols := Columns{
		{Name: "id", FieldType: FieldTypeInt, DataType: DataTypeInt},
		{Name: "name", FieldType: FieldTypeVarchar, DataType: DataTypeString},
	}

	displays := cols.ToColumnDisplay()

	if len(displays) != 2 {
		t.Fatalf("len = %d, want 2", len(displays))
	}

	if displays[0].Name != "id" || displays[1].Name != "name" {
		t.Errorf("names = %q, %q; want \"id\", \"name\"", displays[0].Name, displays[1].Name)
	}

	// display widths should be computed
	if displays[0].NameDisplayWidth != 2 {
		t.Errorf("displays[0].NameDisplayWidth = %d, want 2", displays[0].NameDisplayWidth)
	}
	if displays[1].NameDisplayWidth != 4 {
		t.Errorf("displays[1].NameDisplayWidth = %d, want 4", displays[1].NameDisplayWidth)
	}
}
