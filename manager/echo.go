package manager

import (
	"fmt"

	"github.com/bushubdegefu/blue-rest/temps"
	"github.com/bushubdegefu/blue-rest/temps/generator"
	"github.com/spf13/cobra"
)

var (
	echocli = &cobra.Command{
		Use:   "echo",
		Short: "generate the basic structure file to start app using echo",
		Long:  `generate the basic structure file to start app using echo`,
		Run: func(cmd *cobra.Command, args []string) {

			// Initialize the project settings
			temps.InitProjectJSON()
			temps.RenderData.ProjectName = temps.ProjectSettings.ProjectName
			temps.RenderData.AppNames = temps.ProjectSettings.AppNames
			temps.RenderData.AuthAppType = temps.ProjectSettings.AuthAppType

			appName, _ := cmd.Flags().GetString("app")
			globalName, _ := cmd.Flags().GetBool("global")

			if appName != "" {
				handleAppDirectory(appName)
				if err := temps.LoadData(config_file); err != nil {
					fmt.Printf("Error loading data: %v\n", err)
					return
				}
				generator.GenerateFiberAppMiddleware(temps.RenderData)
				generator.GenerateFiberSetup(temps.RenderData)
				if appName == temps.ProjectSettings.AuthAppName {
					temps.ProjectSettings.CurrentAppName = appName
					generator.GenerateJWTUtils(temps.ProjectSettings)
					generator.GenerateUtilsApp(temps.ProjectSettings)
					loginFrame(appName, "echo")
				}

			} else if globalName {
				generator.GenerateGlobalEchoAppMiddleware(temps.RenderData)
				generator.GenerateAppEchoGlobal(temps.RenderData)
				runSwagInitForApps()
			} else {
				fmt.Println("No app name specified")
			}
			temps.CommonCMD()
		},
	}
)

func loginFrame(appName, frame string) {
	temps.ProjectSettings.CurrentAppName = appName
	generator.GenerateJWTUtils(temps.ProjectSettings)
	generator.GenerateUtilsApp(temps.ProjectSettings)
	temps.ProjectSettings.Models = temps.RenderData.Models
	// Generate login frame
	if frame == "fiber" {
		generator.GenerateFiberLogin(temps.ProjectSettings)
	} else {
		generator.GenerateEchoLogin(temps.ProjectSettings)

	}
}

func init() {
	echocli.Flags().StringP("app", "a", "", "Specify the app name, so that echo app will be generated")
	echocli.Flags().BoolP("global", "g", false, "basic echo app with for global, creates app.go( in manager package) and middleware.go on the main module takes true or false")
	echocli.Flags().StringVarP(&config_file, "config", "c", "config.json", "Specify the data file to load")
	goFrame.AddCommand(echocli)
}
