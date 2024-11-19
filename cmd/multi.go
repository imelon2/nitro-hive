/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	hlog "github.com/imelon2/nitro-hive/common/hLog"
	"github.com/imelon2/nitro-hive/common/utils"
	"github.com/imelon2/nitro-hive/config"
	"github.com/imelon2/nitro-hive/simulate"
	"github.com/imelon2/nitro-hive/transaction"
	"github.com/spf13/cobra"
)

// multiCmd represents the multi command
var multiCmd = &cobra.Command{
	Use:   "multi",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		value := big.NewInt(int64(config.GlobalConfig.MultiOptions.TransactionOptions.Value))
		gasPrice := big.NewInt(int64(config.GlobalConfig.MultiOptions.TransactionOptions.GasPrice))
		gasLimit := uint64(config.GlobalConfig.MultiOptions.TransactionOptions.Gas)
		data := common.Hex2Bytes(utils.Unhexlify(config.GlobalConfig.MultiOptions.TransactionOptions.Data))
		to := common.HexToAddress(config.GlobalConfig.MultiOptions.TransactionOptions.To)
		perTx := config.GlobalConfig.MultiOptions.PerTx

		signerIndex := config.GlobalConfig.MultiOptions.PrivateKeyRange.StartIndex
		signerCount := config.GlobalConfig.MultiOptions.PrivateKeyRange.Total

		simulation := simulate.NewSimulateContext()
		signersPk := (simulation.PrivateKey)[signerIndex : signerIndex+signerCount]

		simulateSigners := make([]simulate.SimulateSigner, 0)
		for _, key := range signersPk {
			signer, err := transaction.NewSginerContext(key)
			if err != nil {
				log.Fatal(err)
			}
			txFuncs := make([]func() (*types.Transaction, error), 0)
			for i := 0; i < perTx; i++ {
				txFunc := signer.TransaferLegacyTx(&to, gasPrice, gasLimit, data, value)
				txFuncs = append(txFuncs, txFunc)
			}

			simulateSigners = append(simulateSigners, simulate.SimulateSigner{
				Signer: signer,
				TxFunc: txFuncs,
			})
		}

		// intro console log
		hlog.SimulateLog(simulateSigners)
		SubStart := time.Now()
		simulation.SimulateWithThread(&simulateSigners)
		Subduration := time.Since(SubStart) // 종료 시점에서 경과 시간 계산
		fmt.Printf("\n\nExecution distribute to sub account simualter time: %v\n\n", Subduration)
	},
}

func init() {
	rootCmd.AddCommand(multiCmd)
}
