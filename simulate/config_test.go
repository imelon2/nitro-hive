package simulate_test

import (
	"fmt"
	"testing"

	"github.com/imelon2/nitro-hive/simulate"
)

func Test_(t *testing.T) {
	config := simulate.NewConfig()
	fmt.Printf("Config.providers.parent                             : %s\n", config.Providers.Parent)
	fmt.Printf("Config.providers.Main                               : %s\n\n", config.Providers.Main)
	fmt.Printf("config.SimulateOptions.Total                        : %d\n", config.SimulateOptions.Total)
	fmt.Printf("config.SimulateOptions.TransactionOptions.Value     : %d\n", config.SimulateOptions.TransactionOptions.Value)
	fmt.Printf("config.SimulateOptions.TransactionOptions.Gas       : %d\n", config.SimulateOptions.TransactionOptions.Gas)
	fmt.Printf("config.SimulateOptions.TransactionOptions.GasPrice  : %d\n\n", config.SimulateOptions.TransactionOptions.GasPrice)
	fmt.Printf("config.SimulateOption.SingleSigner.PrivateKey       : %s\n", config.SimulateOptions.SingleSigner.PrivateKey)
	fmt.Printf("config.SimulateOption.MultiSigner.StartIndex        : %d\n", config.SimulateOptions.MultiSigner.StartIndex)
	fmt.Printf("config.SimulateOption.MultiSigner.Total             : %d\n", config.SimulateOptions.MultiSigner.Total)

}
