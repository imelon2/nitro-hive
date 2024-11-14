/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	c "github.com/imelon2/nitro-hive/common"

	"github.com/imelon2/nitro-hive/common/utils"
	"github.com/imelon2/nitro-hive/config"
	"github.com/imelon2/nitro-hive/simulate"
	"github.com/imelon2/nitro-hive/transaction"
	"github.com/spf13/cobra"
)

// DistributeCmd represents the Distribute command
var DistributeCmd = &cobra.Command{
	Use:   "distribute",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		value, err := cmd.Flags().GetInt64("value")
		if err != nil {
			log.Fatal(err)
		}
		gasPrice, err := cmd.Flags().GetInt64("gasPrice")
		if err != nil {
			log.Fatal(err)
		}
		gasLimit, err := cmd.Flags().GetUint64("gasLimit")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("value per address: %d\n", value)
		fmt.Printf("gasPrice per tx: %d\n", gasPrice)
		fmt.Printf("gasLimit per tx: %d\n", gasLimit)

		simulation := simulate.NewSimulateContext()

		pk := utils.Unhexlify(config.GlobalConfig.SimulateOptions.SingleSigner.PrivateKey)
		key, _ := crypto.HexToECDSA(pk)
		signer, err := transaction.NewSginerContext(key)
		if err != nil {
			log.Fatal(err)
		}

		multicallMaxCall := c.MULTICALL_MAX_TX_COUNT
		simulateAccount := len(simulation.Address)
		subAccountCount := simulateAccount / multicallMaxCall
		remained := simulateAccount % multicallMaxCall
		if remained != 0 {
			subAccountCount++
		}
		fmt.Printf("Send Multicall Tx Total %d count | remained %d \n\n", subAccountCount, remained)

		subAccount := (simulation.Address)[:subAccountCount]

		for len(subAccount) > 0 {
			receiver := subAccount
			if len(subAccount) > multicallMaxCall {
				receiver = (receiver)[:multicallMaxCall]
			}

			fmt.Printf("Sub Account: %d\n", len(receiver))
			amountPerAccount := make([]*big.Int, 0)
			totalValue := big.NewInt(value * int64(multicallMaxCall))
			totalValue = totalValue.Add(totalValue, c.MULTICALL_FEE)

			amountPerAccount = append(amountPerAccount, totalValue)

			for simulateAccount >= multicallMaxCall {

				txFuncs := make([]func() (*types.Transaction, error), 0)
				txFunc := signer.Distribute(receiver, big.NewInt(gasPrice), gasLimit, amountPerAccount)
				txFuncs = append(txFuncs, txFunc)

				Start := time.Now()
				simulation.SimulateWait(txFuncs)
				duration := time.Since(Start) // 종료 시점에서 경과 시간 계산
				fmt.Printf("\n\nExecution SimulateWait time: %v\n", duration)

				simulateAccount -= multicallMaxCall
			}

			subAccount = (subAccount)[len(receiver):]
			fmt.Print("RUN\n")
		}

		return

		subAccountPk := (simulation.PrivateKey)[:subAccountCount]
		txSubFuncs := make([]func() (*types.Transaction, error), 0)
		for _, _pk := range subAccountPk {
			subSigner, err := transaction.NewSginerContext(_pk)
			if err != nil {
				log.Fatal(err)
			}

			receivers := len(simulation.Address)
			if receivers > multicallMaxCall {
				receivers = multicallMaxCall // 마지막에 남은 요소가 chunkSize보다 적은 경우
			}
			// 현재 추출할 chunk 설정
			to := (simulation.Address)[:receivers]
			simulation.Address = (simulation.Address)[receivers:]

			txFunc := subSigner.Distribute(to, big.NewInt(gasPrice), gasLimit, utils.FillBigIntArray(receivers, big.NewInt(value)))
			txSubFuncs = append(txSubFuncs, txFunc)
		}

		SubStart := time.Now()
		simulation.SimulateWithThread(txSubFuncs)
		Subduration := time.Since(SubStart) // 종료 시점에서 경과 시간 계산
		fmt.Printf("\n\nExecution SimulateWait time: %v\n", Subduration)
	},
}

func init() {
	rootCmd.AddCommand(DistributeCmd)

	DistributeCmd.Flags().Int64P("value", "", 10000000000000000, "Number of node created data to retrieve(DEFAULT 1 ETH)")
	DistributeCmd.Flags().Int64P("gasPrice", "", 300000000, "Number of node created data to retrieve(0.5 GWEI)")
	DistributeCmd.Flags().Uint64P("gasLimit", "", 0, "Number of node created data to retrieve")
}
