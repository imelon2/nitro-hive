package simulate

import (
	"context"
	"log"
	"math/big"
)

func (signer *SignerContext) UpdateNonce() {
	nonce, err := signer.MainClient.PendingNonceAt(context.Background(), *signer.Account)
	if err != nil {
		log.Fatalf("fail update nonce: %v", err)
	}
	signer.SignerOpt.Nonce = big.NewInt(int64(nonce))
}
