package generator

import (
	"os"

	"github.com/bushubdegefu/blue-rest/temps"
)

func GenerateGoogleOauth(data temps.Data) {
	tmpl := temps.LoadTemplate("googleOauth")
	err := os.MkdirAll("oauth", os.ModePerm)
	if err != nil {
		panic(err)
	}
	temps.WriteTemplateToFile("oauth/google.go", tmpl, data)
}
