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

	hlog "github.com/imelon2/nitro-hive/common/hLog"
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
	Run: func(cmd *cobra.Command, args []string) {
		isSubA, err := cmd.Flags().GetBool("sub-account")
		if err != nil {
			log.Fatal(err)
		}
		isSiA, err := cmd.Flags().GetBool("simulate-account")
		if err != nil {
			log.Fatal(err)
		}

		value := big.NewInt(int64(config.GlobalConfig.DistributeOptions.Value))
		gasPrice := big.NewInt(int64(config.GlobalConfig.DistributeOptions.GasPrice))
		gasLimit := uint64(config.GlobalConfig.DistributeOptions.Gas)
		multicallMaxCall := c.MULTICALL_MAX_TX_COUNT

		/*
		* The signer sends the sub accounts funds to distribute to the simulate account.
		* Each transaction can send funds to a maximum of 250 simulate accounts.
		* All transfers are done via `Multicall3`, and you can perform up to 250 transfers.
		 */
		if isSubA {
			simulation := simulate.NewSimulateContext()

			simulateAccountCount := len(simulation.Address)
			subAccountCount := simulateAccountCount / multicallMaxCall
			remained := simulateAccountCount % multicallMaxCall
			if remained != 0 {
				subAccountCount++
			}

			pk := utils.Unhexlify(config.GlobalConfig.DistributeOptions.PrivateKey)
			key, _ := crypto.HexToECDSA(pk)
			signer, err := transaction.NewSginerContext(key)
			if err != nil {
				log.Fatal(err)
			}
			subAccount := (simulation.Address)[:subAccountCount]

			subTx := subAccountCount / multicallMaxCall
			if subAccountCount%multicallMaxCall != 0 {
				subTx++
			}

			txFuncs := make([]func() (*types.Transaction, error), 0)

			for len(subAccount) > 0 {
				subAccountReceiver := subAccount
				//  multicall transfer max count := multicallMaxCall
				if len(subAccount) > multicallMaxCall {
					subAccountReceiver = (subAccountReceiver)[:multicallMaxCall]
				}

				amountPerAccount := make([]*big.Int, 0)
				for simulateAccountCount > 0 {
					simulateReceiver := simulateAccountCount
					//  multicall transfer max count := multicallMaxCall
					if simulateReceiver > multicallMaxCall {
						simulateReceiver = multicallMaxCall
					}

					totalValue := big.NewInt(1).Mul(value, big.NewInt(int64(simulateReceiver)))
					totalValue = totalValue.Add(totalValue, c.MULTICALL_FEE)

					amountPerAccount = append(amountPerAccount, totalValue)

					simulateAccountCount -= simulateReceiver

					if len(amountPerAccount) == multicallMaxCall {
						break
					}
				}

				txFunc := signer.Distribute(subAccountReceiver, gasPrice, gasLimit, amountPerAccount)
				txFuncs = append(txFuncs, txFunc)
				subAccount = (subAccount)[len(subAccountReceiver):]
			}

			txSubFuncs := []simulate.SimulateSigner{
				{
					TxFunc: txFuncs,
					Signer: signer,
				},
			}
			// intro console log
			hlog.SimulateLog(txSubFuncs)

			SubStart := time.Now()
			simulation.SimulateWait(&txSubFuncs)
			Subduration := time.Since(SubStart) // 종료 시점에서 경과 시간 계산
			fmt.Printf("\n\nExecution distribute to sub account simualter time: %v\n\n", Subduration)
		}

		/*
		* Each sub account will send `value` per simulate account and sub account can send funds to a maximum of 250 simulate accounts.
		* All transfers are done via `Multicall3`, and you can perform up to 250 transfers.
		 */
		if isSiA {
			simulation := simulate.NewSimulateContext()
			simulateAccountCount := len(simulation.Address)
			subAccountCount := simulateAccountCount / multicallMaxCall

			if simulateAccountCount%multicallMaxCall != 0 {
				subAccountCount++
			}

			txSubFuncs := make([]simulate.SimulateSigner, 0)
			subAccountPk := (simulation.PrivateKey)[:subAccountCount]
			for _, _pk := range subAccountPk {
				subSigner, err := transaction.NewSginerContext(_pk)
				if err != nil {
					log.Fatal(err)
				}

				receivers := len(simulation.Address)
				if receivers > multicallMaxCall {
					receivers = multicallMaxCall
				}

				to := (simulation.Address)[:receivers]
				simulation.Address = (simulation.Address)[receivers:]

				txFunc := subSigner.Distribute(to, gasPrice, gasLimit, utils.FillBigIntArray(receivers, value))
				txSubFuncs = append(txSubFuncs, simulate.SimulateSigner{
					TxFunc: []func() (*types.Transaction, error){
						txFunc,
					},
					Signer: subSigner,
				})
			}

			// intro console log
			hlog.SimulateLog(txSubFuncs)

			SubStart := time.Now()
			simulation.SimulateWithThread(&txSubFuncs)
			Subduration := time.Since(SubStart) // 종료 시점에서 경과 시간 계산
			fmt.Printf("\n\nExecution distribute to simulate accounts simualter time: %v\n\n", Subduration)
		}
	},
}

func init() {
	rootCmd.AddCommand(DistributeCmd)

	DistributeCmd.Flags().BoolP("sub-account", "u", true, "Distribute to sub account funds from signer")
	DistributeCmd.Flags().BoolP("simulate-account", "i", true, "Distribute to simualter account funds from sub accounts")
}
