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

type ColumnShow struct {
	Column

	NameZhLen     int
	NameShowWidth int
	FieldTypeLen  int
}

func NewColumnShow(c Column) ColumnShow {
	cs := ColumnShow{
		Column: c,
	}

	cs.zhShowWidth()

	return cs
}

func (c *ColumnShow) zhShowWidth() {
	byteLen := len(c.Name)
	runeLen := utf8.RuneCountInString(c.Name)
	zhLen := (byteLen - runeLen) / 2

	c.NameZhLen = zhLen
	c.NameShowWidth = runeLen + zhLen
	c.FieldTypeLen = len(c.FieldType)
}

type Columns []Column

func NewColumnsFromsimpledb(columns []schema.Column) Columns {
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

func (c Columns) ToColumnShow() []ColumnShow {
	var cs []ColumnShow
	for _, c := range c {
		cs = append(cs, NewColumnShow(c))
	}

	return cs
}
