package simulate

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func (context *SimulateContext) Simulate(txFunc func(*common.Address) (*types.Transaction, error)) {
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
