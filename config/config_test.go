package config_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/imelon2/nitro-hive/config"
)

func Test_Config(t *testing.T) {
	config := config.GlobalConfig
	jsonData, _ := json.MarshalIndent(config, "", "  ")
	fmt.Println(string(jsonData))
}

func Test_GetCpu(t *testing.T) {
	config.GetCpu()
}
