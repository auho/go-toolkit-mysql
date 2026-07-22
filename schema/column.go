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
	cs := ColumnDisplay{
		Column: c,
	}

	cs.calculateDisplayWidth()

	return cs
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
	var cs Columns
	for _, c := range columns {
		ft := FieldType(strings.ToLower(c.FieldType))
		nc := Column{
			Name:      c.Name,
			FieldType: ft,
			DataType:  FieldTypeToDataType(ft),
		}

		cs = append(cs, nc)
	}

	return cs
}

func (c Columns) ToColumnDisplay() []ColumnDisplay {
	var cs []ColumnDisplay
	for _, c := range c {
		cs = append(cs, NewColumnDisplay(c))
	}

	return cs
}
