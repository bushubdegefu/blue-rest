package generator

import (
	"os"

	"github.com/bushubdegefu/blue-rest/temps"
)

func GenerateMicrosoftOauth(data temps.Data) {
	tmpl := temps.LoadTemplate("microsoftOauth")
	err := os.MkdirAll("oauth", os.ModePerm)
	if err != nil {
		panic(err)
	}
	temps.WriteTemplateToFile("oauth/microsoft.go", tmpl, data)
}
