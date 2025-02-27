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
			// Load data from the config file
			if err := temps.LoadData(config_file); err != nil {
				fmt.Printf("Error loading data: %v\n", err)
				return
			}

			// Check if tests flag is enabled
			if testFlag, _ := cmd.Flags().GetString("tests"); testFlag == "on" {
				generateFiberTestFiles()
			}

			// Generate main fiber framework files
			generateFiberFiles()
		},
	}
)

func generateFiberFiles() {
	// Generate the Fiber framework structure
	temps.FiberFrame()
	// running go mod tidy finally
	if err := exec.Command("swag", "init").Run(); err != nil {
		fmt.Printf("error generating swagger: %v \n", err)
	}
	temps.CommonCMD() // Common command structure
}

func generateFiberTestFiles() {
	// Generate the test structure for Fiber
	temps.TestFrameFiber()
	temps.CommonCMD() // Common command structure
}

func init() {
	// // Register flags for the fiber command
	// fibercli.Flags().StringP("tests", "t", "", "Enable test generation by specifying \"on\". Defaults to off.")
	goFrame.AddCommand(fibercli)
}
