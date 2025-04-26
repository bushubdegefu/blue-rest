package generator

import (
	"os"

	"github.com/bushubdegefu/blue-rest/temps"
)

func GenerateCommon(data temps.Data) {
	tmpl := temps.LoadTemplate("common")
	err := os.MkdirAll("common", os.ModePerm)
	if err != nil {
		panic(err)
	}
	data.SetBackTick()
	temps.WriteTemplateToFile("common/common.go", tmpl, data)
}

func GenerateDjangoAuth(data temps.Data) {
	tmpl := temps.LoadTemplate("django")

	temps.WriteTemplateToFile("config.json", tmpl, data)
}

func GenerateSSOAuth(data temps.Data) {
	tmpl := temps.LoadTemplate("sso")

	temps.WriteTemplateToFile("config.json", tmpl, data)
}
