package temps

import (
	"os"
	"text/template"
)

func FiberTracerFrame(frame string) {
	// ############################################################
	RenderData.FrameName = frame
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
	// ####################################
	prom_tmpl, err := template.New("RenderData").Parse(promethusMetricsTemplate)
	if err != nil {
		panic(err)
	}

	// Create the models directory if it does not exist
	err = os.MkdirAll("observe", os.ModePerm)
	if err != nil {
		panic(err)
	}

	prom_file, err := os.Create("prometheus.yml")
	if err != nil {
		panic(err)
	}
	defer prom_file.Close()

	err = prom_tmpl.Execute(prom_file, RenderData)
	if err != nil {
		panic(err)
	}

}

func StandardTracerFrame(frame string) {
	// ############################################################
	RenderData.ProjectName = ProjectSettings.ProjectName
	RenderData.FrameName = frame
	var common_tmpl *template.Template
	var err error
	if frame == "echo" {
		common_tmpl, err = template.New("RenderData").Funcs(FuncMap).Parse(standardTracerTemplate)
		if err != nil {
			panic(err)
		}
	} else {
		common_tmpl, err = template.New("RenderData").Funcs(FuncMap).Parse(tracerTemplate)
		if err != nil {
			panic(err)
		}
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

	// ############################################################
	prom_tmpl, err := template.New("RenderData").Parse(promethusMetricsTemplate)
	if err != nil {
		panic(err)
	}

	// Create the models directory if it does not exist
	err = os.MkdirAll("observe", os.ModePerm)
	if err != nil {
		panic(err)
	}

	prom_file, err := os.Create("prometheus.yml")
	if err != nil {
		panic(err)
	}
	defer prom_file.Close()

	err = prom_tmpl.Execute(prom_file, RenderData)
	if err != nil {
		panic(err)
	}

}

func PrometheusTracerFrame(frame string) {
	// ############################################################
	common_tmpl, err := template.New("RenderData").Parse(promSettingTemplate)
	if err != nil {
		panic(err)
	}

	// Create the models directory if it does not exist
	err = os.MkdirAll("observe", os.ModePerm)
	if err != nil {
		panic(err)
	}

	common_file, err := os.Create("observe/prometheus_bucket.go")
	if err != nil {
		panic(err)
	}
	defer common_file.Close()
	if frame == "" {
		frame = "standard"
	}

	var frame_input = struct {
		Frame       string
		ProjectName string
	}{Frame: frame, ProjectName: RenderData.ProjectName}
	err = common_tmpl.Execute(common_file, frame_input)
	if err != nil {
		panic(err)
	}

}

func GenericTracerTemplate() {
	// ############################################################
	common_tmpl, err := template.New("RenderData").Parse(genericAppTracerTemplate)
	if err != nil {
		panic(err)
	}

	// Create the models directory if it does not exist
	err = os.MkdirAll("observe", os.ModePerm)
	if err != nil {
		panic(err)
	}

	common_file, err := os.Create("observe/generic.go")
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

var AppTracer = otel.Tracer(fmt.Sprintf("cli-server-%v", configs.AppConfig.GetOrDefault("APP_NAME", "blue-app")))

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

	const traceExporterFiber = "jaeger"
	exportedTo := fmt.Sprintf("%s:%s", tracerHost, tracerPort)
	fmt.Println("Exporting traces to:", exportedTo)

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


	"{{.ProjectName}}/configs"
	"go.opentelemetry.io/otel"
	{{- if  eq .FrameName  "echo" }}
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/attribute"
	{{- end }}
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

var AppTracer = otel.Tracer(fmt.Sprintf("cli-server-%v", configs.AppConfig.GetOrDefault("APP_NAME", "blue-app")))

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

	const traceExporterEcho = "jaeger"

	if (traceExporter != "" && tracerHost != "") || traceExporter == traceExporterEcho {
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

{{- if  eq .FrameName  "echo" }}

func EchoAppSpanner(ctx echo.Context, span_name string) (context.Context, oteltrace.Span) {
	gen, _ := uuid.NewV7()
	id := gen.String()

	//  getting request body
	trace, span := AppTracer.Start(ctx.Request().Context(), span_name,
		oteltrace.WithAttributes(attribute.String("id", id)),

	)
	return trace, span
}
{{- end}}

type RouteTracer struct {
	Tracer context.Context
	Span   oteltrace.Span
}

`

var promSettingTemplate = `
package observe

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"{{.ProjectName}}/configs"
{{- if  eq .Frame  "fiber" }}
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
{{- else if  eq .Frame  "echo" }}
	"github.com/labstack/echo/v4"
{{- end}}
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

var (
	// CPU usage gauge
	cpuUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cpu_usage_percentage",
			Help: "Current CPU usage percentage",
		},
		[]string{"cpu", "service"},
	)

	// Memory usage gauge
	memUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "memory_usage_percentage",
			Help: "Current memory usage percentage",
		},
		[]string{"memory", "service"},
	)

	// Register HTTP request duration metrics
	httpDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of HTTP request duration (seconds) by method, path, and status code.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status_code", "service"},
	)
)

