package simulate

import (
	"bufio"
	"context"
	"crypto/ecdsa"
	"log"
	"os"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

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
	Start      time.Time
	Ctx        context.Context
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

		// if len(simulateContext.Address) >= com.MAX_ACCOUNT_COUNT {
		// 	break
		// }
	}

	simulateContext.MainClient = mainClient
	simulateContext.Ctx = context.Background()
	simulateContext.Total = config.GlobalConfig.SimulateOptions.Total

	return &simulateContext
}
