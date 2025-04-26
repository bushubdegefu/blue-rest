package generator

import (
	"fmt"
	"os"
	"strings"

	"github.com/bushubdegefu/blue-rest/temps"
)

func GenerateModels(data temps.Data) {
	tmpl := temps.LoadTemplate("model")
	migrationTmpl := temps.LoadTemplate("migration")
	helperTmpl := temps.LoadTemplate("helperModels")

	_ = os.MkdirAll("models", os.ModePerm)

	for _, model := range data.Models {
		filePath := fmt.Sprintf("models/%s.go", strings.ToLower(model.Name))
		temps.WriteTemplateToFileModel(filePath, tmpl, model)
	}

	temps.WriteTemplateToFile("models/init.go", migrationTmpl, data)
	temps.WriteTemplateToFile("models/helper.go", helperTmpl, data)
}
