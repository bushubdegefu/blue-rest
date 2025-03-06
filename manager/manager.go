package manager

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	goFrame = &cobra.Command{
		Use:           "Blue Rest",
		Short:         "Blue Rest – command-line tool to aid structure your golang backend projects with gorm and fiber/echo",
		Long:          "Blue Rest – command-line tool to aid structure your golang backend projects with gorm and fiber/echo for SQL based projects",
		Version:       "0.2.4",
		SilenceErrors: true,
		SilenceUsage:  true,
	}
)

func Execute() {
	if err := goFrame.Execute(); err != nil {

		fmt.Println(err)
		os.Exit(1)
	}
}
