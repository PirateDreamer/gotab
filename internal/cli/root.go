package cli

import (
	"fmt"
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
			case "tableName":
				TableName = true
			case "rmPrefix":
				RmPrefix = true
			case "jsonTag":
				JsonTag = true
			case "gormTag":
				GormTag = true
			case "comment":
				Comment = true
			case "all":
				TableName = true
				RmPrefix = true
				JsonTag = true
				GormTag = true
				Comment = true
			}
		}
		if Tables == "all" {
			GenAllTable = true
		}
	},
}

func init() {
	model.Flags().StringVarP(&Dsn, "dsn", "d", "", "database dsn")
	model.Flags().StringVarP(&Tables, "tables", "t", "", "table name")
	model.Flags().StringVarP(&OutDir, "outDir", "o", "", "output directory")
	model.MarkFlagRequired("dsn")
	model.MarkFlagRequired("tables")
	model.MarkFlagRequired("outDir")
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
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
