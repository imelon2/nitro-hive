/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/imelon2/nitro-hive/account"
	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		count, err := cmd.Flags().GetInt("count")
		if err != nil {
			log.Fatal(err)
		}

		account.Generate(count)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().IntP("count", "c", 10000, "Number of accounts")
}
