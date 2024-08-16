package logic

type Field struct {
	Name    string
	Type    string
	JSONTag string
	GormTag string
	COMMENT string
}

type ModelData struct {
	PackageName string
	ModelName   string
	Fields      []Field
	TableName   string
}

type Option struct {
	// this is gen model config
	HasTableName bool
	HasRmPrefix  bool

	// this is gen field config
	HasJsonTag bool
	HasGormTag bool
	HasNote    bool

	// 模式
	// IsCli bool //cli、web模式（web模式获取sql语句、映射、生成配置）

	ConnType    int    // 0: 无 1: 命令 2: 配置文件，默认是命令
	Dsn         string // 数据库连接地址
	ConnPath    string // 配置文件
	SqlFilePath string // 当连接类型为0，根据sql文件内容，默认是当前项目的model.sql
	AllTable    bool   // 指定表，所有表，默认是所有表
	TableNames  string // 指定生成的表名称
	OutDir      string // 生成model文件目录，默认是当前项目的model目录

	// 字段类型映射
	MapType     int    // 1: 默认是读取map 2: 配置文件
	MapTypeFile string // 配置文件地址,默认是当前项目的type.json
}

func HandlerOption(option *Option) {
	if option.SqlFilePath == "" {
		if option.Dsn != "" {
			option.ConnType = 1
		} else {
			if option.ConnPath != "" {
				option.ConnType = 2
			}
		}
	}
	if option.TableNames == "" {
		option.AllTable = true
	}
	if option.MapTypeFile != "" {
		option.MapType = 2
	}
}

type Options struct {
	F func(*Option)
}
