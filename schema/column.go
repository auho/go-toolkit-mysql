package schema

import (
	"strings"
	"unicode/utf8"

	"github.com/auho/go-simple-db/v2/schema"
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
	_len := len(c.Name)
	_runeLen := utf8.RuneCountInString(c.Name)
	_zhLen := (_len - _runeLen) / 2

	c.NameZhLen = _zhLen
	c.NameShowWidth = _runeLen + _zhLen
	c.FieldTypeLen = len(c.FieldType)
}

type Columns []Column

func NewColumnsFromSimpleDb(columns []schema.Column) Columns {
	var cs Columns
	for _, _c := range columns {
		ft := FieldType(strings.ToLower(_c.FieldType))
		nc := Column{
			Name:      _c.Name,
			FieldType: ft,
			DataType:  FileTypeToDataType(ft),
		}

		cs = append(cs, nc)
	}

	return cs
}

func (c Columns) ToColumnShow() []ColumnShow {
	var cs []ColumnShow
	for _, _c := range c {
		cs = append(cs, NewColumnShow(_c))
	}

	return cs
}
