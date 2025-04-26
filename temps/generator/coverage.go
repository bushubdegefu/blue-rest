package generator

import (
	"fmt"
	"os"
	"strings"

	"github.com/bushubdegefu/blue-rest/temps"
)

// For Echo coverage testing, we need to generate a test file for each model
func GenerateEchoCoverage(data temps.Data) {
	tmplSetting := temps.LoadTemplate("echoCoverSetting")
	tmplTests := temps.LoadTemplate("echoCover")
	err := os.MkdirAll("testsetting", os.ModePerm)
	if err != nil {
		panic(err)
	}

	_ = os.MkdirAll("tests", os.ModePerm)

	for _, model := range data.Models {
		filePath := fmt.Sprintf("tests/%s_controller_test.go", strings.ToLower(model.Name))
		temps.WriteTemplateToFileModel(filePath, tmplTests, model)
	}
	temps.WriteTemplateToFile("testsetting/settings.go", tmplSetting, data)
}

// For Fiber coverage testing, we need to generate a test file for each model
func GenerateFiberCoverage(data temps.Data) {
	tmplSetting := temps.LoadTemplate("fiberCoverSetting")
	tmplTests := temps.LoadTemplate("fiberCover")
	err := os.MkdirAll("testsetting", os.ModePerm)
	if err != nil {
		panic(err)
	}

	_ = os.MkdirAll("tests", os.ModePerm)

	for _, model := range data.Models {
		filePath := fmt.Sprintf("tests/%s_controller_test.go", strings.ToLower(model.Name))
		temps.WriteTemplateToFileModel(filePath, tmplTests, model)
	}
	temps.WriteTemplateToFile("testsetting/settings.go", tmplSetting, data)
}
