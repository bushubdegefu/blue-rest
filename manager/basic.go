package manager

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bushubdegefu/blue-rest/temps"
	"github.com/bushubdegefu/blue-rest/temps/generator"
	"github.com/spf13/cobra"
)

var (
	initalizemodule = &cobra.Command{
		Use:   "init",
		Short: "Initialize the module with name",
		Long:  `Provide name to initialize the project using the "name" flag.`,
		Run: func(cmd *cobra.Command, args []string) {
			moduleName, _ := cmd.Flags().GetString("name")
			appName, _ := cmd.Flags().GetString("app")
			authAppName, _ := cmd.Flags().GetString("auth")
			authAppType, _ := cmd.Flags().GetString("type")

			if appName == "" && moduleName == "" {

				fmt.Println("Please provide app name with app name flag or module name with  name flag")

			} else {
				// Initialize the module
				if moduleName != "" {
					temps.CommonProjectName(moduleName, authAppName, authAppType)
					temps.CommonModInit(moduleName)
					// temps.CommonCMDInit()
				}

				// If no module name, fetch the project name
				if moduleName == "" {
					temps.InitProjectJSON()

				}

				// Get current working directory
				currentDir, _ := os.Getwd()
				temps.RenderData.ProjectName = moduleName
				generator.GenerateMainAndManager(temps.RenderData)
				generator.GenerateConfig(temps.RenderData)

				// Handle appName if provided
				if appName != "" {
					handleAppInitialization(appName, currentDir, authAppName)
				}
			}
			temps.CommonCMD()
		},
	}
	configcli = &cobra.Command{
		Use:   "config",
		Short: "Template Configuration Variables need for the apps registerd to run",
		Long:  `Template Configuration Variables need for the apps registerd to run.`,
		Run: func(cmd *cobra.Command, args []string) {
			temps.InitProjectJSON()
			temps.RenderData.ProjectName = temps.ProjectSettings.ProjectName
			temps.RenderData.AppNames = temps.ProjectSettings.AppNames
			generator.GenerateConfig(temps.RenderData)
			generator.GenerateConfigEnv(temps.RenderData)
			generator.GenerateConfigAppEnv(temps.RenderData)

		},
	}

	basicCommand = &cobra.Command{
		Use:   "basic",
		Short: "Generate a basic folder structure for a project",
		Long:  `This command generates a basic folder structure for a project. The type flag determines the specific setup.`,
		Run: func(cmd *cobra.Command, args []string) {
			projectType, _ := cmd.Flags().GetString("type")
			appName, _ := cmd.Flags().GetString("app")
			frame, _ := cmd.Flags().GetString("frame")

			temps.InitProjectJSON()

			// Handle appName if provided
			if appName != "" {
				handleAppDirectory(appName)
				if err := temps.LoadData(config_file); err != nil {
					fmt.Printf("Error loading data: %v\n", err)
					return
				}
				temps.RenderData.ProjectName = temps.ProjectSettings.ProjectName
				temps.RenderData.AppName = appName
			}

			// Generate structure based on project type
			handleProjectType(projectType, frame, cmd)
		},
	}
)

func handleAppInitialization(appName, currentDir, authAppName string) {

	temps.RenderData.AppName = appName
	temps.ProjectSettings.AppendAppName(appName, authAppName)
	// Create app directory and switch to it
	os.Mkdir(appName, os.ModePerm)
	newDir := filepath.Join(currentDir, appName)
	_ = os.Chdir(newDir)
	if temps.ProjectSettings.AuthAppType == "sso" {
		// Generate the SSO schema app
		generator.GenerateSSOAuth(temps.RenderData)
	} else {
		// Generate the Django auth Schema app
		generator.GenerateDjangoAuth(temps.RenderData)
	}
}

func handleAppDirectory(appName string) {
	currentDir, _ := os.Getwd()
	newDir := filepath.Join(currentDir, appName)
	_ = os.Chdir(newDir)
}

