package internal

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/xwb1989/sqlparser"
	"html/template"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// 获取sql的途径，生成映射的途径

// 获取createSql、解析sql、结构化、生成模版

func Run() {

}

func getCreateTableSql(option *Option) (createSQLs []string, err error) {
	var db *sql.DB
	if db, err = sql.Open("mysql", option.Dsn); err != nil {
		return
	}
	defer db.Close()

	// get tables name
	tables := make([]string, 0)
	if !option.AllTable {
		tables = option.TableNames
	} else {
		query := "SHOW TABLES"
		var rows *sql.Rows
		rows, err = db.Query(query)
		defer rows.Close()
		if err != nil {
			return
		}

		for rows.Next() {
			var tableName string
			err = rows.Scan(&tableName)
			if err != nil {
				return
			}
			tables = append(tables, tableName)
		}
	}

	// get tables create sql
	for _, tableName := range tables {
		query := fmt.Sprintf("SHOW CREATE TABLE `%s`", tableName)
		var createSQL string
		err = db.QueryRow(query).Scan(&tableName, &createSQL)
		if err != nil {
			return
		}
		createSQLs = append(createSQLs, createSQL)
	}
	return
}

func parseSql(sqlStatement string) (*sqlparser.DDL, error) {
	stmt, err := sqlparser.Parse(sqlStatement)
	if err != nil {
		log.Printf("Error parsing SQL: %s\n", err)
		return nil, err
	}
	switch stmt.(type) {
	case *sqlparser.DDL:
	default:
		return nil, errors.New("sql type is error")
	}
	ddl := stmt.(*sqlparser.DDL)
	if ddl.Action != "create" {
		return nil, errors.New("sql type is error")
	}
	return ddl, err
}

func convertSqlStruct(ddl *sqlparser.DDL, option *Option) (data ModelData) {
	files := make([]Field, 0)
	for _, fileInfo := range ddl.TableSpec.Columns {
		re := regexp.MustCompile(`\(\d+\)`)
		file := Field{
			Name:    convertToCamelCase(fileInfo.Name.String(), false),
			Type:    TypeMap[strings.ToUpper(re.ReplaceAllString(fileInfo.Type.Type, ""))],
			JSONTag: fileInfo.Name.String(),
			GormTag: "column:" + fileInfo.Name.String(),
		}
		if file.Name == "DeletedAt" {
			file.Type = "gorm.DeletedAt"
		}
		if fileInfo.Type.Comment != nil {
			file.COMMENT = string(fileInfo.Type.Comment.Val)
		}
		if fileInfo.Name.String() == ddl.TableSpec.Indexes[0].Columns[0].Column.String() {
			file.GormTag += ";primary_key"
		}
		// 类型
		if fileInfo.Type.NotNull {
			file.GormTag += ";not null"
		}
		if fileInfo.Type.Default != nil && string(fileInfo.Type.Default.Val) != "" {
			file.GormTag += ";default:" + string(fileInfo.Type.Default.Val)
		}
		if fileInfo.Type.Autoincrement {
			file.GormTag += ";auto_increment"
		}

		files = append(files, file)
	}
	data = ModelData{
		ModelName:   convertToCamelCase(ddl.NewName.Name.String(), option.HasRmPrefix),
		TableName:   ddl.NewName.Name.String(),
		Fields:      files,
		PackageName: filepath.Base(option.OutDir),
	}
	return
}

func convertToCamelCase(s string, rmPrefix bool) (result string) {
	words := strings.Split(s, "_")

	for i, word := range words {
		words[i] = strings.Title(word)
	}

	if rmPrefix && len(words) > 1 {
		result = strings.Join(words[1:], "")
		return
	}
	result = strings.Join(words, "")

	return result
}

func generateTemplate(data ModelData, option *Option) (err error) {
	// 加载模板文件
	fileTemplate := `{{.Name}} {{.Type}}`
	if option.HasJsonTag && option.HasGormTag {
		fileTemplate = fileTemplate + " `gorm:\"{{.GormTag}}\" json:\"{{.JSONTag}}\"`"
	} else if option.HasJsonTag && !option.HasGormTag {
		fileTemplate = fileTemplate + " `json:\"{{.JSONTag}}\"`"
	} else if option.HasGormTag && !option.HasJsonTag {
		fileTemplate = fileTemplate + "`gorm:\"{{.GormTag}}\"`"
	}

	if option.HasNote {
		fileTemplate = fileTemplate + ` //{{.COMMENT}}`
	}

	modelTemapte := `
package {{.PackageName}}

import (
	"gorm.io/gorm"
	"time"
)

type {{.ModelName}} struct {
	{{range .Fields -}}
	` + fileTemplate + `
	{{end -}}
}
`
	if option.HasTableName {
		modelTemapte = modelTemapte + `
func (m *{{.ModelName}}) TableName() string {
	return "{{.TableName}}"
}`
	}
	tmpl, err := template.New("model").Parse(modelTemapte)
	if err != nil {
		panic(err)
	}

	// 生成文件
	var codeStr strings.Builder
	err = tmpl.Execute(&codeStr, data)
	if err != nil {
		return
	}

	if err = CreateDir(filepath.Base(option.OutDir)); err != nil {
		log.Printf("error creating a out dir: %s", err.Error())
		return
	}

	newFile, err := os.Create(option.OutDir + "/" + data.ModelName + ".go")
	if err != nil {
		log.Printf("error creating a model file: %s", err.Error())
		return
	}
	defer newFile.Close()

	// 写入生成的代码
	_, err = newFile.WriteString(string(codeStr.String()))
	if err != nil {
		return
	}

	// 规范代码结构
	gofmtCmd := exec.Command("gofmt", "-w", option.OutDir+"/"+data.ModelName+".go")
	err = gofmtCmd.Run()
	if err != nil {
		log.Fatalf("Failed to run gofmt command: %s", err)
		return
	}

	// 导入依赖包
	goimportsCmd := exec.Command("goimports", "-w", option.OutDir+"/"+data.ModelName+".go")
	err = goimportsCmd.Run()
	if err != nil {
		log.Fatalf("Failed to run goimports command: %s", err)
	}
	return
}

func CreateDir(folderPath string) error {
	_, err := os.Stat(folderPath)
	if os.IsNotExist(err) {
		// 文件夹不存在，创建文件夹
		err := os.Mkdir(folderPath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}
