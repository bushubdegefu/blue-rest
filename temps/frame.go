package temps

import (
	"fmt"
	"os"
	"text/template"
)

func Frame() {

	//  this is creating manger file inside the manager folder
	// ############################################################
	manager_tmpl, err := template.New("RenderData").Parse(managerTemplate)
	if err != nil {
		fmt.Printf("Frame - 1: %v\n", err)
	}

	// Create the models directory if it does not exist
	err = os.MkdirAll("manager", os.ModePerm)
	if err != nil {
		fmt.Printf("Frame - 2: %v\n", err)
	}

	manager_file, err := os.Create("manager/manager.go")
	if err != nil {
		fmt.Printf("Frame - 3: %v\n", err)
	}
	defer manager_file.Close()

	err = manager_tmpl.Execute(manager_file, RenderData)
	if err != nil {
		fmt.Printf("Frame - 4: %v\n", err)
	}

	//this is creating the main.go file
	// ############################################################
	main_tmpl, err := template.New("RenderData").Parse(mainTemplate)
	if err != nil {
		fmt.Printf("Frame - 5: %v\n", err)
	}

	// Create the models directory if it does not exist
	err = os.MkdirAll("manager", os.ModePerm)
	if err != nil {
		fmt.Printf("Frame - 6: %v\n", err)
	}

	main_file, err := os.Create("main.go")
	if err != nil {
		fmt.Printf("Frame - 7: %v\n", err)
	}
	defer main_file.Close()

	err = main_tmpl.Execute(main_file, RenderData)
	if err != nil {
		fmt.Printf("Frame - 8: %v\n", err)
	}

	//  this is creating .env and configuration file
	// ############################################################
	config_tmpl, err := template.New("RenderData").Parse(configsTemplate)
	if err != nil {
		fmt.Printf("Frame - 9: %v\n", err)
	}

	// Create the models directory if it does not exist
	err = os.MkdirAll("configs", os.ModePerm)
	if err != nil {
		fmt.Printf("Frame - 10: %v\n", err)
	}
	err = os.MkdirAll("tests", os.ModePerm)
	if err != nil {
		fmt.Printf("Frame - 11: %v\n", err)
	}

	config_file, err := os.Create("configs/configs.go")
	if err != nil {
		fmt.Printf("Frame - 12: %v\n", err)
	}
	defer config_file.Close()

	err = config_tmpl.Execute(config_file, RenderData)
	if err != nil {
		fmt.Printf("Frame - 13: %v\n", err)
	}

	//  this is creating .env and configuration file
	// ############################################################
	env_tmpl, err := template.New("RenderData").Parse(envTemplate)
	if err != nil {
		fmt.Printf("Frame - 14: %v\n", err)
	}

	env_file, err := os.Create("configs/.env")
	if err != nil {
		fmt.Printf("Frame - 15: %v\n", err)
	}
	defer env_file.Close()

	tenv_file, err := os.Create("tests/.env")
	if err != nil {
		fmt.Printf("Frame - 16: %v\n", err)
	}
	defer tenv_file.Close()

	err = env_tmpl.Execute(tenv_file, RenderData)
	if err != nil {
		fmt.Printf("Frame - 17: %v\n", err)
	}

	err = env_tmpl.Execute(env_file, RenderData)
	if err != nil {
		fmt.Printf("Frame - 18: %v\n", err)
	}

	//  this is creating .env and configuration file
	// ############################################################
	devenv_tmpl, err := template.New("RenderData").Parse(devenvTemplate)
	if err != nil {
		fmt.Printf("Frame - 19: %v\n", err)
	}

	devenv_file, err := os.Create("configs/.dev.env")
	if err != nil {
		fmt.Printf("Frame - 20: %v\n", err)
	}
	defer devenv_file.Close()

	err = devenv_tmpl.Execute(devenv_file, RenderData)
	if err != nil {
		fmt.Printf("Frame - 21: %v\n", err)
	}

	// ############################################################
	prodenv_tmpl, err := template.New("RenderData").Parse(devenvTemplate)
	if err != nil {
		fmt.Printf("Frame - 22: %v\n", err)
	}

	prodenv_file, err := os.Create("configs/.prod.env")
	if err != nil {
		fmt.Printf("Frame - 23: %v\n", err)
	}
	defer prodenv_file.Close()

	err = prodenv_tmpl.Execute(prodenv_file, RenderData)
	if err != nil {
		fmt.Printf("Frame - 24: %v\n", err)
	}

	// ##########################################################
	testenv_tmpl, err := template.New("RenderData").Parse(devenvTemplate)
	if err != nil {
		fmt.Printf("Frame - 25: %v\n", err)
	}

	testenv_file, err := os.Create("tests/.test.env")
	if err != nil {
		fmt.Printf("Frame - 26: %v\n", err)
	}
	defer testenv_file.Close()

	err = testenv_tmpl.Execute(testenv_file, RenderData)
	if err != nil {
		fmt.Printf("Frame - 27: %v\n", err)
	}

}

