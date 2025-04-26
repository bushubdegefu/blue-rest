package generator

import (
	"github.com/bushubdegefu/blue-rest/temps"
)

// ############################################
//
//	for Fiber App Generat functions
//
// ############################################
func GenerateFiberSetup(data temps.Data) {
	tmpl := temps.LoadTemplate("fiberSetup")

	temps.WriteTemplateToFile("setup.go", tmpl, data)
}

func GenerateFiberAppMiddleware(data temps.Data) {
	tmpl := temps.LoadTemplate("fiberAppMiddleware")

	temps.WriteTemplateToFile("middleware.go", tmpl, data)
}

func GenerateGlobalFiberAppMiddleware(data temps.Data) {
	tmpl := temps.LoadTemplate("globalFiberMiddleware")

	temps.WriteTemplateToFile("manager/middleware.go", tmpl, data)
}

func GenerateAppFiberGlobal(data temps.Data) {
	tmpl := temps.LoadTemplate("globalFiberApp")
	data.SetBackTick()
	temps.WriteTemplateToFile("manager/app.go", tmpl, data)
}

// ############################################
//  for Fiber App Generat functions
// ############################################

func GenerateEchoSetup(data temps.Data) {
	tmpl := temps.LoadTemplate("echoSetup")

	temps.WriteTemplateToFile("setup.go", tmpl, data)
}

func GenerateEchoAppMiddleware(data temps.Data) {
	tmpl := temps.LoadTemplate("echoAppMiddleware")

	temps.WriteTemplateToFile("middleware.go", tmpl, data)
}

func GenerateGlobalEchoAppMiddleware(data temps.Data) {
	tmpl := temps.LoadTemplate("globalEchoMiddleware")

	temps.WriteTemplateToFile("manager/middleware.go", tmpl, data)
}

func GenerateAppEchoGlobal(data temps.Data) {
	tmpl := temps.LoadTemplate("globalEchoApp")

	data.SetBackTick()
	temps.WriteTemplateToFile("manager/app.go", tmpl, data)
}
