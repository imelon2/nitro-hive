package simulate

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Providers      Providers      `yaml:"providers"`
	Signer         Signer         `yaml:"signer"`
	SimulateOption SimulateOption `yaml:"simulate-option"`
}

type Providers struct {
	Parent string `yaml:"parent"`
	Main   string `yaml:"main"`
}

type Signer struct {
	PrivateKey string `yaml:"privateKey"`
}

type SimulateOption struct {
	Total        int          `yaml:"total"`
	AccountRange AccountRange `yaml:"account-range"`
}

type AccountRange struct {
	StartIndex int `yaml:"start-index"`
	Total      int `yaml:"total"`
}

var (
	GlobalConfig Config
)

func NewConfig() *Config {
	return &GlobalConfig
}

func init() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("bad path")
	}
	root := filepath.Dir(filename)
	projectDir := filepath.Dir(root)
	configPath := filepath.Join(projectDir, "config.yml")

	file, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}

	err = yaml.Unmarshal(file, &GlobalConfig)
	if err != nil {
		log.Fatalf("failed to unmarshal YAML: %v", err)
	}
}
