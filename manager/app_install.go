package manager

import (
	"fmt"
	"os"

	"github.com/bushubdegefu/blue-rest/dist"
	"github.com/bushubdegefu/blue-rest/temps"
	"github.com/bushubdegefu/blue-rest/temps/generator"
	"github.com/spf13/cobra"
)

var (
	appinstallcli = &cobra.Command{
		Use:   "ui",
		Short: "Create React Django Admin UI dist files",
		Long:  `Create React Django Admin UI dist files`,
		Run: func(cmd *cobra.Command, args []string) {
			temps.InitProjectJSON()
			if temps.ProjectSettings.AuthAppType == "sso" {
				dist.SSOAdminUI()
			} else {
				dist.DjangoAdminUI()

			}
		},
	}

	installauthcli = &cobra.Command{
		Use:   "auth",
		Short: "Install Authentication Managment App With the UI",
		Long:  `Install Authentication Managment App With the UI`,
		Run: func(cmd *cobra.Command, args []string) {
			authAppName, _ := cmd.Flags().GetString("app")
			projectName, _ := cmd.Flags().GetString("project")
			userName, _ := cmd.Flags().GetString("user")
			frame, _ := cmd.Flags().GetString("frame")
			authType, _ := cmd.Flags().GetString("sso")
			if authType == "sso" {
				InstallSSOhApp(userName, projectName, authAppName, frame)
			} else {
				InstallAuthApp(userName, projectName, authAppName, frame)
			}
			runSwagInitForApps()
			temps.CommonCMDInit()
			temps.CommonCMD()
		},
	}
)

func InstallAuthApp(userName, projectName, authAppName, frame string) {
	if userName == "" {
		fmt.Println("Should Provide Github Username")
		return
	}
	if authAppName == "" {
		authAppName = "django-auth"
	}

	moduleName := fmt.Sprintf("github.com/%v/%v", userName, projectName)

	temps.CommonProjectName(moduleName, projectName, "standalone")
	temps.CommonModInit(moduleName)
	// Get current working directory
	temps.InitProjectJSON()

	temps.RenderData.ProjectName = moduleName
	generator.GenerateMainAndManager(temps.RenderData)
	generator.GenerateLogs(temps.RenderData)
	generator.GenerateCommon(temps.RenderData)
	generator.GenerateDBConn(temps.ProjectSettings)
	handleOtelFrame(frame)

	currentDir, _ := os.Getwd()
	handleAppInitialization(authAppName, currentDir, authAppName)
	os.Chdir(currentDir)
	_ = handleAppDirectoryAndLoadConfig(authAppName)

	temps.RenderData.AuthAppName = temps.ProjectSettings.AuthAppName
	temps.RenderData.AppName = temps.ProjectSettings.AuthAppName
	temps.RenderData.AppNames = temps.ProjectSettings.AppNames
	temps.ProjectSettings.CurrentAppName = authAppName

	generator.GenerateTasks(temps.RenderData)
	generator.GenerateConfigTestEnv(temps.RenderData)
	generateTests(frame)

	generator.GenerateJWTUtils(temps.ProjectSettings)
	generator.GenerateUtilsApp(temps.ProjectSettings)
	generator.GenerateModels(temps.RenderData)
	gengorm(frame)

	if frame == "echo" {
		generator.GenerateEchoAppMiddleware(temps.RenderData)
		generator.GenerateEchoSetup(temps.RenderData)
		loginFrame(authAppName, "echo")
	} else {
		generator.GenerateFiberAppMiddleware(temps.RenderData)
		generator.GenerateFiberSetup(temps.RenderData)
		loginFrame(authAppName, "fiber")
	}

	// Go back to root directory
	os.Chdir(currentDir)
	generator.GenerateAppDatabaseMigration(temps.RenderData)
	if frame == "echo" {
		generator.GenerateGlobalEchoAppMiddleware(temps.RenderData)
		generator.GenerateAppEchoGlobal(temps.RenderData)
	} else {
		generator.GenerateGlobalFiberAppMiddleware(temps.RenderData)
		generator.GenerateAppFiberGlobal(temps.RenderData)
	}
	temps.RenderData.ProjectName = temps.ProjectSettings.ProjectName
	temps.RenderData.AppNames = temps.ProjectSettings.AppNames
	generator.GenerateConfig(temps.RenderData)
	generator.GenerateConfigEnv(temps.RenderData)
	generator.GenerateConfigAppEnv(temps.RenderData)

	dist.DjangoAdminUI()
	// After all commands are successfully executed
	fmt.Println("App Installed successfully.")
}

