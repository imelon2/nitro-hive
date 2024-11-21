package simulate

import (
	"bufio"
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"

	"github.com/imelon2/nitro-hive/common/path"
	"github.com/imelon2/nitro-hive/common/utils"
	config "github.com/imelon2/nitro-hive/config"
)

type SimulateContext struct {
	MainClient *ethclient.Client
	Address    []*common.Address
	PrivateKey []*ecdsa.PrivateKey
	Total      int
	Wait       sync.WaitGroup
	FailCount  int
	Ctx        context.Context
	Progress   *mpb.Progress
}

func NewSimulateContext() *SimulateContext {
	simulateContext := SimulateContext{}

	mainClient, err := ethclient.Dial(config.GlobalConfig.Providers.Main)
	if err != nil {
		log.Fatalf("main client: %v", err)
	}

	privateKeyFilePath := path.PrivateKeyPath()
	privateKeyFile, err := os.Open(privateKeyFilePath)
	if err != nil {
		log.Fatalf("Failed to open private key file: %v", err)
	}
	defer privateKeyFile.Close()

	scanner := bufio.NewScanner(privateKeyFile)
	for scanner.Scan() {
		line := scanner.Text()
		pk := utils.Unhexlify(line)
		key, err := crypto.HexToECDSA(pk)
		if err != nil {
			log.Fatalf("Failed to HexToECDSA : %v", err)
		}

		publicKey := key.Public()
		publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)
		address := crypto.PubkeyToAddress(*publicKeyECDSA)
		simulateContext.Address = append(simulateContext.Address, &address)
		simulateContext.PrivateKey = append(simulateContext.PrivateKey, key)
	}

	simulateContext.MainClient = mainClient
	simulateContext.Ctx = context.Background()
	// simulateContext.Total = config.GlobalConfig.SimulateOptions.Total

	// for multi log
	simulateContext.Progress = nil
	if config.GlobalConfig.CommonOptions.ProgressLog {
		simulateContext.Progress = mpb.New(mpb.WithWaitGroup(&simulateContext.Wait), mpb.WithWidth(40))
	}

	return &simulateContext
}

/**
*	@perNow per task start time
*	@task last task duration
*	@taskAverage all task time of average
 */
func (context *SimulateContext) AddProgress(index int, total int64, perNow *time.Time, task *time.Duration, taskAverage *time.Duration) *mpb.Bar {
	bar := context.Progress.AddBar(total,
		mpb.PrependDecorators(
			decor.Name(fmt.Sprintf("Signer#%d: ", index)),
			decor.Name(" ("),
			decor.Counters("", ""),
			decor.Name(")"),
		),
		mpb.AppendDecorators(
			decor.Percentage(),
			decor.Name(" |"),
			decor.Any(func(s decor.Statistics) string {
				if s.Completed {
					return fmt.Sprintf(" %0.3fs", task.Seconds())
				}

				taskDuration := time.Since(*perNow)
				elapsed := taskDuration.Seconds()
				return fmt.Sprintf(" %0.3fs", elapsed)
			}),
			decor.Any(func(s decor.Statistics) string {
				if s.Completed {
					return fmt.Sprintf(" | %.3fs/opt", taskAverage.Seconds()/float64(s.Current))
				}

				if s.Current == 0 {
					return " | 0.00s/pot"
				}

				avgTimePerTask := taskAverage.Seconds() / float64(s.Current)
				return fmt.Sprintf(" | %.3fs/opt", avgTimePerTask)
			}),
			decor.Name(" | "),
			decor.Elapsed(decor.ET_STYLE_MMSS), // progressed time
			decor.Name(" | ETA: "),
			decor.OnComplete(decor.AverageETA(decor.ET_STYLE_GO), "DONE"), // Average Estimated Time of Arrival
		),
	)
	return bar
}
