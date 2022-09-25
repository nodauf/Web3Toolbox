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
	"encoding/json"
	"fmt"
	"github.com/nodauf/web3Toolbox/utilsCmd/storageChanged"
	"github.com/spf13/cobra"
	"log"
)

type debugTx struct {
	ID      int64  `json:"id"`
	JsonRPC string `json:"jsonrpc"`
	Result  struct {
		Gas         int64    `json:"gas"`
		ReturnValue string   `json:"returnValue"`
		Storage     struct{} `json:"storage"`
		StructLogs  []struct {
			Depth   int64             `json:"depth"`
			Error   string            `json:"error"`
			Gas     int64             `json:"gas"`
			GasCost int64             `json:"gasCost"`
			Memory  []string          `json:"memory"`
			Op      string            `json:"op"`
			Pc      int64             `json:"pc"`
			Stack   []string          `json:"stack"`
			Storage map[string]string `json:"storage"`
		} `json:"structLogs"`
	} `json:"result"`
}

// downloadSourcesCmd represents the downloadSources command
var storageChangedCmd = &cobra.Command{
	Use:   "storageChanged transactionID",
	Short: "Detect changes in smart contract storage for a specific TX",
	Long:  `Use the call rpc debug_traceTransaction to trace the execution of the TX. Only endpoint that supports this method can be used. Ganache is one of them.`,
	Run: func(cmd *cobra.Command, args []string) {
		res, err := storageChanged.DebugTransaction(blockchainURL, args[0])
		if err != nil {
			log.Fatal(err)
		}
		var resDebugTx debugTx
		err = json.Unmarshal(res, &resDebugTx)
		if err != nil {
			log.Fatal(err)
		}
		var initialStorage = make(map[string]string)
		var finalStorage = make(map[string]string)
		for _, data := range resDebugTx.Result.StructLogs {
			if data.Op == "SSTORE" {
				for key, value := range data.Storage {
					if _, ok := initialStorage[key]; !ok {
						initialStorage[key] = value
					}
					finalStorage[key] = value
				}
			}
		}
		if len(initialStorage) == 0 {
			fmt.Println("No storage changes")
		}
		for key, value := range initialStorage {
			if value != finalStorage[key] {
				fmt.Println("----------------------------------------------------")
				fmt.Println(key)
				fmt.Println(value, "->", finalStorage[key])
			}
		}
	},
}

func init() {
	EthCmd.AddCommand(storageChangedCmd)

	storageChangedCmd.Args = cobra.ExactArgs(1)
	storageChangedCmd.Flags().StringVarP(&blockchainURL, "blockchainURL", "b", "http://localhost:8545", "Blockchain URL")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downloadSourcesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downloadSourcesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
