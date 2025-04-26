package generator

import (
	"os"

	"github.com/bushubdegefu/blue-rest/temps"
)

func GenerateAppDatabaseMigration(data temps.Data) {
	tmpl := temps.LoadTemplate("migrationApp")
	err := os.MkdirAll("manager", os.ModePerm)
	if err != nil {
		panic(err)
	}
	temps.WriteTemplateToFile("manager/migration.go", tmpl, data)
}
