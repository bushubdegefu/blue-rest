package generator

import "github.com/bushubdegefu/blue-rest/temps"

func GenerateFiberLogin(data temps.ProjectSetting) {
	tmpl := temps.LoadTemplate("fiberLogin")

	temps.WriteTemplateToFileSetting("controllers/app.go", tmpl, data)
}
