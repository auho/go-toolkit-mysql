package schema

func FileTypeToDataType(ft FieldType) DataType {
	switch ft {
	case FieldTypeTinyint, FieldTypeSmallint, FieldTypeMediumint, FieldTypeInt, FieldTypeInteger, FieldTypeBigint:
		return DataTypeInt
	case FieldTypeDecimal, FieldTypeFloat, FieldTypeDouble:
		return DataTypeFloat
	case FieldTypeChar, FieldTypeVarchar:
		return DataTypeString
	default:
		return DataTypeUnknown
	}
}
