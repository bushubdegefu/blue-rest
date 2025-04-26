package generator

import (
	"os"

	"github.com/bushubdegefu/blue-rest/temps"
)

func GenerateJWTUtils(data temps.ProjectSetting) {
	tmpl := temps.LoadTemplate("jwtUtils")
	err := os.MkdirAll("utils", os.ModePerm)
	if err != nil {
		panic(err)
	}
	temps.WriteTemplateToFileSetting("utils/utils.go", tmpl, data)
}
