package manager

import (
	"github.com/bushubdegefu/blue-rest/dist"
	"github.com/spf13/cobra"
)

var (
	appinstallcli = &cobra.Command{
		Use:   "ui",
		Short: "Create React Django Admin UI dist files",
		Long:  `Create React Django Admin UI dist files`,
		Run: func(cmd *cobra.Command, args []string) {
			dist.DjangoAdminUI()
		},
	}
)

func init() {
	goFrame.AddCommand(appinstallcli)
}
