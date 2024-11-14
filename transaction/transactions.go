package transaction

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	constant "github.com/imelon2/nitro-hive/common"
	"github.com/imelon2/nitro-hive/common/utils"
	config "github.com/imelon2/nitro-hive/config"
	multicall3 "github.com/imelon2/nitro-hive/solgen"
)

func (signer *SignerContext) NativeTransafer() func(*common.Address) (*types.Transaction, error) {
	txFunc := func(to *common.Address) (*types.Transaction, error) {
		signer.NonceMutex.Lock()
		defer signer.NonceMutex.Unlock()
		signer.UpdateNonce()

		var err error
		GasPrice := big.NewInt(int64(config.GlobalConfig.TransactionOptions.GasPrice))
		if GasPrice.Cmp(big.NewInt(0)) == 0 {
			GasPrice, err = signer.MainClient.SuggestGasPrice(*signer.Ctx)
			if err != nil {
				return nil, err
			}
		}
		Data := common.Hex2Bytes(utils.Unhexlify(config.GlobalConfig.TransactionOptions.Data))
		if to == nil {
			_to := common.HexToAddress(config.GlobalConfig.TransactionOptions.To)
			to = &_to
		}

		tx := types.NewTx(&types.LegacyTx{
			Nonce:    signer.SignerOpt.Nonce.Uint64(),
			To:       to,
			Value:    big.NewInt(int64(config.GlobalConfig.TransactionOptions.Value)),
			Gas:      uint64(config.GlobalConfig.TransactionOptions.Gas),
			GasPrice: GasPrice,
			Data:     Data,
		})

		signedTx, err := signer.SignerOpt.Signer(*signer.Account, tx)
		if err != nil {
			return nil, err
		}
		err = signer.MainClient.SendTransaction(*signer.Ctx, signedTx)
		return signedTx, err
	}

	return txFunc
}

func (signer *SignerContext) Distribute(to []*common.Address, gasPrice *big.Int, gasLimit uint64, amountPerAccount []*big.Int) func() (*types.Transaction, error) {

	txFunc := func() (*types.Transaction, error) {
		signer.NonceMutex.Lock()
		signer.UpdateNonce()
		defer signer.NonceMutex.Unlock()

		Multicall3, err := multicall3.NewMulticall3(common.HexToAddress(constant.MULTICALL_ADDRESS), signer.MainClient)
		if err != nil {
			return nil, err
		}

		calls := make([]multicall3.Multicall3Call3Value, len(to))

		Value := big.NewInt(0)
		for i, account := range to {
			calls[i] = multicall3.Multicall3Call3Value{
				Target:       *account,
				AllowFailure: true,
				Value:        amountPerAccount[i],
				CallData:     nil,
			}

			Value.Add(Value, amountPerAccount[i])
		}

		signer.SignerOpt.GasPrice = gasPrice
		signer.SignerOpt.GasLimit = gasLimit
		signer.SignerOpt.Value = Value

		return Multicall3.Aggregate3Value(signer.SignerOpt, calls)
	}

	return txFunc
}
