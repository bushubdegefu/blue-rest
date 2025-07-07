package manager

import (
	"fmt"
	"os/exec"

	"github.com/bushubdegefu/blue-rest/temps"
	"github.com/bushubdegefu/blue-rest/temps/generator"
	"github.com/spf13/cobra"
)

var (
	fibercli = &cobra.Command{
		Use:   "fiber",
		Short: "Generate basic structure files to start an app using Fiber",
		Long:  `Generate the basic structure files to start an app using Fiber.`,
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
				}
				generator.GenerateUtilsApp(temps.ProjectSettings)
				loginFrame(appName, "fiber")

			} else if globalName {
				generator.GenerateGlobalFiberAppMiddleware(temps.RenderData)
				generator.GenerateAppFiberGlobal(temps.RenderData)
				runSwagInitForApps()

			} else {
				fmt.Println("No app name specified")
			}
			temps.CommonCMD()
		},
	}
)

func runSwagInitForApps() {
	temps.InitProjectJSON()
	// swag init --generalInfo setup.go --output  blue-auth/docs --dir=blue-auth,common
	for _, appName := range temps.ProjectSettings.AppNames {
		// Construct paths for generalInfo, output, and dir
		// generalInfo := filepath.Join(appName, "setup.go")
		// outputDir := filepath.Join(appName, "docs")
		dirArg := fmt.Sprintf("%s,common", appName)
		outputDir := fmt.Sprintf("%s/docs", appName)

		// Prepare the swag init command
		cmd := exec.Command(
			"swag", "init",
			"--generalInfo", "setup.go",
			"--output", outputDir,
			"--dir", dirArg,
		)

		// Run the command and handle errors
		if err := cmd.Run(); err != nil {
			fmt.Printf("Error generating swagger for app '%s': %v\n", appName, err)
		} else {
			fmt.Printf("Swagger generated for app '%s'\n", appName)
		}
	}
}

func init() {
	// // Register flags for the fiber command
	// fibercli.Flags().StringP("tests", "t", "", "Enable test generation by specifying \"on\". Defaults to off.")
	fibercli.Flags().StringP("app", "a", "", "Specify the app name, so that echo app will be generated")
	fibercli.Flags().BoolP("global", "g", false, "basic echo app with for global, creates app.go( in manager package) and middleware.go on the main module takes true or false")
	fibercli.Flags().StringVarP(&config_file, "config", "c", "config.json", "Specify the data file to load")
	goFrame.AddCommand(fibercli)
}
