package schema

// FieldTypeToDataType maps a MySQL field type to its coarse DataType category.
// Unrecognised types return DataTypeUnknown.
func FieldTypeToDataType(ft FieldType) DataType {
	switch ft {
	case FieldTypeTinyint, FieldTypeSmallint, FieldTypeMediumint, FieldTypeInt, FieldTypeInteger, FieldTypeBigint, FieldTypeBit:
		return DataTypeInt
	case FieldTypeDecimal, FieldTypeFloat, FieldTypeDouble:
		return DataTypeFloat
	case FieldTypeBool, FieldTypeBoolean:
		return DataTypeBool
	case FieldTypeDate, FieldTypeTime, FieldTypeDatetime, FieldTypeTimestamp, FieldTypeYear:
		return DataTypeTime
	case FieldTypeChar, FieldTypeVarchar, FieldTypeText, FieldTypeEnum, FieldTypeSet:
		return DataTypeString
	case FieldTypeBinary, FieldTypeVarbinary, FieldTypeBlob:
		return DataTypeBytes
	default:
		return DataTypeUnknown
	}
}
