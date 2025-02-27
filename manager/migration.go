package manager

import (
	"fmt"

	"github.com/bushubdegefu/blue-rest/temps"
	"github.com/spf13/cobra"
)

var (
	migrationdcli = &cobra.Command{
		Use:   "migration",
		Short: "Generate Data Models based on the GORM using provided spec on the config.json file ",
		Long:  `Generate Data Models based on the GORM using provided spec on the config.json file`,
		Run: func(cmd *cobra.Command, args []string) {

			// Load data with the specified filename
			if err := temps.LoadData(config_file); err != nil {
				fmt.Printf("Error loading data: %v\n", err)
				return
			}
			migrationmod()
		},
	}
)

func migrationmod() {

	temps.MigrationFrame()
	temps.CommonCMD()
}

func init() {
	goFrame.AddCommand(migrationdcli)
}
