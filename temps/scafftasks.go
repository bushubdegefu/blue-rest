package temps

import (
	"os"
	"text/template"
)

func TasksFrame() {
	// ####################################################
	//  rabbit template
	tasks_tmpl, err := template.New("RenderData").Funcs(FuncMap).Parse(tasksFileTemplate)
	if err != nil {
		panic(err)
	}

	// Create the models directory if it does not exist
	err = os.MkdirAll("bluetasks", os.ModePerm)
	if err != nil {
		panic(err)
	}

	tasks_file, err := os.Create("bluetasks/tasks.go")
	if err != nil {
		panic(err)
	}
	defer tasks_file.Close()

	err = tasks_tmpl.Execute(tasks_file, RenderData)
	if err != nil {
		panic(err)
	}
}

func LogFilesFrame() {
	// ####################################################
	// Log file For the APP
	log_file_tmpl, err := template.New("RenderData").Funcs(FuncMap).Parse(logFileTemplate)
	if err != nil {
		panic(err)
	}

	// Create the models directory if it does not exist
	err = os.MkdirAll("logs", os.ModePerm)
	if err != nil {
		panic(err)
	}

	log_file_file, err := os.Create("logs/logfile.go")
	if err != nil {
		panic(err)
	}
	defer log_file_file.Close()

	err = log_file_tmpl.Execute(log_file_file, RenderData)
	if err != nil {
		panic(err)
	}
}

var logFileTemplate = `
package logs

import (
	"log"
	"os"
	"fmt"
)

func Logfile(app_name string) (*os.File, error) {

	// Custom File Writer for logging
	log_file_name := fmt.Sprintf("%s_blue.log", app_name)
	file, err := os.OpenFile(log_file_name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
		return nil, err
	}
	return file, nil
}
`
var tasksFileTemplate = `
package bluetasks

import (
	"fmt"
	"os/exec"
	"strconv"
	"time"

	"{{.ProjectName}}/configs"
	"{{.ProjectName}}/logs"
	"{{.ProjectName}}/database"

	"github.com/madflojo/tasks"
)

func ScheduledTasks() *tasks.Scheduler {

	//  initalizing scheduler for regullarly running tasks
	scheduler := tasks.New()


	// // Add a task to move to Logs Directory Every Interval, Interval to Be Provided From Configuration File
	gormLoggerfile, _ := database.GormLoggerFile("{{ .AppName | replaceString }}")
	//  App should not start
	log_file, _ := logs.Logfile("{{ .AppName | replaceString }}")
	// Getting clear log interval from env
	clearIntervalLogs, _ := strconv.Atoi(configs.AppConfig.GetOrDefault("CLEAR_LOGS_INTERVAL", "1440"))
	if _, err := scheduler.Add(&tasks.Task{
		Interval: time.Duration(clearIntervalLogs) * time.Minute,
		TaskFunc: func() error {
			currentTime := time.Now()
			FileName := fmt.Sprintf("%v-%v-%v-%v-%v", currentTime.Year(), currentTime.Month(), currentTime.Day(), currentTime.Hour(), currentTime.Minute())
			//  make sure to replace the names of log files correctly here
			Command := fmt.Sprintf("cp {{ .AppName | replaceString }}_blue.log logs/{{ .AppName | replaceString }}_blue-%v.log", FileName)
			Command2 := fmt.Sprintf("cp {{ .AppName | replaceString }}_gorm.log logs/{{ .AppName | replaceString }}_gorm-%v.log", FileName)
			if _, err := exec.Command("bash", "-c", Command).Output(); err != nil {
				fmt.Printf("error: %v\n", err)
			}

			if _, err := exec.Command("bash", "-c", Command2).Output(); err != nil {
				fmt.Printf("error: %v\n", err)
			}

			err := gormLoggerfile.Truncate(0)
			if err != nil {
				fmt.Println("Error truncating gorm logger file:", err)
			}
			lerr := log_file.Truncate(0)
			if lerr != nil {
				fmt.Println("Error truncating log file:", err)
			}
			return nil
		},
	}); err != nil {
		fmt.Println(err)
	}

	return scheduler
}

`
