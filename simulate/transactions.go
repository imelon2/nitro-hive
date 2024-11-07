package simulate

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/imelon2/nitro-hive/common/utils"
)

func (signer *SignerContext) NativeTransafer() func(*common.Address) (*types.Transaction, error) {
	txFunc := func(to *common.Address) (*types.Transaction, error) {
		signer.NonceMutex.Lock()
		defer signer.NonceMutex.Unlock()
		signer.UpdateNonce()

		var err error
		GasPrice := big.NewInt(int64(GlobalConfig.TransactionOptions.GasPrice))
		if GasPrice.Cmp(big.NewInt(0)) == 0 {
			GasPrice, err = signer.MainClient.SuggestGasPrice(signer.Ctx)
			if err != nil {
				return nil, err
			}
		}
		Data := common.Hex2Bytes(utils.Unhexlify(GlobalConfig.TransactionOptions.Data))
		if to == nil {
			_to := common.HexToAddress(GlobalConfig.TransactionOptions.To)
			to = &_to
		}

		tx := types.NewTx(&types.LegacyTx{
			Nonce:    signer.SignerOpt.Nonce.Uint64(),
			To:       to,
			Value:    big.NewInt(int64(GlobalConfig.TransactionOptions.Value)),
			Gas:      uint64(GlobalConfig.TransactionOptions.Gas),
			GasPrice: GasPrice,
			Data:     Data,
		})

		signedTx, err := signer.SignerOpt.Signer(*signer.Account, tx)
		if err != nil {
			return nil, err
		}
		err = signer.MainClient.SendTransaction(signer.Ctx, signedTx)
		return signedTx, err
	}

	return txFunc
}

// func NativeTransafer1() {
// 	pk := utils.Unhexlify(GlobalConfig.SimulateOptions.SingleSigner.PrivateKey)
// 	key, _ := crypto.HexToECDSA(pk)
// 	signer, _ := NewSginerContext(key)

// 	simulate := NewSimulateContext()

// 	txFunc := func(id int) (*types.Transaction, error) {
// 		signer.NonceMutex.Lock()
// 		defer signer.NonceMutex.Unlock()
// 		signer.UpdateNonce()

// 		transferAmount := new(big.Int).Mul(big.NewInt(params.GWei), big.NewInt(1))

// 		tx := types.NewTx(&types.LegacyTx{
// 			Nonce:    signer.SignerOpt.Nonce.Uint64(),
// 			To:       simulate.Address[id],
// 			Value:    transferAmount,
// 			Gas:      25000,
// 			GasPrice: big.NewInt(100000000),
// 		})

// 		signedTx, err := signer.SignerOpt.Signer(*signer.Account, tx)
// 		if err != nil {
// 			return nil, err
// 		}
// 		err = signer.MainClient.SendTransaction(signer.Ctx, signedTx)
// 		return signedTx, err
// 	}

// 	simulate.Simulate(txFunc)
// }
