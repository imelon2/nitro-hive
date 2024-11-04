package simulate

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Providers       Providers       `yaml:"providers"`
	SimulateOptions SimulateOptions `yaml:"simulate-options"`
}

type Providers struct {
	Parent string `yaml:"parent"`
	Main   string `yaml:"main"`
}

type SimulateOptions struct {
	Total              int                `yaml:"total"`
	TransactionOptions TransactionOptions `yaml:"transaction-options"`
	SingleSigner       SingleSigner       `yaml:"single-signer"`
	MultiSigner        MultiSigner        `yaml:"multi-signer"`
}

type TransactionOptions struct {
	Value    int `yaml:"value"`
	Gas      int `yaml:"gas"`
	GasPrice int `yaml:"gas-price"`
}

type SingleSigner struct {
	PrivateKey string `yaml:"private-key"`
}

type MultiSigner struct {
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