func InitProm(prom *prometheus.Registry) {
	// Register metrics
	prom.MustRegister(cpuUsage)
	prom.MustRegister(memUsage)
	prom.MustRegister(httpDuration)

	// Start collecting system metrics in a goroutine
	go collectSystemMetrics()
}

func collectSystemMetrics() {
	for {
		// Get CPU usage
		cpuPercent, err := cpu.Percent(0, false)
		if err != nil {
			log.Println("Error getting CPU usage:", err)
		} else {
			// Set CPU usage in the metric
			cpuUsage.WithLabelValues("total", configs.AppConfig.Get("APP_NAME")).Set(cpuPercent[0])
		}

		// Get memory usage
		v, err := mem.VirtualMemory()
		if err != nil {
			log.Println("Error getting memory usage:", err)
		} else {
			// Set memory usage in the metric
			memUsage.WithLabelValues("total", configs.AppConfig.Get("APP_NAME")).Set(float64(v.Used) / float64(v.Total) * 100)
		}

		// Sleep for 10 seconds before checking again
		time.Sleep(5 * time.Second)
	}
}
{{- if  eq .Frame  "fiber" }}
func RegisterMetricsHandler(prom *prometheus.Registry) http.Handler {
	// Return the handler for Prometheus scraping using the custom registry
	return promhttp.HandlerFor(prom, promhttp.HandlerOpts{})
}

func SetupPrometheusMetrics(app *fiber.App) {
	// Initialize Prometheus registry
	prom := prometheus.NewRegistry()
	InitProm(prom)

	// Middleware to track HTTP request duration
	app.Use(func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start).Seconds()

		httpDuration.WithLabelValues(c.Method(), c.Path(), fmt.Sprintf("%d", c.Response().StatusCode()), configs.AppConfig.Get("APP_NAME")).Observe(duration)
		return err
	})

	// Expose the /metrics endpoint for Prometheus to scrape
	app.Get("/metrics", adaptor.HTTPHandler(RegisterMetricsHandler(prom)))
}
{{- else if eq .Frame "echo"}}
func RegisterMetricsHandler(prom *prometheus.Registry) http.Handler {
	// Return the handler for Prometheus scraping using the custom registry
	return promhttp.HandlerFor(prom, promhttp.HandlerOpts{})
}

// WrapHTTPHandler takes an http.Handler and returns an echo.HandlerFunc
func WrapHTTPHandler(h http.Handler) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Use the http.Handler to serve the HTTP request
		h.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}

func SetupPrometheusMetrics(app *echo.Echo) {
	// Initialize Prometheus registry
	prom := prometheus.NewRegistry()
	InitProm(prom)

	app.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c) // Call the next handler
			duration := time.Since(start).Seconds()

			// Record the duration (add logic to log or track metrics here)
			httpDuration.WithLabelValues(c.Request().Method, // HTTP Method (GET, POST, etc.)
				c.Request().URL.Path,                   // Request Path (e.g., "/api/v1/resource")
				fmt.Sprintf("%d", c.Response().Status), // Status code
				configs.AppConfig.Get("APP_NAME")).Observe(duration)

			return err
		}
	})

	// return err
	// Expose the /metrics endpoint for Prometheus to scrape
	app.GET("/metrics", WrapHTTPHandler(RegisterMetricsHandler(prom)))
}
{{- end}}
`

var promethusMetricsTemplate = `
global:
  scrape_interval: 15s  # How frequently to scrape targets

scrape_configs:
  - job_name: 'example-job'
    static_configs:
      - targets: ['localhost:7500']  # Scrape metrics from localhost:7500
    metrics_path: '/metrics'  # Default path is /metrics, but can be set explicitly


`

var genericAppTracerTemplate = `
package observe

import (
	"context"

	"{{.ProjectName}}/configs"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func AppSpanner(ctx context.Context) (context.Context, oteltrace.Span) {
	gen, _ := uuid.NewV7()
	id := gen.String()

	trace, span := AppTracer.Start(ctx, configs.AppConfig.GetOrDefault("APP_NAME", "blue-app"),
		oteltrace.WithAttributes(attribute.String("id", id)),
	)

	return trace, span
}

`
