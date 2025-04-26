package generator

import (
	"os"

	"github.com/bushubdegefu/blue-rest/temps"
)

func GenerateHelperOauth(data temps.Data) {
	tmpl := temps.LoadTemplate("helperOauth")
	err := os.MkdirAll("oauth", os.ModePerm)
	if err != nil {
		panic(err)
	}
	temps.WriteTemplateToFile("oauth/helper.go", tmpl, data)
}
