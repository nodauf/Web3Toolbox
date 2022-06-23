package ethcmd

/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
	"log"
)

var number int

// downloadSourcesCmd represents the downloadSources command
var newAccountCmd = &cobra.Command{
	Use:   "newAccount",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		for i := 0; i < number; i++ {
			privateKey, err := crypto.GenerateKey()
			if err != nil {
				log.Fatal(err)
			}
			publicKeyECDSA, ok := privateKey.Public().(*ecdsa.PublicKey)
			if !ok {
				log.Fatal("error casting public key to ECDSA")
			}

			address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
			fmt.Printf("Address: %s\n", address)
			fmt.Printf("Private key: %s\n", hexutil.Encode(crypto.FromECDSA(privateKey))[2:])
			if i != number-1 {
				fmt.Println("-----------------------------------------------------")
			}

		}
	},
}

func init() {
	EthCmd.AddCommand(newAccountCmd)

	newAccountCmd.Flags().IntVarP(&number, "number", "n", 1, "number of accounts to create")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downloadSourcesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downloadSourcesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
