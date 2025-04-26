package temps

import (
	"fmt"
	"os/exec"
	"time"
)

func CommonCMDInit() {
	time.Sleep(2 * time.Second)
	// running go mod tidy finally
	if err := exec.Command("go", "get", "-u", ".").Run(); err != nil {
		fmt.Printf("error go get: %v \n", err)
	}
}

func CommonCMD() {

	time.Sleep(2 * time.Second)
	// running go mod tidy finally
	if err := exec.Command("go", "mod", "tidy").Run(); err != nil {
		fmt.Printf("error tidy: %v \n", err)
	}
}

func CommonModInit(project_module string) {
	// running go mod tidy finally
	if err := exec.Command("go", "mod", "init", project_module).Run(); err != nil {
		fmt.Printf("error: %v \n", err)
	}
}
