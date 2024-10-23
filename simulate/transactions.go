package simulate

import (
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
		signer.NonceMutex.Lock()
		defer signer.NonceMutex.Unlock()
		signer.UpdateNonce()

		transferAmount := new(big.Int).Mul(big.NewInt(params.GWei), big.NewInt(1))

		tx := types.NewTx(&types.LegacyTx{
			Nonce:    signer.SignerOpt.Nonce.Uint64(),
			To:       simulate.Address[id],
			Value:    transferAmount,
			Gas:      25000,
			GasPrice: big.NewInt(100000000),
		})

		signedTx, err := signer.SignerOpt.Signer(*signer.Account, tx)
		if err != nil {
			return nil, err
		}
		err = signer.MainClient.SendTransaction(signer.Ctx, signedTx)
		return signedTx, err
	}

	simulate.Simulate(txFunc)
}
