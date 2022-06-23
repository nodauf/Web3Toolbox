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
	"fmt"
	utilsDownloadSources "github.com/nodauf/web3Toolbox/utilsCmd/downloadSources"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var apiKey string
var apiDomain string
var contractAddress string

// downloadSourcesCmd represents the downloadSources command
var downloadSourcesCmd = &cobra.Command{
	Use:   "downloadSources contractAddress [...contractAddress]",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		sources := utilsDownloadSources.DownloadSources(apiDomain, apiKey, contractAddress)
		directoryParent := "downloadedSourceCode/" + contractAddress + "/" + sources.Result[0].ContractName
		os.MkdirAll(directoryParent, 0755)
		fmt.Println(contractAddress + " - " + sources.Result[0].ContractName)
		// If the struct is populated that means there may have multiple contracts and etherscan provides much information
		if len(sources.Result[0].SourceCode.Sources) != 0 {
			for contractName, subContracts := range sources.Result[0].SourceCode.Sources {

				directory := filepath.Dir(contractName)
				name := filepath.Base(contractName)
				os.MkdirAll(directoryParent+"/"+directory, 0755)
				err := ioutil.WriteFile(directoryParent+"/"+directory+"/"+name, []byte(subContracts.Content), os.ModePerm)
				if err != nil {
					log.Fatal("writefile1 ", err)
				}
			}
			// Otherwise we only have one contract string
		} else {
			err := ioutil.WriteFile(directoryParent+"/"+sources.Result[0].ContractName+".sol", []byte(sources.Result[0].SourceCodeString), os.ModePerm)
			if err != nil {
				log.Fatal("writefile2 ", err)
			}
		}
	},
}

func init() {
	EthCmd.AddCommand(downloadSourcesCmd)
	downloadSourcesCmd.Flags().StringVarP(&apiKey, "apiKey", "k", "", "API key")
	downloadSourcesCmd.Flags().StringVarP(&apiDomain, "apiDomain", "d", "https://api.etherscan.io", "API domain")
	downloadSourcesCmd.Flags().StringVarP(&contractAddress, "contractAddress", "c", "", "Contract address")

	//downloadSourcesCmd.Args = cobra.MinimumNArgs(1)
	downloadSourcesCmd.MarkFlagRequired("apiKey")
	downloadSourcesCmd.MarkFlagRequired("contractAddress")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downloadSourcesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downloadSourcesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
