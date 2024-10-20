package simulate_test

import (
	"fmt"
	"testing"

	"github.com/imelon2/nitro-hive/simulate"
)

func Test_NewSimulateContext(t *testing.T) {
	context := simulate.NewSimulateContext(10)

	fmt.Printf("Simulate Private Key Length : %d", len(context.PrivateKey))
}
