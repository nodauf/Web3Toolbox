package utilscmd

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
	"github.com/nodauf/solc-go"
	"github.com/nodauf/web3Toolbox/utilsCmd/contract2Interface"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var contractFile string

// contract2InterfaceCmd represents the contract2Interface command
var contract2InterfaceCmd = &cobra.Command{
	Use:   "contract2Interface",
	Short: "Convert a solidity contract to an interface",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		contractData, err := os.ReadFile(contractFile)
		if err != nil {
			log.Fatal("Fail to open file", err)
		}
		downloadedFile := contract2Interface.DownloadSolcJS(contract2Interface.GetSolidityVersion(contractData))
		compiler, err := solc.NewFromFile(downloadedFile)
		if err != nil {
			log.Fatal(err)
		}
		input := &solc.Input{
			Language: "Solidity",
			Sources: map[string]solc.SourceIn{
				"target.sol": solc.SourceIn{Content: string(contractData)},
			},
			Settings: solc.Settings{
				EVMVersion: "byzantium",
				OutputSelection: map[string]map[string][]string{
					"*": map[string][]string{
						"*": []string{
							//"abi",
							//"evm.bytecode.object",
							//"evm.bytecode.sourceMap",
							//"evm.deployedBytecode.object",
							//"evm.deployedBytecode.sourceMap",
							//"evm.methodIdentifiers",
						},
						"": []string{
							"ast",
						},
					},
				},
			},
		}
		output, err := compiler.Compile(input)
		if err != nil {
			log.Fatal(err)
		}
		//x, _ := json.MarshalIndent(output.Sources["target.sol"].AST, "", "	")
		//fmt.Printf("%s ", x)
		fmt.Println(contract2Interface.InterfaceFromAST(output.Sources["target.sol"].AST))
		//var interface string

	},
}

func init() {
	UtilsCmd.AddCommand(contract2InterfaceCmd)
	contract2InterfaceCmd.Flags().StringVarP(&contractFile, "file", "f", "", "contract file")
	contract2InterfaceCmd.MarkFlagRequired("file")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// contract2InterfaceCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// contract2InterfaceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
