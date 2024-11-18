package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

// func main() {
// 	var total int64 = 100
// 	var taskStartTime time.Time
// 	var firstTaskStartTime time.Time

// 	p := mpb.New(
// 		mpb.WithWidth(40),
// 	)
// 	bar := p.New(total,
// 		mpb.BarStyle().Rbound("|"),
// 		mpb.PrependDecorators(
// 			decor.Name("Signer#1: "),
// 			decor.Name(" ("),
// 			decor.Counters("", ""),
// 			decor.Name(")"),
// 		),
// 		mpb.AppendDecorators(
// 			decor.Percentage(),
// 			decor.Name(" ]"),
// 			decor.Any(func(s decor.Statistics) string {
// 				elapsed := time.Since(taskStartTime).Seconds()
// 				return fmt.Sprintf(" %0.2fs", elapsed)
// 			}),
// 			decor.Any(func(s decor.Statistics) string {
// 				if s.Current == 0 {
// 					return " | 0.00s/pot"
// 				}
// 				totalElapsed := time.Since(firstTaskStartTime).Seconds()

// 				avgTimePerTask := totalElapsed / float64(s.Current)

// 				return fmt.Sprintf(" | %.2fs/opt", avgTimePerTask)
// 			}),
// 			decor.Name(" | ETA: "),
// 			decor.OnComplete(decor.AverageETA(decor.ET_STYLE_GO), "DONE"), // Average Estimated Time of Arrival
// 		),
// 	)

// 	firstTaskStartTime = time.Now()
// 	for i := 0; i < int(total); i++ {
// 		// start := time.Now()
// 		// time.Sleep(time.Second / 1) // 100ms 단위로 진행
// 		// fmt.Printf("Execution %d time: %v\n", i, duration)

// 		// if i == 5 {
// 		// 	time.Sleep(time.Second * 5) // 100ms 단위로 진행
// 		// } else {

// 		// }
// 		taskStartTime = time.Now()
// 		time.Sleep(1200 * time.Millisecond)
// 		// duration := time.Since(start)
// 		// lastIncrementTime += duration.Microseconds()
// 		bar.Increment()

// 		// cur++
// 	}

// 	p.Wait()
// }

func main() {
	var wg sync.WaitGroup
	// passed wg will be accounted at p.Wait() call
	p := mpb.New(mpb.WithWaitGroup(&wg))
	total, numBars := 100, 3
	wg.Add(numBars)

	for i := 0; i < numBars; i++ {
		name := fmt.Sprintf("Bar#%d:", i)
		bar := p.AddBar(int64(total),
			mpb.PrependDecorators(
				// simple name decorator
				decor.Name(name),
				// decor.DSyncWidth bit enables column width synchronization
				decor.Percentage(decor.WCSyncSpace),
			),
			mpb.AppendDecorators(
				// replace ETA decorator with "done" message, OnComplete event
				decor.OnComplete(
					// ETA decorator with ewma age of 30
					decor.EwmaETA(decor.ET_STYLE_GO, 30, decor.WCSyncWidth), "done",
				),
			),
		)
		// simulating some work
		go func() {
			defer wg.Done()
			rng := rand.New(rand.NewSource(time.Now().UnixNano()))
			max := 100 * time.Millisecond
			for i := 0; i < total; i++ {
				// start variable is solely for EWMA calculation
				// EWMA's unit of measure is an iteration's duration
				start := time.Now()
				time.Sleep(time.Duration(rng.Intn(10)+1) * max / 10)
				// we need to call EwmaIncrement to fulfill ewma decorator's contract
				bar.EwmaIncrement(time.Since(start))
			}
		}()
	}
	// wait for passed wg and for all bars to complete and flush
	p.Wait()
}
