package schema

type DataType byte
type FieldType string
type TypeFlag string

const (
	DataTypeUnknown DataType = iota
	DataTypeBool
	DataTypeInt
	DataTypeUint
	DataTypeFloat
	DataTypeString
	DataTypeTime
	DataTypeBytes
)

const (
	FieldTypeBit       FieldType = "bit"
	FieldTypeTinyint   FieldType = "tinyint"
	FieldTypeBool      FieldType = "bool"
	FieldTypeBoolean   FieldType = "boolean"
	FieldTypeSmallint  FieldType = "smallint"
	FieldTypeMediumint FieldType = "mediumint"
	FieldTypeInt       FieldType = "int"
	FieldTypeInteger   FieldType = "integer"
	FieldTypeBigint    FieldType = "bigint"
	FieldTypeDecimal   FieldType = "decimal"
	FieldTypeFloat     FieldType = "float"
	FieldTypeDouble    FieldType = "double"
	FieldTypeDate      FieldType = "date"      // '0000-00-00'
	FieldTypeTime      FieldType = "time"      // '00:00:00'
	FieldTypeDatetime  FieldType = "datetime"  // '0000-00-00 00:00:00'
	FieldTypeTimestamp FieldType = "timestamp" // '0000-00-00 00:00:00'
	FieldTypeYear      FieldType = "year"      // 0000
	FieldTypeChar      FieldType = "char"
	FieldTypeVarchar   FieldType = "varchar"
	FieldTypeBinary    FieldType = "binary"
	FieldTypeVarbinary FieldType = "varbinary"
	FieldTypeBlob      FieldType = "blob"
	FieldTypeText      FieldType = "text"
	FieldTypeEnum      FieldType = "enum"
	FieldTypeSet       FieldType = "set"
)
