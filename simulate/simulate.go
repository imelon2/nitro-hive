package simulate

import (
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

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

func (context *SimulateContext) SimulateWait(txFuncs []func() (*types.Transaction, error)) {
	for i := 0; i < len(txFuncs); i++ {
		func(i int) {
			context.Start = time.Now()
			defer func() { // 함수 종료 시 실행될 작업 정의
				duration := time.Since(context.Start) // 종료 시점에서 경과 시간 계산
				fmt.Printf("Execution time: %v\n", duration)
			}()

			tx, err := txFuncs[i]()
			if err != nil {

			}
			receipt, err := bind.WaitMined(context.Ctx, context.MainClient, tx)
			if err != nil {
				log.Fatalf("WaitMined: %v", err)
			}
			fmt.Printf("HASH : %s\n", receipt.TxHash.Hex())
		}(i)
	}
}

func (context *SimulateContext) SimulateWithThread(txFuncs []func() (*types.Transaction, error)) {
	for i := 0; i < len(txFuncs); i++ {
		context.Wait.Add(1)
		go func(i int) {
			defer context.Wait.Done()
			var start *time.Time
			_start := time.Now()
			start = &_start

			fmt.Printf("Start Send %d Tx at %s \n", i, start.Format("2006-01-02 15:04:05.000"))

			defer func() { // 함수 종료 시 실행될 작업 정의
				duration := time.Since(*start) // 종료 시점에서 경과 시간 계산
				fmt.Printf("Execution %d time: %v\n", i, duration)
			}()

			tx, err := txFuncs[i]()
			if err != nil {
				log.Fatalf("txFuncs: %v", err)
			}
			fmt.Printf("HASH : %s\n", tx.Hash().Hex())
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
