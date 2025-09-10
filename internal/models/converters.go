package models

// ConvertValueTypeToModel ...
func ConvertValueTypeToModel(valueType string) ValueType {
	switch valueType {
	case "string":
		return ValueTypeString
	case "bool":
		return ValueTypeBool
	case "int":
		return ValueTypeInt
	case "int64":
		return ValueTypeInt64
	case "uint":
		return ValueTypeUint
	case "uint64":
		return ValueTypeUint64
	case "float":
		return ValueTypeFloat
	case "float64":
		return ValueTypeFloat64
	default:
		return ValueTypeUnknown
	}
}

// ConvertAuditActionToModel ...
func ConvertAuditActionToModel(auditAction string) AuditAction {
	switch auditAction {
	case "config_updated":
		return AuditActionConfigUpdated
	default:
		return AuditActionUnknown
	}
}
