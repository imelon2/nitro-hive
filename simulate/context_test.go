package simulate_test

import (
	"fmt"
	"testing"

	"github.com/imelon2/nitro-hive/simulate"
)

func Test_NewSimulateContext(t *testing.T) {
	context := simulate.NewSimulateContext()

	fmt.Printf("Simulate Private Key Length : %d\n", len(context.PrivateKey))
	fmt.Printf("Simulate Address Length     : %d\n", len(context.Address))

	fmt.Printf("Simulate Address Length     : %s\n", context.Address[0].Hex())
}
