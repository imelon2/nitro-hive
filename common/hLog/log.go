package hlog

import (
	"fmt"
	"math/big"

	"github.com/imelon2/nitro-hive/common"
	"github.com/imelon2/nitro-hive/simulate"
	"github.com/pterm/pterm"
)

type DistributeIntroLogParam struct {
	TotalAccount int
	PerAmount    int64
}

func SimulateLog(si []simulate.SimulateSigner) {
	accountCount := 0
	txCount := 0
	for _, simulater := range si {
		accountCount++
		txCount += len(simulater.TxFunc)
	}
	msg := fmt.Sprintf("Total account: %d\n", accountCount)
	msg += fmt.Sprintf("Total transaction count: %d\n", txCount)
	msg += fmt.Sprintf("ㄴ account per transaction count: %d", txCount/accountCount)
	pterm.DefaultBox.Println(msg)
	fmt.Print("\n")
}

func DistributeIntroLog(param DistributeIntroLogParam) {
	total := param.TotalAccount
	sub := total / common.MULTICALL_MAX_TX_COUNT
	remained := total % common.MULTICALL_MAX_TX_COUNT
	if remained != 0 {
		sub++
	}

	subTx := sub / common.MULTICALL_MAX_TX_COUNT
	remainedSubTx := total % common.MULTICALL_MAX_TX_COUNT
	if remainedSubTx != 0 {
		subTx++
	}

	totalValue := big.NewInt(param.PerAmount * int64(total))
	fee := big.NewInt(common.MULTICALL_FEE.Int64() * int64(sub))
	require := totalValue.Add(totalValue, fee)

	fmt.Printf("Total account count: %d\n", param.TotalAccount)
	fmt.Printf("Sub account count: %d | need to execute %d transaction\n", sub, subTx)
	fmt.Printf("distribute amount: %d | multicall fee: %d | requirement amount: %d\n", totalValue, fee, require)

}
