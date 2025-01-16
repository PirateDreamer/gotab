package logic

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/xwb1989/sqlparser"
)

type ProtoGenInfo struct {
	Table string
	Files []string
}

func SqlToProto(sql string, delPrefix bool) (protoStr string, err error) {
	// 解析sql
	var ddl *sqlparser.DDL
	if ddl, err = ParseDDl(sql); err != nil {
		return
	}
	protoGenInfo := ProtoGenInfo{Files: make([]string, 0), Table: ddl.Table.Name.String()}
	for index, fileInfo := range ddl.TableSpec.Columns {
		var comment string
		if fileInfo.Type.Comment != nil {
			comment = fmt.Sprintf(" // %s", string(fileInfo.Type.Comment.Val))
		}
		if fileInfo.Name.String() == "" {
			continue
		}
		fileGen := fmt.Sprintf("%s %s = %d;%s", GetProtoType(fileInfo.Type.Type), fileInfo.Name.String(), index+1, comment)
		protoGenInfo.Files = append(protoGenInfo.Files, fileGen)
	}
	protoTemplate := `
message {{.Table}} {
	{{- range $item := .Files}}
	{{$item -}}
	{{- end}}
}
	`
	tmpl, err := template.New("proto").Parse(protoTemplate)
	if err != nil {
		panic(err)
	}
	var codeStr strings.Builder
	err = tmpl.Execute(&codeStr, protoGenInfo)
	if err != nil {
		return
	}
	protoStr = codeStr.String()

	return
}
