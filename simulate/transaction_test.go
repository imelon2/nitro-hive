package simulate_test

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/imelon2/nitro-hive/common/utils"
	si "github.com/imelon2/nitro-hive/simulate"
)

func Test_MultiSimulate(t *testing.T) {
	simulate := si.NewSimulateContext()

	start := si.GlobalConfig.SimulateOptions.MultiSigner.StartIndex
	total := si.GlobalConfig.SimulateOptions.MultiSigner.Total

	txFuncs := make([]func(*common.Address) (*types.Transaction, error), total)
	for i := start; i < total; i++ {
		signer, err := si.NewSginerContext(simulate.PrivateKey[i])
		if err != nil {
			t.Fatalf("failed to create signer context: %v", err)
		}
		tx := signer.NativeTransafer()

		txFuncs[i] = tx
	}

	fmt.Printf("txFuns length: %d\n", len(txFuncs))
	simulate.MultiSimulate(txFuncs)
}

func Test_NativeTransafer(t *testing.T) {
	simulate := si.NewSimulateContext()

	pk := utils.Unhexlify(si.GlobalConfig.SimulateOptions.SingleSigner.PrivateKey)
	key, _ := crypto.HexToECDSA(pk)
	signer, _ := si.NewSginerContext(key)
	tx := signer.NativeTransafer()

	simulate.Distribute(tx)
}
