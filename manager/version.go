package manager

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	versioncli = &cobra.Command{
		Use:   "version",
		Short: "Print out the version of the cli app on terminal",
		Long:  `Print out the version of the cli app on terminal`,
		Run: func(cmd *cobra.Command, args []string) {
			version()
		},
	}
)

func version() {
	fmt.Println("Frame version 0.0.2")
}

func init() {
	goFrame.AddCommand(versioncli)

}
