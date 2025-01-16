package logic

import (
	"errors"
	"log"
	"regexp"
	"strings"

	"github.com/xwb1989/sqlparser"
)

// ParseDD; 解析ddl sql
func ParseDDl(sql string) (ddl *sqlparser.DDL, err error) {
	var statement sqlparser.Statement
	if statement, err = sqlparser.Parse(sql); err != nil {
		log.Printf("Error parsing SQL: %s\n", err)
		return
	}
	switch statement.(type) {
	case *sqlparser.DDL:
	default:
		err = errors.New("sql type is error")
		return
	}
	ddl = statement.(*sqlparser.DDL)
	if ddl.Action != "create" {
		err = errors.New("sql type is error")
		return
	}
	return
}

func GetProtoType(sqlFieldType string) string {
	re := regexp.MustCompile(`\(\d+\)`)
	return SqlProtoType[strings.ToUpper(re.ReplaceAllString(sqlFieldType, ""))]
}
