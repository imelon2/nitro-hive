package simulate_test

import (
	"fmt"
	"testing"

	"github.com/imelon2/nitro-hive/simulate"
)

func Test_(t *testing.T) {
	config := simulate.NewConfig()
	fmt.Printf("Config.providers.parent   : %s\n", config.Providers.Parent)
	fmt.Printf("Config.providers.Main    : %s\n", config.Providers.Main)
	fmt.Printf("config.Signer.PrivateKey  : %s\n\n", config.Signer.PrivateKey)

	fmt.Printf("config.SimulateOption.Total  : %d\n", config.SimulateOption.Total)
	fmt.Printf("config.SimulateOption.AccountRange.StartIndex  : %d\n", config.SimulateOption.AccountRange.StartIndex)
	fmt.Printf("config.SimulateOption.AccountRange.Total  : %d\n", config.SimulateOption.AccountRange.Total)

}
