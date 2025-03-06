package temps

import (
	"os"
	"text/template"
)

func MigrationFrame() {
	//  this is creating manger file inside the manager folder
	// ############################################################
	migration_tmpl, err := template.New("RenderData").Parse(migrationTemplate)
	if err != nil {
		panic(err)
	}

	// Create the models directory if it does not exist
	err = os.MkdirAll("manager", os.ModePerm)
	if err != nil {
		panic(err)
	}

	migration_file, err := os.Create("manager/migrate.go")
	if err != nil {
		panic(err)
	}
	defer migration_file.Close()

	err = migration_tmpl.Execute(migration_file, RenderData)
	if err != nil {
		panic(err)
	}

}

var migrationTemplate = `
package manager

import (
	"fmt"
	"{{.ProjectName}}/models"
	"github.com/spf13/cobra"
)

var (
	{{.AppName}}migrate= &cobra.Command{
		Use:   "migrate",
		Short: "Run Database Migration for found in init migration Models",
		Long:  {{.BackTick}}Migrate to init database{{.BackTick}},
		Run: func(cmd *cobra.Command, args []string) {
			init_migrate()
		},
	}

	{{.AppName}}clean= &cobra.Command{
		Use:   "clean",
		Short: "Drop Database Models for found in init migration Models",
		Long:  {{.BackTick}}Drop Models found in the models definition{{.BackTick}},
		Run: func(cmd *cobra.Command, args []string) {
			clean_database()
		},
	}

)

func init_migrate() {
	models.InitDatabase(false)
	fmt.Println("Migrated Database Models sucessfully")
}

func clean_database() {
	models.CleanDatabase(false)
	fmt.Println("Dropped Tables sucessfully")
}


func init() {
	goFrame.AddCommand({{.AppName}}migrate)
	goFrame.AddCommand({{.AppName}}clean)
}

`
