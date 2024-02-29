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
	},
}

func InitModel(option *internal.Option) {
	// 初始化model生成命令
	model.Flags().StringVarP(&option.Dsn, "dsn", "d", "", "database dsn")
	model.Flags().StringVarP(&option.ConnPath, "connPath", "c", "./config.yaml", "table connection file path,key is mysql.dsn")
	model.Flags().StringVarP(&option.TableNames, "tables", "s", "", "table name")

	model.Flags().StringVarP(&option.SqlFilePath, "sqlPath", "f", "", "if dns is exist,sql file path")
	model.Flags().StringVarP(&option.OutDir, "outDir", "o", "./model", "output directory")
	model.Flags().StringVarP(&option.MapTypeFile, "mapPath", "m", "", "sql type to go type map file path")

	model.Flags().BoolVarP(&option.HasTableName, "tableName", "t", true, "gen model with table name")
	model.Flags().BoolVarP(&option.HasRmPrefix, "rmPrefix", "r", true, "gen model remove table prefix")

	model.Flags().BoolVarP(&option.HasGormTag, "gormTag", "g", true, "gen model with gorm tag")
	model.Flags().BoolVarP(&option.HasJsonTag, "jsonTag", "j", true, "gen model with json tag")
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

func Execute(option *internal.Option) {
	InitModel(option)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
