
package configs

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

const (
	defaultFileName         = "/.env"
	defaultOverrideFileName = "/.local.env"
)

type EnvConfig struct {
	defaultPath string
	prodFlag    string
}

type Config interface {
	Get(string) string
	GetOrDefault(string, string) string
}

var AppConfig EnvConfig

func NewEnvFile(configFolder string) {
	AppConfig = EnvConfig{
		defaultPath: configFolder,
	}
	AppConfig.read()
}

func (e *EnvConfig) read() {
	var (
		defaultFile  = e.defaultPath + defaultFileName
		overrideFile = e.defaultPath + defaultOverrideFileName
	)
	err := godotenv.Overload(defaultFile)
	env := e.Get("APP_ENV")
	if err != nil {
		fmt.Printf("WARNING: Failed to load config from file: %v, Err: %v \n", defaultFile, err)
	} else {
		fmt.Printf("INFO: Loaded config from file: %v\n", defaultFile)
	}

	// If 'APP_ENV' is set to x, then App will read '.env' from configs directory, and then it will be overwritten
	// by configs present in file '.x.env'
	overrideFile = fmt.Sprintf("%s/.%s.env", e.defaultPath, env)
	if env == "" && e.prodFlag == "" {
		overrideFile = fmt.Sprintf("%s/dev.env", e.defaultPath)
	}

	if e.prodFlag != "" {
		overrideFile = fmt.Sprintf("%s/.%s.env", e.defaultPath, e.prodFlag)
	}
	err = godotenv.Overload(overrideFile)
	if err != nil {
		fmt.Printf("WARNING: to load config from file: %v, Err: %v \n", overrideFile, err)
	} else {
		fmt.Printf("INFO: config from file: %v \n", overrideFile)
	}
}

func (e *EnvConfig) Get(key string) string {
	return os.Getenv(key)
}

func (e *EnvConfig) SetEnv(key string) {
	AppConfig = EnvConfig{
		defaultPath: "./configs",
		prodFlag: key,
	}
	AppConfig.read()
}

func (e *EnvConfig) GetOrDefault(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return defaultValue
}
