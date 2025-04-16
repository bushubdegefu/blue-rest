package temps

import (
	"os"
	"text/template"
)

func MigrationFrame() {
	//  this is creating manger file inside the manager folder
	RenderData.ProjectName = ProjectSettings.ProjectName
	RenderData.AppNames = ProjectSettings.AppNames
	RenderData.BackTick = "`"
	// ############################################################
	migration_tmpl, err := template.New("RenderData").Funcs(FuncMap).Parse(migrationTemplate)
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
	{{- range .AppNames}}
	{{ . | replaceString }} "{{$.ProjectName}}/{{ . }}/models"
	{{- end }}
	"github.com/spf13/cobra"
)

var (
	{{.AppName | replaceString }}migrate= &cobra.Command{
		Use:   "migrate",
		Short: "Run Database Migration for found in init migration Models",
		Long:  "Migrate to init database",
		Run: func(cmd *cobra.Command, args []string) {
			// Check for the type flag
			migrateType, _ := cmd.Flags().GetString("type")
			if migrateType == "create" {
				init_migrate()
			} else if migrateType == "stats" {
				create_views()
			} else {
				populate_migrate()
			}
		},
	}

	{{.AppName | replaceString }}clean= &cobra.Command{
		Use:   "clean",
		Short: "Drop Database Models for found in init migration Models",
		Long:  "Drop Models found in the models definition",
		Run: func(cmd *cobra.Command, args []string) {
			clean_database()
		},
	}

	{{ if eq .AuthAppName .AppName }}
	createsuperuser = &cobra.Command{
		Use:   "superuser",
		Short: "Create super user",
		Long:  "Create super user",
		Run: func(cmd *cobra.Command, args []string) {
			{{ .AuthAppName | replaceString }}.CreateSuperUser()
			fmt.Println("Super user created successfully")
		},
	}
	{{- end}}

)

func init_migrate() {

    {{- range .AppNames}}
    {{ . | replaceString }}.InitDatabase(false)
    {{- end }}
	fmt.Println("Migrated Database Models sucessfully")
}

func populate_migrate() {
{{- if eq .AuthAppType "standalone" }}
    {{- range .AppNames}}
	{{ . | replaceString }}.Populate(false)
	{{- end}}
{{- end}}
	fmt.Println("Auth Populated Sucessfuly Database Models sucessfully")
}

func create_views() {
	{{- range .AppNames}}
	{{ . | replaceString }}.CreateStatsDatabase()
	{{- end}}
	fmt.Println("Auth Created App stat views")
}

func clean_database() {
	{{- range .AppNames}}
	{{ . | replaceString }}.CleanDatabase(false)
	{{- end}}
	fmt.Println("Dropped Tables sucessfully")
}

func init() {
	{{.AppName | replaceString }}migrate.Flags().StringP("type", "t", "", "Specify create to \"create\" the models to database, and \"populate\" to populate default permissions and content types")
	goFrame.AddCommand({{.AppName | replaceString }}migrate)
	goFrame.AddCommand({{.AppName | replaceString }}clean)
	goFrame.AddCommand(createsuperuser)
}

`
