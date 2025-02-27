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

			// Load data with the specified filename
			if err := temps.LoadData(config_file); err != nil {
				fmt.Printf("Error loading data: %v\n", err)
				return
			}

			echogen()
		},
	}
)

func echogen() {
	temps.EchoFrame()
	// running go mod tidy finally
	if err := exec.Command("swag", "init").Run(); err != nil {
		fmt.Printf("error generating swagger: %v \n", err)
	}
	temps.CommonCMD()
}

func init() {
	goFrame.AddCommand(echocli)
}
