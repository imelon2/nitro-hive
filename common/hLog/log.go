package hlog

import (
	"fmt"
	"math/big"
	"time"

	"github.com/imelon2/nitro-hive/common"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

type DistributeIntroLogParam struct {
	TotalAccount int
	PerAmount    int64
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

func NewSignerProgress(index int, total int64) (*mpb.Bar, *mpb.Progress) {
	var taskStartTime time.Time
	var firstTaskStartTime time.Time

	p := mpb.New(
		mpb.WithWidth(40),
	)

	bar := p.New(total,
		mpb.BarStyle().Rbound("|"),
		mpb.PrependDecorators(
			decor.Name(fmt.Sprintf("Signer#%d: ", index)),
			decor.Name(" ("),
			decor.Counters("", ""),
			decor.Name(")"),
		),
		mpb.AppendDecorators(
			decor.Percentage(),
			decor.Name(" ]"),
			decor.Any(func(s decor.Statistics) string {
				elapsed := time.Since(taskStartTime).Seconds()
				return fmt.Sprintf(" %0.2fs", elapsed)
			}),
			decor.Any(func(s decor.Statistics) string {
				if s.Current == 0 {
					return " | 0.00s/pot"
				}
				totalElapsed := time.Since(firstTaskStartTime).Seconds()

				avgTimePerTask := totalElapsed / float64(s.Current)

				return fmt.Sprintf(" | %.2fs/opt", avgTimePerTask)
			}),
			decor.Name(" | ETA: "),
			decor.OnComplete(decor.AverageETA(decor.ET_STYLE_GO), "DONE"), // Average Estimated Time of Arrival
		),
	)

	firstTaskStartTime = time.Now()
	taskStartTime = time.Now()

	return bar, p
}
