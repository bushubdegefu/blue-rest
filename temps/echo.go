package temps

import (
	"fmt"
	"os"
	"text/template"
)

func EchoFrameSetupAndMiddleware(appName string) {
	RenderData.PackageAppName = replaceString(RenderData.AppName)
	// ############################################################
	app_middleware_tmpl, err := template.New("RenderData").Funcs(FuncMap).Parse(middlwareTemplate)
	if err != nil {
		fmt.Println("Error creating template:", err)
	}

	app_middleware_file, err := os.Create("middleware.go")
	if err != nil {
		panic(err)
	}
	defer app_middleware_file.Close()

	err = app_middleware_tmpl.Execute(app_middleware_file, RenderData)
	if err != nil {
		panic(err)
	}
	// // ############################################################
	app_setup_render, err := template.New("RenderData").Funcs(FuncMap).Parse(setupFileTemplate)
	if err != nil {
		panic(err)
	}
	app_setup_file, err := os.Create("setup.go")
	if err != nil {
		panic(err)
	}
	defer app_setup_file.Close()

	err = app_setup_render.Execute(app_setup_file, RenderData)
	if err != nil {
		panic(err)
	}

}

func EchoAppAndMiddleware() {

	RenderData.ProjectName = ProjectSettings.ProjectName
	RenderData.SetBackTick()
	RenderData.AppNames = ProjectSettings.AppNames

	// ############################################################
	app_middleware_tmpl, err := template.New("RenderData").Funcs(FuncMap).Parse(golbalMiddlwareTemplate)
	if err != nil {
		fmt.Println("Error creating template:", err)
	}

	app_middleware_file, err := os.Create("manager/middleware.go")
	if err != nil {
		panic(err)
	}
	defer app_middleware_file.Close()

	err = app_middleware_tmpl.Execute(app_middleware_file, RenderData)
	if err != nil {
		panic(err)
	}

	// // ############################################################
	app_tmpl, err := template.New("RenderData").Funcs(FuncMap).Parse(globalAppTemplate)
	if err != nil {
		panic(err)
	}
	app_setup_file, err := os.Create("manager/app.go")
	if err != nil {
		panic(err)
	}
	defer app_setup_file.Close()

	err = app_tmpl.Execute(app_setup_file, RenderData)
	if err != nil {
		panic(err)
	}

}

var middlwareTemplate = `
package {{.PackageAppName}}

import (
	"fmt"
	"strings"
	"{{.ProjectName}}/database"
	"{{.ProjectName}}/observe"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/attribute"
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
		db, err := database.ReturnSession("{{ .AppName | replaceString }}")
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
	if ctx.Path() == "/api/v1/blue_auth/login" || ctx.Path() == "/api/v1/blue_auth/stats" {
		return true, nil
	}


//  You have to fix the NextAuthValidator function, it will let all values pass
	// using required role access logic
	return true, nil
}


// AddAppTokenIfMissing is a middleware that checks if the x-app-token header is present in the request. so that the login route can work
func AddAppTokenIfMissing(next echo.HandlerFunc) echo.HandlerFunc {
	return func(contx echo.Context) error {
		// Check if x-app-token header exists
		appToken := contx.Request().Header.Get("x-app-token")

		// If the x-app-token header is missing, set a default value
		if appToken == "" {
			contx.Request().Header.Set("x-app-token", "login")
		}

		// Continue processing the request
		return next(contx)
	}
}
// Custom Middlewares can be added here specfic to the app

`

