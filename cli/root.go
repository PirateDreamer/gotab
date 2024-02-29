package cli

import (
	"fmt"
	"gotab/internal"
	"os"

	"github.com/spf13/cobra"
)

var model = cobra.Command{
	Use:   "model",
	Short: "Generate model",
	Long: `You can add model tags, JSON tags, and table tags to GORM, 
												and remove table prefixes`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, v := range args {
			switch v {
			case "dsn":
				OptionInfo = 
			case "connPath":
			case "tables":
			case "sqlPath":
			case "outDir":
			case "mapPath":
			case "tableName":
			case "rmPrefix":
			case "gormTag":
			case "jsonTag":
			case "note":
			}
		}
	},
}

func InitModel(option *internal.Option) {
	// 初始化model生成命令
	model.Flags().StringVarP(&option.Dsn, "dsn", "d", "", "database dsn")
	model.Flags().StringVarP(&option.ConnPath, "connPath", "cp", "./config.yaml", "table name")
	model.Flags().StringVarP(&option.TableNames, "tables", "t", "", "table name")

	model.Flags().StringVarP(&option.SqlFilePath, "sqlPath", "s", "", "if dns is exist,sql file path")
	model.Flags().StringVarP(&option.OutDir, "outDir", "o", "./model", "output directory")
	model.Flags().StringVarP(&option.MapTypeFile, "mapPath", "m", "", "sql type to go type map file path")

	model.Flags().BoolVarP(&option.HasTableName, "tableName", "tm", true, "gen model with table name")
	model.Flags().BoolVarP(&option.HasRmPrefix, "rmPrefix", "rm", true, "gen model remove table prefix")

	model.Flags().BoolVarP(&option.HasGormTag, "gormTag", "gt", true, "gen model with gorm tag")
	model.Flags().BoolVarP(&option.HasJsonTag, "jsonTag", "jt", true, "gen model with json tag")
	model.Flags().BoolVarP(&option.HasNote, "note", "n", true, "gen model with note")

	rootCmd.AddCommand(&model)
}

var rootCmd = &cobra.Command{
	Use:   "gotab",
	Short: "gotab is a very fast go-web site generator",
	Long: `gotab is a go-web code generator 
					that supports generating backend data models and interfaces`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please fill in the model, after goalt")
	},
}

func Execute() {
	InitModel()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
