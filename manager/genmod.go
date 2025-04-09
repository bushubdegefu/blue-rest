package manager

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bushubdegefu/blue-rest/temps"
	"github.com/spf13/cobra"
)

var (
	config_file string

	// CRUD command for generating CRUD handlers
	curdcli = &cobra.Command{
		Use:   "crud",
		Short: "Generate CRUD handlers based on GORM for the specified framework",
		Long:  `This command generates only the CRUD handlers using GORM, based on the provided spec, for either the Echo or Fiber framework.`,
		Run:   runCrudCommand,
	}

	// Models command for generating data models
	modelscli = &cobra.Command{
		Use:   "models",
		Short: "Generate data models based on GORM with annotations",
		Long:  `This command generates data models using GORM, based on the provided spec in the config.json file, along with GORM relationship annotations.`,
		Run:   runModelsCommand,
	}
)

// runCrudCommand handles the execution of the 'crud' command
func runCrudCommand(cmd *cobra.Command, args []string) {
	frame, _ := cmd.Flags().GetString("frame")
	appName, _ := cmd.Flags().GetString("app")
	temps.InitProjectJSON()
	temps.RenderData.AppName = appName
	if frame == "" {
		fmt.Println("Error: --frame flag is required. Use --frame=echo or --frame=fiber.")
		return
	}

	if appName == "" {
		fmt.Println("Error: --app flag is required.")
		return
	}

	temps.RenderData.AppName = appName
	// Change to the app's directory and load the config data
	if err := handleAppDirectoryAndLoadConfig(appName); err != nil {
		fmt.Println(err)
		return
	}

	// Generate CRUD based on the frame type
	if frame == "echo" || frame == "fiber" {
		temps.AuthUtilsFrame(appName)
		gengorm(frame)

	} else {
		fmt.Println("Error: Invalid frame value. Use --frame=echo or --frame=fiber.")
	}
}

// runModelsCommand handles the execution of the 'models' command
func runModelsCommand(cmd *cobra.Command, args []string) {
	modelsType, _ := cmd.Flags().GetString("type")
	appName, _ := cmd.Flags().GetString("app")
	appAuth, _ := cmd.Flags().GetBool("auth")

	if appName == "" {
		fmt.Println("Error: --app flag is required.")
		return
	}

	// Change to the app's directory and load the config data
	if err := handleAppDirectoryAndLoadConfig(appName); err != nil {
		fmt.Println(err)
		return
	}
	temps.RenderData.AuthApp = false
	if appAuth == true {
		temps.RenderData.AuthApp = true
	}
	// Generate models and migrations
	if modelsType == "init" {
		temps.ModelDataFrame()
		temps.MigrationInit()
	} else {
		temps.ModelDataFrame()
		temps.CommonCMD()
	}
}

// handleAppDirectoryAndLoadConfig changes the working directory to the app's directory and loads the config data
func handleAppDirectoryAndLoadConfig(appName string) error {
	currentDir, _ := os.Getwd()
	newDir := filepath.Join(currentDir, appName)
	if err := os.Chdir(newDir); err != nil {
		return fmt.Errorf("error changing directory: %v", err)
	}

	if err := temps.LoadData(config_file); err != nil {
		return fmt.Errorf("error loading data: %v", err)
	}
	return nil
}

// gengorm generates GORM-related handlers for either Echo or Fiber
func gengorm(frame string) {

	if frame == "echo" {
		temps.CurdFrameEcho()
	} else if frame == "fiber" {
		temps.CurdFrameFiber()

	}
	temps.CommonCMD()
}

func init() {
	// Register flags for CRUD and Models commands
	modelscli.Flags().StringVarP(&config_file, "config", "c", "config.json", "Specify the data file to load")
	modelscli.Flags().StringP("type", "t", "", "Rerender the migration function by setting type to \"init\"")
	modelscli.Flags().StringP("app", "a", "", "Set app name, e.g., \"blue-auth\"")
	modelscli.Flags().BoolP("auth", "au", false, "Tell if generating models for auth app true or false")

	curdcli.Flags().StringVarP(&config_file, "config", "c", "config.json", "Specify the data file to load")
	curdcli.Flags().StringP("frame", "f", "", "Specify the framework to use (echo or fiber)")
	curdcli.Flags().StringP("app", "a", "", "Specify the app name using the app flag")

	// Register commands to the root (goFrame)
	goFrame.AddCommand(curdcli)
	goFrame.AddCommand(modelscli)
}
