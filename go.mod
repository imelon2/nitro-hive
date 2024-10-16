module github.com/imelon2/arbload

go 1.23.1

replace github.com/ethereum/go-ethereum v1.13.13 => ./go-ethereum

require (
	github.com/ethereum/go-ethereum v1.13.13
	github.com/spf13/cobra v1.8.1
)

require (
	github.com/btcsuite/btcd/btcec/v2 v2.2.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.1 // indirect
	github.com/holiman/uint256 v1.2.4 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/crypto v0.17.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
)