func InstallSSOhApp(userName, projectName, authAppName, frame string) {
	if userName == "" {
		fmt.Println("Should Provide Github Username")
		return
	}
	if authAppName == "" {
		authAppName = "django-auth"
	}

	moduleName := fmt.Sprintf("github.com/%v/%v", userName, projectName)
	temps.CommonProjectName(moduleName, projectName, "sso")
	temps.CommonModInit(moduleName)
	// Get current working directory
	temps.InitProjectJSON()

	temps.RenderData.ProjectName = moduleName
	generator.GenerateMainAndManager(temps.RenderData)
	generator.GenerateLogs(temps.RenderData)
	generator.GenerateCommon(temps.RenderData)
	generator.GenerateDBConn(temps.ProjectSettings)
	handleOtelFrame(frame)

	currentDir, _ := os.Getwd()
	handleAppInitialization(authAppName, currentDir, authAppName)
	os.Chdir(currentDir)
	_ = handleAppDirectoryAndLoadConfig(authAppName)

	temps.RenderData.AuthAppName = temps.ProjectSettings.AuthAppName
	temps.RenderData.AppName = temps.ProjectSettings.AuthAppName
	temps.RenderData.AppNames = temps.ProjectSettings.AppNames
	temps.ProjectSettings.CurrentAppName = authAppName

	generator.GenerateTasks(temps.RenderData)
	generator.GenerateConfigTestEnv(temps.RenderData)
	generateTests(frame)

	generator.GenerateJWTUtils(temps.ProjectSettings)
	generator.GenerateUtilsApp(temps.ProjectSettings)
	generator.GenerateModels(temps.RenderData)
	gengorm(frame)
	if frame == "echo" {
		generator.GenerateEchoAppMiddleware(temps.RenderData)
		generator.GenerateEchoSetup(temps.RenderData)
		loginFrame(authAppName, "echo")
	} else {
		generator.GenerateFiberAppMiddleware(temps.RenderData)
		generator.GenerateFiberSetup(temps.RenderData)
		loginFrame(authAppName, "fiber")
	}

	// Go back to root directory
	os.Chdir(currentDir)
	generator.GenerateAppDatabaseMigration(temps.RenderData)
	if frame == "echo" {
		generator.GenerateGlobalEchoAppMiddleware(temps.RenderData)
		generator.GenerateAppEchoGlobal(temps.RenderData)
	} else {
		generator.GenerateGlobalFiberAppMiddleware(temps.RenderData)
		generator.GenerateAppFiberGlobal(temps.RenderData)
	}
	temps.RenderData.ProjectName = temps.ProjectSettings.ProjectName
	temps.RenderData.AppNames = temps.ProjectSettings.AppNames
	generator.GenerateConfig(temps.RenderData)
	generator.GenerateConfigEnv(temps.RenderData)
	generator.GenerateConfigAppEnv(temps.RenderData)

	dist.SSOAdminUI()
	// After all commands are successfully executed

	fmt.Println("App Installed successfully.")
}

func init() {
	installauthcli.Flags().StringP("frame", "f", "", "Specify the framework for the template (echo or fiber)")
	installauthcli.Flags().StringP("project", "p", "", "Specify the project name using the app flag")
	installauthcli.Flags().StringP("app", "a", "", "Specify the app name using the app flag")
	installauthcli.Flags().StringP("user", "u", "", "Specify the githubrepo user name using the user flag")
	installauthcli.Flags().StringP("sso", "s", "", "Specify the authentication app type sso or standalone")
	goFrame.AddCommand(appinstallcli)
	goFrame.AddCommand(installauthcli)
}
