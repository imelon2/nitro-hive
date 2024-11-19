package main

import (
	"fmt"
	"time"

	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

func main() {
	var total int64 = 100
	var taskStartTime time.Time
	var firstTaskStartTime time.Time

	p := mpb.New(
		mpb.WithWidth(40),
	)
	bar := p.New(total,
		mpb.BarStyle().Rbound("|"),
		mpb.PrependDecorators(
			decor.Name("Signer#1: "),
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
	for i := 0; i < int(total); i++ {
		// start := time.Now()
		// time.Sleep(time.Second / 1) // 100ms 단위로 진행
		// fmt.Printf("Execution %d time: %v\n", i, duration)

		// if i == 5 {
		// 	time.Sleep(time.Second * 5) // 100ms 단위로 진행
		// } else {

		// }
		taskStartTime = time.Now()
		time.Sleep(1200 * time.Millisecond)
		// duration := time.Since(start)
		// lastIncrementTime += duration.Microseconds()
		bar.Increment()

		// cur++
	}

	p.Wait()
}
