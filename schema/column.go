package schema

import (
	"strings"
	"unicode/utf8"

	"github.com/auho/go-simple-db/v3/schema"
)

type Column struct {
	Name      string
	FieldType FieldType
	DataType  DataType
}

type ColumnDisplay struct {
	Column

	NameZhLen        int
	NameDisplayWidth int
	FieldTypeLen     int
}

func NewColumnDisplay(c Column) ColumnDisplay {
	display := ColumnDisplay{
		Column: c,
	}

	display.calculateDisplayWidth()

	return display
}

func (c *ColumnDisplay) calculateDisplayWidth() {
	byteLen := len(c.Name)
	runeLen := utf8.RuneCountInString(c.Name)
	zhLen := (byteLen - runeLen) / 2

	c.NameZhLen = zhLen
	c.NameDisplayWidth = runeLen + zhLen
	c.FieldTypeLen = len(c.FieldType)
}

type Columns []Column

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

func (c Columns) ToColumnDisplay() []ColumnDisplay {
	var displays []ColumnDisplay
	for _, c := range c {
		displays = append(displays, NewColumnDisplay(c))
	}

	return displays
}
