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
	"github.com/ethereum/go-ethereum/crypto"
	hlog "github.com/imelon2/nitro-hive/common/hLog"
	"github.com/imelon2/nitro-hive/common/utils"
	"github.com/imelon2/nitro-hive/config"
	"github.com/imelon2/nitro-hive/simulate"
	"github.com/imelon2/nitro-hive/transaction"
	"github.com/spf13/cobra"
)

// singleCmd represents the single command
var singleCmd = &cobra.Command{
	Use:   "single",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		value := big.NewInt(int64(config.GlobalConfig.SingleOptions.TransactionOptions.Value))
		gasPrice := big.NewInt(int64(config.GlobalConfig.SingleOptions.TransactionOptions.GasPrice))
		gasLimit := uint64(config.GlobalConfig.SingleOptions.TransactionOptions.Gas)
		data := common.Hex2Bytes(utils.Unhexlify(config.GlobalConfig.SingleOptions.TransactionOptions.Data))
		to := common.HexToAddress(config.GlobalConfig.SingleOptions.TransactionOptions.To)
		total := config.GlobalConfig.SingleOptions.PerTx

		simulation := simulate.NewSimulateContext()

		pk := utils.Unhexlify(config.GlobalConfig.SingleOptions.PrivateKey)
		key, _ := crypto.HexToECDSA(pk)
		signer, err := transaction.NewSginerContext(key)
		if err != nil {
			log.Fatal(err)
		}

		txFuncs := make([]func() (*types.Transaction, error), 0)
		for i := 0; i < total; i++ {
			txFunc := signer.TransaferLegacyTx(&to, gasPrice, gasLimit, data, value)
			txFuncs = append(txFuncs, txFunc)
		}

		simulateSigner := []simulate.SimulateSigner{
			{
				TxFunc: txFuncs,
				Signer: signer,
			},
		}
		// intro console log
		hlog.SimulateLog(simulateSigner)

		SubStart := time.Now()
		simulation.SimulateWait(&simulateSigner)
		Subduration := time.Since(SubStart) // 종료 시점에서 경과 시간 계산
		fmt.Printf("\n\nExecution distribute to sub account simualter time: %v\n\n", Subduration)
	},
}

func init() {
	rootCmd.AddCommand(singleCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// singleCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// singleCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
