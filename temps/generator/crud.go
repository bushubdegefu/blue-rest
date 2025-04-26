package generator

import (
	"fmt"
	"os"
	"strings"

	"github.com/bushubdegefu/blue-rest/temps"
)

func GenerateCrudFiber(data temps.Data) {
	tmpl := temps.LoadTemplate("crudFiber")
	_ = os.MkdirAll("controllers", os.ModePerm)

	for _, model := range data.Models {
		filePath := fmt.Sprintf("controllers/%v_controller.go", strings.ToLower(model.Name))

		temps.WriteTemplateToFileModel(filePath, tmpl, model)
	}

}

func GenerateCrudEcho(data temps.Data) {
	tmpl := temps.LoadTemplate("crudEcho")

	_ = os.MkdirAll("controllers", os.ModePerm)

	for _, model := range data.Models {
		filePath := fmt.Sprintf("controllers/%v_controller.go", strings.ToLower(model.Name))
		temps.WriteTemplateToFileModel(filePath, tmpl, model)
	}

}
