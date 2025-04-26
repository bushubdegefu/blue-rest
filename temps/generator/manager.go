package generator

import (
	"os"

	"github.com/bushubdegefu/blue-rest/temps"
)

func GenerateMainAndManager(data temps.Data) {
	tmplMain := temps.LoadTemplate("main")
	tmplManager := temps.LoadTemplate("manager")
	err := os.MkdirAll("manager", os.ModePerm)
	if err != nil {
		panic(err)
	}
	temps.WriteTemplateToFile("main.go", tmplMain, data)
	temps.WriteTemplateToFile("manager/manager.go", tmplManager, data)
}
