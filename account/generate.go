package main

import (
	"crypto/ecdsa"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	constant "github.com/imelon2/nitro-hive/common"
	"github.com/imelon2/nitro-hive/common/path"
)

func main() {
	accountFilePath := path.AccountPath()
	privateKeyFilePath := path.PrivateKeyPath()

	accountFile, err := os.Create(accountFilePath)
	if err != nil {
		log.Fatal(err)
	}
	privateKeyFile, err := os.Create(privateKeyFilePath)
	if err != nil {
		log.Fatal(err)
	}

	defer accountFile.Close()
	defer privateKeyFile.Close()

	for i := 0; i < constant.MAX_ACCOUNT_COUNT; i++ {
		privateKey, err := crypto.GenerateKey()
		if err != nil {
			log.Fatal(err)
		}

		privateKeyBytes := crypto.FromECDSA(privateKey)
		privateKeyHex := hexutil.Encode(privateKeyBytes)
		publicKey := privateKey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			log.Fatal("Error casting public key to ECDSA")
		}

		address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

		if i < constant.MAX_ACCOUNT_COUNT-1 {
			_, err = accountFile.WriteString(fmt.Sprintf("%s\n", address))
		}
		if err != nil {
			log.Fatal(err)
		}

		if i < constant.MAX_ACCOUNT_COUNT-1 {
			_, err = privateKeyFile.WriteString(fmt.Sprintf("%s\n", privateKeyHex))
		}

		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println(constant.MAX_ACCOUNT_COUNT, "addresses and private key generated and saved to account_100k, privateKey_100k")
}