var setupFileTemplate = `
package {{.PackageAppName}}

import (
	"{{.ProjectName}}/{{.AppName}}/controllers"
	"{{.ProjectName}}/logs"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)


//  Please Note the sequence you mount the middlewares
func SetupRoutes(app *echo.Echo) {
	logOutput, _ := logs.Logfile("{{ .AppName | replaceString }}")

	// the Otel spanner middleware
	app.Use(otelechospanstarter)

	// db session injection
	app.Use(dbsessioninjection)

	app.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: {{.BackTick}}{"time":"${time_rfc3339_nano}","id":"${id}","remote_ip":"${remote_ip}",{{.BackTick}} +
			{{.BackTick}}"host":"${host}","method":"${method}","uri":"${uri}","user_agent":"${user_agent}",{{.BackTick}} +
			{{.BackTick}}"status":${status},"error":"${error}","latency":${latency},"latency_human":"${latency_human}"{{.BackTick}} +
			{{.BackTick}},"bytes_in":${bytes_in},"bytes_out":${bytes_out}}{{.BackTick}} + "\n",
		Output: logOutput,
	}))

	// then authentication middlware
	gapp := app.Group("/api/v1/{{.AppName | replaceString}}")

	gapp.POST("/login", controllers.Login).Name = "login"
	gapp.GET("/stats", controllers.DbStatEndpoint).Name = "db_stat"
	{{- range .Models}}
	gapp.GET("/{{.LowerName}}", controllers.Get{{.Name}}s).Name = "{{.AppName | replaceString}}_can_view_{{.LowerName}}"
	gapp.GET("/{{.LowerName}}/:{{.LowerName}}_id", controllers.Get{{.Name}}ByID).Name = "{{.AppName | replaceString}}_can_view_{{.LowerName}}"
	gapp.POST("/{{.LowerName}}", controllers.Post{{.Name}}).Name = "{{.AppName | replaceString}}_can_add_{{.LowerName}}"
	gapp.PATCH("/{{.LowerName}}/:{{.LowerName}}_id", controllers.Patch{{.Name}}).Name = "{{.AppName | replaceString}}_can_change_{{.LowerName}}"
	gapp.DELETE("/{{.LowerName}}/:{{.LowerName}}_id", controllers.Delete{{.Name}}).Name = "{{.AppName | replaceString}}_can_delete_{{.LowerName}}"
	{{range .Relations}}
	gapp.GET("/{{.LowerFieldName}}{{.LowerParentName}}/:{{.LowerParentName}}_id",controllers.Get{{.FieldName}}{{.ParentName}}s).Name = "{{.AppName | replaceString}}_can_view_{{.LowerFieldName}}"
	gapp.GET("/{{.LowerFieldName}}complement{{.LowerParentName}}/:{{.LowerParentName}}_id",controllers.Get{{.FieldName}}Complement{{.ParentName}}s).Name = "{{.AppName | replaceString}}_can_view_{{.LowerFieldName}}complement"
	gapp.GET("/{{.LowerFieldName}}noncomplement{{.LowerParentName}}/:{{.LowerParentName}}_id",controllers.Get{{.FieldName}}NonComplement{{.ParentName}}s).Name = "{{.AppName | replaceString}}_can_view_{{.LowerFieldName}}complement"
	gapp.POST("/{{.LowerFieldName}}{{.LowerParentName}}/:{{.LowerFieldName}}_id/:{{.LowerParentName}}_id",controllers.Add{{.FieldName}}{{.ParentName}}s).Name = "{{.AppName | replaceString}}_can_add_{{.LowerFieldName}}"
	gapp.DELETE("/{{.LowerFieldName}}{{.LowerParentName}}/:{{.LowerFieldName}}_id/:{{.LowerParentName}}_id",controllers.Delete{{.FieldName}}{{.ParentName}}s).Name = "{{.AppName | replaceString}}_can_delete_{{.LowerFieldName}}"
	{{end}}
	{{- end}}


}
`

