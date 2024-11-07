package simulate_test

import (
	"fmt"
	"testing"

	"github.com/imelon2/nitro-hive/simulate"
)

func Test_Config(t *testing.T) {
	config := simulate.NewConfig()
	fmt.Printf("Config.providers.Main                               : %s\n\n", config.Providers.Main)
	fmt.Printf("config.SimulateOptions.Total                        : %d\n\n", config.SimulateOptions.Total)
	fmt.Printf("config.SimulateOptions.SingleSigner.PrivateKey      : %s\n\n", config.SimulateOptions.SingleSigner.PrivateKey)
	fmt.Printf("config.SimulateOptions.MultiSigner.Total            : %d\n", config.SimulateOptions.MultiSigner.Total)
	fmt.Printf("config.SimulateOptions.MultiSigner.StartIndex       : %d\n\n", config.SimulateOptions.MultiSigner.StartIndex)
	fmt.Printf("config.SimulateOptions.TransactionOptions.Gas       : %d\n", config.TransactionOptions.Gas)
	fmt.Printf("config.SimulateOptions.TransactionOptions.GasPrice  : %d\n", config.TransactionOptions.GasPrice)
	fmt.Printf("config.SimulateOptions.TransactionOptions.Value     : %d\n", config.TransactionOptions.Value)
	fmt.Printf("config.SimulateOptions.TransactionOptions.Data      : %s\n", config.TransactionOptions.Data)
	fmt.Printf("config.SimulateOptions.TransactionOptions.To        : %s\n", config.TransactionOptions.To)

}

func Test_GetCpu(t *testing.T) {
	simulate.GetCpu()
}