var mainTemplate = `
package main

import (
	"{{.ProjectName}}/manager"
)

//	@title			Swagger {{.ProjectName}} API
//	@version		0.1
//	@description	This is {{.ProjectName}} API OPENAPI Documentation.
//	@termsOfService	http://swagger.io/terms/
//  @BasePath  /api/v1

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						X-APP-TOKEN
//	@description				Description for what is this security definition being used

//	@securityDefinitions.apikey Refresh
//	@in							header
//	@name						X-REFRESH-TOKEN
//	@description				Description for what is this security definition being used

func main() {
	manager.Execute()
}
`

var managerTemplate = `
package manager

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	goFrame = &cobra.Command{
		Use:           "AppMan",
		Short:         "AppMan – command-line tool to aid structure you fiber backend projects with gorm",
		Long:          "Manager File Framed by go frame",
		Version:       "0.0.0",
		SilenceErrors: true,
		SilenceUsage:  true,
	}
)

func Execute() {
	if err := goFrame.Execute(); err != nil {

		fmt.Println(err)
		os.Exit(1)
	}
}
`

var configsTemplate = `
package configs

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

const (
	defaultFileName         = "/.env"
	defaultOverrideFileName = "/.local.env"
)

type EnvConfig struct {
	defaultPath string
	prodFlag    string
}

type Config interface {
	Get(string) string
	GetOrDefault(string, string) string
}

var AppConfig EnvConfig

func NewEnvFile(configFolder string) {
	AppConfig = EnvConfig{
		defaultPath: configFolder,
	}
	AppConfig.read()
}

func (e *EnvConfig) read() {
	var (
		defaultFile  = e.defaultPath + defaultFileName
		overrideFile = e.defaultPath + defaultOverrideFileName
	)
	err := godotenv.Overload(defaultFile)
	env := e.Get("APP_ENV")
	if err != nil {
		fmt.Printf("WARNING: Failed to load config from file: %v, Err: %v \n", defaultFile, err)
	} else {
		fmt.Printf("INFO: Loaded config from file: %v\n", defaultFile)
	}

	// If 'APP_ENV' is set to x, then App will read '.env' from configs directory, and then it will be overwritten
	// by configs present in file '.x.env'
	overrideFile = fmt.Sprintf("%s/.%s.env", e.defaultPath, env)
	if env == "" && e.prodFlag == "" {
		overrideFile = fmt.Sprintf("%s/dev.env", e.defaultPath)
	}

	if e.prodFlag != "" {
		overrideFile = fmt.Sprintf("%s/.%s.env", e.defaultPath, e.prodFlag)
	}
	err = godotenv.Overload(overrideFile)
	if err != nil {
		fmt.Printf("WARNING: to load config from file: %v, Err: %v \n", overrideFile, err)
	} else {
		fmt.Printf("INFO: config from file: %v \n", overrideFile)
	}
}

func (e *EnvConfig) Get(key string) string {
	return os.Getenv(key)
}

func (e *EnvConfig) SetEnv(key string) {
	AppConfig = EnvConfig{
		defaultPath: "./configs",
		prodFlag: key,
	}
	AppConfig.read()
}

func (e *EnvConfig) GetOrDefault(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return defaultValue
}
`

var envTemplate = `
APP_ENV=dev
`
var devenvTemplate = `
APP_NAME=dev
HTTP_PORT=7500
TEST_NAME="Development Development"
BODY_LIMIT=70
READ_BUFFER_SIZE=40
RATE_LIMIT_PER_SECOND=5000

#Interval in minutes
CLEAR_LOGS_INTERVAL=1

#Database config settings
#DB_TYPE=postgres
#POSTGRES_URI="host=localhost user=blueuser password=default dbname=learning_one port=5432 sslmode=disable"
DB_TYPE="sqlite"
SQLLITE_URI="blue-rest.db"
#DB_TYPE="mysql"
#MYSQL_URI="yenefivy_beimnet:bluenet%402025@tcp(109.70.148.37:3306)/gorm?charset=utf8&parseTime=True&loc=Local"

#Messeage qeue settings specifically rabbit
RABBIT_URI="amqps://xrqlluoo:4hAUYGqztMsWyFdT5r65j4xudTw-AWl1@puffin.rmq2.cloudamqp.com/xrqlluoo"

# JWT token settings
JWT_SALT_LIFE_TIME=60 #in minutes
JWT_SALT_LENGTH=25

#RPC settings
RPC_PORT=6500


#Observability settings
TRACE_EXPORTER=jaeger
TRACER_HOST=localhost
TRACER_PORT=14317

`
