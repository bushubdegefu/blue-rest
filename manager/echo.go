package manager

import (
	"fmt"
	"os/exec"

	"github.com/bushubdegefu/blue-rest/temps"
	"github.com/spf13/cobra"
)

var (
	echocli = &cobra.Command{
		Use:   "echo",
		Short: "generate the basic structure file to start app using echo",
		Long:  `generate the basic structure file to start app using echo`,
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
				temps.EchoFrameSetupAndMiddleware(appName)
				temps.AuthUtilsFrame(appName)
				temps.AuthLoginFrame(appName, "echo")

			} else if globalName {
				temps.EchoAppAndMiddleware()
				echogen()
			} else {
				fmt.Println("No app name specified")
			}
			temps.CommonCMD()
		},
	}
)

func echogen() {
	// running go mod tidy finally
	if err := exec.Command("swag", "init").Run(); err != nil {
		fmt.Printf("error generating swagger: %v \n", err)
	}
}

func init() {
	echocli.Flags().StringP("app", "a", "", "Specify the app name, so that echo app will be generated")
	echocli.Flags().BoolP("global", "g", false, "basic echo app with for global, creates app.go( in manager package) and middleware.go on the main module takes true or false")
	echocli.Flags().StringVarP(&config_file, "config", "c", "config.json", "Specify the data file to load")
	goFrame.AddCommand(echocli)
}
