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
	Start      *time.Time
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
	simulateContext.Progress = mpb.New(mpb.WithWaitGroup(&simulateContext.Wait), mpb.WithWidth(40))

	return &simulateContext
}

func (context *SimulateContext) AddProgress(index int, total int64, initTime time.Time, task time.Duration, taskAvergage time.Duration) *mpb.Bar {
	bar := context.Progress.AddBar(total,
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
				if s.Completed {
					return fmt.Sprintf(" %0.2fs", task.Seconds())
				}
				task = time.Since(initTime)

				elapsed := task.Seconds()
				return fmt.Sprintf(" %0.2fs", elapsed)
			}),
			decor.Any(func(s decor.Statistics) string {
				if s.Completed {
					return fmt.Sprintf(" | %.2fs/opt", taskAvergage.Seconds()/float64(s.Current))
				}

				taskAvergage = time.Since(initTime)
				if s.Current == 0 {
					return " | 0.00s/pot"
				}

				avgTimePerTask := taskAvergage.Seconds() / float64(s.Current)
				return fmt.Sprintf(" | %.2fs/opt", avgTimePerTask)
			}),
			decor.Name(" | ETA: "),
			decor.OnComplete(decor.AverageETA(decor.ET_STYLE_GO), "DONE"), // Average Estimated Time of Arrival
		),
	)
	return bar
}
