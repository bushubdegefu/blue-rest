package manager

import (
	"fmt"

	"github.com/bushubdegefu/blue-rest/temps"
	"github.com/spf13/cobra"
)

var (
	testscli = &cobra.Command{
		Use:   "test",
		Short: "Generate basic coverage test code for Fiber/Echo for the generated crud endpoints.",
		Long:  `Generate basic coverage test code for Fiber/Echo for the generated crud endpoints.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Load data from the config file
			if err := temps.LoadData(config_file); err != nil {
				fmt.Printf("Error loading data: %v\n", err)
				return
			}

			frame, _ := cmd.Flags().GetString("frame")
			if frame == "" {
				fmt.Println("Error: --frame flag is required. Use --frame=echo or --frame=fiber.")
				return
			}

			if frame == "echo" || frame == "fiber" {
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
		temps.TestFrameFiber()
	} else if frame == "echo" {
		temps.TestFrameEcho()
	} else {
		fmt.Println("Error: Invalid frame value. Use --frame=echo or --frame=fiber.")
	}
	temps.CommonCMD() // Common command structure
}

func init() {
	// // Register flags for the fiber command
	testscli.Flags().StringP("frame", "f", "", "Specify the framework to use (echo or fiber) for the tests")
	goFrame.AddCommand(testscli)
}
