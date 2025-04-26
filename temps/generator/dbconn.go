package generator

import (
	"os"

	"github.com/bushubdegefu/blue-rest/temps"
)

func GenerateDBConn(data temps.ProjectSetting) {
	tmpl := temps.LoadTemplate("database")
	err := os.MkdirAll("database", os.ModePerm)
	if err != nil {
		panic(err)
	}
	temps.WriteTemplateToFileSetting("database/database.go", tmpl, data)
}
