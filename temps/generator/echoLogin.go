package generator

import (
	"github.com/bushubdegefu/blue-rest/temps"
)

func GenerateEchoLogin(data temps.ProjectSetting) {
	tmpl := temps.LoadTemplate("echoLogin")

	temps.WriteTemplateToFileSetting("controllers/login.go", tmpl, data)
}
