package schema

import (
	"strings"
	"unicode/utf8"

	"github.com/auho/go-simple-db/v3/schema"
)

// Column holds the toolkit-level metadata for a single MySQL column.
type Column struct {
	Name      string
	FieldType FieldType
	DataType  DataType
}

// ColumnDisplay wraps a Column with pre-computed display metrics used to
// align column names in tabular output. The metrics account for the display
// width of multi-byte characters (e.g. CJK ideographs occupy two cells).
type ColumnDisplay struct {
	Column

	NameZhLen        int // number of double-width characters in Name
	NameDisplayWidth int // terminal display width of Name (rune count + extra width for CJK)
	FieldTypeLen     int // byte length of FieldType
}

// NewColumnDisplay creates a ColumnDisplay from a Column and computes its
// display metrics.
func NewColumnDisplay(c Column) ColumnDisplay {
	display := ColumnDisplay{
		Column: c,
	}

	display.calculateDisplayWidth()

	return display
}

// calculateDisplayWidth computes NameZhLen, NameDisplayWidth and FieldTypeLen.
// A character is treated as double-width when its UTF-8 byte length exceeds
// its rune length (i.e. multi-byte characters such as CJK ideographs).
func (c *ColumnDisplay) calculateDisplayWidth() {
	byteLen := len(c.Name)
	runeLen := utf8.RuneCountInString(c.Name)
	zhLen := (byteLen - runeLen) / 2

	c.NameZhLen = zhLen
	c.NameDisplayWidth = runeLen + zhLen
	c.FieldTypeLen = len(c.FieldType)
}

// Columns is an ordered collection of Column values.
type Columns []Column

// NewColumnsFromSimpleDB converts columns from the go-simple-db schema package
// into toolkit Columns, deriving the DataType for each field type.
func NewColumnsFromSimpleDB(columns []schema.Column) Columns {
	var cols Columns
	for _, c := range columns {
		ft := FieldType(strings.ToLower(c.FieldType))
		col := Column{
			Name:      c.Name,
			FieldType: ft,
			DataType:  FieldTypeToDataType(ft),
		}

		cols = append(cols, col)
	}

	return cols
}

// ToColumnDisplay converts all Columns to ColumnDisplay values with
// pre-computed display metrics.
func (c Columns) ToColumnDisplay() []ColumnDisplay {
	var displays []ColumnDisplay
	for _, c := range c {
		displays = append(displays, NewColumnDisplay(c))
	}

	return displays
}
