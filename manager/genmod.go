package manager

import (
	"fmt"

	"github.com/bushubdegefu/blue-rest/temps"
	"github.com/spf13/cobra"
)

var (
	config_file string
	curdcli     = &cobra.Command{
		Use:   "crud",
		Short: "Generate CRUD handlers based on GORM for the specified framework",
		Long:  `This command generates only the CRUD handlers using GORM, based on the provided spec, for either the Echo or Fiber framework.`,
		Run: func(cmd *cobra.Command, args []string) {
			frame, _ := cmd.Flags().GetString("frame")
			if frame == "" {
				fmt.Println("Error: --frame flag is required. Use --frame=echo or --frame=fiber.")
				return
			}

			// Load data with the specified filename
			if err := temps.LoadData(config_file); err != nil {
				fmt.Printf("Error loading data: %v\n", err)
				return
			}

			if frame == "echo" || frame == "fiber" {
				gengorm(frame)
			} else {
				fmt.Println("Error: Invalid frame value. Use --frame=echo or --frame=fiber.")
			}
		},
	}

	modelscli = &cobra.Command{
		Use:   "models",
		Short: "Generate data models based on GORM with annotations",
		Long:  `This command generates data models using GORM, based on the provided spec on the config.json file along with gorm relationship annotations.`,
		Run: func(cmd *cobra.Command, args []string) {

			// Load data with the specified filename
			if err := temps.LoadData(config_file); err != nil {
				fmt.Printf("Error loading data: %v\n", err)
				return
			}
			temps.ModelDataFrame()
			temps.CommonCMD()
		},
	}
)

// Function to generate CRUD for Fiber
func genfibercurd() {
	temps.CurdFrameFiber()
	temps.CommonCMD()
}

// Function to generate Models for Echo or Fiber
func gengorm(frame string) {

	temps.DbConnDataFrame()
	if frame == "echo" {
		temps.ServiceFrame()
		temps.CurdFrameEcho()
	} else if frame == "fiber" {
		temps.FiberTracerFrame()
		temps.CurdFrameFiber()

	}
	temps.CommonCMD()
}

// Init function to add commands to the root
func init() {
	// Add flag for data-file
	curdcli.Flags().StringVarP(&config_file, "config", "c", "config.json", "Specify the data file to load")
	modelscli.Flags().StringVarP(&config_file, "config", "c", "config.json", "Specify the data file to load")

	curdcli.Flags().StringP("frame", "f", "", "Specify the framework to use (echo or fiber)")

	// Register commands to the root (goFrame)
	goFrame.AddCommand(curdcli)
	goFrame.AddCommand(modelscli)
}
