// Package schema defines MySQL column and table metadata types used across
// the toolkit. It provides field-type classification (FieldType), a coarser
// data-type category (DataType), and display-width calculation helpers that
// account for multi-byte characters such as CJK ideographs.
package schema

// DataType is a coarse category that groups MySQL field types into a small
// set of kinds useful for statistical analysis (e.g. deciding whether "empty"
// means 0 for numbers or '' for strings).
type DataType byte

// FieldType is the raw MySQL column type as returned by the database
// (e.g. "int", "varchar", "datetime"). Values are always lower-case.
type FieldType string

// TypeFlag is reserved for future use.
type TypeFlag string

// DataType values group MySQL field types into analysis-friendly categories.
const (
	DataTypeUnknown DataType = iota // unrecognised or unsupported field type
	DataTypeBool                    // boolean
	DataTypeInt                     // integer types (bit, tinyint, int, bigint, ...)
	DataTypeUint                    // unsigned integer (currently unused)
	DataTypeFloat                   // floating-point and fixed-point types
	DataTypeString                  // character string types
	DataTypeTime                    // date/time types
	DataTypeBytes                   // binary string types
)

// FieldType constants mirror the MySQL column type names. Each constant is
// the lower-case form returned by INFORMATION_SCHEMA or SHOW COLUMNS.
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
