package simulate

import (
	"log"
	"math/big"
)

func (signer *SignerContext) UpdateNonce() {
	nonce, err := signer.MainClient.PendingNonceAt(signer.Ctx, *signer.Account)
	if err != nil {
		log.Fatalf("fail update nonce: %v", err)
	}
	signer.SignerOpt.Nonce = big.NewInt(int64(nonce))
}
