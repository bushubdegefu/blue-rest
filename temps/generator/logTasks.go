package generator

import (
	"os"

	"github.com/bushubdegefu/blue-rest/temps"
)

func GenerateTasks(data temps.Data) {
	tmpl := temps.LoadTemplate("tasks")
	err := os.MkdirAll("bluetasks", os.ModePerm)
	if err != nil {
		panic(err)
	}
	temps.WriteTemplateToFile("bluetasks/tasks.go", tmpl, data)
}

func GenerateLogs(data temps.Data) {
	tmpl := temps.LoadTemplate("logs")
	err := os.MkdirAll("logs", os.ModePerm)
	if err != nil {
		panic(err)
	}
	temps.WriteTemplateToFile("logs/logfile.go", tmpl, data)
}
