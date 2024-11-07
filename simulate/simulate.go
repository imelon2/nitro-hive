package simulate

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func (context *SimulateContext) Distribute(txFunc func(*common.Address) (*types.Transaction, error)) {
	start := GlobalConfig.SimulateOptions.MultiSigner.StartIndex
	total := GlobalConfig.SimulateOptions.MultiSigner.Total
	for i := start; i < total; i++ {
		func(i int) {
			tx, err := txFunc(context.Address[i])
			if err != nil {
				log.Fatalf("txFunc: %v", err)
			}

			fmt.Printf("HASH : %s\n", tx.Hash().Hex())
		}(i)
	}
}

func (context *SimulateContext) SimulateWait(txFunc func(*common.Address) (*types.Transaction, error)) {
	for i := 0; i < context.Total; i++ {
		context.Wait.Add(1)
		go func(i int) {
			defer context.Wait.Done()
			tx, err := txFunc(context.Address[i])
			if err != nil {
				log.Fatalf("txFunc: %v", err)
			}

			receipt, err := bind.WaitMined(context.Ctx, context.MainClient, tx)
			if err != nil {
				log.Fatalf("WaitMined: %v", err)
			}
			fmt.Printf("HASH : %s\n", receipt.TxHash.Hex())
		}(i)
	}

	context.Wait.Wait()
}

func Simulate(index int, total int, txFunc func(*common.Address) (*types.Transaction, error)) {
	for i := 0; i < total; i++ {
		func(i int) {
			tx, err := txFunc(nil)
			if err != nil {
				log.Fatalf("txFunc: %v", err)
			}

			fmt.Printf("HASH-%d : %s\n", index, tx.Hash().Hex())
		}(index)
	}
}

func (context *SimulateContext) MultiSimulate(txFuncs []func(*common.Address) (*types.Transaction, error)) {
	signers := len(txFuncs)
	signerPerTx := context.Total / signers
	fmt.Printf("per tx: %d\n", signerPerTx)
	fmt.Printf("total: %d\n", context.Total)
	for i := 0; i < signers; i++ {
		context.Wait.Add(1)
		go func(in int, total int) {
			defer context.Wait.Done()
			Simulate(in, total, txFuncs[in])
		}(i, signerPerTx)
	}

	context.Wait.Wait()
}
