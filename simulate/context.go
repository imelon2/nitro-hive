package simulate

import (
	"bufio"
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"
	"os"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/imelon2/nitro-hive/common/path"
	"github.com/imelon2/nitro-hive/common/utils"
)

type SimulateContext struct {
	PrivateKey   []*ecdsa.PrivateKey
	ParentClient *ethclient.Client
	SignerOpt    *bind.TransactOpts
	NonceMutex   *sync.Mutex
}

type SignerContext struct {
	MainClient   *ethclient.Client
	ParentClient *ethclient.Client
	SignerOpt    *bind.TransactOpts
	NonceMutex   *sync.Mutex
}

func NewSimulateContext(count int) *SimulateContext {
	// func NewSimulateContext() (*SimulateContext, error) {
	context := SimulateContext{}

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
		context.PrivateKey = append(context.PrivateKey, key)
		if len(context.PrivateKey) >= count {
			break
		}
	}

	return &context
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

	gasPrice, _ := mainClient.SuggestGasPrice(context.Background())
	opt.GasPrice = gasPrice

	// chain.GasLimit = config.GasLimit => Estimate

	nonce, err := mainClient.PendingNonceAt(context.Background(), address)
	if err != nil {
		log.Fatalf("Nonce: %v", err)
	}
	opt.Nonce = big.NewInt(int64(nonce))

	return &SignerContext{
		MainClient:   mainClient,
		ParentClient: parentClient,
		SignerOpt:    opt,
		NonceMutex:   new(sync.Mutex),
	}, nil
}
