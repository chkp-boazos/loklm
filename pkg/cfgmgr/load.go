package cfgmgr

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/BurntSushi/toml"
)

const DIRNAME = "loklm"

type ConfigLoader interface {
	configLoad() (Configuration, error)
}

type TomlFileConfigLoader struct {
	filePath string
}

func (fcl TomlFileConfigLoader) configLoad() (Configuration, error) {
	var config Configuration
	var err error = nil
	_, err = toml.DecodeFile(fcl.filePath, &config)
	return config, err
}

func LoadToml(configPath string) (Configuration, error) {
	if configPath == "" {
		stateDir := fmt.Sprintf("/tmp/%s", DIRNAME)
		if runtime.GOOS == "windows" {
			programData := os.Getenv("ProgramData")
			if programData == "" {
				programData = "C:\\ProgramData"
			}
			stateDir = filepath.Join(programData, DIRNAME)

		}
		return Configuration{
			General{
				Version:  "v1",
				Name:     "",
				Network:  "cortex-net",
				StateDir: stateDir,
			},
			ExposedContainer{
				Name:        "jupyter",
				Port:        8888,
				ExposedPort: "8888/tcp",
				Hostname:    "jupyter",
				Image:       "jupyter/base-notebook",
				Dir:         "/home/jovyan",
			},
			ExposedContainer{
				Name:        "ollama",
				Port:        11434,
				ExposedPort: "11434/tcp",
				Hostname:    "llm",
				Image:       "ollama/ollama",
				Dir:         "/root/.ollama",
			},
			nil,
		}, nil
	}
	tomlCfgLoader := TomlFileConfigLoader{
		filePath: configPath,
	}
	return tomlCfgLoader.configLoad()
}
