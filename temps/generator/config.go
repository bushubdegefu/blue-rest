package generator

import (
	"os"

	"github.com/bushubdegefu/blue-rest/temps"
)

func GenerateConfig(data temps.Data) {
	tmpl := temps.LoadTemplate("config")
	err := os.MkdirAll("configs", os.ModePerm)
	if err != nil {
		panic(err)
	}
	temps.WriteTemplateToFile("configs/configs.go", tmpl, data)
}

func GenerateConfigEnv(data temps.Data) {
	tmpl := temps.LoadTemplate("env")

	temps.WriteTemplateToFile("configs/.env", tmpl, data)
}

func GenerateConfigAppEnv(data temps.Data) {
	tmpl := temps.LoadTemplate("projectEnv")

	temps.WriteTemplateToFile("configs/.dev.env", tmpl, data)
	temps.WriteTemplateToFile("configs/.prod.env", tmpl, data)
}

func GenerateConfigTestEnv(data temps.Data) {
	tmpl := temps.LoadTemplate("testEnv")
	err := os.MkdirAll("tests", os.ModePerm)
	if err != nil {
		panic(err)
	}
	temps.WriteTemplateToFile("tests/.test.env", tmpl, data)
}
