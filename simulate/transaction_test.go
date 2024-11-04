package simulate_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/imelon2/nitro-hive/common/utils"
	si "github.com/imelon2/nitro-hive/simulate"
)

func Test_NativeTransafer(t *testing.T) {
	simulate := si.NewSimulateContext()

	pk := utils.Unhexlify(si.GlobalConfig.SimulateOptions.SingleSigner.PrivateKey)
	key, _ := crypto.HexToECDSA(pk)
	signer, _ := si.NewSginerContext(key)
	tx := signer.NativeTransafer()

	simulate.Simulate(tx)
}
