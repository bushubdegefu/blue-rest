package manager

import (
	"fmt"
	"os/exec"

	"github.com/bushubdegefu/blue-rest/temps"
	"github.com/spf13/cobra"
)

var (
	fibercli = &cobra.Command{
		Use:   "fiber",
		Short: "Generate basic structure files to start an app using Fiber",
		Long:  `Generate the basic structure files to start an app using Fiber.`,
		Run: func(cmd *cobra.Command, args []string) {

			temps.InitProjectJSON()
			appName, _ := cmd.Flags().GetString("app")
			globalName, _ := cmd.Flags().GetBool("global")

			if appName != "" {
				handleAppDirectory(appName)
				if err := temps.LoadData(config_file); err != nil {
					fmt.Printf("Error loading data: %v\n", err)
					return
				}
				temps.FiberFrameSetupAndMiddleware(appName)
				temps.AuthUtilsFrame(appName)
				temps.AuthLoginFrame(appName, "fiber")

			} else if globalName {
				temps.FiberAppAndMiddleware()
				generateFiberFiles()
			} else {
				fmt.Println("No app name specified")
			}
			temps.CommonCMD()
		},
	}
)

func generateFiberFiles() {
	// running go mod tidy finally
	if err := exec.Command("swag", "init").Run(); err != nil {
		fmt.Printf("error generating swagger: %v \n", err)
	}
	// Common command structure
}

func init() {
	// // Register flags for the fiber command
	// fibercli.Flags().StringP("tests", "t", "", "Enable test generation by specifying \"on\". Defaults to off.")
	fibercli.Flags().StringP("app", "a", "", "Specify the app name, so that echo app will be generated")
	fibercli.Flags().BoolP("global", "g", false, "basic echo app with for global, creates app.go( in manager package) and middleware.go on the main module takes true or false")
	fibercli.Flags().StringVarP(&config_file, "config", "c", "config.json", "Specify the data file to load")
	goFrame.AddCommand(fibercli)
}
