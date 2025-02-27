package manager

import (
	"fmt"

	"github.com/bushubdegefu/blue-rest/temps"
	"github.com/spf13/cobra"
)

var (
	initalizemodule = &cobra.Command{
		Use:   "init",
		Short: "Initalize the module with name ",
		Long:  `Provide name to initalize the project ussing the \"name\" tag.`,
		Run: func(cmd *cobra.Command, args []string) {
			moduleName, _ := cmd.Flags().GetString("name")
			temps.CommonModInit(moduleName)
			temps.CommonTargetJSON(moduleName)
		},
	}
	// Parent 'basic' command that uses a type flag
	basicCommand = &cobra.Command{
		Use:   "basic",
		Short: "Generate a basic folder structure for a project",
		Long:  `This command generates a basic folder structure for a project. The type flag determines the specific setup.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Check for the type flag
			projectType, _ := cmd.Flags().GetString("type")

			// Load data with the specified filename
			if err := temps.LoadData(config_file); err != nil {
				fmt.Printf("Error loading data: %v\n", err)
				return
			}

			// Run corresponding command based on the project type
			switch projectType {
			case "git":
				basiccmd()
			case "rsa":
				rsa_basic()
			case "db":
				standarddatabase()
				temps.CommonCMD()
			case "consumer":
				standardrabbit()
				temps.CommonCMD()
			case "producer":
				standpublish()
				temps.CommonCMD()
			case "tasks":
				standtasks()
				temps.CommonCMD()
			case "pagination":
				commongormpagination()
				temps.CommonCMD()
			default:
				fmt.Println("Unknown type specified. Valid types are: rsa, db, rabbit, tasks, pagination, git.")
			}
		},
	}
)

func projectintialize() {

}

func basiccmd() {
	temps.GitDockerFrame()
	temps.HaproxyFrame()
}

func standtasks() {
	temps.TasksFrame()
	temps.LogFilesFrame()
}

func rsa_basic() {
	temps.RSAHelper()
}

func commongormpagination() {
	temps.CommonFrame()
}

func standardcmd() {
	temps.Frame()
	temps.DbConnDataFrame()
	temps.StandardTracerFrame()
	temps.CommonFrame()
	temps.RabbitFrame()
}

func standardrabbit() {
	temps.Frame()
	temps.StandardTracerFrame()
	temps.CommonRabbitFrame()
	temps.RabbitFrame()
	temps.PublishFrame()
	temps.ConsumeFrame()
	temps.RunConsumeFrame()

}

func standpublish() {
	temps.Frame()
	temps.StandardTracerFrame()
	temps.CommonRabbitFrame()
	temps.RabbitFrame()
	temps.PublishFrame()

}

func standarddatabase() {
	temps.Frame()
	temps.DbConnDataFrame()
}

func init() {
	// Register flags for all commands
	initalizemodule.Flags().StringP("name", "n", "", "Specify the module name  (github.com/someuser/someproject)")
	// Register flags for the 'basic' command
	basicCommand.Flags().StringP("type", "t", "", "Specify the type of folder structure to generate: rsa, db, producer, consumer, tasks, pagination")

	goFrame.AddCommand(basicCommand)
	goFrame.AddCommand(initalizemodule)

}
