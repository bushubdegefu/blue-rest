package temps

import (
	"fmt"
	"os"
	"text/template"
)

func FiberTracerFrame() {
	// Open the JSON file
	file, err := os.Open("config.json")
	if err != nil {
		fmt.Println("Error opening JSON file:", err)
		return
	}
	defer file.Close() // Defer closing the file until the function returns

	// Decode the JSON content into the data structure

	// ############################################################
	common_tmpl, err := template.New("RenderData").Parse(tracerTemplate)
	if err != nil {
		panic(err)
	}

	// Create the models directory if it does not exist
	err = os.MkdirAll("observe", os.ModePerm)
	if err != nil {
		panic(err)
	}

	common_file, err := os.Create("observe/tracer.go")
	if err != nil {
		panic(err)
	}
	defer common_file.Close()

	err = common_tmpl.Execute(common_file, RenderData)
	if err != nil {
		panic(err)
	}

}

func StandardTracerFrame() {
	// Open the JSON file
	file, err := os.Open("config.json")
	if err != nil {
		fmt.Println("Error opening JSON file:", err)
		return
	}
	defer file.Close() // Defer closing the file until the function returns

	// Decode the JSON content into the data structure

	// ############################################################
	common_tmpl, err := template.New("RenderData").Parse(standardTracerTemplate)
	if err != nil {
		panic(err)
	}

	// Create the models directory if it does not exist
	err = os.MkdirAll("observe", os.ModePerm)
	if err != nil {
		panic(err)
	}

	common_file, err := os.Create("observe/tracer.go")
	if err != nil {
		panic(err)
	}
	defer common_file.Close()

	err = common_tmpl.Execute(common_file, RenderData)
	if err != nil {
		panic(err)
	}

}

var tracerTemplate = `
package observe

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"{{.ProjectName}}/configs"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

var AppTracer = otel.Tracer(fmt.Sprintf("fiber-server-%v", configs.AppConfig.Get("APP_NAME")))

func InitTracer() *sdktrace.TracerProvider {
	traceExporter := configs.AppConfig.Get("TRACE_EXPORTER")
	tracerHost := configs.AppConfig.Get("TRACER_HOST")
	tracerPort := configs.AppConfig.GetOrDefault("TRACER_PORT", "9411")

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(configs.AppConfig.Get("APP_NAME")),
		)),
	)
	// app logger with jager

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	// otel.SetErrorHandler(&otelErrorHandler{logger: app_logger})

	const traceExporterFiber = "fiber"

	if (traceExporter != "" && tracerHost != "") || traceExporter == traceExporterFiber {
		var (
			exporter sdktrace.SpanExporter
			// err      error
		)

		switch strings.ToLower(traceExporter) {
		case "jaeger":
			// app_logger.Log("Exporting traces to jaeger.")

			exporter, _ = otlptracegrpc.New(context.Background(), otlptracegrpc.WithInsecure(),
				otlptracegrpc.WithEndpoint(fmt.Sprintf("%s:%s", tracerHost, tracerPort)))

			batcher := sdktrace.NewBatchSpanProcessor(exporter)
			tp.RegisterSpanProcessor(batcher)
		}
	}
	return tp
}


func FiberAppSpanner(ctx *fiber.Ctx, span_name string ) (context.Context, oteltrace.Span) {
	gen, _ := uuid.NewV7()
	id := gen.String()

	trace, span := AppTracer.Start(ctx.UserContext(), span_name,
		oteltrace.WithAttributes(attribute.String("id", id)),
		oteltrace.WithAttributes(attribute.String("request", ctx.Request().String())),
	)
	return trace, span
}

type RouteTracer struct {
	Tracer context.Context
	Span   oteltrace.Span
}

`

var standardTracerTemplate = `
package observe

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"{{.ProjectName}}/configs"
	"go.opentelemetry.io/otel"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

var AppTracer = otel.Tracer(fmt.Sprintf("fiber-server-%v", configs.AppConfig.Get("APP_NAME")))

func InitTracer() *sdktrace.TracerProvider {
	traceExporter := configs.AppConfig.Get("TRACE_EXPORTER")
	tracerHost := configs.AppConfig.Get("TRACER_HOST")
	tracerPort := configs.AppConfig.GetOrDefault("TRACER_PORT", "9411")

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(configs.AppConfig.Get("APP_NAME")),
		)),
	)
	// app logger with jager

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	// otel.SetErrorHandler(&otelErrorHandler{logger: app_logger})

	const traceExporterFiber = "fiber"

	if (traceExporter != "" && tracerHost != "") || traceExporter == traceExporterFiber {
		var (
			exporter sdktrace.SpanExporter
			// err      error
		)

		switch strings.ToLower(traceExporter) {
		case "jaeger":
			// app_logger.Log("Exporting traces to jaeger.")

			exporter, _ = otlptracegrpc.New(context.Background(), otlptracegrpc.WithInsecure(),
				otlptracegrpc.WithEndpoint(fmt.Sprintf("%s:%s", tracerHost, tracerPort)))

			batcher := sdktrace.NewBatchSpanProcessor(exporter)
			tp.RegisterSpanProcessor(batcher)
		}
	}
	return tp
}

func EchoAppSpanner(ctx echo.Context, span_name string) (context.Context, oteltrace.Span) {
	gen, _ := uuid.NewV7()
	id := gen.String()

	//  getting request body
	trace, span := AppTracer.Start(ctx.Request().Context(), span_name,
		oteltrace.WithAttributes(attribute.String("id", id)),

	)
	return trace, span
}


type RouteTracer struct {
	Tracer context.Context
	Span   oteltrace.Span
}


`
