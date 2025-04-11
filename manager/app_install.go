package manager

import (
	"fmt"
	"os"
	"os/exec"

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

	installauthcli = &cobra.Command{
		Use:   "auth",
		Short: "Install Authentication Managment App With the UI",
		Long:  `Install Authentication Managment App With the UI`,
		Run: func(cmd *cobra.Command, args []string) {
			appName, _ := cmd.Flags().GetString("app")
			userName, _ := cmd.Flags().GetString("user")
			InstallAuthApp(userName, appName)
		},
	}
)

// Executes a command and returns any error
func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout // Pipe the command output to the terminal
	cmd.Stderr = os.Stderr // Pipe errors to the terminal

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("command %s %v failed: %w", name, args, err)
	}
	return nil
}

func InstallAuthApp(userName, appName string) {
	if userName == "" {
		fmt.Println("Should Provide Github Username")
		return
	}
	if appName == "" {
		appName = "django-auth"
	}

	modleName := fmt.Sprintf("gitthub.com/%v/%v", userName, appName)
	// ame = fmt.Sprintf(format string,"")
	// List of commands to run sequentially
	commands := []struct {
		cmd  string
		args []string
	}{
		{"brest", []string{"init", "-n", modleName}},
		{"brest", []string{"init", "-a", appName, "-p", appName}},
		{"brest", []string{"basic", "-t", "pagination"}},
		{"brest", []string{"basic", "-t", "db"}},
		{"brest", []string{"basic", "-t", "logs"}},
		{"brest", []string{"basic", "-t", "otel", "-f", "fiber"}},
		{"brest", []string{"basic", "-t", "tasks", "-a", appName}},
		{"brest", []string{"models", "-a", appName}},
		{"brest", []string{"crud", "-a", appName}},
		{"brest", []string{"fiber", "-a", appName}},
		{"brest", []string{"fiber", "-g", "true"}},
		{"brest", []string{"config"}},
		{"brest", []string{"basic", "-t", "migration"}},
		{"brest", []string{"ui"}},
	}

	// Run each command sequentially
	for _, cmd := range commands {
		err := runCommand(cmd.cmd, cmd.args...)
		if err != nil {
			fmt.Printf("Error running command: %v\n", err)
			return
		}
	}

	// After all commands are successfully executed
	fmt.Println("All CLI commands executed successfully.")
}

func init() {
	installauthcli.Flags().StringP("app", "a", "", "Specify the app name using the app flag")
	installauthcli.Flags().StringP("user", "u", "", "Specify the githubrepo user name using the user flag")
	goFrame.AddCommand(appinstallcli)
}
