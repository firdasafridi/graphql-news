package platform

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type Environment struct {
	AppName    string `json:"app_name"`
	Database   string `json:"database"`
	Host       string `json:"host"`
	Port       string `json:"port"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	SSLMode    string `json:"ssl_mode"`
	SearchPath string `json:"search_path"`
	Debug      bool   `json:"debug"`
	DirWork    string `json:"dir_word"`
}

func NewEnvironment() (env *Environment) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("NewEnvironment: " + err.Error())
	}

	env = &Environment{
		DirWork: wd,
	}
	env.loadConfig()

	return env
}

func (env *Environment) loadConfig() {

	wd := env.DirWork
	configPath := filepath.Join(wd, pathEnv)

	b, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal("platform: Environment.load: cannot find " + configPath)
	}

	fmt.Printf("platform: loading configuration from %q ...\n", configPath)
	err = json.Unmarshal(b, env)
	if err != nil {
		log.Fatalf("platform: Environment.load: %s: %s", configPath, err.Error())
	}

	if env.Debug {
		log.Println("Debug is active.")
		log.Println(string(b))
	}
}
