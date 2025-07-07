package generator

import (
	"github.com/bushubdegefu/blue-rest/temps"
)

func GenerateEchoLogin(data temps.ProjectSetting) {
	tmpl := temps.LoadTemplate("echoLogin")

	temps.WriteTemplateToFileSetting("controllers/app.go", tmpl, data)
}
