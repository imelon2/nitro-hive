package simulate

import (
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/imelon2/nitro-hive/transaction"
)

type SimulateSigner struct {
	TxFunc []func() (*types.Transaction, error)
	Signer *transaction.SignerContext
}

func (context *SimulateContext) SimulateWaitOne(txFunc func() (*types.Transaction, error)) {
	tx, err := txFunc()
	if err != nil {
		log.Fatalf("txFunc: %v", err)
	}

	receipt, err := bind.WaitMined(context.Ctx, context.MainClient, tx)
	if err != nil {
		log.Fatalf("WaitMined: %v", err)
	}
	// fmt.Printf("HASH : %s\n", tx.Hash().Hex())
	fmt.Printf("HASH : %s\n", receipt.TxHash.Hex())
}

func (context *SimulateContext) SimulateWait(simulate *[]SimulateSigner) {
	for signerIndex, s := range *simulate {
		context.Wait.Add(1)

		now := time.Now()
		bar := context.AddProgress(signerIndex, int64(len(s.TxFunc)), now, s.Signer.Task, s.Signer.TaskAvergage)
		for i := 0; i < len(s.TxFunc); i++ {
			func(i int) {
				defer context.Wait.Done()
				tx, err := s.TxFunc[i]()
				if err != nil {
					log.Fatalf("txFunc: %v", err)
				}
				_, err = bind.WaitMined(context.Ctx, context.MainClient, tx)
				if err != nil {
					log.Fatalf("WaitMined: %v", err)
				}
				bar.EwmaIncrement(time.Since(now))
			}(i)
		}
		context.Progress.Wait()
	}
}

func (context *SimulateContext) SimulateWithThread(simulate *[]SimulateSigner) {
	for signerIndex, s := range *simulate {
		context.Wait.Add(1)

		now := time.Now()
		bar := context.AddProgress(signerIndex, int64(len(s.TxFunc)), now, s.Signer.Task, s.Signer.TaskAvergage)
		go func() {
			defer context.Wait.Done()
			for i := 0; i < len(s.TxFunc); i++ {
				_, err := s.TxFunc[i]()
				if err != nil {
					log.Fatalf("txFunc: %v", err)
				}

				bar.EwmaIncrement(time.Since(now))
			}
		}()
	}
	context.Progress.Wait()
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
