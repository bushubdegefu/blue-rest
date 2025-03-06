package temps

import (
	"os"
	"text/template"
)

func EchoFrame() {
	//  this is creating manger file inside the manager folder
	// ############################################################
	echo_tmpl, err := template.New("RenderData").Parse(devechoTemplate)
	if err != nil {
		panic(err)
	}

	// Create the models directory if it does not exist
	err = os.MkdirAll("manager", os.ModePerm)
	if err != nil {
		panic(err)
	}

	devecho_file, err := os.Create("manager/devecho.go")
	if err != nil {
		panic(err)
	}
	defer devecho_file.Close()

	err = echo_tmpl.Execute(devecho_file, RenderData)
	if err != nil {
		panic(err)
	}

	// ##########################################

}

var devechoTemplate = `
package manager

import (
	"time"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"


	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/madflojo/tasks"
	"github.com/labstack/gommon/log"
	"go.opentelemetry.io/otel/attribute"
	"{{.ProjectName}}/configs"
	"{{.ProjectName}}/controllers"
	"github.com/spf13/cobra"
	"golang.org/x/time/rate"

	echoSwagger "github.com/swaggo/echo-swagger"

	"{{.ProjectName}}/database"
	_ "{{.ProjectName}}/docs"
	"{{.ProjectName}}/observe"
	"{{.ProjectName}}/bluetasks"
)

var (
	env string
	app_otel    string
	app_tls     string
	{{.AppName}}devechocli= &cobra.Command{
		Use:   "run",
		Short: "Run Development server ",
		Long:  {{.BackTick}}Run Gofr development server{{.BackTick}},
		Run: func(cmd *cobra.Command, args []string) {
		switch env {
		case "":
			echo_run("dev")
		default:
			echo_run(env)
		}
		},
	}
)

func otelechospanstarter(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		routeName := ctx.Path() + "_" + strings.ToLower(ctx.Request().Method)
		tracer, span := observe.EchoAppSpanner(ctx, fmt.Sprintf("%v-root", routeName))
		ctx.Set("tracer", &observe.RouteTracer{Tracer: tracer, Span: span})

		// Process request
		err := next(ctx)
		if err != nil {
			return err
		}

		span.SetAttributes(attribute.String("response", fmt.Sprintf("%v", ctx.Response().Status)))
		span.End()
		return nil
	}
}

func dbsessioninjection(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		db, err := database.ReturnSession()
		if err != nil {
			return err
		}

		ctx.Set("db", db)
		nerr := next(ctx)
		if nerr != nil {
			return nerr
		}
		return nil
	}
}

func NextAuthValidator(key string, ctx echo.Context) (bool, error) {
	//  You have to fix the NextAuthValidator function, it will let all values pass
	// using required role access logic
	return true, nil
}

func echo_run(env string) {
	//  loading dev env file first
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


	// starting the app
	app := echo.New()

	//  prometheus metrics middleware
	app.Use(echoprometheus.NewMiddleware("echo_blue"))

	logOutput, _ := bluetasks.Logfile()

	app.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: {{.BackTick}}{"time":"${time_rfc3339_nano}","id":"${id}","remote_ip":"${remote_ip}",{{.BackTick}} +
			{{.BackTick}}"host":"${host}","method":"${method}","uri":"${uri}","user_agent":"${user_agent}",{{.BackTick}} +
			{{.BackTick}}"status":${status},"error":"${error}","latency":${latency},"latency_human":"${latency_human}"{{.BackTick}} +
			{{.BackTick}},"bytes_in":${bytes_in},"bytes_out":${bytes_out}}{{.BackTick}} + "\n",
		Output: logOutput,
	}))


	// Middleware stack
	configLimit, _ := strconv.ParseFloat(configs.AppConfig.GetOrDefault("RATE_LIMIT_PER_SECOND", "50000"), 64)
	rateLimit := rate.Limit(configLimit)

	// Rate Limiting to throttle overload
	app.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rateLimit)))

	// Recover incase of panic attacks
	app.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10, // 1 KB
		LogLevel:  log.ERROR,
	}))

	app.GET("/docs/*", echoSwagger.WrapHandler)

	// Setting up Endpoints
	SetupRoutes(app, false)

	// Start the server
	go startServer(app)


	// Starting Task Scheduler ( Running task that run regularly based on the provided configs)
	scd_tasks := bluetasks.ScheduledTasks()

	// Create a context that listens for interrupt signals (e.g., Ctrl+C).
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	// Ensure the stop function is called when the function exits to clean up resources.
	defer stop()

	// Graceful shutdown
	waitForShutdown(app, scd_tasks, ctx)

}

// waitForShutdown listens for an interrupt signal (such as SIGINT) and gracefully shuts down the Echo app.
func waitForShutdown(app *echo.Echo, sccheduledTasks *tasks.Scheduler, ctx context.Context) {

	// Block and wait for an interrupt signal (this will block until the signal is received).
	<-ctx.Done()
	fmt.Println("Gracefully shutting down...")

	// Once the interrupt signal is received, create a new context with a 10-second timeout.
	// This will allow time for any active requests to complete before forcing shutdown.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel() // Ensure the cancel function is called when the context is no longer needed.

	// Attempt to gracefully shut down the Echo server.
	// If an error occurs during the shutdown process, log the fatal error.
	if err := app.Shutdown(ctx); err != nil {
		fmt.Println(err)
	}

	// Stop scheduler after shutdown.
	sccheduledTasks.Stop()

	// Log a message indicating the server is being shut down gracefully.
	fmt.Println("Gracefully shutting down...")
}

func startServer(app *echo.Echo) {
	HTTP_PORT := configs.AppConfig.Get("HTTP_PORT")
	if app_tls == "on" {
		CERT_FILE := "./server.pem"
		KEY_FILE := "./server-key.pem"
		app.Logger.Fatal(app.StartTLS("0.0.0.0:"+HTTP_PORT, CERT_FILE, KEY_FILE))
	} else {
		app.Logger.Fatal(app.Start("0.0.0.0:" + HTTP_PORT))
	}
}


func init() {
	{{.AppName}}devechocli.Flags().StringVar(&env, "env", "help", "Which environment to run for example prod or dev")
	{{.AppName}}devechocli.Flags().StringVar(&app_otel, "otel", "help", "Turn on/off OpenTelemetry tracing")
	{{.AppName}}devechocli.Flags().StringVar(&app_tls, "tls", "help", "Turn on/off tls, \"on\" for auto on and \"off\" for auto off")
	goFrame.AddCommand({{.AppName}}devechocli)
}


func SetupRoutes(app *echo.Echo, test bool) {
	// the Otel spanner middleware
	app.Use(otelechospanstarter)

	// db session injection
	app.Use(dbsessioninjection)

	gapp := app.Group("/api/v1")

	if !test {
		// Authentication middleware
		gapp.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
			KeyLookup: "header:x-app-token",
			Validator: NextAuthValidator,
		}))
	}

	{{range .Models}}
	gapp.GET("/{{.LowerName}}", controllers.Get{{.Name}}s).Name = "get_all_{{.LowerName}}s"
	gapp.GET("/{{.LowerName}}/:{{.LowerName}}_id", controllers.Get{{.Name}}ByID).Name = "get_one_{{.LowerName}}s"
	gapp.POST("/{{.LowerName}}", controllers.Post{{.Name}}).Name = "post_{{.LowerName}}"
	gapp.PATCH("/{{.LowerName}}/:{{.LowerName}}_id", controllers.Patch{{.Name}}).Name = "patch_{{.LowerName}}"
	gapp.DELETE("/{{.LowerName}}/:{{.LowerName}}_id", controllers.Delete{{.Name}}).Name = "delete_{{.LowerName}}"
	{{range .Relations}}
	gapp.POST("/{{.LowerFieldName}}{{.LowerParentName}}/{{ "{" }}{{.LowerFieldName}}_id{{ "}" }}/{{ "{" }}{{.LowerParentName}}_id{{ "}" }}",controllers.Add{{.FieldName}}{{.ParentName}}s).Name = "add_{{.LowerFieldName}}{{.LowerParentName}}"
	gapp.DELETE("/{{.LowerFieldName}}{{.LowerParentName}}/{{ "{" }}{{.LowerFieldName}}_id{{ "}" }}/{{ "{" }}{{.LowerParentName}}_id{{ "}" }}",controllers.Delete{{.FieldName}}{{.ParentName}}s).Name = "delete_{{.LowerFieldName}}{{.LowerParentName}}"
	{{end}}
	{{end}}


}

`
