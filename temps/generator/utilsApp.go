package generator

import (
	"os"

	"github.com/bushubdegefu/blue-rest/temps"
)

func GenerateUtilsApp(data temps.ProjectSetting) {
	tmpl := temps.LoadTemplate("utilsApp")
	err := os.MkdirAll("utils", os.ModePerm)
	if err != nil {
		panic(err)
	}
	temps.WriteTemplateToFileSetting("utils/jwt_utils.go", tmpl, data)
}
