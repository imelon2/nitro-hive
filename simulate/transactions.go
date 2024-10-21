package simulate

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	"github.com/imelon2/nitro-hive/common/utils"
)

func NativeTransafer() {
	pk := utils.Unhexlify(GlobalConfig.Signer.PrivateKey)
	key, _ := crypto.HexToECDSA(pk)
	signer, _ := NewSginerContext(key)

	simulate := NewSimulateContext()

	txFunc := func(id int) (*types.Transaction, error) {

		fmt.Printf("IS? : %d\n", id)
		transferAmount := new(big.Int).Mul(big.NewInt(params.Ether), big.NewInt(1))

		tx := types.NewTx(&types.LegacyTx{
			Nonce: signer.SignerOpt.Nonce.Uint64(),
			To:    simulate.Address[id],
			Value: transferAmount,
			// Gas:      config.GasLimit,
			// GasPrice: gasPrice,
		})

		signedTx, err := signer.SignerOpt.Signer(*signer.Account, tx)
		if err != nil {
			return nil, err
		}
		return signedTx, signer.MainClient.SendTransaction(signer.Ctx, signedTx)
	}

	simulate.Simulate(txFunc)
}
