package simulate

import (
	"bufio"
	"context"
	"crypto/ecdsa"
	"log"
	"os"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	com "github.com/imelon2/nitro-hive/common"
	"github.com/imelon2/nitro-hive/common/path"
	"github.com/imelon2/nitro-hive/common/utils"
)

type SimulateContext struct {
	MainClient *ethclient.Client
	// PrivateKey []*ecdsa.PrivateKey
	Address   []*common.Address
	Total     int
	Wait      sync.WaitGroup
	FailCount int
	Ctx       context.Context
}

type SignerContext struct {
	MainClient   *ethclient.Client
	ParentClient *ethclient.Client
	Account      *common.Address
	SignerOpt    *bind.TransactOpts
	NonceMutex   *sync.Mutex
	Ctx          context.Context
}

func NewSimulateContext() *SimulateContext {
	// func NewSimulateContext() (*SimulateContext, error) {
	simulateContext := SimulateContext{}

	mainClient, err := ethclient.Dial(GlobalConfig.Providers.Main)
	if err != nil {
		log.Fatalf("main client: %v", err)
	}
	simulateContext.MainClient = mainClient

	privateKeyFilePath := path.PrivateKeyPath()
	file, err := os.Open(privateKeyFilePath)
	if err != nil {
		log.Fatalf("Failed to open private key file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
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

		if len(simulateContext.Address) >= com.MAX_ACCOUNT_COUNT {
			break
		}
	}

	//@ TODO : Multi call
	// simulateContext.PrivateKey = append(simulateContext.PrivateKey, key)
	// index := GlobalConfig.SimulateOption.AccountRange.StartIndex

	// for scanner.Scan() {
	// 	if index != 0 {
	// 		index--
	// 		continue
	// 	}

	// 	line := scanner.Text()
	// 	pk := utils.Unhexlify(line)
	// 	key, err := crypto.HexToECDSA(pk)
	// 	if err != nil {
	// 		log.Fatalf("Failed to HexToECDSA : %v", err)
	// 	}
	// 	simulateContext.PrivateKey = append(simulateContext.PrivateKey, key)
	// 	publicKey := key.Public()

	// 	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)
	// 	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	// 	simulateContext.Address = append(simulateContext.Address, &address)
	// 	simulateContext.Ctx = context.Background()

	// 	if len(simulateContext.PrivateKey) >= GlobalConfig.SimulateOption.AccountRange.Total {
	// 		break
	// 	}
	// }

	simulateContext.Ctx = context.Background()
	simulateContext.Total = GlobalConfig.SimulateOption.Total

	return &simulateContext
}

func NewSginerContext(pk *ecdsa.PrivateKey) (*SignerContext, error) {

	mainClient, err := ethclient.Dial(GlobalConfig.Providers.Main)
	if err != nil {
		log.Fatalf("main client: %v", err)
	}

	parentClient, err := ethclient.Dial(GlobalConfig.Providers.Parent)
	if err != nil {
		log.Fatalf("parent client: %v", err)
	}

	mainChainID, err := mainClient.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// pk := utils.Unhexlify(GlobalConfig.Signer.PrivateKey)
	// key, _ := crypto.HexToECDSA(pk)

	publicKey := pk.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA)

	opt, err := bind.NewKeyedTransactorWithChainID(pk, mainChainID)
	if err != nil {
		log.Fatalf("NewKeyedTransactorWithChainID: %s", err)
	}

	// gasPrice, _ := mainClient.SuggestGasPrice(context.Background())
	// opt.GasPrice = gasPrice

	// chain.GasLimit = config.GasLimit => Estimate

	// nonce, err := mainClient.PendingNonceAt(context.Background(), address)
	// if err != nil {
	// 	log.Fatalf("Nonce: %v", err)
	// }
	// opt.Nonce = big.NewInt(int64(nonce))

	return &SignerContext{
		MainClient:   mainClient,
		ParentClient: parentClient,
		Account:      &address,
		SignerOpt:    opt,
		NonceMutex:   new(sync.Mutex),
		Ctx:          context.Background(),
	}, nil
}