var globalAppTemplate = `
package manager

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"

	{{- range .AppNames}}
	{{ . | replaceString }} "{{$.ProjectName}}/{{ . }}"
	{{- end }}
	"{{.ProjectName}}/configs"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/madflojo/tasks"
	"github.com/spf13/cobra"
	"golang.org/x/time/rate"

	echoSwagger "github.com/swaggo/echo-swagger"

	{{- range .AppNames}}
	{{ . | replaceString }}_tasks "{{$.ProjectName}}/{{ . }}/bluetasks"
	{{- end }}

	_ "{{.ProjectName}}/docs"
	"{{.ProjectName}}/observe"
)

var (
	env                   string
	app_otel              string
	app_tls               string
	echocli = &cobra.Command{
		Use:   "run",
		Short: "Run Development server ",
		Long:  "Run development server",
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

	// enable cross origin requests
	app.Use(middleware.CORS())


	// setup prom monitoring
	observe.SetupPrometheusMetrics(app)

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

	// Mounting Global Middleware
	MountGlobalMiddleware(app)

	app.GET("/docs/*", echoSwagger.WrapHandler)

	// Serve static files from the "dist/django_admin_ui" folder
	app.Static("/", "./dist/django_admin_ui")

	// Handle "/admin/*" route and send the "index.html" file
	app.GET("/admin/*", func(c echo.Context) error {
		return c.File("./dist/django_admin_ui/index.html")
	}).Name = "Admin_UI"


	// Setting up Endpoints
	{{- range .AppNames}}
	{{ . | replaceString }}.SetupRoutes(app)
	{{- end }}

	// building path route name path for authentication middleware
	GetApplicationRoutes(app)

	// Starting Task Scheduler ( Running task that run regularly based on the provided configs)
	{{- range .AppNames}}
	scd_tasks_{{ . | replaceString }} := {{ . | replaceString }}_tasks.ScheduledTasks()
	{{- end }}

	scd_tasks := []*tasks.Scheduler{
	{{- range .AppNames}}
		scd_tasks_{{ . | replaceString }},
	{{- end }}
		}


	// Start the server
	go startServer(app)


	// Create a context that listens for interrupt signals (e.g., Ctrl+C).
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	// Ensure the stop function is called when the function exits to clean up resources.
	defer stop()

	// Graceful shutdown
	waitForShutdown(app, scd_tasks, ctx)

}

// waitForShutdown listens for an interrupt signal (such as SIGINT) and gracefully shuts down the Echo app.
func waitForShutdown(app *echo.Echo, scheduledTasks []*tasks.Scheduler, ctx context.Context) {

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

	// Iterate through scheduledTasks and stop each one
		for _, task := range scheduledTasks {
			task.Stop()
		}

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
	echocli.Flags().StringVar(&env, "env", "help", "Which environment to run for example prod or dev")
	echocli.Flags().StringVar(&app_otel, "otel", "help", "Turn on/off OpenTelemetry tracing")
	echocli.Flags().StringVar(&app_tls, "tls", "help", "Turn on/off tls, \"on\" for auto on and \"off\" for auto off")
	goFrame.AddCommand(echocli)
}

`

var golbalMiddlwareTemplate = `
package manager

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var AppRouteNames map[string]string

func GetApplicationRoutes(app *echo.Echo) {
	// Lock the Mutex to ensure safe access to AppRouteNames

	AppRouteNames = make(map[string]string)
	for _, route := range app.Routes() {
		routeName := route.Name
		if route.Name == "" {
			// Skip routes without a name
			continue
		}
		AppRouteNames[route.Path] = routeName
	}
}

// SetRouteName header based on path
func SetRouteNameHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(contx echo.Context) error {

		routeName, exists := AppRouteNames[contx.Path()]

		// If the route name doesn't exist in the map, set it to "not-set"
		if !exists {
			routeName = "not-set"
		}

		// If the x-app-token header is missing, set a default value
		contx.Request().Header.Set("route-name", routeName)

		// Continue processing the request
		return next(contx)
	}
}

func NextAuthValidator(key string, ctx echo.Context) (bool, error) {
	// fmt.Println(ctx.Path())
	if ctx.Path() == "/api/v1/blue_auth/login" || ctx.Path() == "/api/v1/blue_auth/stats" {
		return true, nil
	}
	fmt.Println("Key: ", key)
	fmt.Println("Route Name: ", ctx.Request().Header.Get("route-name"))

	//  You have to fix the NextAuthValidator function, it will let all values pass
	// using required role access logic
	return true, nil
}

// AddAppTokenIfMissing is a middleware that checks if the x-app-token header is present in the request. so that the login route can work
func AddAppTokenIfMissing(next echo.HandlerFunc) echo.HandlerFunc {
	return func(contx echo.Context) error {
		// Check if x-app-token header exists
		appToken := contx.Request().Header.Get("x-app-token")

		// If the x-app-token header is missing, set a default value
		if appToken == "" {
			contx.Request().Header.Set("x-app-token", "login")
		}

		// Continue processing the request
		return next(contx)
	}
}

func MountGlobalMiddleware(app *echo.Echo) {
	// Mount the middleware
	app.Use(SetRouteNameHeader)
	app.Use(AddAppTokenIfMissing)
	app.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "header:x-app-token",
		Validator: NextAuthValidator,
	}))

}

`
