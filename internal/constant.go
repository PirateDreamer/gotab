package internal

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