func handleProjectType(projectType, frame string, cmd *cobra.Command) {
	switch projectType {
	case "git":
		appName, _ := cmd.Flags().GetString("app")
		if appName == "" {
			basiccmd()
		} else {
			fmt.Println("git does not need appName flag")
		}
	case "api":
		appName, _ := cmd.Flags().GetString("app")
		if appName == "" {
			fmt.Println("api flag need additional flag app with ")
		} else {
			temps.ClientAuthAndIndexJSFrame()
			temps.ClientJSFrame()
		}
	case "oauth":
		appName, _ := cmd.Flags().GetString("app")
		if appName == "" {
			fmt.Println("api flag need additional flag app with ")
		} else {
			generator.GenerateHelperOauth(temps.RenderData)
			generator.GenerateGoogleOauth(temps.RenderData)
			generator.GenerateMicrosoftOauth(temps.RenderData)
			temps.ClientJSFrame()
		}
	case "otel":
		appName, _ := cmd.Flags().GetString("app")
		if appName == "" {
			handleOtelFrame(frame)
		} else {
			fmt.Println("Otel does not need app flag.")
			return
		}
		temps.CommonCMD()
	case "rsa":
		appName, _ := cmd.Flags().GetString("app")
		if appName == "" {
			rsa_basic()
		} else {
			fmt.Println("Does not require app flag")
		}
	case "logs":
		appName, _ := cmd.Flags().GetString("app")
		if appName == "" {
			temps.InitProjectJSON()
			generator.GenerateLogs(temps.RenderData)
		} else {
			fmt.Println("Does not require app flag")
		}
	case "db":
		appName, _ := cmd.Flags().GetString("app")
		if appName == "" {
			standarddatabase()
			temps.CommonCMD()
		}
	case "consumer":
		standardrabbit()
		handleFrame(frame)
		temps.GenericTracerTemplate()
		temps.CommonCMD()
	case "producer":
		standpublish()
		handleFrame(frame)
		temps.GenericTracerTemplate()
		temps.CommonCMD()
	case "tasks":
		appName, _ := cmd.Flags().GetString("app")
		if appName == "" {
			fmt.Println("tasks flag need additional flag app")
		} else {
			generator.GenerateTasks(temps.RenderData)

			temps.CommonCMD()
		}
	case "pagination":
		appName, _ := cmd.Flags().GetString("app")
		if appName != "" {
			fmt.Println("pagination type does not need app flag")
		} else {
			generator.GenerateCommon(temps.RenderData)
			temps.CommonCMD()
		}
	case "migration":
		appName, _ := cmd.Flags().GetString("app")
		if appName != "" {
			fmt.Println("migration type does not need app flag")
		} else {
			temps.InitProjectJSON()
			temps.RenderData.AuthAppName = temps.ProjectSettings.AuthAppName
			temps.RenderData.AppName = temps.ProjectSettings.AuthAppName
			temps.RenderData.AppNames = temps.ProjectSettings.AppNames
			generator.GenerateAppDatabaseMigration(temps.RenderData)
			temps.CommonCMD()
		}
	default:
		fmt.Println("Unknown type specified. Valid types are: rsa, db, rabbit, tasks, pagination, git.")
	}
}

func handleOtelFrame(frame string) {
	if frame == "echo" || frame == "fiber" {
		temps.StandardTracerFrame(frame)
		temps.PrometheusTracerFrame(frame)
	} else {
		fmt.Println("Unknown frame specified. Valid frames are: echo, fiber.")
	}
}

func handleFrame(frame string) {
	if frame == "echo" || frame == "fiber" {
		temps.StandardTracerFrame(frame)
		temps.PrometheusTracerFrame(frame)
	} else {
		temps.StandardTracerFrame(frame)
	}
}

func basiccmd() {
	temps.GitDockerFrame()
	temps.HaproxyFrame()
}

func rsa_basic() {
	temps.RSAHelper()
}

func standardrabbit() {
	temps.CommonRabbitFrame()
	temps.RabbitFrame()
	temps.PublishFrame()
	temps.ConsumeFrame()
	temps.RunConsumeFrame()
}

func standpublish() {
	temps.CommonRabbitFrame()
	temps.RabbitFrame()
	temps.PublishFrame()
}

func standarddatabase() {
	temps.InitProjectJSON()
	generator.GenerateDBConn(temps.ProjectSettings)
}

func init() {
	// Register flags for all commands
	initalizemodule.Flags().StringP("name", "n", "", "Specify the module name  (github.com/someuser/someproject)")
	initalizemodule.Flags().StringP("app", "a", "", "Specify the application name  like auth-app,hrm-app")
	initalizemodule.Flags().StringP("auth", "p", "", "Specify the authentication application name  defaults to django_auth")
	initalizemodule.Flags().StringP("type", "t", "", "specify if you are using standalone authentication like django admin or sso like solution")

	// Register flags for the 'basic' command
	basicCommand.Flags().StringP("type", "t", "", "Specify the type of folder structure to generate: rsa, db, producer,logs, consumer, tasks, pagination, otel,migration")
	basicCommand.Flags().StringP("frame", "f", "", "Specify the Spanner function you want for the tracer, echo/fiber, meant to be used with otel flag")
	basicCommand.Flags().StringP("name", "n", "", "Specify the project module name as in github.com/someuser/someproject for the json template generation")
	basicCommand.Flags().StringP("app", "a", "", "Specify the app name, all it will try to generate for all jsons")

	goFrame.AddCommand(basicCommand)
	goFrame.AddCommand(configcli)
	goFrame.AddCommand(initalizemodule)
}
