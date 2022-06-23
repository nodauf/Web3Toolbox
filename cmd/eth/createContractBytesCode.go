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
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	utilsCreateContractBytesCode "github.com/nodauf/web3Toolbox/utilsCmd/createContractBytesCode"
	"github.com/spf13/cobra"
	"log"
	"strings"
)

var contractBytesCode string

// downloadSourcesCmd represents the downloadSources command
var createContractBytesCodeCmd = &cobra.Command{
	Use:   "createContractBytesCode",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := ethclient.Dial(blockchainURL)
		if err != nil {
			log.Fatal(err)
		}
		if strings.HasPrefix(privateKey, "0x") {
			privateKey = privateKey[2:]
		}

		if strings.HasPrefix(contractBytesCode, "0x") {
			contractBytesCode = contractBytesCode[2:]
		}

		signedTx := utilsCreateContractBytesCode.CreateContract(client, privateKey, common.Hex2Bytes(contractBytesCode))
		fmt.Printf("Tx sent: %s\n", signedTx.Hash().Hex())
		fmt.Println("Waiting for confirmation...")
		receipt, err := bind.WaitMined(context.Background(), client, signedTx)
		if err != nil {
			log.Fatalf("Tx failed: %s", err.Error())
		}
		var status string
		if receipt.Status == 0 {
			status = "failed ❌"
		} else if receipt.Status == 1 {
			status = "success ✔"
		} else {
			status = "unknown"
		}
		fmt.Printf("Contract address: %s\nStatus: %s\n", receipt.ContractAddress, status)
	},
}

func init() {
	EthCmd.AddCommand(createContractBytesCodeCmd)
	createContractBytesCodeCmd.Flags().StringVarP(&privateKey, "privateKey", "k", "", "Private key")
	createContractBytesCodeCmd.Flags().StringVarP(&blockchainURL, "blockchainURL", "b", "http://localhost:8545", "Blockchain URL")
	createContractBytesCodeCmd.Flags().StringVarP(&contractBytesCode, "contractBytesCode", "c", "", "Contract bytes code")

	createContractBytesCodeCmd.MarkFlagRequired("privateKey")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downloadSourcesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downloadSourcesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
