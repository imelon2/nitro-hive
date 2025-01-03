package transaction

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	config "github.com/imelon2/nitro-hive/config"
	"github.com/vbauerster/mpb/v8"
)

type ProgressClass struct {
	Bar      *mpb.Bar
	Progress *mpb.Progress
}

type SignerContext struct {
	MainClient  *ethclient.Client
	Account     *common.Address
	SignerOpt   *bind.TransactOpts
	NonceMutex  *sync.Mutex
	Task        *time.Duration
	TaskAverage *time.Duration
	PerNow      *time.Time
	ProgressBar *mpb.Bar
	Ctx         *context.Context
}

var PROVIDER_INDEX = 0

func getCUrrentProviderUrl() string {
	PROVIDER_INDEX = (PROVIDER_INDEX + 1) % len(config.GlobalConfig.Providers)
	return config.GlobalConfig.Providers[PROVIDER_INDEX]
}

func NewSginerContext(pk *ecdsa.PrivateKey) (*SignerContext, error) {
	mainClient, err := ethclient.Dial(getCUrrentProviderUrl())
	if err != nil {
		log.Fatalf("main client: %v", err)
	}

	mainChainID, err := mainClient.NetworkID(context.Background())
	if err != nil {
		fmt.Printf("here\n\n")
		log.Fatal(err)
	}

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
	ctx := context.Background()
	return &SignerContext{
		MainClient:  mainClient,
		Account:     &address,
		SignerOpt:   opt,
		NonceMutex:  new(sync.Mutex),
		Task:        new(time.Duration),
		TaskAverage: new(time.Duration),
		PerNow:      new(time.Time),
		Ctx:         &ctx,
	}, nil
}

func (signer *SignerContext) UpdateNonce() {
	nonce, err := signer.MainClient.PendingNonceAt(*signer.Ctx, *signer.Account)
	if err != nil {
		log.Fatalf("fail update nonce: %v", err)
	}
	signer.SignerOpt.Nonce = big.NewInt(int64(nonce))
}
