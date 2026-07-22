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
	FieldTypeBit       FieldType = "bit"       //
	FieldTypeTinyint             = "tinyint"   //
	FieldTypeBool                = "bool"      //
	FieldTypeBoolean             = "boolean"   //
	FieldTypeSmallint            = "smallint"  //
	FieldTypeMediumint           = "mediumint" //
	FieldTypeInt                 = "int"       //
	FieldTypeInteger             = "integer"   //
	FieldTypeBigint              = "bigint"    //
	FieldTypeDecimal             = "decimal"   //
	FieldTypeFloat               = "float"     //
	FieldTypeDouble              = "double"    //
	FieldTypeDate                = "date"      // '0000-00-00'
	FieldTypeTime                = "time"      // '00:00:00'
	FieldTypeDatetime            = "datetime"  // '0000-00-00 00:00:00'
	FieldTypeTimestamp           = "timestamp" // '0000-00-00 00:00:00'
	FieldTypeYear                = "year"      // 0000
	FieldTypeChar                = "char"      //
	FieldTypeVarchar             = "varchar"   //
	FieldTypeBinary              = "binary"    //
	FieldTypeVarbinary           = "varbinary" //
	FieldTypeBlob                = "blob"      //
	FieldTypeText                = "text"      //
	FieldTypeEnum                = "enum"      //
	FieldTypeSet                 = "set"       //
)
