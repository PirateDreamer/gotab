package logic

// TypeMap Mapping relationships between database type and golang type
var TypeMap = map[string]string{
	"INT":      "int",
	"VARCHAR":  "string",
	"TINYINT":  "int",
	"DATETIME": "time.Time",
	"TEXT":     "string",
	"LONGTEXT": "string",
	"BIGINT":   "int64",
}

var SqlProtoType = map[string]string{
	"INT":      "int32",
	"VARCHAR":  "string",
	"TINYINT":  "int32",
	"DATETIME": "int64",
	"TEXT":     "string",
	"LONGTEXT": "string",
	"BIGINT":   "uint64",
}
