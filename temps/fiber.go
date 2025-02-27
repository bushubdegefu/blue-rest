package temps

import (
	"fmt"
	"os"
	"text/template"
)

func FiberFrame() {
	//  this is creating manger file inside the manager folder
	// ############################################################
	devf_tmpl, err := template.New("RenderData").Parse(devfTemplate)
	if err != nil {
		panic(err)
	}

	// Create the models directory if it does not exist
	err = os.MkdirAll("manager", os.ModePerm)
	if err != nil {
		panic(err)
	}

	devf_file, err := os.Create("manager/devfiber.go")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer devf_file.Close()

	err = devf_tmpl.Execute(devf_file, RenderData)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

}

// https://help.sumologic.com/docs/apm/traces/get-started-transaction-tracing/opentelemetry-instrumentation/go/

var devfTemplate = `
package manager

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"net/http"
	"fmt"
	"regexp"
	"log"
	"strconv"
	"time"


	"os"
	"os/signal"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/contrib/otelfiber"
	"github.com/madflojo/tasks"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.opentelemetry.io/otel/attribute"
	"github.com/gofiber/swagger"
	"{{.ProjectName}}/configs"
	"{{.ProjectName}}/bluetasks"
	"{{.ProjectName}}/observe"
	"{{.ProjectName}}/common"
	"{{.ProjectName}}/controllers"
	_ "{{.ProjectName}}/docs"
	"{{.ProjectName}}/database"
	"github.com/spf13/cobra"
)

var (
	env string
	app_otel string
	app_tls string
	{{.AppName}}cli= &cobra.Command{
		Use:   "run",
		Short: "Run Development or Production server Based on Provided --env=dev flag. Defaults to dev ",
		Long:  {{.BackTick}}Run {{.AppName}} development server{{.BackTick}},
		Run: func(cmd *cobra.Command, args []string) {
			switch env {
			case "":
				fiber_run("dev")
			default:
				fiber_run(env)
			}
		},
	}
	protectedURLs = []*regexp.Regexp{
		regexp.MustCompile("^/api/v1/login"),
		regexp.MustCompile("^/api/v1/checklogin"),
		regexp.MustCompile("^/lmetrics"),
		regexp.MustCompile("^/docs"),
		regexp.MustCompile("^/metrics$"),
	}
)

func otelspanstarter(ctx *fiber.Ctx) error {
	//  creating trace context from span if they exist
	route_name := ctx.Path() + "_" + strings.ToLower(ctx.Route().Method)
	tracer, span := observe.FiberAppSpanner(ctx, fmt.Sprintf("%v-root", route_name))
	ctx.Locals("tracer", &observe.RouteTracer{Tracer: tracer, Span: span})
	if err := ctx.Next(); err != nil {
		return err
	}
	span.SetAttributes(attribute.String("response", ctx.Response().String()))
	span.End()
	return nil
}

func dbsessioninjection(ctx *fiber.Ctx) error {
	db, err := database.ReturnSession()
	if err != nil {
		return ctx.Status(http.StatusServiceUnavailable).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}
	ctx.Locals("db", db)
	return ctx.Next()
}

func NextFunc(contx *fiber.Ctx) error {
	return contx.Next()
}

// this is path filter which wavies token requirement for provided paths
func authFilter(c *fiber.Ctx) bool {
	originalURL := strings.ToLower(c.OriginalURL())

	for _, pattern := range protectedURLs {
		if pattern.MatchString(originalURL) {
			c.Request().Header.Add("X-APP-TOKEN", "allowed")
			return true
		}
	}
	return false
}

func NextRoute(contx *fiber.Ctx, key string) (bool, error) {
	contx.Next()
	// fix the below return line based on logic
	// currently will pass unautheticated
	return true, nil
}

func fiber_run(env string) {
	prefork := false
	if env == "prod" {
		prefork = true
	}

	//  Loading Configuration
	configs.AppConfig.SetEnv(env)

	if app_otel == "on" {
		// Starting Otel Global tracer
		tp := observe.InitTracer()
		defer func() {
			if err := tp.Shutdown(context.Background()); err != nil {
				log.Printf("Error shutting down tracer provider: %v", err)
			}
		}()
	}

	// Basic App Configs
	body_limit, _ := strconv.Atoi(configs.AppConfig.GetOrDefault("BODY_LIMIT", "70"))
	read_buffer_size, _ := strconv.Atoi(configs.AppConfig.GetOrDefault("READ_BUFFER_SIZE", "70"))
	//load config file
	app := fiber.New(fiber.Config{
		Prefork: prefork,
		// Network:     fiber.NetworkTCP,
		// Immutable:   true,
		JSONEncoder:    json.Marshal,
		JSONDecoder:    json.Unmarshal,
		BodyLimit:      body_limit * 1024 * 1024,
		ReadBufferSize: read_buffer_size * 1024,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			// Status code defaults to 500
			code := fiber.StatusInternalServerError
			// Retrieve the custom status code if it's a *fiber.Error
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}
			// Send custom error page
			err = ctx.Status(code).SendFile(fmt.Sprintf("./%d.html", code))
			if err != nil {
				// In case the SendFile fails
				return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
			}
			// Return from handler
			return nil
		},
	})

	// recover from panic attacks middlerware
	app.Use(recover.New())

	// allow cross origin request
	app.Use(cors.New())

	//  rate limiting middleware
	rate_limit_per_second, _ := strconv.Atoi(configs.AppConfig.GetOrDefault("RATE_LIMIT_PER_SECOND", "50000"))
	app.Use(limiter.New(limiter.Config{
		Max:               rate_limit_per_second,
		Expiration:        1 * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
	}))

	// idempotency middleware
	app.Use(idempotency.New(idempotency.Config{
		Lifetime: 10 * time.Second,
			}))

	log_file, err := bluetasks.Logfile()
	if err != nil {
		fmt.Printf("Error Creating Logfile %v\n", err)
		panic(err)
	}

	// logger middle ware with the custom file writer object
	app.Use(logger.New(logger.Config{
		Format:     "\n${cyan}-[${time}]-[${ip}] -${white}${pid} ${red}${status} ${blue}[${method}] ${white}-${path}\n [${body}]\n[${error}]\n[${resBody}]\n[${reqHeaders}]\n[${queryParams}]\n",
		TimeFormat: "15:04:05",
		TimeZone:   "Local",
		Output:     log_file,
	}))

	// prometheus middleware concrete instance
	prometheus := fiberprometheus.New("gobluefiber")
	prometheus.RegisterAt(app, "/metrics")

	// prometheus monitoring middleware
	app.Use(prometheus.Middleware)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	// swagger docs
	app.Get("/docs/*", swagger.HandlerDefault)
	app.Get("/docs/*", swagger.New()).Name("swagger_routes")

	// fiber native monitoring metrics endpoint
	app.Get("/lmetrics", monitor.New(monitor.Config{Title: "goBlue Metrics Page"})).Name("custom_metrics_route")

	// Starting the APP
	go startServer(app)

	//  starting scheduler files
	sccheduledTasks := bluetasks.ScheduledTasks()

	// signal for graceful shutdown
	shutdownChan := make(chan os.Signal, 1)

	// Gracefully shuting Down
	waitForShutdown(app, sccheduledTasks, shutdownChan)
}

// waitForShutdown listens for an interrupt signal (such as SIGINT) and gracefully shuts down the Echo app.
func waitForShutdown(app *fiber.App, sccheduledTasks *tasks.Scheduler, shutdownChan chan os.Signal) {
	// Create channel to signify a signal being sent
	signal.Notify(shutdownChan, os.Interrupt) // When an interrupt or termination signal is sent, notify the channel

	<-shutdownChan // This blocks the main thread until an interrupt is received
	fmt.Println("Gracefully shutting down...")
	app.Shutdown()

	fmt.Println("Running cleanup tasks...")

	// Your cleanup tasks go here
	sccheduledTasks.Stop()
	fmt.Println("Gracefully shutting down...")
	fmt.Println("App was successfully shutdown.")
}

// start the app function
func startServer(app *fiber.App) {
	HTTP_PORT := configs.AppConfig.Get("HTTP_PORT")

	if app_tls == "on" {
		CERT_FILE := "./server.pem"
		KEY_FILE := "./server-key.pem"
		listen_port := "0.0.0.0:" + HTTP_PORT
		if err := app.ListenTLS(":"+HTTP_PORT, CERT_FILE, KEY_FILE); err != nil {
			app.ListenTLS(listen_port, CERT_FILE, KEY_FILE)
		}
	} else {
		if err := app.Listen(":" + HTTP_PORT); err != nil {
			app.Listen("0.0.0.0:" + HTTP_PORT)
		}
	}
}


func init() {
	{{.AppName}}cli.Flags().StringVar(&env, "env", "help", "Which environment to run for example prod or dev")
	{{.AppName}}cli.Flags().StringVar(&app_otel, "otel", "help", "Turn on/off OpenTelemetry tracing")
	{{.AppName}}cli.Flags().StringVar(&app_tls, "tls", "help", "Turn on/off tls, \"on\" for auto on and \"off\" for auto off")
	goFrame.AddCommand({{.AppName}}cli)

}


func SetupRoutes(app *fiber.App) {

	//app logging open telemetery
	app.Use(otelfiber.Middleware())
	app.Use(otelspanstarter)

	// database session injection to local context
	app.Use(dbsessioninjection)

	// Role Middleware
	gapp := app.Group("/api/v1", keyauth.New(keyauth.Config{
		Next:      authFilter,
		KeyLookup: "header:X-APP-TOKEN",
		Validator: NextRoute,
	}))

	{{range .Models}}
	gapp.Get("/{{.LowerName}}",NextFunc).Name("get_all_{{.LowerName}}s").Get("/{{.LowerName}}", controllers.Get{{.Name}}s)
	gapp.Get("/{{.LowerName}}/:{{.LowerName}}_id",NextFunc).Name("get_one_{{.LowerName}}s").Get("/{{.LowerName}}/:{{.LowerName}}_id", controllers.Get{{.Name}}ByID)
	gapp.Post("/{{.LowerName}}",NextFunc).Name("post_{{.LowerName}}").Post("/{{.LowerName}}", controllers.Post{{.Name}})
	gapp.Patch("/{{.LowerName}}/:{{.LowerName}}_id",NextFunc).Name("patch_{{.LowerName}}").Patch("/{{.LowerName}}/:{{.LowerName}}_id", controllers.Patch{{.Name}})
	gapp.Delete("/{{.LowerName}}/:{{.LowerName}}_id",NextFunc).Name("delete_{{.LowerName}}").Delete("/{{.LowerName}}/:{{.LowerName}}_id", controllers.Delete{{.Name}}).Name("delete_{{.LowerName}}")
	{{range .Relations}}{{if .OtM}}
	gapp.Patch("/{{.LowerFieldName}}{{.LowerParentName}}/:{{.LowerFieldName}}_id",NextFunc).Name("add_{{.LowerFieldName}}{{.LowerParentName}}").Patch("/{{.LowerFieldName}}{{.LowerParentName}}/:{{.LowerFieldName}}_id",controllers.Add{{.FieldName}}{{.ParentName}}s)
	gapp.Delete("/{{.LowerFieldName}}{{.LowerParentName}}/:{{.LowerFieldName}}_id",NextFunc).Name("delete_{{.LowerFieldName}}{{.LowerParentName}}").Delete("/{{.LowerFieldName}}{{.LowerParentName}}/:{{.LowerFieldName}}_id",controllers.Delete{{.FieldName}}{{.ParentName}}s){{end}}
	{{if .MtM}}gapp.Post("/{{.LowerFieldName}}{{.LowerParentName}}/:{{.LowerFieldName}}_id/:{{.LowerParentName}}_id",NextFunc).Name("add_{{.LowerFieldName}}{{.LowerParentName}}").Post("/{{.LowerFieldName}}{{.LowerParentName}}/:{{.LowerFieldName}}_id/:{{.LowerParentName}}_id",controllers.Add{{.FieldName}}{{.ParentName}}s)
	gapp.Delete("/{{.LowerFieldName}}{{.LowerParentName}}/:{{.LowerFieldName}}_id/:{{.LowerParentName}}_id",NextFunc).Name("delete_{{.LowerFieldName}}{{.LowerParentName}}").Delete("/{{.LowerFieldName}}{{.LowerParentName}}/:{{.LowerFieldName}}_id/:{{.LowerParentName}}_id",controllers.Delete{{.FieldName}}{{.ParentName}}s){{end}}{{end}}
	{{end}}
}

`
