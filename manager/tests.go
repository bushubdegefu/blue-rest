package manager

import (
	"fmt"

	"github.com/bushubdegefu/blue-rest/temps"
	"github.com/bushubdegefu/blue-rest/temps/generator"
	"github.com/spf13/cobra"
)

var (
	testscli = &cobra.Command{
		Use:   "test",
		Short: "Generate basic coverage test code for Fiber/Echo for the generated crud endpoints.",
		Long:  `Generate basic coverage test code for Fiber/Echo for the generated crud endpoints.`,
		Run: func(cmd *cobra.Command, args []string) {

			frame, _ := cmd.Flags().GetString("frame")
			appName, _ := cmd.Flags().GetString("app")
			if frame == "" || appName == "" {
				fmt.Println("Error: --frame and --app flags are required. Use --frame=echo or --frame=fiber and --app=app_name.")
				return
			}

			handleAppDirectoryAndLoadConfig(appName)
			if frame == "echo" || frame == "fiber" {

				generator.GenerateConfigTestEnv(temps.RenderData)
				generateTests(frame)
			} else {
				fmt.Println("Error: Invalid frame value. Use --frame=echo or --frame=fiber.")
			}

		},
	}
)

func generateTests(frame string) {
	// Generate the test structure for Fiber

	if frame == "fiber" {
		// temps.TestFrameFiber()
		generator.GenerateFiberCoverage(temps.RenderData)
	} else if frame == "echo" {
		// temps.TestFrameEcho()
		generator.GenerateEchoCoverage(temps.RenderData)
	} else {
		fmt.Println("Error: Invalid frame value. Use --frame=echo or --frame=fiber.")
	}
	temps.CommonCMD() // Common command structure
}

func init() {
	// // Register flags for the fiber command
	testscli.Flags().StringP("frame", "f", "", "Specify the framework to use (echo or fiber) for the tests")
	testscli.Flags().StringP("app", "a", "", "Specify the application name")
	goFrame.AddCommand(testscli)
}
