package logic

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Houserqu/ginc"

	_ "github.com/go-sql-driver/mysql"

	"github.com/spf13/viper"
	"github.com/xwb1989/sqlparser"
)

// 获取sql的途径，生成映射的途径

// 获取createSql、解析sql、结构化、生成模版

func Run(option *Option) (err error) {
	var createSqls []string
	switch option.ConnType {
	case 0:
		// 读取sql文件
		var sqlByte []byte
		if sqlByte, err = ioutil.ReadFile(option.SqlFilePath); err != nil {
			return
		}
		createSqls = strings.Split(string(sqlByte), ";")
		createSqls = createSqls[:len(createSqls)-1]
	case 1:
		if createSqls, err = getCreateTableSql(option); err != nil {
			return
		}
	case 2:
		ginc.InitConfig()
		viper.SetConfigName(strings.ReplaceAll(filepath.Base(option.ConnPath), filepath.Ext(option.ConnPath), ""))
		viper.SetConfigType(strings.ReplaceAll(filepath.Ext(option.ConnPath), ".", ""))
		viper.AddConfigPath(filepath.Dir(option.ConnPath))
		if err = viper.ReadInConfig(); err != nil {
			return
		}
		option.Dsn = viper.GetString("mysql.dsn")
		// 解析读取配置文件中的yaml
		if createSqls, err = getCreateTableSql(option); err != nil {
			return
		}
	default:
		log.Println("connType is error,support 0、1、2")
		return
	}

	for _, v := range createSqls {
		var paraseDDL *sqlparser.DDL
		if paraseDDL, err = parseSql(v); err != nil {
			return
		}
		modelData := convertSqlStruct(paraseDDL, option)
		if err = generateTemplate(modelData, option); err != nil {
			return
		}
	}
	log.Println("gen model is end")
	return
}

func getCreateTableSql(option *Option) (createSQLs []string, err error) {
	var db *sql.DB
	if db, err = sql.Open("mysql", option.Dsn); err != nil {
		log.Printf("open mysql error: %v\n", err)
		return
	}
	defer db.Close()

	// get tables name
	tables := make([]string, 0)
	if !option.AllTable {
		tables = strings.Split(option.TableNames, ",")
	} else {
		query := "SHOW TABLES"
		var rows *sql.Rows
		rows, err = db.Query(query)
		if err != nil {
			log.Printf("query mysql error: %v\n", err)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var tableName string
			err = rows.Scan(&tableName)
			if err != nil {
				log.Printf("scan mysql error: %v\n", err)
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
			log.Printf("query mysql error: %v\n", err)
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
		file.GormTag += ";type:" + fileInfo.Type.DescribeType()
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
		log.Printf("error creating a out dir: %s\n", err.Error())
		return
	}

	newFile, err := os.Create(option.OutDir + "/" + data.TableName + ".go")
	if err != nil {
		log.Printf("error creating a model file: %s\n", err.Error())
		return
	}
	defer newFile.Close()

	// 写入生成的代码
	_, err = newFile.WriteString(string(codeStr.String()))
	if err != nil {
		return
	}

	// 规范代码结构
	gofmtCmd := exec.Command("gofmt", "-w", option.OutDir+"/"+data.TableName+".go")
	err = gofmtCmd.Run()
	if err != nil {
		log.Fatalf("Failed to run gofmt command: %s\n", err)
		return
	}

	// 导入依赖包
	goimportsCmd := exec.Command("goimports", "-w", option.OutDir+"/"+data.TableName+".go")
	err = goimportsCmd.Run()
	if err != nil {
		log.Fatalf("Failed to run goimports command: %s\n", err)
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
