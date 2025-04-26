package temps

import (
	"embed"
	"fmt"
	"os"
	"text/template"
)

//go:embed templates/*.tmpl
var TemplateFS embed.FS

func WriteTemplateToFile(filePath string, tmpl *template.Template, data Data) {
	f, err := os.Create(filePath)
	if err != nil {
		panic(fmt.Errorf("failed to create file %s: %w", filePath, err))
	}
	defer f.Close()

	if err := tmpl.Execute(f, data); err != nil {
		panic(fmt.Errorf("failed to execute template on %s: %w", filePath, err))
	}
}

func WriteTemplateToFileSetting(filePath string, tmpl *template.Template, data ProjectSetting) {
	f, err := os.Create(filePath)
	if err != nil {
		panic(fmt.Errorf("failed to create file %s: %w", filePath, err))
	}
	defer f.Close()

	if err := tmpl.Execute(f, data); err != nil {
		panic(fmt.Errorf("failed to execute template on %s: %w", filePath, err))
	}
}

func WriteTemplateToFileModel(filePath string, tmpl *template.Template, data Model) {
	f, err := os.Create(filePath)
	if err != nil {
		panic(fmt.Errorf("failed to create file %s: %w", filePath, err))
	}
	defer f.Close()

	if err := tmpl.Execute(f, data); err != nil {
		panic(fmt.Errorf("failed to execute template on %s: %w", filePath, err))
	}
}

func LoadTemplate(name string) *template.Template {
	tmplContent, err := TemplateFS.ReadFile("templates/" + name + ".tmpl")
	if err != nil {
		panic(fmt.Errorf("failed to read embedded template: %w", err))
	}
	tmpl, err := template.New(name).Funcs(FuncMap).Parse(string(tmplContent))
	if err != nil {
		panic(fmt.Errorf("failed to parse template: %w", err))
	}
	return tmpl
}
