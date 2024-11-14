package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Providers          Providers          `yaml:"providers"`
	SimulateOptions    SimulateOptions    `yaml:"simulate-options"`
	TransactionOptions TransactionOptions `yaml:"transaction-options"`
}

type Providers struct {
	Main string `yaml:"main"`
}

type SimulateOptions struct {
	Total        int          `yaml:"total"`
	SingleSigner SingleSigner `yaml:"single-signer"`
	MultiSigner  MultiSigner  `yaml:"multi-signer"`
}

type TransactionOptions struct {
	Value    int    `yaml:"value"`
	Gas      int    `yaml:"gas"`
	GasPrice int    `yaml:"gas-price"`
	Data     string `yaml:"data"`
	To       string `yaml:"to"`
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

func GetCpu() {
	currCPU := runtime.NumCPU() // 내 PC CPU 개수
	// runtime.GOMAXPROCS(currCPU)
	fmt.Println("Max System CPU : ", currCPU)                   // 설정값 출력
	fmt.Println("Current System CPU : ", runtime.GOMAXPROCS(0)) // 설정값 출력
}
