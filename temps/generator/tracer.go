package generator

import (
	"os"

	"github.com/bushubdegefu/blue-rest/temps"
)

func GenerateTracerFiberSetup(data temps.Data) {
	tmpl := temps.LoadTemplate("fiberTracer")
	tmplMetric := temps.LoadTemplate("promyml")
	err := os.MkdirAll("observe", os.ModePerm)
	if err != nil {
		panic(err)
	}
	temps.WriteTemplateToFile("observe/tracer.go", tmpl, data)
	temps.WriteTemplateToFile("prometheus.yml", tmplMetric, data)
}

func GenerateTracerEchoSetup(data temps.Data) {
	tmpl := temps.LoadTemplate("echoTracer")
	tmplMetric := temps.LoadTemplate("promyml")
	err := os.MkdirAll("observe", os.ModePerm)
	if err != nil {
		panic(err)
	}
	temps.WriteTemplateToFile("observe/tracer.go", tmpl, data)
	temps.WriteTemplateToFile("prometheus.yml", tmplMetric, data)
}

func GenerateTracerGenericSetup(data temps.Data) {
	tmpl := temps.LoadTemplate("genericTracer")

	err := os.MkdirAll("observe", os.ModePerm)
	if err != nil {
		panic(err)
	}
	temps.WriteTemplateToFile("observe/generic.go", tmpl, data)

}

func GeneratePromMetricsSetup(data temps.Data, frame string) {
	tmpl := temps.LoadTemplate("prometheus")

	err := os.MkdirAll("observe", os.ModePerm)
	if err != nil {
		panic(err)
	}
	data.Frame = frame
	temps.WriteTemplateToFile("observe/prometheus_bucket.go", tmpl, data)

}
